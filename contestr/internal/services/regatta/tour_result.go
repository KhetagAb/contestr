package regatta

import (
	"contestr/pkg/regatta"
	"time"
)

const (
	OvertakePoints    = 5
	SolvePoints       = 5
	SolveInTimePoints = 5
)

const (
	SubmissionStatusOK string = "OK"

	GiveBonusOnlyForFirstSubmission bool = true
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
	group := t.Groups[participant]
	score := 0
	participantResults := t.Results[participant]
	result := make(ParticipantResult, len(t.Problems))

	for _, problem := range t.Problems {
		problemCode := t.ProblemsMapping[problem]

		participantSolveTime, participantSolved := participantResults[problem]
		if participantSolved {
			score += SolvePoints
			if time.Duration(participantSolveTime)*time.Second < t.Duration {
				score += SolveInTimePoints
			}
		}

		overtookCount := 0
		for _, opponent := range group {
			if participant == opponent {
				continue
			}
			opponentResults := t.Results[opponent]
			opponentSolveTime, opponentSolved := opponentResults[problem]

			if participantSolved {
				if !opponentSolved || participantSolveTime < opponentSolveTime {
					overtookCount++
				}
			}
		}

		if GiveBonusOnlyForFirstSubmission {
			if overtookCount == len(group)-1 {
				score += OvertakePoints
			}
		} else {
			score += OvertakePoints * overtookCount
		}

		result[problemCode] = ProblemResult{
			problemCode:        problemCode,
			score:              score,
			lastSubmissionTime: participantSolveTime,
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
