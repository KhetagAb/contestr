package regatta

// RebuildChain recalculates start_time for not-started tours so each tour starts
// when the previous one ends. Started tours keep their start_time unchanged.
func RebuildChain(tourTimes []TourConfig) {
	if len(tourTimes) == 0 {
		return
	}

	anchor := len(tourTimes)
	for i, tour := range tourTimes {
		if !tour.Started {
			anchor = i
			break
		}
	}

	if anchor >= len(tourTimes) {
		return
	}

	if anchor == 0 {
		tourTimes[0].StartTime = 0
	} else {
		prev := tourTimes[anchor-1]
		tourTimes[anchor].StartTime = prev.StartTime + prev.Duration
	}

	for i := anchor + 1; i < len(tourTimes); i++ {
		prev := tourTimes[i-1]
		tourTimes[i].StartTime = prev.StartTime + prev.Duration
	}
}

// ApplyTimetableSchedule applies tour durations while preserving started flags
// and start times of already started tours, then rebuilds the chain.
func ApplyTimetableSchedule(existing *ToursTimetable, durations []int) []TourConfig {
	result := make([]TourConfig, len(durations))
	for i, duration := range durations {
		result[i] = TourConfig{Duration: duration}
		if existing != nil && i < len(existing.TourTimes) {
			result[i].Started = existing.TourTimes[i].Started
			if existing.TourTimes[i].Started {
				result[i].StartTime = existing.TourTimes[i].StartTime
			}
		}
	}
	RebuildChain(result)
	return result
}

// BuildToursMeta computes UI status metadata for each tour.
func BuildToursMeta(tourTimes []TourConfig, elapsedSeconds int) ([]TourMeta, *int) {
	meta := make([]TourMeta, len(tourTimes))
	var nextTourNumber *int

	firstPending := -1
	for i, tour := range tourTimes {
		if !tour.Started && firstPending < 0 {
			firstPending = i
			n := i + 1
			nextTourNumber = &n
		}
	}

	for i, tour := range tourTimes {
		scheduledStart := tour.StartTime
		status := TourStatusPlanned

		switch {
		case tour.Started:
			status = TourStatusStarted
		case firstPending == i && elapsedSeconds < scheduledStart:
			status = TourStatusNext
		case firstPending == i && elapsedSeconds >= scheduledStart:
			status = TourStatusStarting
		}

		meta[i] = TourMeta{
			TourNumber: i + 1,
			Status:     status,
		}
	}

	return meta, nextTourNumber
}
