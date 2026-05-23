package regatta

import (
	"contestr/pkg/logger"
	"contestr/pkg/regatta"
	"contestr/pkg/util"
	"context"
	"fmt"
	"slices"
	"strings"
)

type StartTourOptions struct {
	IsPause bool
}

// competitiveSetup содержит параметры конкурентного (не паузного) тура.
type competitiveSetup struct {
	round     int
	groups    map[regatta.Participant]regatta.Group
	groupNums map[regatta.Participant]int
	problems  []int
	groupSize int
}

func (s *Regatta) StartTour(ctx context.Context, contestId int, durationSeconds int, opts StartTourOptions) (string, error) {
	if contestId <= 0 {
		return "", fmt.Errorf("invalid contest ID: %d", contestId)
	}
	if durationSeconds <= 0 {
		return "", fmt.Errorf("invalid duration: %d seconds", durationSeconds)
	}

	logger.Infof(ctx, "Starting segment on contest %d, duration=%ds pause=%v", contestId, durationSeconds, opts.IsPause)

	tours, err := s.loadToursSorted(ctx, contestId)
	if err != nil {
		return "", err
	}

	var setup *competitiveSetup
	if !opts.IsPause {
		setup, err = s.buildCompetitiveSetup(ctx, contestId, tours)
		if err != nil {
			return "", err
		}
	}

	tour := buildTourRecord(contestId, durationSeconds, len(tours)+1, opts.IsPause, setup)

	id, err := s.tourRepository.Create(ctx, &tour)
	if err != nil {
		return "", fmt.Errorf("failed to create tour: %w", err)
	}

	return id.Hex(), nil
}

func (s *Regatta) buildCompetitiveSetup(ctx context.Context, contestId int, completedTours []regatta.Tour) (*competitiveSetup, error) {
	contest, err := s.contestRepo.GetByContestID(ctx, contestId)
	if err != nil {
		return nil, fmt.Errorf("failed to get contest: %w", err)
	}

	participantsMap := contestParticipantsMap(contest.Participants)
	if len(participantsMap) == 0 {
		return nil, fmt.Errorf("contest %d has no participants", contestId)
	}

	tourSettings := regatta.NormalizeTourSettings(contest.TourSettings)

	ratedParticipants, err := s.participantsOrderedByRating(ctx, contestId, participantsMap, completedTours, contest)
	if err != nil {
		return nil, err
	}

	formed := util.FormGroupsWithSwapProbability(
		ratedParticipants,
		tourSettings.GroupSize,
		tourSettings.GroupShuffleProbability(),
	)

	return &competitiveSetup{
		round:     regatta.CompetitiveRoundCount(completedTours) + 1,
		groups:    convertGroups(formed),
		groupNums: regatta.ParticipantsToGroupNumbersMapping(formed),
		groupSize: tourSettings.GroupSize,
		problems:  nextTourProblems(completedTours, tourSettings.ProblemsPerTour),
	}, nil
}

func buildTourRecord(contestId, durationSeconds, sequence int, isPause bool, setup *competitiveSetup) regatta.Tour {
	tour := regatta.Tour{
		ContestID:         contestId,
		Sequence:          sequence,
		IsPause:           isPause,
		DurationInSeconds: durationSeconds,
	}
	if isPause {
		tour.Name = fmt.Sprintf("Pause of contest %d (seq %d)", contestId, sequence)
	} else {
		tour.Round = setup.round
		tour.Name = fmt.Sprintf("Tour №%d of contest %d", setup.round, contestId)
		tour.Groups = setup.groups
		tour.GroupNumbers = setup.groupNums
		tour.GroupSize = setup.groupSize
		tour.Problems = setup.problems
	}
	return tour
}

// participantsOrderedByRating возвращает участников по убыванию суммы очков.
// FormGroupsWithSwapProbability использует этот порядок для формирования столов.
func (s *Regatta) participantsOrderedByRating(
	ctx context.Context,
	contestID int,
	participantsMap map[string]string,
	completedTours []regatta.Tour,
	contest *regatta.Contest,
) ([]regatta.Participant, error) {
	ids := make([]regatta.Participant, 0, len(participantsMap))
	for id := range participantsMap {
		ids = append(ids, id)
	}

	totals := make(map[regatta.Participant]int, len(ids))

	// Первый тур: у всех 0 очков, сортировка только по имени.
	if regatta.CompetitiveRoundCount(completedTours) > 0 {
		submissions, err := s.contestRepo.GetSubmissions(ctx, contestID)
		if err != nil {
			return nil, fmt.Errorf("failed to get submissions for grouping: %w", err)
		}
		runs := convertSubmissionsToRuns(submissions)
		scoringSettings := regatta.NormalizeScoringSettings(contest.ScoringSettings)
		offsets := regatta.SegmentOffsets(completedTours)

		for _, tour := range completedTours {
			if tour.IsPause {
				continue
			}
			segmentStart := offsets[tour.Sequence].Start
			result := CalculateResultWithSettings(tour, segmentStart, runs, scoringSettings).Export()
			for participant, problemResults := range result {
				for _, pr := range problemResults {
					if pr.score > 0 {
						totals[participant] += pr.score
					}
				}
			}
		}
	}

	slices.SortFunc(ids, func(a, b regatta.Participant) int {
		if totals[a] != totals[b] {
			return totals[b] - totals[a]
		}
		if c := strings.Compare(participantsMap[a], participantsMap[b]); c != 0 {
			return c
		}
		return strings.Compare(a, b)
	})

	return ids, nil
}

func contestParticipantsMap(participants []regatta.ContestParticipant) map[string]string {
	result := make(map[string]string, len(participants))
	for _, participant := range participants {
		result[participant.ID] = participant.DisplayName
	}
	return result
}

func nextTourProblems(tours []regatta.Tour, count int) []int {
	maxProblem := 0
	for _, tour := range tours {
		for _, problem := range tour.Problems {
			if problem > maxProblem {
				maxProblem = problem
			}
		}
	}
	problems := make([]int, 0, count)
	for i := 1; i <= count; i++ {
		problems = append(problems, maxProblem+i)
	}
	return problems
}

func convertGroups(groups [][]string) map[regatta.Participant]regatta.Group {
	result := make(map[regatta.Participant]regatta.Group)
	for _, group := range groups {
		for _, participantID := range group {
			groupCopy := make([]string, len(group))
			copy(groupCopy, group)
			result[participantID] = groupCopy
		}
	}
	return result
}
