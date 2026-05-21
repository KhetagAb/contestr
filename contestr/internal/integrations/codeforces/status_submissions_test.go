package codeforces

import (
	"testing"

	"contestr/pkg/regatta"
)

func TestStatusSubmissionsIncludeFailedWithTime(t *testing.T) {
	t.Parallel()

	problemIDs := map[string]int{"A": 1}
	allowed := map[string]bool{"alice": true}
	status := []StatusSubmission{
		{
			Author:              StatusParty{Members: []StatusMember{{Handle: "alice"}}},
			Problem:             StatusProblem{Index: "A"},
			RelativeTimeSeconds: 150,
			Verdict:             "WRONG_ANSWER",
		},
		{
			Author:              StatusParty{Members: []StatusMember{{Handle: "alice"}}},
			Problem:             StatusProblem{Index: "A"},
			RelativeTimeSeconds: 420,
			Verdict:             "OK",
		},
	}

	submissions := appendStatusSubmissions(nil, status, problemIDs, allowed)

	if len(submissions) != 2 {
		t.Fatalf("submissions len = %d, want 2: %+v", len(submissions), submissions)
	}
	if submissions[0].Time != 150 || submissions[0].Status != "WRONG_ANSWER" || submissions[0].Points != 0 {
		t.Fatalf("first submission = %+v, want WA at 150", submissions[0])
	}
	if submissions[1].Time != 420 || submissions[1].Status != "OK" || submissions[1].Points != 100 {
		t.Fatalf("second submission = %+v, want OK at 420 with 100 points", submissions[1])
	}
}

func TestNormalizeBinarySubmissions(t *testing.T) {
	t.Parallel()

	submissions := []regatta.ContestSubmission{
		{Status: "OK", Points: 100},
		{Status: "WRONG_ANSWER", Points: 0},
	}
	normalizeBinarySubmissions(submissions)
	if submissions[0].Points != 100 {
		t.Fatalf("OK points = %d, want 100", submissions[0].Points)
	}
	if submissions[1].Points != 0 {
		t.Fatalf("WA points = %d, want 0", submissions[1].Points)
	}
}
