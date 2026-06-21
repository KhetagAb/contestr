package contest_sync

import (
	"time"
)

const codeforcesPhaseFinished = "FINISHED"
const codeforcesPhaseBefore = "BEFORE"

func ShouldSyncCf(phase string, lastUpdated, now time.Time, interval, intervalBeforeStart time.Duration) bool {
	if phase == codeforcesPhaseFinished {
		return false
	}
	if phase == codeforcesPhaseBefore {
		return now.Sub(lastUpdated) >= intervalBeforeStart
	}
	return now.Sub(lastUpdated) >= interval
}

func ShouldSyncEjudge(lastUpdated, now time.Time, interval time.Duration) bool {
	return now.Sub(lastUpdated) >= interval
}
