package codeforces

import (
	"testing"

	"contestr/pkg/regatta"

	"github.com/togatoga/goforces"
)

func TestNormalizeCFHandle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"tourist", "tourist"},
		{"g44870=Aleksandrova_Sofyya", "Aleksandrova_Sofyya"},
		{"g44870=timofeev-dmitry", "timofeev-dmitry"},
	}

	for _, tt := range tests {
		if got := normalizeCFHandle(tt.in); got != tt.want {
			t.Errorf("normalizeCFHandle(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestBuildParticipantsUsesMappingsWhenPresent(t *testing.T) {
	t.Parallel()

	rows := []goforces.RanklistRow{
		{Party: goforces.Party{Members: []goforces.Member{{Handle: "tourist"}}}},
		{Party: goforces.Party{Members: []goforces.Member{{Handle: "Petr"}}}},
	}
	mappings := map[string]string{"Petr": "Пётр"}

	participants, allowed := buildParticipants(1, rows, mappings)

	if len(participants) != 1 {
		t.Fatalf("participants len = %d, want 1", len(participants))
	}
	if participants[0].ID != "Petr" || participants[0].DisplayName != "Пётр" {
		t.Fatalf("participant = %+v, want mapped Petr", participants[0])
	}
	if !allowed["Petr"] || allowed["tourist"] {
		t.Fatalf("allowed handles = %+v, want only Petr", allowed)
	}
}

func TestBuildParticipantsFallsBackToStandingsRows(t *testing.T) {
	t.Parallel()

	rows := []goforces.RanklistRow{
		{Party: goforces.Party{Members: []goforces.Member{{Handle: "g44870=alice"}}}},
		{Party: goforces.Party{Members: []goforces.Member{{Handle: "bob"}}}},
		{Party: goforces.Party{Members: []goforces.Member{{Handle: "bob"}}}},
	}

	participants, allowed := buildParticipants(1, rows, nil)

	if len(participants) != 2 {
		t.Fatalf("participants len = %d, want 2", len(participants))
	}
	if participants[0].ID != "alice" || participants[0].OriginalID != "g44870=alice" {
		t.Fatalf("first participant = %+v, want normalized alice with original handle", participants[0])
	}
	if participants[1].ID != "bob" {
		t.Fatalf("second participant = %+v, want bob", participants[1])
	}
	if !allowed["alice"] || !allowed["bob"] {
		t.Fatalf("allowed handles = %+v, want alice and bob", allowed)
	}
}

func TestAppendPartialStandingsFallbackAddsVisiblePartialScore(t *testing.T) {
	t.Parallel()

	rows := []goforces.RanklistRow{
		{
			Party: goforces.Party{Members: []goforces.Member{{Handle: "alice"}}},
			ProblemResults: []goforces.ProblemResult{
				{Points: 10, BestSubmissionTimeSeconds: 120},
			},
		},
	}
	problems := []goforces.Problem{{Index: "A"}}
	allowed := map[string]bool{"alice": true}

	submissions := appendPartialStandingsFallback(42, nil, rows, problems, allowed)

	if len(submissions) != 1 {
		t.Fatalf("submissions len = %d, want 1: %+v", len(submissions), submissions)
	}
	got := submissions[0]
	if got.ParticipantID != "alice" || got.ProblemID != 1 || got.Points != 10 || got.Status != "PARTIAL" || got.Time != 120 {
		t.Fatalf("submission = %+v, want alice A with 10 partial points at 120", got)
	}
}

func TestAppendPartialStandingsFallbackDoesNotDuplicateStatusScore(t *testing.T) {
	t.Parallel()

	existing := []regatta.ContestSubmission{
		{
			ParticipantID:     "alice",
			ProblemID:         1,
			Time:              100,
			Status:            "PARTIAL",
			Points:            15,
			OriginalProblemID: "A",
		},
	}
	rows := []goforces.RanklistRow{
		{
			Party: goforces.Party{Members: []goforces.Member{{Handle: "alice"}}},
			ProblemResults: []goforces.ProblemResult{
				{Points: 10, BestSubmissionTimeSeconds: 120},
			},
		},
	}
	problems := []goforces.Problem{{Index: "A"}}
	allowed := map[string]bool{"alice": true}

	submissions := appendPartialStandingsFallback(42, existing, rows, problems, allowed)

	if len(submissions) != 1 {
		t.Fatalf("submissions len = %d, want existing only: %+v", len(submissions), submissions)
	}
	if submissions[0].Points != 15 {
		t.Fatalf("points = %d, want existing status score 15", submissions[0].Points)
	}
}
