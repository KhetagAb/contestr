package codeforces

import (
	"contestr/internal/integrations"
	"contestr/internal/repository"
	"contestr/pkg/logger"
	"contestr/pkg/regatta"
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/togatoga/goforces"
)

// normalizeCFHandle strips gym-group prefix from handles (e.g. "g44870=user" -> "user").
func normalizeCFHandle(handle string) string {
	if i := strings.LastIndex(handle, "="); i >= 0 {
		return handle[i+1:]
	}
	return handle
}

type CodeforcesAdapter struct {
	service    *Service
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

func (a *CodeforcesAdapter) FetchContest(ctx context.Context, contestID int, opts integrations.FetchContestOptions) (*regatta.Contest, error) {
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

	participants, allowedHandles := buildParticipants(contestID, standings.Rows, handleMappings)
	logger.Infof(ctx, "[CF] Created %d participants", len(participants))

	settings := regatta.NormalizeScoringSettings(opts.ScoringSettings)
	submissions := make([]regatta.ContestSubmission, 0)

	if settings.Mode == regatta.ScoringModePartial {
		submissions, err = a.fetchPartialSubmissions(
			ctx,
			contestID,
			standings.Rows,
			standings.Problems,
			allowedHandles,
		)
		if err != nil {
			return nil, err
		}
	} else {
		submissions = a.buildBinarySubmissions(contestID, standings.Rows, standings.Problems, allowedHandles)
	}

	logger.Infof(ctx, "[CF] Parsed %d submissions", len(submissions))

	logger.Infof(ctx, "[CF] Final participants count: %d, submissions count: %d", len(participants), len(submissions))

	return &regatta.Contest{
		ContestID:       int(standings.Contest.ID),
		ContestName:     standings.Contest.Name,
		System:          "codeforces",
		StartTime:       startTime,
		LastUpdated:     time.Now(),
		ScoringSettings: settings,
		TourSettings:    regatta.NormalizeTourSettings(opts.TourSettings),
		Participants:    participants,
		Submissions:     submissions,
	}, nil
}

func buildParticipants(contestID int, rows []goforces.RanklistRow, handleMappings map[string]string) ([]regatta.ContestParticipant, map[string]bool) {
	participants := make([]regatta.ContestParticipant, 0)
	allowedHandles := make(map[string]bool)

	if len(handleMappings) > 0 {
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
		return participants, allowedHandles
	}

	seen := make(map[string]bool)
	for _, row := range rows {
		rawHandle := rawRowHandle(contestID, row)
		handle := normalizeCFHandle(rawHandle)
		if handle == "" || seen[handle] {
			continue
		}
		participants = append(participants, regatta.ContestParticipant{
			ID:          handle,
			DisplayName: handle,
			OriginalID:  rawHandle,
		})
		allowedHandles[handle] = true
		seen[handle] = true
	}

	return participants, allowedHandles
}

func rawRowHandle(contestID int, row goforces.RanklistRow) string {
	if len(row.Party.Members) > 0 {
		return row.Party.Members[0].Handle
	}
	if row.Rank > 0 {
		return fmt.Sprintf("team_%d_%d", contestID, row.Rank)
	}
	return fmt.Sprintf("team_%d", contestID)
}

func (a *CodeforcesAdapter) buildBinarySubmissions(contestID int, rows []goforces.RanklistRow, problems []goforces.Problem, allowedHandles map[string]bool) []regatta.ContestSubmission {
	submissions := make([]regatta.ContestSubmission, 0)
	for _, row := range rows {
		handle := normalizeCFHandle(rawRowHandle(contestID, row))

		if !allowedHandles[handle] {
			continue
		}

		for i, problemResult := range row.ProblemResults {
			problemID := i + 1
			problemIndex := ""
			if len(problems) > i {
				problemIndex = problems[i].Index
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
					Points:            100,
					OriginalProblemID: problemIndex,
				})
			}

			if problemResult.RejectedAttemptCount > 0 {
				for j := 0; j < int(problemResult.RejectedAttemptCount); j++ {
					submissions = append(submissions, regatta.ContestSubmission{
						ParticipantID:     handle,
						ProblemID:         problemID,
						Time:              0,
						Status:            "WRONG_ANSWER", // TODO А почему тут только WA?
						OriginalProblemID: problemIndex,
					})
				}
			}
		}
	}
	return submissions
}

