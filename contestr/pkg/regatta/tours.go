package regatta

import (
	"fmt"
	"strconv"
	"time"
)

type Participant = string

type Problem = int

type Group = []Participant

type Tour struct {
	Name  string `bson:"name"`
	Index int    `bson:"index"`
	// TODO refactor
	StarTime          int `bson:"start_time"`
	EndTime           int `bson:"end_time"`
	DurationInSeconds int `bson:"duration_in_seconds"`

	Groups       map[Participant]Group `bson:"groups"`
	GroupSize    int                   `bson:"group_size"`
	Problems     []Problem             `bson:"problems"`
	ContestID    int                   `bson:"contest_id"`
	GroupNumbers map[Participant]int   `bson:"group_numbers"`
}

func (t *Tour) ProblemsIDsToNameMapping(problems []Problem) map[Problem]string {
	mapping := make(map[Problem]string)

	for index, problem := range problems {
		letter := 'A' + rune(index)
		mapping[problem] = strconv.Itoa(t.Index) + fmt.Sprintf("%c", letter)
	}

	return mapping
}

func ParticipantsToGroupNumbersMapping(groups [][]string) map[Participant]int {
	result := make(map[Participant]int)

	for idx, group := range groups {
		for _, participantID := range group {
			result[participantID] = idx + 1
		}
		idx += 1
	}

	return result
}

type ContestStandings struct {
	ContestId             int          `json:"contest_id,omitempty"`
	ContestName           string       `json:"contest_name,omitempty"`
	Rows                  []ContestRow `json:"rows,omitempty"`
	ContestStartTime      time.Time    `json:"contest_start_time,omitempty"`
	CurrentTime           time.Time    `json:"current_time,omitempty"`
	CurrentTourStartTime  int          `json:"current_tour_start_time,omitempty"`
	CurrentTourDuration   int          `json:"current_tour_duration,omitempty"`
}

type ProblemResult struct {
	ProblemCode        string `json:"problem_code"`
	Score              int    `json:"score"`
	LastSubmissionTime int    `json:"last_submission_time"`
}

type ContestRow struct {
	UserID         string          `json:"user_id"`
	DisplayName    string          `json:"display_name"`
	ProblemResults []ProblemResult `json:"problem_results"`
	SolvedProblems int             `json:"solved_problems"`
	TeamNumber     int             `json:"team_number"`
	TotalScore     int             `json:"total_score"`

}
