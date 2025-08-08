package regatta

import (
	"contestr/pkg/regatta"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	startTime, err := time.Parse("2006-01-02 15:04:05", parsedContest.StartTime)
	if err != nil {
		return regatta.ContestStandings{}, fmt.Errorf("failed to parse contest time %d: %w", contestID, err)
	}

	currentTime, err := time.Parse("2006-01-02 15:04:05", parsedContest.CurrentTime)
	if err != nil {
		return regatta.ContestStandings{}, fmt.Errorf("failed to parse contest time %d: %w", contestID, err)
	}

	displayNameByParticipant := getDisplayNameByParticipant(parsedContest.Users)

	tours, err := s.tourRepository.FindByContestID(ctx, contestID)
	if err != nil {
		return regatta.ContestStandings{}, fmt.Errorf("failed to find tours for contest %d: %w", contestID, err)
	}

	standings := regatta.ContestStandings{
		ContestId:        contestID, // TODO parse int
		ContestName:      parsedContest.Name,
		CurrentTime:      currentTime,
		ContestStartTime: startTime,
		Rows:             []regatta.ContestRow{},
	}

	if len(tours) == 0 {
		contestRows := []regatta.ContestRow{}
		for _, participant := range parsedContest.Users.Users {
			contestRows = append(contestRows, regatta.ContestRow{
				DisplayName: displayNameByParticipant[participant.ID],
				UserID:      participant.ID,
			})
		}
		standings.Rows = contestRows
		return standings, nil
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

	// func sumTimes(row: regatta.ContestRow) {
	// 	// for _, r := range row.ProblemResults {
	// 	// 	max(time1, time2)

	// 	// }
	// }

	slices.SortFunc(contestRows, func(row1, row2 regatta.ContestRow) int {
		if row1.TotalScore != row2.TotalScore {
			return row2.TotalScore - row1.TotalScore
		} else {
			// TODO: take sum of times or time of last sumbit
		}
		return strings.Compare(row1.DisplayName, row2.DisplayName)
	})
	standings.Rows = contestRows

	return standings, nil
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

func getDisplayNameByParticipant(users regatta.Users) map[Participant]string {
	displayNameByParticipant := make(map[Participant]string)
	for _, user := range users.Users {
		displayNameByParticipant[user.ID] = user.Name
	}
	return displayNameByParticipant
}
