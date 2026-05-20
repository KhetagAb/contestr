package regatta

import (
	"os"
	"testing"

	"contestr/pkg/logger"
	"contestr/pkg/regatta"
)

func TestMain(m *testing.M) {
	logger.Init("regatta-test", "error")
	os.Exit(m.Run())
}

func testTour() regatta.Tour {
	return regatta.Tour{
		Name:              "Tour 1",
		Sequence:          1,
		Round:             1,
		DurationInSeconds: 3600,
		Problems:          []int{1, 2},
		Groups: map[string][]string{
			"alice": {"alice", "bob"},
			"bob":   {"alice", "bob"},
		},
		GroupNumbers: map[string]int{
			"alice": 1,
			"bob":   1,
		},
	}
}

func TestParticipantScore_eventsAndScore(t *testing.T) {
	tour := testTour()
	tr := CalculateResult(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusOK},
		{UserID: "bob", ProbID: 1, Time: 500, Status: SubmissionStatusOK},
	})

	out := tr.ParticipantScore("alice")
	wantScore := 100 + SolvePoints + SolveInTimePoints + OvertakePoints
	if out.Results["1A"].score != wantScore {
		t.Fatalf("alice score = %d, want %d", out.Results["1A"].score, wantScore)
	}

	if len(out.Events) != 1 {
		t.Fatalf("expected 1 event per solve, got %d", len(out.Events))
	}
	ev := out.Events[0]
	if ev.Points != wantScore {
		t.Fatalf("event points = %d, want %d", ev.Points, wantScore)
	}
	if !ev.SolvedInTime || !ev.FirstInGroup {
		t.Fatalf("expected in-time and first-in-group flags, got %+v", ev)
	}
}

func TestParticipantScore_singleParticipantFullSolveGetsProblemAndSolveBonuses(t *testing.T) {
	tour := testTour()
	tour.Groups = map[string][]string{
		"alice": {"alice"},
	}
	tour.GroupNumbers = map[string]int{
		"alice": 1,
	}

	tr := CalculateResult(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusOK},
	})

	want := 100 + SolvePoints + SolveInTimePoints
	out := tr.ParticipantScore("alice")
	if got := out.Results["1A"].score; got != want {
		t.Fatalf("single participant score = %d, want %d", got, want)
	}
	if len(out.Events) != 1 {
		t.Fatalf("expected solve event, got %d events", len(out.Events))
	}
	if out.Events[0].FirstInGroup {
		t.Fatalf("single participant should not receive overtake bonus: %+v", out.Events[0])
	}
}

func TestBuildContestEvents_chronological(t *testing.T) {
	tour := testTour()
	tours := []regatta.Tour{tour}
	runs := []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusOK},
		{UserID: "bob", ProbID: 1, Time: 500, Status: SubmissionStatusOK},
	}
	names := map[string]string{"alice": "Alice A", "bob": "Bob B"}

	events := BuildContestEvents(tours, runs, names)
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].TimeSec != 100 || events[0].ParticipantID != "alice" {
		t.Fatalf("first event should be alice at 100, got %+v", events[0])
	}
	if events[0].DisplayName != "Alice A" {
		t.Fatalf("display name = %q", events[0].DisplayName)
	}
	if events[0].Points != 100+SolvePoints+SolveInTimePoints+OvertakePoints {
		t.Fatalf("alice points = %d", events[0].Points)
	}

	for _, e := range events {
		if e.ParticipantID == "bob" && e.FirstInGroup {
			t.Fatal("bob should not get first_in_group when alice solved earlier")
		}
	}
}

func TestScoreForSolve_noMultiOvertakeBonus(t *testing.T) {
	tour := testTour()
	tour.Groups = map[string][]string{
		"alice": {"alice", "bob", "carol"},
		"bob":   {"alice", "bob", "carol"},
		"carol": {"alice", "bob", "carol"},
	}
	tr := CalculateResult(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusOK},
		{UserID: "bob", ProbID: 1, Time: 200, Status: SubmissionStatusOK},
		{UserID: "carol", ProbID: 1, Time: 300, Status: SubmissionStatusOK},
	})

	want := 100 + SolvePoints + SolveInTimePoints + OvertakePoints
	got := scoreForProblem(tr, "alice", 1)
	if got != want {
		t.Fatalf("score = %d, want %d (no +%d per opponent)", got, want, OvertakePoints)
	}
}

