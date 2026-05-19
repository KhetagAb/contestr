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

)

type Interval time.Duration

type RegattaService interface {
	AdvanceTimetable(ctx context.Context, contestID int, mode regattasvc.AdvanceMode, opts regattasvc.TimetableViewOptions) error
	LoadTimetable(ctx context.Context, contestID int) (*regattapkg.ToursTimetable, error)
}

type TimetableSyncService struct {
	registry     contest_registry.ContestRegistry
	contestRepo  repository.ContestRepository
	regatta      RegattaService
	syncInterval time.Duration
	viewOpts     regattasvc.TimetableViewOptions
}

func NewTimetableSyncService(
	registry contest_registry.ContestRegistry,
	contestRepo repository.ContestRepository,
	regatta RegattaService,
	syncInterval Interval,
	autoStartAvailable bool,
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
		viewOpts: regattasvc.TimetableViewOptions{
			ServerAutoStartAvailable: autoStartAvailable,
		},
	}
}

func (s *TimetableSyncService) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()

	s.AdvanceDue(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			s.AdvanceDue(ctx)
		}
	}
}

func (s *TimetableSyncService) AdvanceDue(ctx context.Context) {
	contestsBySystem := s.registry.GetAllContests()

	for system, contestIDs := range contestsBySystem {
		for _, contestID := range contestIDs {
			if err := s.advanceContest(ctx, contestID); err != nil {
				logger.Errorf(ctx, "failed to advance timetable for contest %d (%s): %v", contestID, system, err)
			}
		}
	}
}

func (s *TimetableSyncService) advanceContest(ctx context.Context, contestID int) error {
	timetable, err := s.regatta.LoadTimetable(ctx, contestID)
	if err != nil {
		if errors.Is(err, regattasvc.ErrTimetableNotFound) {
			return nil
		}
		return err
	}

	if !timetable.AutoStartEnabled {
		return nil
	}

	err = s.regatta.AdvanceTimetable(ctx, contestID, regattasvc.AdvanceAuto, s.viewOpts)
	if errors.Is(err, regattasvc.ErrNothingToAdvance) {
		return nil
	}
	if err != nil {
		return err
	}

	logger.Infof(ctx, "advanced timetable for contest %d", contestID)
	return nil
}
