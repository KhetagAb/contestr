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
	Sequence          int    `bson:"sequence"`
	Round             int    `bson:"round"`
	IsPause           bool   `bson:"is_pause"`
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
		mapping[problem] = strconv.Itoa(t.Round) + fmt.Sprintf("%c", letter)
	}

	return mapping
}

func ParticipantsToGroupNumbersMapping(groups [][]string) map[Participant]int {
	result := make(map[Participant]int)

	for idx, group := range groups {
		for _, participantID := range group {
			result[participantID] = idx + 1
		}
	}

	return result
}

type ContestStandings struct {
	ContestId            int            `json:"contest_id,omitempty"`
	ContestName          string         `json:"contest_name,omitempty"`
	Rows                 []ContestRow   `json:"rows,omitempty"`
	Events               []RegattaEvent `json:"events,omitempty"`
	ContestStartTime     time.Time      `json:"contest_start_time,omitempty"`
	CurrentTime          time.Time      `json:"current_time,omitempty"`
	CurrentTourStartTime int            `json:"current_tour_start_time,omitempty"`
	CurrentTourDuration  int            `json:"current_tour_duration,omitempty"`
	IsPauseBreak         bool           `json:"is_pause_break,omitempty"`
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

type ScheduleSlot struct {
	Duration int    `bson:"duration" json:"duration"`
	Kind     string `bson:"kind" json:"kind"`
}

type ToursTimetable struct {
	ContestId        int            `bson:"contest_id" json:"contest_id"`
	PendingSlots     []ScheduleSlot `bson:"pending_slots" json:"pending_slots"`
	AutoStartEnabled bool           `bson:"auto_start_enabled" json:"auto_start_enabled"`
}

func (t *ToursTimetable) HeadSlot() (ScheduleSlot, bool) {
	if t == nil || len(t.PendingSlots) == 0 {
		return ScheduleSlot{}, false
	}
	return t.PendingSlots[0], true
}

func (t *ToursTimetable) PopHead() {
	if t == nil || len(t.PendingSlots) == 0 {
		return
	}
	t.PendingSlots = t.PendingSlots[1:]
}
