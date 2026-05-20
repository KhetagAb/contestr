package regatta

const (
	EventTypeProblemSolved   = "problem_solved"
	EventTypeProblemOvertake = "problem_overtake"
)

type RegattaEvent struct {
	Type          string `json:"type"`
	TimeSec       int    `json:"time_sec"`
	ParticipantID string `json:"participant_id"`
	DisplayName   string `json:"display_name"`
	ProblemCode   string `json:"problem_code"`
	TeamNumber    int    `json:"team_number,omitempty"`
	Points        int    `json:"points"`
	SolvedInTime  bool `json:"solved_in_time"`
	FirstInGroup  bool `json:"first_in_group,omitempty"`
}
