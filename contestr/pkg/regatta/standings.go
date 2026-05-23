package regatta

import "time"

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
	LastSubmissionTime int    `json:"last_submission_time,omitempty"`
}

type ContestRow struct {
	UserID         string          `json:"user_id"`
	DisplayName    string          `json:"display_name"`
	ProblemResults []ProblemResult `json:"problem_results"`
	SolvedProblems int             `json:"solved_problems"`
	TeamNumber     int             `json:"team_number"`
	TotalScore     int             `json:"total_score"`
}
