package regata

import (
	"contestr/internal/lib/ejudge"
	"time"
)

const (
	OVERTAKE_POINTS = 5
	SOLVE_POINTS    = 5
)

const (
	SubmissionStatusOK string = "OK"
)

// Problem int ejudge ID
type Problem int

// Participant int ejudge ID
type Participant int

// ContestResult problem -> submission_time
type ContestResult map[Problem]int

// Group participant group
type Group []Participant

// Tour regata tour
type Tour struct {
	Name      string                        `json:"name"`
	StartTime time.Time                     `json:"start_time"`
	Duration  time.Duration                 `json:"duration"`
	Groups    map[Participant]Group         `json:"groups"` // participant -> group
	Problems  []Problem                     `json:"problems"`
	ContestID int                           `json:"contest_id"`
	Results   map[Participant]ContestResult `json:"results"` // participant -> result
}

func (t *Tour) CalcSubmissions(submissions []ejudge.Submission) map[Participant]ContestResult {
	results := make(map[Participant]ContestResult)

	for _, submission := range submissions {
		if submission.Status != SubmissionStatusOK {
			continue
		}

		participant := Participant(submission.UserID)
		problem := Problem(submission.ProbID)

		_, ok := results[participant]
		if !ok {
			results[participant] = map[Problem]int{
				problem: submission.Time,
			}
		} else {
			results[participant][problem] = submission.Time
		}
	}

	return results
}

func (t *Tour) ParticipantScore(participant Participant) int {
	group := t.Groups[participant]
	score := 0

	for _, opponent := range group {
		if participant == opponent {
			continue
		}
		participantResults := t.Results[participant]
		opponentResults := t.Results[opponent]

		for _, problem := range t.Problems {
			participantProblemResult, participantSolved := participantResults[problem]
			opponentProblemResult, opponentSolved := opponentResults[problem]

			if participantSolved {
				score += SOLVE_POINTS
				if !opponentSolved || participantProblemResult < opponentProblemResult {
					score += OVERTAKE_POINTS
				}
			}
		}
	}

	return score
}
