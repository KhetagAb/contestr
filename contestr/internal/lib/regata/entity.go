package regata

import (
	"contestr/internal/lib/ejudge"
	"time"
)

// ContestResult problem -> submission_time
type ContestResult map[string]int

// Group participant group
type Group []string

type Tour struct {
	Name      string                   `json:"name"`
	StartTime time.Time                `json:"start_time"`
	Duration  time.Duration            `json:"duration"`
	Groups    map[string]Group         `json:"groups"` // participant -> group
	Problems  []string                 `json:"problems"`
	ContestID int                      `json:"contest_id"`
	Results   map[string]ContestResult `json:"results"`
}

func (t *Tour) CalcSubmissions(submissions []ejudge.Submission) map[string]ContestResult {
	return make(map[string]ContestResult)
}

func (t *Tour) ParticipantScore(participant string) int {
	return 0
}
