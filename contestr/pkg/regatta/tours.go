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
	Name              string `bson:"name"`
	Index             int    `bson:"index"`
	StartTime         int    `bson:"start_time"`
	DurationInSeconds int    `bson:"duration_in_seconds"`

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
	ContestId            int          `json:"contest_id,omitempty"`
	ContestName          string       `json:"contest_name,omitempty"`
	Rows                 []ContestRow `json:"rows,omitempty"`
	ContestStartTime     time.Time    `json:"contest_start_time,omitempty"`
	CurrentTime          time.Time    `json:"current_time,omitempty"`
	CurrentTourStartTime int          `json:"current_tour_start_time,omitempty"`
	CurrentTourDuration  int          `json:"current_tour_duration,omitempty"`
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

type TourConfig struct {
	StartTime int  `bson:"start_time" json:"start_time"`
	Duration  int  `bson:"duration" json:"duration"`
	Started   bool `bson:"started" json:"started"`
}

type ToursTimetable struct {
	ContestId        int          `bson:"contest_id" json:"contest_id"`
	TourTimes        []TourConfig `bson:"tour_times" json:"tour_times"`
	AutoStartEnabled bool         `bson:"auto_start_enabled" json:"auto_start_enabled"`
}

func (t *ToursTimetable) FirstNotStartedTour() (int, TourConfig, bool) {
	if t == nil {
		return 0, TourConfig{}, false
	}

	for i, tour := range t.TourTimes {
		if !tour.Started {
			return i + 1, tour, true
		}
	}
	return 0, TourConfig{}, false
}
