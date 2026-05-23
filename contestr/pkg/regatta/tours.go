package regatta

import (
	"fmt"
	"strconv"
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
