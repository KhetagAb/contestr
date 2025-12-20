package codeforces

import (
	"contestr/pkg/regatta"
	"context"
	"fmt"
	"time"
)

type CodeforcesAdapter struct {
	service *Service
}

func NewCodeforcesAdapter(service *Service) *CodeforcesAdapter {
	return &CodeforcesAdapter{
		service: service,
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

	startTime := time.Unix(standings.Contest.StartTimeSeconds, 0)

	participants := make([]regatta.ContestParticipant, 0)
	participantMap := make(map[string]bool)

	submissions := make([]regatta.ContestSubmission, 0)

	for _, row := range standings.Rows {
		var handle string
		var displayName string

		if len(row.Party.Members) > 0 {
			handle = row.Party.Members[0].Handle
			displayName = row.Party.Members[0].Handle
		} else {
			handle = fmt.Sprintf("team_%d", contestID)
			displayName = "Team"
		}

		if !participantMap[handle] {
			participants = append(participants, regatta.ContestParticipant{
				ID:          handle,
				DisplayName: displayName,
				OriginalID:  handle,
			})
			participantMap[handle] = true
		}

		for i, problemResult := range row.ProblemResults {
			if problemResult.Points > 0 {
				problemID := i + 1
				submissionTime := 0
				if problemResult.BestSubmissionTimeSeconds > 0 {
					submissionTime = int(problemResult.BestSubmissionTimeSeconds)
				}

				problemIndex := ""
				if len(standings.Problems) > i {
					problemIndex = standings.Problems[i].Index
				} else {
					problemIndex = string(rune('A' + i))
				}

				submissions = append(submissions, regatta.ContestSubmission{
					ParticipantID:     handle,
					ProblemID:         problemID,
					Time:              submissionTime,
					Status:            "OK",
					OriginalProblemID: problemIndex,
				})
			}
		}
	}

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
