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
	wantScore := 100 + SolveInTimePoints + OvertakePoints
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

	want := 100 + SolveInTimePoints
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
	if events[0].Points != 100+SolveInTimePoints+OvertakePoints {
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

	want := 100 + SolveInTimePoints + OvertakePoints
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
		{UserID: "alice", ProbID: 3, Time: 5000, Status: SubmissionStatusOK},
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
		{UserID: "alice", ProbID: 3, Time: 5000, Status: SubmissionStatusOK},
		{UserID: "bob", ProbID: 2, Time: 500, Status: SubmissionStatusOK},
		{UserID: "bob", ProbID: 4, Time: 5500, Status: SubmissionStatusOK},
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

func TestOvertake_uniqueFirstFullSolveGetsOvertake(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	tr := CalculateResultWithSettings(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusOK, Points: 100},
		{UserID: "bob", ProbID: 1, Time: 200, Status: SubmissionStatusOK, Points: 100},
	}, settings)

	want := 100 + SolveInTimePoints + OvertakePoints
	if got := tr.ParticipantScore("alice").Results["1A"].score; got != want {
		t.Fatalf("alice score = %d, want %d", got, want)
	}
	if got := tr.ParticipantScore("bob").Results["1A"].score; got != 100+SolveInTimePoints {
		t.Fatalf("bob score = %d, want no overtake", got)
	}
}

func TestOvertake_tiedFirstFullSolveGetsNoOvertake(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	tr := CalculateResultWithSettings(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusOK, Points: 100},
		{UserID: "bob", ProbID: 1, Time: 100, Status: SubmissionStatusOK, Points: 100},
	}, settings)

	want := 100 + SolveInTimePoints
	for _, participant := range []string{"alice", "bob"} {
		if got := tr.ParticipantScore(participant).Results["1A"].score; got != want {
			t.Fatalf("%s score = %d, want %d", participant, got, want)
		}
	}
}

func TestOvertake_tourEndHighestPartialGetsOvertake(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
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

	var bobOvertake *regatta.RegattaEvent
	var bobAt200 int
	for i := range events {
		e := &events[i]
		if e.ParticipantID != "bob" || e.ProblemCode != "1A" {
			continue
		}
		if e.TimeSec == 200 {
			bobAt200++
		}
		if e.Type == regatta.EventTypeProblemOvertake {
			bobOvertake = &events[i]
		}
	}
	if bobAt200 != 1 {
		t.Fatalf("expected single tour-end overtake at 200, got %d events", bobAt200)
	}
	if bobOvertake == nil {
		t.Fatalf("expected bob tour-end overtake event, got %+v", events)
	}
	if bobOvertake.TimeSec != 200 || bobOvertake.Points != 60+OvertakePoints {
		t.Fatalf("expected bob overtake at 200 with bonus points, got %+v", bobOvertake)
	}
}

func TestOvertake_activeTourHighestPartialGetsNoTourEndOvertake(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	runs := []Run{
		{UserID: "alice", ProbID: 1, Time: 900, Status: SubmissionStatusPartial, Points: 20},
		{UserID: "bob", ProbID: 1, Time: 700, Status: SubmissionStatusPartial, Points: 10},
	}

	tr := CalculateResultWithSettingsAt(tour, 0, runs, settings, 1100)

	if got := tr.ParticipantScore("alice").Results["1A"].score; got != 20 {
		t.Fatalf("alice score = %d, want raw points only while tour is active", got)
	}
	events := BuildContestEventsAt([]regatta.Tour{tour}, runs, nil, 1100, settings)
	for _, e := range events {
		if e.Type == regatta.EventTypeProblemOvertake {
			t.Fatalf("expected no tour-end overtake event while tour is active, got %+v", events)
		}
	}
	if len(events) != 2 {
		t.Fatalf("expected two partial score events while tour is active, got %+v", events)
	}
}

func TestOvertake_tiedHighestPartialGetsNoOvertake(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
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

func TestOvertake_zeroNeverGetsOvertake(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	tr := CalculateResultWithSettings(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusPartial, Points: 0},
	}, settings)

	if got := tr.ParticipantScore("alice").Results["1A"].score; got != 0 {
		t.Fatalf("alice score = %d, want 0", got)
	}
}