func (a *CodeforcesAdapter) fetchPartialSubmissions(
	ctx context.Context,
	contestID int,
	rows []goforces.RanklistRow,
	problems []goforces.Problem,
	allowedHandles map[string]bool,
) ([]regatta.ContestSubmission, error) {
	statusSubmissions, err := a.service.GetContestStatusWithPoints(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch codeforces status submissions: %w", err)
	}

	problemIDs := make(map[string]int, len(problems))
	for i, problem := range problems {
		problemIDs[problem.Index] = i + 1
	}

	submissions := make([]regatta.ContestSubmission, 0, len(statusSubmissions))
	for _, sub := range statusSubmissions {
		if len(sub.Author.Members) == 0 {
			continue
		}
		handle := normalizeCFHandle(sub.Author.Members[0].Handle)
		if !allowedHandles[handle] {
			continue
		}

		problemID, ok := problemIDs[sub.Problem.Index]
		if !ok {
			continue
		}

		points := 0
		if sub.Points != nil {
			// Codeforces IOI-style submissions expose scored points. Current scoring assumes
			// contest problems are normalized to 0..100; other ranges need per-problem max support.
			points = int(math.Round(*sub.Points))
		} else if sub.Verdict == "OK" {
			points = 100
		}

		if points <= 0 && sub.Verdict != "OK" {
			continue
		}

		submissions = append(submissions, regatta.ContestSubmission{
			ParticipantID:     handle,
			ProblemID:         problemID,
			Time:              sub.RelativeTimeSeconds,
			Status:            sub.Verdict,
			Points:            points,
			OriginalProblemID: sub.Problem.Index,
		})
	}

	return appendPartialStandingsFallback(contestID, submissions, rows, problems, allowedHandles), nil
}

func appendPartialStandingsFallback(
	contestID int,
	submissions []regatta.ContestSubmission,
	rows []goforces.RanklistRow,
	problems []goforces.Problem,
	allowedHandles map[string]bool,
) []regatta.ContestSubmission {
	bestPoints := make(map[string]int, len(submissions))
	for _, submission := range submissions {
		key := submissionKey(submission.ParticipantID, submission.ProblemID)
		if submission.Points > bestPoints[key] {
			bestPoints[key] = submission.Points
		}
	}

	for _, row := range rows {
		handle := normalizeCFHandle(rawRowHandle(contestID, row))
		if !allowedHandles[handle] {
			continue
		}

		for i, problemResult := range row.ProblemResults {
			points := int(math.Round(problemResult.Points))
			if points <= 0 {
				continue
			}

			problemID := i + 1
			key := submissionKey(handle, problemID)
			if points <= bestPoints[key] {
				continue
			}

			status := "PARTIAL"
			if points >= 100 {
				status = "OK"
			}

			submissions = append(submissions, regatta.ContestSubmission{
				ParticipantID:     handle,
				ProblemID:         problemID,
				Time:              int(problemResult.BestSubmissionTimeSeconds),
				Status:            status,
				Points:            points,
				OriginalProblemID: problemIndexAt(problems, i),
			})
			bestPoints[key] = points
		}
	}

	return submissions
}

func submissionKey(participant string, problemID int) string {
	return fmt.Sprintf("%s:%d", participant, problemID)
}

func problemIndexAt(problems []goforces.Problem, index int) string {
	if len(problems) > index {
		return problems[index].Index
	}
	return string(rune('A' + index))
}
