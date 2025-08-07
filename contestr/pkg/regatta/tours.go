package regatta

import (
	"fmt"
	"strconv"
	"time"
)

type Participant = int

type Problem = int

type Group = []Participant

type Tour struct {
	Name      string                `bson:"name"`
	Index     int                   `bson:"index"`
	StartTime time.Time             `bson:"start_time"`
	Duration  time.Duration         `bson:"duration"`
	Groups    map[Participant]Group `bson:"groups"`
	GroupSize int                   `bson:"group_size"`
	Problems  []Problem             `bson:"problems"`
	ContestID int                   `bson:"contest_id"`
}

func (t *Tour) ProblemsIDsToNameMapping(problems []Problem) map[Problem]string {
	mapping := make(map[Problem]string)

	for index, problem := range problems {
		letter := 'A' + rune(index)
		mapping[problem] = strconv.Itoa(t.Index) + fmt.Sprintf("%c", letter)
	}

	return mapping
}
