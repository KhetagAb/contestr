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

// ShortSubmissionVerdict maps platform submission status to a short label for the event log.
func ShortSubmissionVerdict(status string) string {
	switch status {
	case "WRONG_ANSWER":
		return "WA"
	case "TIME_LIMIT_EXCEEDED":
		return "TL"
	case "RUNTIME_ERROR":
		return "RE"
	case "COMPILATION_ERROR":
		return "CE"
	case "MEMORY_LIMIT_EXCEEDED":
		return "ML"
	case "PRESENTATION_ERROR":
		return "PE"
	case "IDLENESS_LIMIT_EXCEEDED":
		return "IL"
	case "PARTIAL":
		return "PT"
	case "FAILED":
		return "FAIL"
	case "REJECTED":
		return "RJ"
	default:
		if status == "" {
			return "?"
		}
		return status
	}
}