func TestBuildContestEvents_byProblemTourNotSubmissionTime(t *testing.T) {
	tour1 := testTour()
	tour2 := testTour()
	tour2.Name = "Tour 2"
	tour2.Sequence = 3
	tour2.Round = 2
	tour2.Problems = []int{3, 4}
	tours := []regatta.Tour{
		tour1,
		{Sequence: 2, Round: 0, IsPause: true, DurationInSeconds: 1000},
		tour2,
	}

	runs := []Run{
		{UserID: "alice", ProbID: 3, Time: 200, Status: SubmissionStatusOK},
	}
	events := BuildContestEvents(tours, runs, map[string]string{"alice": "Alice"})
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].ProblemCode != "2A" {
		t.Fatalf("problem code = %q, want 2A", events[0].ProblemCode)
	}
}

func TestBuildContestEvents_countMatchesTable(t *testing.T) {
	tour1 := testTour()
	tour2 := testTour()
	tour2.Name = "Tour 2"
	tour2.Sequence = 3
	tour2.Round = 2
	tour2.Problems = []int{3, 4}
	tours := []regatta.Tour{
		tour1,
		{Sequence: 2, Round: 0, IsPause: true, DurationInSeconds: 1000},
		tour2,
	}

	runs := []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusOK},
		{UserID: "alice", ProbID: 3, Time: 4200, Status: SubmissionStatusOK},
		{UserID: "bob", ProbID: 2, Time: 500, Status: SubmissionStatusOK},
		{UserID: "bob", ProbID: 4, Time: 5000, Status: SubmissionStatusOK},
	}

	events := BuildContestEvents(tours, runs, nil)
	offsets := regatta.SegmentOffsets(tours)
	tableSolves := 0
	for _, tour := range tours {
		if tour.IsPause {
			continue
		}
		start := offsets[tour.Sequence].Start
		for _, pr := range CalculateResult(tour, start, runs).Export() {
			for _, p := range pr {
				if p.score > 0 {
					tableSolves++
				}
			}
		}
	}
	if len(events) != tableSolves {
		t.Fatalf("events %d != table solves %d", len(events), tableSolves)
	}
}

func TestBuildContestEvents_sortedByTime(t *testing.T) {
	tour1 := testTour()
	tour2 := testTour()
	tour2.Sequence = 2
	tour2.Round = 2
	tour2.Problems = []int{3, 4}
	tours := []regatta.Tour{tour1, tour2}

	runs := []Run{
		{UserID: "alice", ProbID: 1, Time: 4100, Status: SubmissionStatusOK},
		{UserID: "bob", ProbID: 1, Time: 200, Status: SubmissionStatusOK},
	}
	events := BuildContestEvents(tours, runs, map[string]string{})
	if len(events) < 2 {
		t.Fatal("expected events")
	}
	if events[0].TimeSec > events[len(events)-1].TimeSec {
		t.Fatalf("events not sorted: first %d last %d", events[0].TimeSec, events[len(events)-1].TimeSec)
	}
	if events[0].ParticipantID != "bob" {
		t.Fatalf("first should be bob at t=200, got %+v", events[0])
	}
}

func TestPartialMode_uniqueFirstFullSolveGetsOvertake(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	settings.Mode = regatta.ScoringModePartial
	tr := CalculateResultWithSettings(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusOK, Points: 100},
		{UserID: "bob", ProbID: 1, Time: 200, Status: SubmissionStatusOK, Points: 100},
	}, settings)

	want := 100 + SolvePoints + SolveInTimePoints + OvertakePoints
	if got := tr.ParticipantScore("alice").Results["1A"].score; got != want {
		t.Fatalf("alice score = %d, want %d", got, want)
	}
	if got := tr.ParticipantScore("bob").Results["1A"].score; got != 100+SolvePoints+SolveInTimePoints {
		t.Fatalf("bob score = %d, want no overtake", got)
	}
}

func TestPartialMode_tiedFirstFullSolveGetsNoOvertake(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	settings.Mode = regatta.ScoringModePartial
	tr := CalculateResultWithSettings(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusOK, Points: 100},
		{UserID: "bob", ProbID: 1, Time: 100, Status: SubmissionStatusOK, Points: 100},
	}, settings)

	want := 100 + SolvePoints + SolveInTimePoints
	for _, participant := range []string{"alice", "bob"} {
		if got := tr.ParticipantScore(participant).Results["1A"].score; got != want {
			t.Fatalf("%s score = %d, want %d", participant, got, want)
		}
	}
}

