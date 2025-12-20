package regatta

import (
	"contestr/pkg/regatta"
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourRepository interface {
	Create(ctx context.Context, tour *regatta.Tour) (primitive.ObjectID, error)
	FindByContestID(ctx context.Context, contestID int) ([]regatta.Tour, error)
}

type ContestRepository interface {
	GetByContestID(ctx context.Context, contestID int) (*regatta.Contest, error)
	GetParticipants(ctx context.Context, contestID int) (map[string]string, error)
	GetSubmissions(ctx context.Context, contestID int) ([]regatta.ContestSubmission, error)
}

type Regatta struct {
	tourRepository TourRepository
	contestRepo    ContestRepository
}

func NewRegatta(
	tourRepository TourRepository,
	contestRepo ContestRepository,
) *Regatta {
	return &Regatta{
		tourRepository: tourRepository,
		contestRepo:    contestRepo,
	}
}

func (s *Regatta) GetContestResult(ctx context.Context, contestID int) (regatta.ContestStandings, error) {
	contest, err := s.contestRepo.GetByContestID(ctx, contestID)
	if err != nil {
		return regatta.ContestStandings{}, fmt.Errorf("failed to get contest %d: %w", contestID, err)
	}

	participantsMap, err := s.contestRepo.GetParticipants(ctx, contestID)
	if err != nil {
		return regatta.ContestStandings{}, fmt.Errorf("failed to get participants: %w", err)
	}

	slugToIntMap := makeSlugToIntMapping(contest.Participants)
	intToSlugMap := makeIntToSlugMapping(contest.Participants)

	tours, err := s.tourRepository.FindByContestID(ctx, contestID)
	if err != nil {
		return regatta.ContestStandings{}, fmt.Errorf("failed to find tours for contest %d: %w", contestID, err)
	}

	currentTime := time.Now()
	standings := regatta.ContestStandings{
		ContestId:        contestID,
		ContestName:      contest.ContestName,
		CurrentTime:      currentTime,
		ContestStartTime: contest.StartTime,
		Rows:             []regatta.ContestRow{},
	}

	if len(tours) == 0 {
		contestRows := []regatta.ContestRow{}
		for _, participant := range contest.Participants {
			participantInt := slugToIntMap[participant.ID]
			contestRows = append(contestRows, regatta.ContestRow{
				DisplayName: participant.DisplayName,
				UserID:      participantInt,
			})
		}
		standings.Rows = contestRows
		return standings, nil
	}

	submissions, err := s.contestRepo.GetSubmissions(ctx, contestID)
	if err != nil {
		return regatta.ContestStandings{}, fmt.Errorf("failed to get submissions: %w", err)
	}

	runs := convertSubmissionsToRuns(submissions, slugToIntMap)

	var contestRows []regatta.ContestRow
	contestStandingsByParticipants := make(ResultsByParticipant)
	participantTotal := make(map[Participant]int)

	for _, tour := range tours {
		result := CalculateResult(tour, runs).Export()

		for participant, participantResult := range result {
			_, was := contestStandingsByParticipants[participant]
			if !was {
				contestStandingsByParticipants[participant] = make(ParticipantResult)
			}

			for problem, problemResult := range participantResult {
				contestStandingsByParticipants[participant][problem] = problemResult
				participantTotal[participant] += problemResult.score
			}
		}
	}

	for participant, participantResult := range contestStandingsByParticipants {
		participantSlug := intToSlugMap[participant]
		displayName := participantsMap[participantSlug]
		if displayName == "" {
			displayName = participantSlug
		}

		contestRows = append(contestRows, regatta.ContestRow{
			DisplayName:    displayName,
			ProblemResults: getProblemResults(participantResult),
			SolvedProblems: getSolvedProblemsCount(participantResult),
			TeamNumber:     tours[len(tours)-1].GroupNumbers[participant],
			TotalScore:     participantTotal[participant],
			UserID:         participant,
		})
	}

	slices.SortFunc(contestRows, func(row1, row2 regatta.ContestRow) int {
		if row1.TotalScore != row2.TotalScore {
			return row2.TotalScore - row1.TotalScore
		}
		return strings.Compare(row1.DisplayName, row2.DisplayName)
	})
	standings.Rows = contestRows

	return standings, nil
}

func makeSlugToIntMapping(participants []regatta.ContestParticipant) map[string]Participant {
	result := make(map[string]Participant)
	for idx, participant := range participants {
		result[participant.ID] = idx + 1
	}
	return result
}

func makeIntToSlugMapping(participants []regatta.ContestParticipant) map[Participant]string {
	result := make(map[Participant]string)
	for idx, participant := range participants {
		result[Participant(idx+1)] = participant.ID
	}
	return result
}

func convertSubmissionsToRuns(submissions []regatta.ContestSubmission, slugToIntMap map[string]Participant) []regatta.Run {
	runs := make([]regatta.Run, 0, len(submissions))
	for _, sub := range submissions {
		participantInt, ok := slugToIntMap[sub.ParticipantID]
		if !ok {
			continue
		}
		runs = append(runs, regatta.Run{
			UserID: participantInt,
			ProbID: sub.ProblemID,
			Time:   sub.Time,
			Status: sub.Status,
		})
	}
	return runs
}

func getProblemResults(result ParticipantResult) []regatta.ProblemResult {
	var results []regatta.ProblemResult
	for _, problemResult := range result {
		results = append(results, regatta.ProblemResult{
			ProblemCode:        problemResult.problemCode,
			Score:              problemResult.score,
			LastSubmissionTime: problemResult.lastSubmissionTime,
		})
	}
	slices.SortFunc(results, func(result1, result2 regatta.ProblemResult) int {
		return strings.Compare(result1.ProblemCode, result2.ProblemCode)
	})
	return results
}

func getSolvedProblemsCount(result ParticipantResult) int {
	count := 0
	for _, problemResult := range result {
		if problemResult.score > 0 {
			count++
		}
	}
	return count
}
