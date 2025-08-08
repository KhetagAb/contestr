package regatta

import (
	"contestr/pkg/regatta"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"slices"
	"strings"
	"time"
)

type TourRepository interface {
	Create(ctx context.Context, tour *regatta.Tour) (primitive.ObjectID, error)
	FindByContestID(ctx context.Context, contestID int) ([]regatta.Tour, error)
}

type EjudgeParser interface {
	FetchAndParseXML(ctx context.Context, contestId int) (*regatta.RunLog, error)
}

type Regatta struct {
	tourRepository TourRepository
	ejudgeParser   EjudgeParser
}

func NewRegatta(
	tourRepository TourRepository,
	ejudgeParser EjudgeParser,
) *Regatta {
	return &Regatta{
		tourRepository: tourRepository,
		ejudgeParser:   ejudgeParser,
	}
}

// TODO рефакторинг
func (s *Regatta) GetContestResult(ctx context.Context, contestID int) (regatta.ContestStandings, error) {
	parsedContest, err := s.ejudgeParser.FetchAndParseXML(ctx, contestID)
	if err != nil {
		return regatta.ContestStandings{}, fmt.Errorf("failed to parse contest %d: %w", contestID, err)
	}

	displayNameByParticipant := getDisplayNameByParticipant(parsedContest.Users)

	tours, err := s.tourRepository.FindByContestID(ctx, contestID)
	if err != nil {
		return regatta.ContestStandings{}, fmt.Errorf("failed to find tours for contest %d: %w", contestID, err)
	}

	if len(tours) == 0 {
		contestRows := []regatta.ContestRow{}
		for _, participant := range parsedContest.Users.Users {
			contestRows = append(contestRows, regatta.ContestRow{
				DisplayName: displayNameByParticipant[participant.ID],
				UserID:      participant.ID,
			})
		}
		return regatta.ContestStandings{
			ContestId:           contestID, // TODO parse int
			ContestName:         parsedContest.Name,
			CurrentTourDuration: math.MaxInt,
			Rows:                contestRows,
		}, nil
	}

	var contestRows []regatta.ContestRow
	contestStandingsByParticipants := make(ResultsByParticipant)
	participantTotal := make(map[Participant]int)

	for _, tour := range tours {
		result := CalculateResult(tour, parsedContest.Runs.Runs).Export()

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
		contestRows = append(contestRows, regatta.ContestRow{
			DisplayName:    displayNameByParticipant[participant],
			ProblemResults: getProblemResults(participantResult),
			SolvedProblems: getSolvedProblemsCount(participantResult),
			TeamNumber:     tours[len(tours)-1].GroupNumbers[participant],
			TotalScore:     participantTotal[participant],
			UserID:         participant,
		})
	}

	slices.SortFunc(contestRows, func(row1, row2 regatta.ContestRow) int {
		if row1.TotalScore == row2.TotalScore {
			return strings.Compare(row1.DisplayName, row2.DisplayName)
		}
		return row1.TotalScore - row2.TotalScore
	})

	contestStandings := regatta.ContestStandings{
		ContestId:            contestID, // TODO parse int
		ContestName:          parsedContest.Name,
		Rows:                 contestRows,
		CurrentTourStartTime: tours[len(tours)-1].StartTime,
		CurrentTourDuration:  tours[len(tours)-1].Duration,
	}

	return contestStandings, nil
}

func getProblemResults(result ParticipantResult) []regatta.ProblemResult {
	var results []regatta.ProblemResult
	for _, problemResult := range result {
		results = append(results, regatta.ProblemResult{
			ProblemCode:        problemResult.problemCode,
			Score:              problemResult.score,
			LastSubmissionTime: time.Unix(int64(problemResult.lastSubmissionTime), 0),
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

func getDisplayNameByParticipant(users regatta.Users) map[Participant]string {
	displayNameByParticipant := make(map[Participant]string)
	for _, user := range users.Users {
		displayNameByParticipant[user.ID] = user.Name
	}
	return displayNameByParticipant
}
