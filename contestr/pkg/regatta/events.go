package regatta

const (
	EventTypeProblemSolved   = "problem_solved"
	EventTypeProblemOvertake = "problem_overtake"
	EventTypeProblemRejected = "problem_rejected"
)

const SubmissionStatusTesting = "TESTING"

// IsIgnorableSubmissionStatus reports interim verdicts that must not be stored or exposed.
func IsIgnorableSubmissionStatus(status string) bool {
	return status == SubmissionStatusTesting
}

type RegattaEvent struct {
	Type          string `json:"type"`
	TimeSec       int    `json:"time_sec"`
	ParticipantID string `json:"participant_id"`
	DisplayName   string `json:"display_name"`
	ProblemCode   string `json:"problem_code"`
	TeamNumber    int    `json:"team_number,omitempty"`
	Points        int    `json:"points"`
	SolvedInTime  bool   `json:"solved_in_time"`
	FirstInGroup  bool   `json:"first_in_group,omitempty"`
	Verdict       string `json:"verdict,omitempty"`
}

var verdictShortcodes = map[string]string{
	"WRONG_ANSWER":            "WA",
	"TIME_LIMIT_EXCEEDED":     "TL",
	"RUNTIME_ERROR":           "RE",
	"COMPILATION_ERROR":       "CE",
	"MEMORY_LIMIT_EXCEEDED":   "ML",
	"PRESENTATION_ERROR":      "PE",
	"IDLENESS_LIMIT_EXCEEDED": "IL",
	"PARTIAL":                 "PT",
	"FAILED":                  "FAIL",
	"REJECTED":                "RJ",
}

// ShortSubmissionVerdict maps platform submission status to a short label for the event log.
func ShortSubmissionVerdict(status string) string {
	if s, ok := verdictShortcodes[status]; ok {
		return s
	}
	if status == "" {
		return "?"
	}
	return status
}
