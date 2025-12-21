package codeforces

import (
	"contestr/internal/repository"
	"contestr/pkg/logger"
	"contestr/pkg/regatta"
	"context"
	"fmt"
	"time"
)

type CodeforcesAdapter struct {
	service   *Service
	handleRepo repository.CodeforcesHandleRepository
}

func NewCodeforcesAdapter(service *Service, handleRepo repository.CodeforcesHandleRepository) *CodeforcesAdapter {
	return &CodeforcesAdapter{
		service:    service,
		handleRepo: handleRepo,
	}
}

func (a *CodeforcesAdapter) GetSystem() string {
	return "codeforces"
}

func (a *CodeforcesAdapter) FetchContest(ctx context.Context, contestID int) (*regatta.Contest, error) {
	standings, err := a.service.GetContest(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch codeforces contest: %w", err)
	}

	logger.Infof(ctx, "[CF] Fetched contest %d: %s, rows count: %d", contestID, standings.Contest.Name, len(standings.Rows))

	handleMappings, err := a.handleRepo.GetAllByContestID(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get handle mappings: %w", err)
	}

	logger.Infof(ctx, "[CF] Loaded handle mappings for contest %d: %d mappings", contestID, len(handleMappings))
	for handle, name := range handleMappings {
		logger.Infof(ctx, "[CF] Mapping: handle=%s, name=%s", handle, name)
	}

	startTime := time.Unix(standings.Contest.StartTimeSeconds, 0)

	participants := make([]regatta.ContestParticipant, 0)
	allowedHandles := make(map[string]bool)

	for handle, mappedName := range handleMappings {
		displayName := handle
		if mappedName != "" {
			displayName = mappedName
		}
		participants = append(participants, regatta.ContestParticipant{
			ID:          handle,
			DisplayName: displayName,
			OriginalID:  handle,
		})
		allowedHandles[handle] = true
	}

	logger.Infof(ctx, "[CF] Created %d participants from handle mappings", len(participants))

	submissions := make([]regatta.ContestSubmission, 0)

	for _, row := range standings.Rows {
		var handle string
		if len(row.Party.Members) > 0 {
			handle = row.Party.Members[0].Handle
		} else {
			handle = fmt.Sprintf("team_%d", contestID)
		}

		if !allowedHandles[handle] {
			continue
		}

		for i, problemResult := range row.ProblemResults {
			problemID := i + 1
			problemIndex := ""
			if len(standings.Problems) > i {
				problemIndex = standings.Problems[i].Index
			} else {
				problemIndex = string(rune('A' + i))
			}

			if problemResult.Points > 0 {
				submissionTime := 0
				if problemResult.BestSubmissionTimeSeconds > 0 {
					submissionTime = int(problemResult.BestSubmissionTimeSeconds)
				}

				submissions = append(submissions, regatta.ContestSubmission{
					ParticipantID:     handle,
					ProblemID:         problemID,
					Time:              submissionTime,
					Status:            "OK",
					OriginalProblemID: problemIndex,
				})
			}

			if problemResult.RejectedAttemptCount > 0 {
				for j := 0; j < int(problemResult.RejectedAttemptCount); j++ {
					submissions = append(submissions, regatta.ContestSubmission{
						ParticipantID:     handle,
						ProblemID:         problemID,
						Time:              0,
						Status:            "WRONG_ANSWER",
						OriginalProblemID: problemIndex,
					})
				}
			}
		}
	}

	logger.Infof(ctx, "[CF] Parsed %d submissions from standings", len(submissions))


	logger.Infof(ctx, "[CF] Final participants count: %d, submissions count: %d", len(participants), len(submissions))

	return &regatta.Contest{
		ContestID:    int(standings.Contest.ID),
		ContestName:  standings.Contest.Name,
		System:       "codeforces",
		StartTime:    startTime,
		LastUpdated:  time.Now(),
		Participants: participants,
		Submissions:  submissions,
	}, nil
}