func TestPartialMode_tourEndHighestPartialGetsOvertake(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	settings.Mode = regatta.ScoringModePartial
	tr := CalculateResultWithSettings(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusPartial, Points: 40},
		{UserID: "alice", ProbID: 1, Time: 3700, Status: SubmissionStatusPartial, Points: 80},
		{UserID: "bob", ProbID: 1, Time: 200, Status: SubmissionStatusPartial, Points: 60},
	}, settings)

	if got := tr.ParticipantScore("alice").Results["1A"].score; got != 80 {
		t.Fatalf("alice score = %d, want latest raw points without overtake", got)
	}
	if got := tr.ParticipantScore("bob").Results["1A"].score; got != 60+OvertakePoints {
		t.Fatalf("bob score = %d, want raw points plus overtake", got)
	}

	events := BuildContestEvents([]regatta.Tour{tour}, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusPartial, Points: 40},
		{UserID: "alice", ProbID: 1, Time: 3700, Status: SubmissionStatusPartial, Points: 80},
		{UserID: "bob", ProbID: 1, Time: 200, Status: SubmissionStatusPartial, Points: 60},
	}, nil, settings)
	if len(events) != 1 {
		t.Fatalf("expected one bob scoring event with overtake badge, got %+v", events)
	}
	if events[0].Type != regatta.EventTypeProblemSolved ||
		events[0].ParticipantID != "bob" ||
		!events[0].FirstInGroup ||
		events[0].TimeSec != 200 ||
		events[0].Points != 60+OvertakePoints {
		t.Fatalf("expected bob scoring event with overtake badge, got %+v", events)
	}
}

func TestPartialMode_activeTourHighestPartialGetsNoTourEndOvertake(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	settings.Mode = regatta.ScoringModePartial
	runs := []Run{
		{UserID: "alice", ProbID: 1, Time: 900, Status: SubmissionStatusPartial, Points: 20},
		{UserID: "bob", ProbID: 1, Time: 700, Status: SubmissionStatusPartial, Points: 10},
	}

	tr := CalculateResultWithSettingsAt(tour, 0, runs, settings, 1100)

	if got := tr.ParticipantScore("alice").Results["1A"].score; got != 20 {
		t.Fatalf("alice score = %d, want raw points only while tour is active", got)
	}
	events := BuildContestEventsAt([]regatta.Tour{tour}, runs, nil, 1100, settings)
	if len(events) != 0 {
		t.Fatalf("expected no tour-end overtake event while tour is active, got %+v", events)
	}
}

func TestPartialMode_tiedHighestPartialGetsNoOvertake(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	settings.Mode = regatta.ScoringModePartial
	tr := CalculateResultWithSettings(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusPartial, Points: 60},
		{UserID: "bob", ProbID: 1, Time: 200, Status: SubmissionStatusPartial, Points: 60},
	}, settings)

	for _, participant := range []string{"alice", "bob"} {
		if got := tr.ParticipantScore(participant).Results["1A"].score; got != 60 {
			t.Fatalf("%s score = %d, want raw points only", participant, got)
		}
	}
}

func TestPartialMode_zeroNeverGetsOvertake(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	settings.Mode = regatta.ScoringModePartial
	tr := CalculateResultWithSettings(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusPartial, Points: 0},
	}, settings)

	if got := tr.ParticipantScore("alice").Results["1A"].score; got != 0 {
		t.Fatalf("alice score = %d, want 0", got)
	}
}

func TestBinaryMode_duringTourOnlyOvertakeOption(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	settings.BinaryOvertakeMode = regatta.BinaryOvertakeModeDuringTour
	tr := CalculateResultWithSettings(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 3700, Status: SubmissionStatusOK},
	}, settings)

	want := 100 + SolvePoints
	if got := tr.ParticipantScore("alice").Results["1A"].score; got != want {
		t.Fatalf("alice score = %d, want %d", got, want)
	}
}

func TestPauseBetweenTours_segmentOffsets(t *testing.T) {
	tours := []regatta.Tour{
		{Sequence: 1, Round: 1, DurationInSeconds: 1000, IsPause: false},
		{Sequence: 2, Round: 0, DurationInSeconds: 300, IsPause: true},
		{Sequence: 3, Round: 2, DurationInSeconds: 2000, IsPause: false},
	}
	offsets := regatta.SegmentOffsets(tours)
	if offsets[3].Start != 1300 {
		t.Fatalf("tour 2 should start after pause, got %d", offsets[3].Start)
	}
}
