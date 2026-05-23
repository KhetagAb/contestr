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
	sorted := SortToursBySequence(tours)
	offsets := SegmentOffsets(sorted)
	activeSeq, hasActive := ActiveSequence(sorted, elapsed)

	anchor := TimelineAnchorEnd(sorted)
	pendingStarts := BuildPendingStarts(anchor, pending)

	segments := buildFactualSegments(sorted, offsets, activeSeq, hasActive)
	segments = append(segments, buildPendingSegments(pending, pendingStarts, elapsed, hasActive, sorted)...)
	return segments
}

func buildFactualSegments(sorted []Tour, offsets map[int]SegmentOffset, activeSeq int, hasActive bool) []TimelineSegment {
	segments := make([]TimelineSegment, 0, len(sorted))
	for _, tour := range sorted {
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
	return segments
}

func buildPendingSegments(pending []ScheduleSlot, pendingStarts []int, elapsed int, hasActive bool, tours []Tour) []TimelineSegment {
	firstOverdue, nextIdx := pendingOverdueIndexes(pendingStarts, elapsed)

	tourRoundBase := CompetitiveRoundCount(tours)
	competitivePendingCount := 0

	segments := make([]TimelineSegment, 0, len(pending))
	for i, slot := range pending {
		kind := NormalizeSlotKind(slot.Kind)
		idx := i

		var roundPtr *int
		if kind == ScheduleSlotKindTour {
			competitivePendingCount++
			r := tourRoundBase + competitivePendingCount
			roundPtr = &r
		}

		segments = append(segments, TimelineSegment{
			PendingIndex: &idx,
			Kind:         kind,
			Round:        roundPtr,
			Duration:     slot.Duration,
			StartTime:    pendingStarts[i],
			Status:       pendingSlotStatus(i, pendingStarts[i], elapsed, hasActive, firstOverdue, nextIdx),
			Editable:     true,
		})
	}
	return segments
}

func pendingOverdueIndexes(pendingStarts []int, elapsed int) (firstOverdue, nextIdx int) {
	firstOverdue, nextIdx = -1, -1
	for i, start := range pendingStarts {
		if elapsed >= start {
			if firstOverdue < 0 {
				firstOverdue = i
			}
		} else if nextIdx < 0 {
			nextIdx = i
		}
	}
	return
}

func pendingSlotStatus(i int, start int, elapsed int, hasActive bool, firstOverdue int, nextIdx int) string {
	if elapsed < start {
		if i == nextIdx {
			return SegmentStatusNext
		}
		return SegmentStatusFuture
	}
	// слот просрочен: активный сегмент ещё идёт — только первый просроченный помечается как starting
	if hasActive {
		if i == firstOverdue {
			return SegmentStatusStarting
		}
		return SegmentStatusFuture
	}
	switch i {
	case firstOverdue:
		return SegmentStatusStarting
	case firstOverdue + 1:
		return SegmentStatusNext
	}
	return SegmentStatusFuture
}
