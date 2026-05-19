package regatta

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"contestr/pkg/regatta"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrTimetableNotFound           = errors.New("timetable not found")
	ErrTimetableAlreadyExists      = errors.New("timetable already exists")
	ErrContestNotFound             = errors.New("contest not found")
	ErrTourNotFound                = errors.New("tour not found in timetable")
	ErrInvalidTimetable            = errors.New("invalid timetable")
	ErrManualStartWithAutostart    = errors.New("manual start disabled while auto start is enabled")
	ErrNothingToAdvance            = errors.New("nothing to advance")
	ErrContestNotStarted           = errors.New("contest has not started yet")
	ErrNoActiveTour                = errors.New("no active tour or pause")
)

type TimetableRepository interface {
	Create(ctx context.Context, timetable *regatta.ToursTimetable) error
	GetByContestID(ctx context.Context, contestID int) (*regatta.ToursTimetable, error)
	Update(ctx context.Context, timetable *regatta.ToursTimetable) error
	DeleteByContestID(ctx context.Context, contestID int) error
}

type AdvanceMode int

const (
	AdvanceManual AdvanceMode = iota
	AdvanceAuto
)

func (s *Regatta) LoadTimetable(ctx context.Context, contestID int) (*regatta.ToursTimetable, error) {
	timetable, err := s.timetableRepository.GetByContestID(ctx, contestID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrTimetableNotFound
		}
		return nil, fmt.Errorf("failed to load timetable for contest %d: %w", contestID, err)
	}
	return timetable, nil
}

func (s *Regatta) insertTimetable(ctx context.Context, timetable regatta.ToursTimetable) (*regatta.ToursTimetable, error) {
	if err := validatePendingSlots(timetable.PendingSlots); err != nil {
		return nil, err
	}

	if _, err := s.timetableRepository.GetByContestID(ctx, timetable.ContestId); err == nil {
		return nil, ErrTimetableAlreadyExists
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("failed to check timetable: %w", err)
	}

	if err := s.timetableRepository.Create(ctx, &timetable); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrTimetableAlreadyExists
		}
		return nil, fmt.Errorf("failed to insert timetable: %w", err)
	}

	return &timetable, nil
}

func (s *Regatta) ReplaceTimetable(ctx context.Context, timetable regatta.ToursTimetable) (*regatta.ToursTimetable, error) {
	if err := validatePendingSlots(timetable.PendingSlots); err != nil {
		return nil, err
	}

	if _, err := s.LoadTimetable(ctx, timetable.ContestId); err != nil {
		return nil, err
	}

	if err := s.timetableRepository.Update(ctx, &timetable); err != nil {
		return nil, fmt.Errorf("failed to replace timetable: %w", err)
	}

	return &timetable, nil
}

func (s *Regatta) RemoveTimetable(ctx context.Context, contestID int) error {
	if contestID <= 0 {
		return fmt.Errorf("%w: contest_id must be positive", ErrInvalidTimetable)
	}

	if err := s.timetableRepository.DeleteByContestID(ctx, contestID); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrTimetableNotFound
		}
		return fmt.Errorf("failed to remove timetable: %w", err)
	}

	return nil
}

func (s *Regatta) contestAdvanceLock(contestID int) *sync.Mutex {
	s.advanceLockMu.Lock()
	defer s.advanceLockMu.Unlock()
	if s.advanceLocks == nil {
		s.advanceLocks = make(map[int]*sync.Mutex)
	}
	mu, ok := s.advanceLocks[contestID]
	if !ok {
		mu = &sync.Mutex{}
		s.advanceLocks[contestID] = mu
	}
	return mu
}

func (s *Regatta) AdvanceTimetable(
	ctx context.Context,
	contestID int,
	mode AdvanceMode,
	opts TimetableViewOptions,
) error {
	if contestID <= 0 {
		return fmt.Errorf("%w: contest_id must be positive", ErrInvalidTimetable)
	}

	if mode == AdvanceManual && opts.ServerAutoStartAvailable {
		timetable, err := s.LoadTimetable(ctx, contestID)
		if err == nil && timetable.AutoStartEnabled {
			return ErrManualStartWithAutostart
		}
	}

	mu := s.contestAdvanceLock(contestID)
	mu.Lock()
	defer mu.Unlock()

	contest, err := s.contestRepo.GetByContestID(ctx, contestID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrContestNotFound
		}
		return fmt.Errorf("failed to get contest %d: %w", contestID, err)
	}

	elapsed := int(time.Now().Sub(contest.StartTime).Seconds())
	if elapsed < 0 {
		return ErrContestNotStarted
	}

	timetable, err := s.LoadTimetable(ctx, contestID)
	if err != nil {
		return err
	}

	tours, err := s.loadToursSorted(ctx, contestID)
	if err != nil {
		return err
	}

	var changed bool

	if mode == AdvanceManual {
		if _, active := regatta.ActiveSequence(tours, elapsed); active {
			if err := s.finishActiveSegment(ctx, contestID, tours, elapsed, true); err != nil {
				return err
			}
			changed = true
			tours, err = s.loadToursSorted(ctx, contestID)
			if err != nil {
				return err
			}
		}
	} else {
		finished, err := s.tryFinishActiveSegment(ctx, contestID, tours, elapsed)
		if err != nil {
			return err
		}
		if finished {
			changed = true
			tours, err = s.loadToursSorted(ctx, contestID)
			if err != nil {
				return err
			}
		}
	}

	started, err := s.tryStartPendingHead(ctx, contestID, timetable, tours, elapsed, mode)
	if err != nil {
		return err
	}
	if started {
		changed = true
		if err := s.timetableRepository.Update(ctx, timetable); err != nil {
			return fmt.Errorf("failed to update timetable: %w", err)
		}
	}

	if !changed {
		return ErrNothingToAdvance
	}

	return nil
}

