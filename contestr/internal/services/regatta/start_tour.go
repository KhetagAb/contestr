package regatta

import (
	"contestr/pkg/logger"
	"contestr/pkg/regatta"
	"contestr/pkg/util"
	"context"
	"fmt"
	"time"
)

const GroupSize = 3

func (s *Regatta) StartTour(ctx context.Context, contestId int, duration time.Duration) (string, error) {
	tours, err := s.tourRepository.FindByContestID(ctx, contestId)
	if err != nil {
		return "", fmt.Errorf("failed to find tours for contest %d: %w", contestId, err)
	}

	contestStandings, err := s.GetContestResult(ctx, contestId)
	if err != nil {
		return "", fmt.Errorf("failed to get contest standings: %w", err)
	}

	logger.Infof(ctx, "starting tour for contest %d", contestId)

	var ratedParticipants []Participant
	for _, row := range contestStandings.Rows {
		logger.Infof(ctx, "participant %s has %d points", row.DisplayName, row.TotalScore)
		ratedParticipants = append(ratedParticipants, row.UserID)
	}

	groups := util.FormGroups(ratedParticipants, GroupSize)
	tourIdx := len(tours) + 1
	convertedGroups := ConvertGroups(groups)

	tour := regatta.Tour{
		Name:         fmt.Sprintf("Tour №%v of contest %v", tourIdx, contestId),
		Index:        tourIdx,
		StartTime:    time.Now(),
		Duration:     duration,
		Groups:       convertedGroups,
		GroupSize:    GroupSize,
		Problems:     []int{2*tourIdx - 1, 2 * tourIdx},
		ContestID:    contestId,
		GroupNumbers: regatta.ParticipantsToGroupNumbersMapping(convertedGroups),
	}

	create, err := s.tourRepository.Create(ctx, &tour)
	if err != nil {
		return "", fmt.Errorf("failed to create tour: %w", err)
	}

	return create.Hex(), nil
}

func ConvertGroups(groups [][]int) map[Participant]Group {
	result := make(map[Participant]Group)

	for _, group := range groups {
		for _, participantID := range group {
			result[participantID] = group
		}
	}

	return result
}
