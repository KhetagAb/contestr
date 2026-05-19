package regatta

import "testing"

func TestRebuildChainFromZero(t *testing.T) {
	tours := []TourConfig{
		{Duration: 1800, Started: false},
		{Duration: 1200, Started: false},
	}
	RebuildChain(tours)
	if tours[0].StartTime != 0 {
		t.Fatalf("tour 1 start: got %d want 0", tours[0].StartTime)
	}
	if tours[1].StartTime != 1800 {
		t.Fatalf("tour 2 start: got %d want 1800", tours[1].StartTime)
	}
}

func TestRebuildChainAfterStarted(t *testing.T) {
	tours := []TourConfig{
		{StartTime: 0, Duration: 1800, Started: true},
		{StartTime: 1800, Duration: 1200, Started: false},
		{StartTime: 9999, Duration: 900, Started: false},
	}
	RebuildChain(tours)
	if tours[0].StartTime != 0 {
		t.Fatalf("started tour start changed: %d", tours[0].StartTime)
	}
	if tours[1].StartTime != 1800 {
		t.Fatalf("tour 2 start: got %d want 1800", tours[1].StartTime)
	}
	if tours[2].StartTime != 3000 {
		t.Fatalf("tour 3 start: got %d want 3000", tours[2].StartTime)
	}
}

func TestApplyTimetableSchedulePreservesStarted(t *testing.T) {
	existing := &ToursTimetable{
		TourTimes: []TourConfig{
			{StartTime: 0, Duration: 1800, Started: true},
			{StartTime: 1800, Duration: 1200, Started: false},
		},
	}
	merged := ApplyTimetableSchedule(existing, []int{2000, 1500})
	if !merged[0].Started || merged[0].StartTime != 0 {
		t.Fatalf("started tour 1: %+v", merged[0])
	}
	if merged[1].StartTime != 2000 {
		t.Fatalf("tour 2 start: got %d want 2000", merged[1].StartTime)
	}
}
