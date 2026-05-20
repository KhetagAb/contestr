package regatta

import (
	"contestr/pkg/logger"
	"contestr/pkg/regatta"
	"contestr/pkg/util"
	"context"
	"fmt"
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

		ratedParticipants := make([]Participant, 0, len(participantsMap))
		for id := range participantsMap {
			ratedParticipants = append(ratedParticipants, id)
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
