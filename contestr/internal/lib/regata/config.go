package regata

import (
	"contestr/internal/lib/util"
	"fmt"
	"strconv"
	"time"
)

type TourConfig struct {
	Name      string  `json:"name"`
	Number    int     `json:"number"`
	StartTime string  `json:"start_time"`
	Duration  int     `json:"duration"`
	Groups    [][]int `json:"groups"`
	Problems  []int   `json:"problems"`
	ContestID int     `json:"contest_id"`
}

func TourFromConfig(tc TourConfig) Tour {
	transformedProblems := util.Transform(tc.Problems, func(problemID int) Problem {
		return Problem(problemID)
	})

	return Tour{
		Name:      tc.Name,
		Number:    tc.Number,
		StartTime: util.ParseTimeOrPanic(tc.StartTime),
		Duration:  time.Duration(tc.Duration) * time.Minute,
		Groups:    ConvertGroups(tc.Groups),
		GroupSize: func() int {
			if len(tc.Groups) == 0 {
				return 0
			}
			return len(tc.Groups[0])
		}(),
		Problems:        transformedProblems,
		ContestID:       tc.ContestID,
		Results:         make(map[Participant]ContestResult),
		ProblemsMapping: tc.ProblemsIDsToNameMapping(transformedProblems),
	}
}

func (t *TourConfig) ProblemsIDsToNameMapping(problems []Problem) map[Problem]string {
	mapping := make(map[Problem]string)

	for index, problem := range problems {
		letter := 'A' + rune(index)
		mapping[problem] = strconv.Itoa(t.Number) + fmt.Sprintf("%c", letter)
	}

	return mapping
}

func ConvertGroups(groups [][]int) map[Participant]Group {
	result := make(map[Participant]Group)

	for _, group := range groups {
		for _, participantID := range group {
			result[Participant(participantID)] = util.Transform(group, func(participantID int) Participant {
				return Participant(participantID)
			})
		}
	}

	return result
}
