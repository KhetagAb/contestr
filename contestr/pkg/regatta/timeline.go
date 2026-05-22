package regatta

import "slices"

const (
	ScheduleSlotKindTour  = "tour"
	ScheduleSlotKindPause = "pause"

	SegmentStatusPast     = "past"
	SegmentStatusActive   = "active"
	SegmentStatusNext     = "next"
	SegmentStatusStarting = "starting"
	SegmentStatusFuture   = "future"
)

type SegmentOffset struct {
	Start int
	End   int
}

func SortToursBySequence(tours []Tour) []Tour {
	out := append([]Tour(nil), tours...)
	slices.SortFunc(out, func(a, b Tour) int {
		return a.Sequence - b.Sequence
	})
	return out
}

func SegmentOffsets(tours []Tour) map[int]SegmentOffset {
	sorted := SortToursBySequence(tours)
	offsets := make(map[int]SegmentOffset, len(sorted))
	cursor := 0
	for _, tour := range sorted {
		start := cursor
		end := start + tour.DurationInSeconds
		offsets[tour.Sequence] = SegmentOffset{Start: start, End: end}
		cursor = end
	}
	return offsets
}

func TimelineAnchorEnd(tours []Tour) int {
	sorted := SortToursBySequence(tours)
	if len(sorted) == 0 {
		return 0
	}
	offsets := SegmentOffsets(tours)
	last := sorted[len(sorted)-1]
	return offsets[last.Sequence].End
}

func ActiveSequence(tours []Tour, elapsed int) (int, bool) {
	if elapsed < 0 {
		return 0, false
	}
	offsets := SegmentOffsets(tours)
	for _, tour := range SortToursBySequence(tours) {
		o := offsets[tour.Sequence]
		if elapsed >= o.Start && elapsed < o.End {
			return tour.Sequence, true
		}
	}
	return 0, false
}

func ElapsedInSegment(tours []Tour, sequence int, elapsed int) int {
	offsets := SegmentOffsets(tours)
	o, ok := offsets[sequence]
	if !ok {
		return 0
	}
	if elapsed < o.Start {
		return 0
	}
	if elapsed >= o.End {
		return o.End - o.Start
	}
	return elapsed - o.Start
}

func BuildPendingStarts(anchor int, slots []ScheduleSlot) []int {
	starts := make([]int, len(slots))
	cursor := anchor
	for i, slot := range slots {
		starts[i] = cursor
		cursor += slot.Duration
	}
	return starts
}

func CompetitiveRoundCount(tours []Tour) int {
	n := 0
	for _, tour := range tours {
		if !tour.IsPause {
			n++
		}
	}
	return n
}

func NextCompetitiveRound(tours []Tour, pending []ScheduleSlot) *int {
	round := CompetitiveRoundCount(tours)
	for _, slot := range pending {
		kind := slot.Kind
		if kind == "" {
			kind = ScheduleSlotKindTour
		}
		if kind == ScheduleSlotKindTour {
			round++
			return &round
		}
	}
	return nil
}

func CompetitiveTourByRound(tours []Tour, round int) *Tour {
	for i := range tours {
		tour := &tours[i]
		if !tour.IsPause && tour.Round == round {
			return tour
		}
	}
	return nil
}

func NormalizeSlotKind(kind string) string {
	if kind == ScheduleSlotKindPause {
		return ScheduleSlotKindPause
	}
	return ScheduleSlotKindTour
}

func BuildTimelineSegments(tours []Tour, pending []ScheduleSlot, elapsed int) []TimelineSegment {
	var segments []TimelineSegment
	offsets := SegmentOffsets(tours)

	activeSeq, hasActive := ActiveSequence(tours, elapsed)
	anchor := TimelineAnchorEnd(tours)
	pendingStarts := BuildPendingStarts(anchor, pending)

	firstOverduePending := -1
	nextPendingIndex := -1
	for i, start := range pendingStarts {
		if elapsed >= start {
			if firstOverduePending < 0 {
				firstOverduePending = i
			}
		} else if nextPendingIndex < 0 {
			nextPendingIndex = i
		}
	}

	for _, tour := range SortToursBySequence(tours) {
		o := offsets[tour.Sequence]
		status := SegmentStatusPast
		if hasActive && tour.Sequence == activeSeq {
			status = SegmentStatusActive
		}
		kind := ScheduleSlotKindTour
		round := tour.Round
		if tour.IsPause {
			kind = ScheduleSlotKindPause
			round = 0
		}
		seq := tour.Sequence
		var roundPtr *int
		if round > 0 {
			roundPtr = &round
		}
		segments = append(segments, TimelineSegment{
			Sequence:  &seq,
			Kind:      kind,
			Round:     roundPtr,
			Duration:  tour.DurationInSeconds,
			StartTime: o.Start,
			Status:    status,
			Editable:  status == SegmentStatusActive,
		})
	}

	for i, slot := range pending {
		kind := NormalizeSlotKind(slot.Kind)
		start := pendingStarts[i]
		status := SegmentStatusFuture
		if elapsed < start {
			if i == nextPendingIndex {
				status = SegmentStatusNext
			}
		} else if !hasActive {
			if firstOverduePending >= 0 {
				switch i {
				case firstOverduePending:
					status = SegmentStatusStarting
				case firstOverduePending + 1:
					status = SegmentStatusNext
				}
			} else if i == 0 {
				status = SegmentStatusNext
			}
		} else if i == 0 {
			status = SegmentStatusStarting
		}

		idx := i
		var roundPtr *int
		if kind == ScheduleSlotKindTour {
			r := CompetitiveRoundCount(tours)
			for j := 0; j <= i; j++ {
				if NormalizeSlotKind(pending[j].Kind) == ScheduleSlotKindTour {
					r++
				}
			}
			roundPtr = &r
		}

		segments = append(segments, TimelineSegment{
			PendingIndex: &idx,
			Kind:         kind,
			Round:        roundPtr,
			Duration:     slot.Duration,
			StartTime:    start,
			Status:       status,
			Editable:     true,
		})
	}

	return segments
}