func TestParticipantScore_rejectedAttemptsNegativeScore(t *testing.T) {
	tour := testTour()
	tr := CalculateResult(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: "WRONG_ANSWER"},
		{UserID: "alice", ProbID: 1, Time: 200, Status: "WRONG_ANSWER"},
	})

	if got := tr.ParticipantScore("alice").Results["1A"].score; got != -2 {
		t.Fatalf("alice score = %d, want -2 (two rejected attempts)", got)
	}
}

func TestParticipantScore_noAttemptsZeroScore(t *testing.T) {
	tour := testTour()
	tr := CalculateResult(tour, 0, []Run{})

	if got := tr.ParticipantScore("alice").Results["1A"].score; got != 0 {
		t.Fatalf("alice score = %d, want 0 (no attempts)", got)
	}
}

func TestParticipantScore_solvedPositiveScoreNoRejected(t *testing.T) {
	tour := testTour()
	tr := CalculateResult(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: "WRONG_ANSWER"},
		{UserID: "alice", ProbID: 1, Time: 500, Status: SubmissionStatusOK},
	})

	got := tr.ParticipantScore("alice").Results["1A"].score
	if got <= 0 {
		t.Fatalf("alice score = %d, want positive solve score", got)
	}
}

func TestBuildContestEvents_includesRejectedSubmissions(t *testing.T) {
	tour := testTour()
	events := BuildContestEvents([]regatta.Tour{tour}, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: "WRONG_ANSWER"},
		{UserID: "alice", ProbID: 1, Time: 200, Status: "WRONG_ANSWER"},
		{UserID: "alice", ProbID: 1, Time: 500, Status: SubmissionStatusOK},
	}, nil)

	var rejected, solved int
	for _, e := range events {
		switch e.Type {
		case regatta.EventTypeProblemRejected:
			rejected++
			if e.ParticipantID != "alice" || e.ProblemCode != "1A" || e.Points != 0 || e.Verdict != "WA" {
				t.Fatalf("unexpected rejected event: %+v", e)
			}
		case regatta.EventTypeProblemSolved:
			solved++
		}
	}

	if rejected != 2 {
		t.Fatalf("rejected events = %d, want 2", rejected)
	}
	if solved != 1 {
		t.Fatalf("solved events = %d, want 1", solved)
	}
}

func TestBuildContestEvents_rejectedTimes(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	events := BuildContestEvents([]regatta.Tour{tour}, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: "WRONG_ANSWER"},
		{UserID: "alice", ProbID: 1, Time: 200, Status: "WRONG_ANSWER"},
	}, nil, settings)

	var times []int
	for _, e := range events {
		if e.Type == regatta.EventTypeProblemRejected {
			times = append(times, e.TimeSec)
		}
	}
	if len(times) != 2 {
		t.Fatalf("rejected events = %d, want 2", len(times))
	}
	if times[0] != 100 || times[1] != 200 {
		t.Fatalf("rejected times = %v, want [100 200]", times)
	}
}

func TestBuildContestEvents_skipsTestingVerdict(t *testing.T) {
	tour := testTour()
	events := BuildContestEvents([]regatta.Tour{tour}, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: regatta.SubmissionStatusTesting},
		{UserID: "alice", ProbID: 1, Time: 200, Status: "WRONG_ANSWER"},
		{UserID: "alice", ProbID: 1, Time: 500, Status: SubmissionStatusOK},
	}, nil)

	var rejected int
	for _, e := range events {
		if e.Type == regatta.EventTypeProblemRejected {
			rejected++
			if e.Verdict == regatta.SubmissionStatusTesting || e.Verdict == "TESTING" {
				t.Fatalf("unexpected TESTING verdict in event: %+v", e)
			}
		}
	}
	if rejected != 1 {
		t.Fatalf("rejected events = %d, want 1 (TESTING skipped)", rejected)
	}
}

