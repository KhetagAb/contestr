package ejudge

import (
	"contestr/internal/integrations"
	"contestr/pkg/regatta"
	"context"
	"fmt"
	"strconv"
	"time"
)

type EjudgeAdapter struct {
	fetcher *ContestXMLFetcher
}

func NewEjudgeAdapter(fetcher *ContestXMLFetcher) *EjudgeAdapter {
	return &EjudgeAdapter{
		fetcher: fetcher,
	}
}

func (a *EjudgeAdapter) GetSystem() string {
	return "ejudge"
}

func (a *EjudgeAdapter) FetchContest(ctx context.Context, contestID int, opts integrations.FetchContestOptions) (*regatta.Contest, error) {
	runLog, err := a.fetcher.FetchAndParseXML(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ejudge contest: %w", err)
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", runLog.StartTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start time: %w", err)
	}

	contestIDInt, err := strconv.Atoi(runLog.ContestID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse contest ID: %w", err)
	}

	participants := make([]regatta.ContestParticipant, 0, len(runLog.Users.Users))
	for _, user := range runLog.Users.Users {
		participants = append(participants, regatta.ContestParticipant{
			ID:          fmt.Sprintf("ejudge:%d", user.ID),
			DisplayName: user.Name,
			OriginalID:  strconv.Itoa(user.ID),
		})
	}

	submissions := make([]regatta.ContestSubmission, 0, len(runLog.Runs.Runs))
	for _, run := range runLog.Runs.Runs {
		submissions = append(submissions, regatta.ContestSubmission{
			ParticipantID:     fmt.Sprintf("ejudge:%d", run.UserID),
			ProblemID:         run.ProbID,
			Time:              run.Time,
			Status:            run.Status,
			Points:            pointsForEjudgeStatus(run.Status),
			OriginalProblemID: strconv.Itoa(run.ProbID),
		})
	}

	return &regatta.Contest{
		ContestID:       contestIDInt,
		ContestName:     runLog.Name,
		System:          "ejudge",
		StartTime:       startTime,
		LastUpdated:     time.Now(),
		ScoringSettings: regatta.NormalizeScoringSettings(opts.ScoringSettings),
		TourSettings:    regatta.NormalizeTourSettings(opts.TourSettings),
		Participants:    participants,
		Submissions:     submissions,
	}, nil
}

func pointsForEjudgeStatus(status string) int {
	if status == "OK" {
		return 100
	}
	return 0
}