func (s *Regatta) tryFinishActiveSegment(
	ctx context.Context,
	contestID int,
	tours []regatta.Tour,
	elapsed int,
) (bool, error) {
	activeSeq, ok := regatta.ActiveSequence(tours, elapsed)
	if !ok {
		return false, nil
	}
	offsets := regatta.SegmentOffsets(tours)
	end := offsets[activeSeq].End
	if elapsed < end {
		return false, nil
	}
	return true, s.shortenActiveToElapsed(ctx, contestID, tours, activeSeq, elapsed)
}

func (s *Regatta) finishActiveSegment(
	ctx context.Context,
	contestID int,
	tours []regatta.Tour,
	elapsed int,
	force bool,
) error {
	activeSeq, ok := regatta.ActiveSequence(tours, elapsed)
	if !ok {
		return nil
	}
	if !force {
		return nil
	}
	return s.shortenActiveToElapsed(ctx, contestID, tours, activeSeq, elapsed)
}

func (s *Regatta) shortenActiveToElapsed(
	ctx context.Context,
	contestID int,
	tours []regatta.Tour,
	activeSeq int,
	elapsed int,
) error {
	elapsedIn := regatta.ElapsedInSegment(tours, activeSeq, elapsed)
	if elapsedIn <= 0 {
		return nil
	}
	return s.tourRepository.UpdateDuration(ctx, contestID, activeSeq, elapsedIn)
}

func (s *Regatta) tryStartPendingHead(
	ctx context.Context,
	contestID int,
	timetable *regatta.ToursTimetable,
	tours []regatta.Tour,
	elapsed int,
	mode AdvanceMode,
) (bool, error) {
	slot, ok := timetable.HeadSlot()
	if !ok {
		return false, nil
	}

	if _, active := regatta.ActiveSequence(tours, elapsed); active {
		return false, nil
	}

	anchor := regatta.TimelineAnchorEnd(tours)
	start := regatta.BuildPendingStarts(anchor, timetable.PendingSlots)[0]

	if mode == AdvanceAuto && elapsed < start {
		return false, nil
	}

	isPause := regatta.NormalizeSlotKind(slot.Kind) == regatta.ScheduleSlotKindPause
	if _, err := s.StartTour(ctx, contestID, slot.Duration, StartTourOptions{IsPause: isPause}); err != nil {
		return false, err
	}

	timetable.PopHead()
	return true, nil
}

func (s *Regatta) loadToursSorted(ctx context.Context, contestID int) ([]regatta.Tour, error) {
	tours, err := s.tourRepository.FindByContestID(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("failed to find tours: %w", err)
	}
	return regatta.SortToursBySequence(tours), nil
}

func (s *Regatta) UpdateActiveTourDuration(
	ctx context.Context,
	contestID int,
	durationSeconds int,
	opts TimetableViewOptions,
) (*regatta.TimetableView, error) {
	if contestID <= 0 {
		return nil, fmt.Errorf("%w: contest_id must be positive", ErrInvalidTimetable)
	}
	if durationSeconds <= 0 {
		return nil, fmt.Errorf("%w: duration must be positive", ErrInvalidTimetable)
	}

	contest, err := s.contestRepo.GetByContestID(ctx, contestID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrContestNotFound
		}
		return nil, fmt.Errorf("failed to get contest %d: %w", contestID, err)
	}

	elapsed := int(time.Now().Sub(contest.StartTime).Seconds())
	if elapsed < 0 {
		return nil, ErrContestNotStarted
	}

	tours, err := s.loadToursSorted(ctx, contestID)
	if err != nil {
		return nil, err
	}

	activeSeq, ok := regatta.ActiveSequence(tours, elapsed)
	if !ok {
		return nil, ErrNoActiveTour
	}

	elapsedIn := regatta.ElapsedInSegment(tours, activeSeq, elapsed)
	if durationSeconds < elapsedIn {
		return nil, fmt.Errorf(
			"%w: duration cannot be less than elapsed time in segment (%ds)",
			ErrInvalidTimetable,
			elapsedIn,
		)
	}

	if err := s.tourRepository.UpdateDuration(ctx, contestID, activeSeq, durationSeconds); err != nil {
		return nil, fmt.Errorf("failed to update tour duration: %w", err)
	}

	timetable, err := s.LoadTimetable(ctx, contestID)
	if errors.Is(err, ErrTimetableNotFound) {
		return s.buildTimetableView(ctx, contestID, &regatta.ToursTimetable{
			ContestId:    contestID,
			PendingSlots: []regatta.ScheduleSlot{},
		}, opts)
	}
	if err != nil {
		return nil, err
	}

	return s.buildTimetableView(ctx, contestID, timetable, opts)
}

func validatePendingSlots(slots []regatta.ScheduleSlot) error {
	for i, slot := range slots {
		if slot.Duration <= 0 {
			return fmt.Errorf("%w: slot %d duration must be positive", ErrInvalidTimetable, i+1)
		}
		kind := regatta.NormalizeSlotKind(slot.Kind)
		if kind != regatta.ScheduleSlotKindTour && kind != regatta.ScheduleSlotKindPause {
			return fmt.Errorf("%w: slot %d invalid kind %q", ErrInvalidTimetable, i+1, slot.Kind)
		}
	}
	return nil
}
