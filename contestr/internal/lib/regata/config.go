package regata

import (
	"contestr/internal/lib/util"
	"time"
)

type TourConfig struct {
	Name      string  `json:"name"`
	StartTime string  `json:"start_time"`
	Duration  int     `json:"duration"`
	Groups    [][]int `json:"groups"`
	Problems  []int   `json:"problems"`
	ContestID int     `json:"contest_id"`
}

func TourFromConfig(tc TourConfig) Tour {
	return Tour{
		Name:      tc.Name,
		StartTime: util.ParseTime(tc.StartTime),
		Duration:  time.Duration(tc.Duration) * time.Minute,
		Groups:    nil,
		Problems: util.Transform(tc.Problems, func(problemID int) Problem {
			return Problem(problemID)
		}),
		ContestID: tc.ContestID,
		Results:   make(map[Participant]ContestResult),
	}
}
