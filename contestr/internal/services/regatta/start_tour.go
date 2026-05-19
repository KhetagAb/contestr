package regatta

import (
	"contestr/pkg/logger"
	"contestr/pkg/regatta"
	"contestr/pkg/util"
	"context"
	"fmt"
)

const GroupSize = 3

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
		participantsMap, err := s.contestRepo.GetParticipants(ctx, contestId)
		if err != nil {
			return "", fmt.Errorf("failed to get participants: %w", err)
		}
		if len(participantsMap) == 0 {
			return "", fmt.Errorf("contest %d has no participants", contestId)
		}

		ratedParticipants := make([]Participant, 0, len(participantsMap))
		for id := range participantsMap {
			ratedParticipants = append(ratedParticipants, id)
		}

		formed := util.FormGroups(ratedParticipants, GroupSize)
		groups = ConvertGroups(formed)
		groupNumbers = regatta.ParticipantsToGroupNumbersMapping(formed)
		groupSize = GroupSize

		round = regatta.CompetitiveRoundCount(tours) + 1
		problems = []int{2*round - 1, 2 * round}
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
