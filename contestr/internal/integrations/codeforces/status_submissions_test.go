package codeforces

import (
	"testing"
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

func TestAppendStatusSubmissions_partialPointsPreserved(t *testing.T) {
	t.Parallel()

	partial := 42.0
	problemIDs := map[string]int{"A": 1}
	allowed := map[string]bool{"alice": true}
	status := []StatusSubmission{
		{
			Author:              StatusParty{Members: []StatusMember{{Handle: "alice"}}},
			Problem:             StatusProblem{Index: "A"},
			RelativeTimeSeconds: 200,
			Verdict:             "PARTIAL",
			Points:              &partial,
		},
	}

	submissions := appendStatusSubmissions(nil, status, problemIDs, allowed)
	if len(submissions) != 1 || submissions[0].Points != 42 || submissions[0].Status != "PARTIAL" {
		t.Fatalf("submissions = %+v, want PARTIAL with 42 points", submissions)
	}
}
