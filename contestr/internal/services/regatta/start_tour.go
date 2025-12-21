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
	if contestId <= 0 {
		return "", fmt.Errorf("invalid contest ID: %d", contestId)
	}
	if duration <= 0 {
		return "", fmt.Errorf("invalid duration: %v", duration)
	}

	logger.Infof(ctx, "Starting tour on contest %d, duration=%v", contestId, duration)

	tours, err := s.tourRepository.FindByContestID(ctx, contestId)
	if err != nil {
		return "", fmt.Errorf("failed to find tours for contest %d: %w", contestId, err)
	}

	contestStandings, err := s.GetContestResult(ctx, contestId)
	if err != nil {
		return "", fmt.Errorf("failed to get contest standings: %w", err)
	}

	var ratedParticipants []Participant
	for _, row := range contestStandings.Rows {
		ratedParticipants = append(ratedParticipants, row.UserID)
	}

	groups := util.FormGroups(ratedParticipants, GroupSize)
	tourIdx := len(tours) + 1

	startTourInSecondsFromStart := int(contestStandings.CurrentTime.Sub(contestStandings.ContestStartTime).Seconds())
	endTourInSecondsFromStart := int((contestStandings.CurrentTime.Sub(contestStandings.ContestStartTime) + duration).Seconds())
	logger.Infof(ctx, "Tour from start: %v, ends: %v (contest start: %v)", startTourInSecondsFromStart, endTourInSecondsFromStart, contestStandings.ContestStartTime)
	tour := regatta.Tour{
		Name:              fmt.Sprintf("Tour №%v of contest %v", tourIdx, contestId),
		Index:             tourIdx,
		StarTime:          startTourInSecondsFromStart,
		EndTime:           endTourInSecondsFromStart,
		DurationInSeconds: int(duration.Seconds()),
		Groups:            ConvertGroups(groups),
		GroupSize:         GroupSize,
		Problems:          []int{2*tourIdx - 1, 2 * tourIdx},
		ContestID:         contestId,
		GroupNumbers:      regatta.ParticipantsToGroupNumbersMapping(groups),
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
