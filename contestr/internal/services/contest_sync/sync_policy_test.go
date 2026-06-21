package contest_sync

import (
	"testing"
	"time"
)

func TestShouldSyncCf(t *testing.T) {
	interval := 15 * time.Second
	beforeStart := time.Minute
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)

	t.Run("FINISHED", func(t *testing.T) {
		if ShouldSyncCf(codeforcesPhaseFinished, now.Add(-time.Hour), now, interval, beforeStart) {
			t.Fatal("expected false for FINISHED")
		}
	})

	t.Run("BEFORE recent", func(t *testing.T) {
		if ShouldSyncCf(codeforcesPhaseBefore, now.Add(-30*time.Second), now, interval, beforeStart) {
			t.Fatal("expected false for recent BEFORE sync")
		}
	})

	t.Run("BEFORE due", func(t *testing.T) {
		if !ShouldSyncCf(codeforcesPhaseBefore, now.Add(-61*time.Second), now, interval, beforeStart) {
			t.Fatal("expected true for BEFORE sync after interval_before_start")
		}
	})

	t.Run("CODING recent", func(t *testing.T) {
		if ShouldSyncCf("CODING", now.Add(-10*time.Second), now, interval, beforeStart) {
			t.Fatal("expected false for recent CODING sync")
		}
	})

	t.Run("CODING due", func(t *testing.T) {
		if !ShouldSyncCf("CODING", now.Add(-16*time.Second), now, interval, beforeStart) {
			t.Fatal("expected true for CODING sync after interval")
		}
	})

	t.Run("empty phase due", func(t *testing.T) {
		if !ShouldSyncCf("", now.Add(-16*time.Second), now, interval, beforeStart) {
			t.Fatal("expected true for empty phase using active interval")
		}
	})
}

func TestShouldSyncEjudge(t *testing.T) {
	interval := 15 * time.Second
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)

	if ShouldSyncEjudge(now.Add(-5*time.Second), now, interval) {
		t.Fatal("expected false when interval not elapsed")
	}
	if !ShouldSyncEjudge(now.Add(-20*time.Second), now, interval) {
		t.Fatal("expected true when interval elapsed")
	}
}
