package regatta

import (
	"contestr/internal/lib/ejudge"
	"time"
)

const (
	OVERTAKE_POINTS      = 5
	SOLVE_POINTS         = 5
	SOLVE_IN_TIME_POINTS = 5
)

const (
	SubmissionStatusOK string = "OK"
)

// Problem int ejudge ID
type Problem = int

// Participant int ejudge ID
type Participant = int

// ContestResult problem -> submission_time
type ContestResult = map[Problem]int

// Group participant group
type Group = []Participant

// Tour regatta tour
type Tour struct {
	Name            string                        `json:"name"`
	Index           int                           `json:"index"`
	StartTime       time.Time                     `json:"start_time"`
	Duration        time.Duration                 `json:"duration"`
	Groups          map[Participant]Group         `json:"groups"` // participant -> group
	GroupSize       int                           `json:"group_size"`
	Problems        []Problem                     `json:"problems"`
	ContestID       int                           `json:"contest_id"`
	Results         map[Participant]ContestResult `json:"-"` // participant -> result
	ProblemsMapping map[Problem]string            `json:"-"`
}

type ProblemCode = string

type ParticipantResult = map[ProblemCode]int

type Table = map[Participant]ParticipantResult

func (t *Tour) CalcSubmissions(submissions []ejudge.Submission) map[Participant]ContestResult {
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

// ParticipantScore problemCode -> score
func (t *Tour) ParticipantScore(participant Participant) ParticipantResult {
	group := t.Groups[participant]
	score := 0
	participantResults := t.Results[participant]
	result := make(ParticipantResult, len(t.Problems))

	for _, problem := range t.Problems {
		problemCode := t.ProblemsMapping[problem]

		participantSolveTime, participantSolved := participantResults[problem]
		if participantSolved {
			score += SOLVE_POINTS
			if time.Duration(participantSolveTime)*time.Second < t.Duration {
				score += SOLVE_IN_TIME_POINTS
			}
		}

		for _, opponent := range group {
			if participant == opponent {
				continue
			}
			opponentResults := t.Results[opponent]
			opponentSolveTime, opponentSolved := opponentResults[problem]

			if participantSolved {
				if !opponentSolved || participantSolveTime < opponentSolveTime {
					score += OVERTAKE_POINTS
				}
			}
		}

		result[problemCode] = score
	}

	return result
}

func (t *Tour) Export() Table {
	result := make(Table)

	for _, group := range t.Groups {
		for _, participant := range group {
			result[participant] = t.ParticipantScore(participant)
		}
	}

	return result
}
