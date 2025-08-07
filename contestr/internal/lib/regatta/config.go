package regatta

import (
	"contestr/internal/lib/util"
	"fmt"
	"strconv"
	"time"
)

type TourConfig struct {
	Name      string  `json:"name"`
	Index     int     `json:"index"`
	StartTime string  `json:"start_time"`
	Duration  int     `json:"duration"`
	Groups    [][]int `json:"groups"`
	Problems  []int   `json:"problems"`
	ContestID int     `json:"contest_id"`
}

func TourFromConfig(tc TourConfig) Tour {
	return Tour{
		Name:      tc.Name,
		Index:     tc.Index,
		StartTime: util.ParseTimeOrPanic(tc.StartTime),
		Duration:  time.Duration(tc.Duration) * time.Minute,
		Groups:    ConvertGroups(tc.Groups),
		GroupSize: func() int {
			if len(tc.Groups) == 0 {
				return 0
			}
			return len(tc.Groups[0])
		}(),
		Problems:        tc.Problems,
		ContestID:       tc.ContestID,
		Results:         make(map[Participant]ContestResult),
		ProblemsMapping: tc.ProblemsIDsToNameMapping(tc.Problems),
	}
}

func (t *TourConfig) ProblemsIDsToNameMapping(problems []Problem) map[Problem]string {
	mapping := make(map[Problem]string)

	for index, problem := range problems {
		letter := 'A' + rune(index)
		mapping[problem] = strconv.Itoa(t.Index) + fmt.Sprintf("%c", letter)
	}

	return mapping
}

func ConvertGroups(groups [][]int) map[Participant]Group {
	result := make(map[Participant]Group)

	for _, group := range groups {
		for _, participantID := range group {
			result[participantID] = group
		}
	}

	return result
}
