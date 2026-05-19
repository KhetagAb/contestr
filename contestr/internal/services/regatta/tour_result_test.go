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
	tr := &TourResult{
		Tour:         tour,
		SegmentStart: 0,
		Results: map[string]map[int]int{
			"alice": {1: 100},
			"bob":   {1: 500},
		},
		ProblemsMapping: tour.ProblemsIDsToNameMapping(tour.Problems),
	}

	out := tr.ParticipantScore("alice")
	wantScore := SolvePoints + SolveInTimePoints + OvertakePoints
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
	if events[0].Points != SolvePoints+SolveInTimePoints+OvertakePoints {
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
	tr := &TourResult{
		Tour:         tour,
		SegmentStart: 0,
		Results: map[string]map[int]int{
			"alice": {1: 100},
			"bob":   {1: 200},
			"carol": {1: 300},
		},
		ProblemsMapping: tour.ProblemsIDsToNameMapping(tour.Problems),
	}

	want := SolvePoints + SolveInTimePoints + OvertakePoints
	got := scoreForSolve(tr, "alice", 1, 100)
	if got != want {
		t.Fatalf("score = %d, want %d (no +5 per opponent)", got, want)
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
