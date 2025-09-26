package regatta

import (
	"contestr/pkg/logger"
	"contestr/pkg/regatta"
	"context"
)

const (
	OvertakePoints    = 5
	SolvePoints       = 5
	SolveInTimePoints = 5
)

const (
	SubmissionStatusOK string = "OK"

	GiveBonusOnlyForFirstSubmission = true
)

type Problem = int

type Participant = int

type SubmissionTime = int

type ContestResult = map[Problem]SubmissionTime

type Group = []Participant

type TourResult struct {
	regatta.Tour
	Results         map[Participant]ContestResult `json:"-"` // participant -> result
	ProblemsMapping map[Problem]string            `json:"-"`
}

type ProblemResult struct {
	problemCode        string
	score              int
	lastSubmissionTime SubmissionTime
}

type ProblemCode = string

type ParticipantResult = map[ProblemCode]ProblemResult

type ResultsByParticipant = map[Participant]ParticipantResult

func (t *TourResult) ParticipantScore(participant Participant) ParticipantResult {
	logger.Infof(context.Background(), "Calculating participant score for participant %v in %v", participant, t.Name)

	group := t.Groups[participant]
	participantResults := t.Results[participant]
	result := make(ParticipantResult, len(t.Problems))

	for _, problem := range t.Problems {
		score := 0
		problemCode := t.ProblemsMapping[problem]

		participantSolveTime, participantSolved := participantResults[problem]
		if !participantSolved {
			logger.Infof(context.Background(), "Participant %v did not solve problem %v", participant, problem)
			result[problemCode] = ProblemResult{
				problemCode: problemCode,
				score:       0,
			}
			continue
		}

		score += SolvePoints
		logger.Infof(context.Background(), "Participant %v solved problem %v: +5 points", participant, problem)
		if participantSolveTime <= t.EndTime {
			score += SolveInTimePoints
			logger.Infof(context.Background(), "Participant %v solved problem %v in time (%v <= %v): +5 points", participant, problem, participantSolveTime, t.EndTime)
		}

		overtookCount := 0
		for _, opponent := range group {
			if participant == opponent {
				continue
			}
			opponentResults := t.Results[opponent]
			opponentSolveTime, opponentSolved := opponentResults[problem]
			if !opponentSolved || participantSolveTime < opponentSolveTime {
				// FIXME GiveBonusOnlyForInTimeSubmission - регулировать
				overtookCount++
			}
		}
		if GiveBonusOnlyForFirstSubmission {
			if overtookCount == len(group)-1 {
				score += OvertakePoints
				logger.Infof(context.Background(), "Participant %v first-solve problem %v: +5 points", participant, problem)
			}
		} else {
			score += OvertakePoints * overtookCount
			logger.Infof(context.Background(), "Participant %v overtook %v other %v participants: +%v points", participant, overtookCount, overtookCount, OvertakePoints*overtookCount)
		}

		result[problemCode] = ProblemResult{
			problemCode:        problemCode,
			score:              score,
			lastSubmissionTime: participantSolveTime - t.StartTime,
		}
	}

	return result
}

func (t *TourResult) Export() ResultsByParticipant {
	result := make(ResultsByParticipant)

	for participant := range t.Groups {
		result[participant] = t.ParticipantScore(participant)
	}

	return result
}

func CalculateResult(tc regatta.Tour, runs []regatta.Run) *TourResult {
	return &TourResult{
		Tour:            tc,
		Results:         calcSubmissions(runs),
		ProblemsMapping: tc.ProblemsIDsToNameMapping(tc.Problems),
	}
}

func calcSubmissions(submissions []regatta.Run) map[Participant]ContestResult {
	results := make(map[Participant]ContestResult)

	for _, submission := range submissions {
		if submission.Status != SubmissionStatusOK {
			continue
		}

		logger.Infof(context.Background(), "OK! UserID: %v, ProbID: %v, Time: %v", submission.UserID, submission.ProbID, submission.Time)

		participant := submission.UserID
		problem := submission.ProbID

		participantResults, ok := results[participant]
		if !ok {
			results[participant] = map[Problem]int{
				problem: submission.Time,
			}
		} else {
			if _, alreadySolved := participantResults[problem]; !alreadySolved {
				results[participant][problem] = submission.Time
			}
		}
	}

	return results
}