func TestBuildContestEvents_partialSubmissionsEmitRawPoints(t *testing.T) {
	tour := testTour()
	events := BuildContestEvents([]regatta.Tour{tour}, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusPartial, Points: 40},
		{UserID: "alice", ProbID: 1, Time: 3700, Status: SubmissionStatusPartial, Points: 80},
		{UserID: "alice", ProbID: 1, Time: 5000, Status: SubmissionStatusOK, Points: 100},
	}, nil)

	var partialPoints []int
	var fullSolve int
	for _, e := range events {
		if e.Type != regatta.EventTypeProblemSolved || e.ParticipantID != "alice" || e.ProblemCode != "1A" {
			continue
		}
		if e.FirstInGroup || e.Points >= 100 {
			fullSolve++
			continue
		}
		partialPoints = append(partialPoints, e.Points)
	}

	if len(partialPoints) != 2 || partialPoints[0] != 40 || partialPoints[1] != 80 {
		t.Fatalf("partial events = %v, want [40 80]", partialPoints)
	}
	if fullSolve != 1 {
		t.Fatalf("full solve events = %d, want 1", fullSolve)
	}
}

func TestBuildContestEvents_rejectedVerdictTL(t *testing.T) {
	tour := testTour()
	events := BuildContestEvents([]regatta.Tour{tour}, []Run{
		{UserID: "alice", ProbID: 1, Time: 300, Status: "TIME_LIMIT_EXCEEDED"},
	}, nil)

	if len(events) != 1 || events[0].Type != regatta.EventTypeProblemRejected || events[0].Verdict != "TL" {
		t.Fatalf("events = %+v, want one rejected TL event", events)
	}
}

func TestUnifiedScoring_okWithZeroPointsFullBonusesInTour(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	tr := CalculateResultWithSettings(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusOK, Points: 0},
		{UserID: "bob", ProbID: 1, Time: 200, Status: SubmissionStatusOK, Points: 0},
	}, settings)

	want := 100 + SolveInTimePoints + OvertakePoints
	if got := tr.ParticipantScore("alice").Results["1A"].score; got != want {
		t.Fatalf("alice score = %d, want %d", got, want)
	}
}

func TestUnifiedScoring_partialWithoutFullSolveNoFullBonus(t *testing.T) {
	tour := testTour()
	tour.Groups = map[string][]string{"alice": {"alice"}}
	tour.GroupNumbers = map[string]int{"alice": 1}
	settings := regatta.DefaultScoringSettings()
	tr := CalculateResultWithSettings(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 100, Status: SubmissionStatusPartial, Points: 40},
	}, settings)

	if got := tr.ParticipantScore("alice").Results["1A"].score; got != 40 {
		t.Fatalf("alice score = %d, want 40 without full-solve bonus", got)
	}
}

func TestBuildContestEvents_noDuplicateWhenFullSolveAfterPartialTourEnd(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	runs := []Run{
		{UserID: "alice", ProbID: 1, Time: 882, Status: SubmissionStatusPartial, Points: 80},
		{UserID: "bob", ProbID: 1, Time: 200, Status: SubmissionStatusPartial, Points: 60},
		{UserID: "alice", ProbID: 1, Time: 2374, Status: SubmissionStatusOK},
	}

	events := BuildContestEvents([]regatta.Tour{tour}, runs, nil, settings)
	var alice1A []regatta.RegattaEvent
	for _, e := range events {
		if e.ParticipantID == "alice" && e.ProblemCode == "1A" && e.Type == regatta.EventTypeProblemSolved {
			alice1A = append(alice1A, e)
		}
	}
	if len(alice1A) != 2 {
		t.Fatalf("alice 1A solved events = %d, want partial + full solve: %+v", len(alice1A), alice1A)
	}
	var fullSolve *regatta.RegattaEvent
	for i := range alice1A {
		if alice1A[i].TimeSec == 2374 {
			fullSolve = &alice1A[i]
			break
		}
	}
	if fullSolve == nil {
		t.Fatalf("want full solve at 2374, got %+v", alice1A)
	}
}

func TestUnifiedScoring_okAfterTourEndNoInTimeOrOvertake(t *testing.T) {
	tour := testTour()
	settings := regatta.DefaultScoringSettings()
	tr := CalculateResultWithSettings(tour, 0, []Run{
		{UserID: "alice", ProbID: 1, Time: 3700, Status: SubmissionStatusOK},
	}, settings)

	want := 100
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
