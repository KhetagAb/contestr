package regatta

import (
	"contestr/pkg/regatta"
	"context"
	"fmt"
	"slices"
)

type TourRepository interface {
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

func (r *Regatta) GetContestResult(ctx context.Context, contestID int) (regatta.ContestStandings, error) {
	parsedContest, err := r.ejudgeParser.FetchAndParseXML(ctx, contestID)
	if err != nil {
		return regatta.ContestStandings{}, fmt.Errorf("failed to parse contest %d: %w", contestID, err)
	}

	displayNameByParticipant := getDisplayNameByParticipant(parsedContest.Users)

	tours, err := r.tourRepository.FindByContestID(ctx, contestID)
	if err != nil {
		return regatta.ContestStandings{}, fmt.Errorf("failed to find tours for contest %d: %w", contestID, err)
	}

	contestRows := []regatta.ContestRow{}
	contestStandingsByParticipants := make(ResultsByParticipant)
	participantTotal := make(map[Participant]int)
	groupNumbers := make(map[Participant]int)

	for _, tour := range tours {
		result := CalculateResult(tour, parsedContest.Runs.Runs).Export()

		for participant, participantResult := range result {
			_, was := contestStandingsByParticipants[participant]
			if !was {
				contestStandingsByParticipants[participant] = make(ParticipantResult)
			}

			for problem, score := range participantResult {
				contestStandingsByParticipants[participant][problem] = score
				participantTotal[participant] += score
			}
		}

		groupNumbers = make(map[Participant]int)
		for number, group := range tour.Groups {
			for _, participant := range group {
				groupNumbers[participant] = number + 1
			}
		}
	}

	for participant, participantResult := range contestStandingsByParticipants {
		contestRows = append(contestRows, regatta.ContestRow{
			DisplayName:    displayNameByParticipant[participant],
			ProblemResults: participantResult,
			SolvedProblems: len(participantResult),
			TeamNumber:     groupNumbers[participant],
			TotalScore:     participantTotal[participant],
		})
	}

	slices.SortFunc(contestRows, func(row1, row2 regatta.ContestRow) int {
		return row1.TotalScore - row2.TotalScore
	})

	contestStandings := regatta.ContestStandings{
		ContestId:   contestID, // TODO parse int
		ContestName: parsedContest.Name,
		Rows:        contestRows,
	}

	return contestStandings, nil
}

func getDisplayNameByParticipant(users regatta.Users) map[Participant]string {
	displayNameByParticipant := make(map[Participant]string)
	for _, user := range users.Users {
		displayNameByParticipant[user.ID] = user.Name
	}
	return displayNameByParticipant
}
