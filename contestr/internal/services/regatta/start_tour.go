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

	sequence := len(tours) + 1
	round := 0
	var groups map[Participant]Group
	var groupNumbers map[Participant]int
	var problems []int
	groupSize := 0

	if !opts.IsPause {
		contest, err := s.contestRepo.GetByContestID(ctx, contestId)
		if err != nil {
			return "", fmt.Errorf("failed to get contest: %w", err)
		}
		participantsMap := contestParticipantsMap(contest.Participants)
		if len(participantsMap) == 0 {
			return "", fmt.Errorf("contest %d has no participants", contestId)
		}
		tourSettings := regatta.NormalizeTourSettings(contest.TourSettings)

		ratedParticipants, err := s.participantsOrderedByRating(
			ctx,
			contestId,
			participantsMap,
			tours,
			contest,
		)
		if err != nil {
			return "", err
		}

		formed := util.FormGroupsWithSwapProbability(
			ratedParticipants,
			tourSettings.GroupSize,
			tourSettings.GroupShuffleProbability(),
		)
		groups = ConvertGroups(formed)
		groupNumbers = regatta.ParticipantsToGroupNumbersMapping(formed)
		groupSize = tourSettings.GroupSize

		round = regatta.CompetitiveRoundCount(tours) + 1
		problems = nextTourProblems(tours, tourSettings.ProblemsPerTour)
	}

	name := fmt.Sprintf("Tour №%d of contest %d", round, contestId)
	if opts.IsPause {
		name = fmt.Sprintf("Pause of contest %d (seq %d)", contestId, sequence)
	}

	tour := regatta.Tour{
		Name:              name,
		Sequence:          sequence,
		Round:             round,
		IsPause:           opts.IsPause,
		DurationInSeconds: durationSeconds,
		Groups:            groups,
		GroupSize:         groupSize,
		Problems:          problems,
		ContestID:         contestId,
		GroupNumbers:      groupNumbers,
	}

	create, err := s.tourRepository.Create(ctx, &tour)
	if err != nil {
		return "", fmt.Errorf("failed to create tour: %w", err)
	}

	return create.Hex(), nil
}

// participantsOrderedByRating returns participant IDs sorted by total score before the new tour
// (descending). FormGroups slices this list into buckets — order must reflect standings.
func (s *Regatta) participantsOrderedByRating(
	ctx context.Context,
	contestID int,
	participantsMap map[string]string,
	completedTours []regatta.Tour,
	contest *regatta.Contest,
) ([]Participant, error) {
	ids := make([]Participant, 0, len(participantsMap))
	for id := range participantsMap {
		ids = append(ids, id)
	}

	totals := make(map[Participant]int, len(ids))
	for _, id := range ids {
		totals[id] = 0
	}

	// Первый тур: у всех 0 очков, сортировка по имени. Дальше — по сумме очков за прошлые туры.
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

	slices.SortFunc(ids, func(a, b Participant) int {
		if totals[a] != totals[b] {
			return totals[b] - totals[a]
		}
		nameA := participantsMap[a]
		nameB := participantsMap[b]
		if c := strings.Compare(nameA, nameB); c != 0 {
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

func ConvertGroups(groups [][]string) map[Participant]Group {
	result := make(map[Participant]Group)

	for _, group := range groups {
		for _, participantID := range group {
			groupCopy := make([]string, len(group))
			copy(groupCopy, group)
			result[participantID] = groupCopy
		}
	}

	return result
}
