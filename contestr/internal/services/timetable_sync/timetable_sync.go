package timetable_sync

import (
	"context"
	"errors"
	"time"

	"contestr/internal/repository"
	"contestr/internal/services/contest_registry"
	regattasvc "contestr/internal/services/regatta"
	"contestr/pkg/logger"
	regattapkg "contestr/pkg/regatta"

	"go.mongodb.org/mongo-driver/mongo"
)

type Interval time.Duration

type RegattaService interface {
	StartTour(ctx context.Context, contestID int, duration time.Duration) (string, error)
	GetTimetable(ctx context.Context, contestID int) (*regattapkg.ToursTimetable, error)
	UpdateTimetable(ctx context.Context, timetable regattapkg.ToursTimetable) (*regattapkg.ToursTimetable, error)
}

type TimetableSyncService struct {
	registry     contest_registry.ContestRegistry
	contestRepo  repository.ContestRepository
	regatta      RegattaService
	syncInterval time.Duration
}

func NewTimetableSyncService(
	registry contest_registry.ContestRegistry,
	contestRepo repository.ContestRepository,
	regatta RegattaService,
	syncInterval Interval,
) *TimetableSyncService {
	interval := time.Duration(syncInterval)
	if interval <= 0 {
		interval = time.Second
	}

	return &TimetableSyncService{
		registry:     registry,
		contestRepo:  contestRepo,
		regatta:      regatta,
		syncInterval: interval,
	}
}

func (s *TimetableSyncService) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()

	s.StartDueTours(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			s.StartDueTours(ctx)
		}
	}
}

func (s *TimetableSyncService) StartDueTours(ctx context.Context) {
	now := time.Now()
	contestsBySystem := s.registry.GetAllContests()

	for system, contestIDs := range contestsBySystem {
		for _, contestID := range contestIDs {
			if err := s.startDueTour(ctx, contestID, now); err != nil {
				logger.Errorf(ctx, "failed to start due timetable tour for contest %d (%s): %v", contestID, system, err)
			}
		}
	}
}

func (s *TimetableSyncService) startDueTour(ctx context.Context, contestID int, now time.Time) error {
	timetable, err := s.regatta.GetTimetable(ctx, contestID)
	if err != nil {
		if errors.Is(err, regattasvc.ErrTimetableNotFound) {
			return nil
		}
		return err
	}

	contest, err := s.contestRepo.GetByContestID(ctx, contestID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return err
	}

	tourNumber, tour, ok := timetable.FirstNotStartedTour()
	if !ok {
		return nil
	}

	scheduledStart := contest.StartTime.Add(time.Duration(tour.StartTime) * time.Second)
	if now.Before(scheduledStart) {
		return nil
	}

	objectID, err := s.regatta.StartTour(ctx, contestID, time.Duration(tour.Duration)*time.Second)
	if err != nil {
		return err
	}

	timetable.TourTimes[tourNumber-1].Started = true
	if _, err := s.regatta.UpdateTimetable(ctx, *timetable); err != nil {
		return err
	}

	logger.Infof(ctx, "started timetable tour %d for contest %d: object_id=%s", tourNumber, contestID, objectID)
	return nil
}
