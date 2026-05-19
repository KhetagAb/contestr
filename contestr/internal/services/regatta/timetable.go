package regatta

import (
	"context"
	"errors"
	"fmt"
	"time"

	"contestr/pkg/regatta"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrTimetableNotFound      = errors.New("timetable not found")
	ErrTimetableAlreadyExists = errors.New("timetable already exists")
	ErrContestNotFound        = errors.New("contest not found")
	ErrTourNotFound           = errors.New("tour not found in timetable")
	ErrTourAlreadyStarted     = errors.New("tour already started")
	ErrInvalidTimetable       = errors.New("invalid timetable")
)

type TimetableRepository interface {
	Create(ctx context.Context, timetable *regatta.ToursTimetable) error
	GetByContestID(ctx context.Context, contestID int) (*regatta.ToursTimetable, error)
	Update(ctx context.Context, timetable *regatta.ToursTimetable) error
	DeleteByContestID(ctx context.Context, contestID int) error
}

// LoadTimetable returns the persisted timetable for a contest.
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

// insertTimetable creates a timetable document; fails if one already exists.
func (s *Regatta) insertTimetable(ctx context.Context, timetable regatta.ToursTimetable) (*regatta.ToursTimetable, error) {
	regatta.RebuildChain(timetable.TourTimes)
	if err := validateTourChain(timetable); err != nil {
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

// ReplaceTimetable overwrites an existing timetable document.
func (s *Regatta) ReplaceTimetable(ctx context.Context, timetable regatta.ToursTimetable) (*regatta.ToursTimetable, error) {
	regatta.RebuildChain(timetable.TourTimes)
	if err := validateTourChain(timetable); err != nil {
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

// RemoveTimetable deletes the timetable for a contest.
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

// StartScheduledTour manually starts the first not-started tour at the current contest time.
func (s *Regatta) StartScheduledTour(ctx context.Context, contestID int, tourNumber int) (*regatta.ToursTimetable, error) {
	if contestID <= 0 {
		return nil, fmt.Errorf("%w: contest_id must be positive", ErrInvalidTimetable)
	}
	if tourNumber <= 0 {
		return nil, fmt.Errorf("%w: tour_number must be positive", ErrInvalidTimetable)
	}

	timetable, err := s.LoadTimetable(ctx, contestID)
	if err != nil {
		return nil, err
	}

	position := tourNumber - 1
	if position >= len(timetable.TourTimes) {
		return nil, ErrTourNotFound
	}

	firstNotStartedTourNumber, _, ok := timetable.FirstNotStartedTour()
	if !ok {
		return nil, ErrTourNotFound
	}
	if firstNotStartedTourNumber != tourNumber {
		return nil, fmt.Errorf("%w: only the first not started tour can be started", ErrInvalidTimetable)
	}

	tour := timetable.TourTimes[position]
	if tour.Started {
		return nil, ErrTourAlreadyStarted
	}

	contest, err := s.contestRepo.GetByContestID(ctx, contestID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrContestNotFound
		}
		return nil, fmt.Errorf("failed to get contest %d: %w", contestID, err)
	}

	startTime := int(time.Since(contest.StartTime).Seconds())
	if startTime < 0 {
		return nil, fmt.Errorf("%w: contest has not started yet", ErrInvalidTimetable)
	}

	shiftToursFromIndex(timetable, position, startTime)
	timetable.TourTimes[position].Started = true

	if err := validateTourChain(*timetable); err != nil {
		return nil, err
	}

	if _, err := s.StartTour(ctx, contestID, time.Duration(tour.Duration)*time.Second); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrContestNotFound
		}
		return nil, err
	}

	if err := s.timetableRepository.Update(ctx, timetable); err != nil {
		return nil, fmt.Errorf("failed to update timetable after starting tour: %w", err)
	}

	return timetable, nil
}

func shiftToursFromIndex(timetable *regatta.ToursTimetable, position int, newStartTime int) {
	delta := newStartTime - timetable.TourTimes[position].StartTime
	for i := position; i < len(timetable.TourTimes); i++ {
		timetable.TourTimes[i].StartTime += delta
	}
}

func validateTourChain(timetable regatta.ToursTimetable) error {
	if timetable.ContestId <= 0 {
		return fmt.Errorf("%w: contest_id must be positive", ErrInvalidTimetable)
	}

	for i, tour := range timetable.TourTimes {
		if tour.StartTime < 0 {
			return fmt.Errorf("%w: tour %d start_time must be non-negative", ErrInvalidTimetable, i+1)
		}
		if tour.Duration <= 0 {
			return fmt.Errorf("%w: tour %d duration must be positive", ErrInvalidTimetable, i+1)
		}
		if i > 0 {
			prevTour := timetable.TourTimes[i-1]
			if prevTour.StartTime+prevTour.Duration > tour.StartTime {
				return fmt.Errorf("%w: tours %d and %d overlap", ErrInvalidTimetable, i, i+1)
			}
		}
	}

	return nil
}

func validateTourDurations(durations []int) error {
	for i, duration := range durations {
		if duration <= 0 {
			return fmt.Errorf("%w: tour %d duration must be positive", ErrInvalidTimetable, i+1)
		}
	}
	return nil
}
