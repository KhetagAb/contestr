package regatta

import "testing"

func TestShortSubmissionVerdict(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status string
		want   string
	}{
		{"WRONG_ANSWER", "WA"},
		{"TIME_LIMIT_EXCEEDED", "TL"},
		{"RUNTIME_ERROR", "RE"},
		{"", "?"},
		{"CUSTOM", "CUSTOM"},
	}

	for _, tt := range tests {
		if got := ShortSubmissionVerdict(tt.status); got != tt.want {
			t.Errorf("ShortSubmissionVerdict(%q) = %q, want %q", tt.status, got, tt.want)
		}
	}
}
