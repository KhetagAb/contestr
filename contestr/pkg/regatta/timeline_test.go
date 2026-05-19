package regatta

import "testing"

func TestSegmentOffsets_withPause(t *testing.T) {
	tours := []Tour{
		{Sequence: 1, Round: 1, DurationInSeconds: 1000, IsPause: false},
		{Sequence: 2, Round: 0, DurationInSeconds: 300, IsPause: true},
		{Sequence: 3, Round: 2, DurationInSeconds: 2000, IsPause: false},
	}
	offsets := SegmentOffsets(tours)

	if offsets[1].Start != 0 || offsets[1].End != 1000 {
		t.Fatalf("tour 1: %+v", offsets[1])
	}
	if offsets[2].Start != 1000 || offsets[2].End != 1300 {
		t.Fatalf("pause: %+v", offsets[2])
	}
	if offsets[3].Start != 1300 || offsets[3].End != 3300 {
		t.Fatalf("tour 2: %+v", offsets[3])
	}
}

func TestActiveSequence(t *testing.T) {
	tours := []Tour{
		{Sequence: 1, DurationInSeconds: 100},
	}
	seq, ok := ActiveSequence(tours, 50)
	if !ok || seq != 1 {
		t.Fatalf("got seq=%d ok=%v", seq, ok)
	}
	_, ok = ActiveSequence(tours, 100)
	if ok {
		t.Fatal("expected no active at boundary end")
	}
}

func TestNextCompetitiveRound_skipsPausePending(t *testing.T) {
	tours := []Tour{{Sequence: 1, Round: 1, IsPause: false}}
	pending := []ScheduleSlot{
		{Duration: 60, Kind: ScheduleSlotKindPause},
		{Duration: 120, Kind: ScheduleSlotKindTour},
	}
	round := NextCompetitiveRound(tours, pending)
	if round == nil || *round != 2 {
		t.Fatalf("got %v want 2", round)
	}
}

func TestBuildPendingStarts(t *testing.T) {
	starts := BuildPendingStarts(1000, []ScheduleSlot{
		{Duration: 100},
		{Duration: 200},
	})
	if starts[0] != 1000 || starts[1] != 1100 {
		t.Fatalf("got %v", starts)
	}
}
