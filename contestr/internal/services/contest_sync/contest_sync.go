package contest_sync

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"contestr/internal/integrations"
	"contestr/internal/repository"
	"contestr/internal/services/contest_registry"
	"contestr/pkg/logger"
	"contestr/pkg/regatta"

	"go.mongodb.org/mongo-driver/mongo"
)

const systemCodeforces = "codeforces"

type ContestSyncService struct {
	registry            contest_registry.ContestRegistry
	adapters            map[string]integrations.ContestAdapter
	contestRepo         repository.ContestRepository
	syncInterval        time.Duration
	intervalBeforeStart time.Duration
}

func NewContestSyncService(
	registry contest_registry.ContestRegistry,
	adapters map[string]integrations.ContestAdapter,
	contestRepo repository.ContestRepository,
	syncInterval time.Duration,
	intervalBeforeStart time.Duration,
) *ContestSyncService {
	if intervalBeforeStart <= 0 {
		intervalBeforeStart = time.Minute
	}
	return &ContestSyncService{
		registry:            registry,
		adapters:            adapters,
		contestRepo:         contestRepo,
		syncInterval:        syncInterval,
		intervalBeforeStart: intervalBeforeStart,
	}
}

func (s *ContestSyncService) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()

	s.SyncPeriodic(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			s.SyncPeriodic(ctx)
		}
	}
}

type SyncResult struct {
	SyncedCount   int
	FailedCount   int
	TotalCount    int
	HasContests   bool
	ErrorMessages []string
}

func (s *ContestSyncService) SyncPeriodic(ctx context.Context) *SyncResult {
	result := &SyncResult{
		ErrorMessages: make([]string, 0),
	}

	contestsBySystem := s.registry.GetAllContests()

	totalContests := 0
	for _, contestIDs := range contestsBySystem {
		totalContests += len(contestIDs)
	}

	if totalContests == 0 {
		logger.Infof(ctx, "no contests to sync")
		result.HasContests = false
		return result
	}

	result.HasContests = true
	result.TotalCount = totalContests

	now := time.Now()

	for system, contestIDs := range contestsBySystem {
		if len(contestIDs) == 0 {
			continue
		}

		for _, contestID := range contestIDs {
			cached, err := s.contestRepo.GetByContestID(ctx, contestID)
			if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
				logger.Errorf(ctx, "failed to load cached contest %d for sync gating: %v", contestID, err)
			}
			if !s.shouldSyncPeriodic(system, cached, now) {
				continue
			}

			if err := s.SyncContest(ctx, contestID); err != nil {
				errMsg := fmt.Sprintf("failed to sync contest %d (%s): %v", contestID, system, err)
				logger.Errorf(ctx, errMsg)
				result.ErrorMessages = append(result.ErrorMessages, errMsg)
				result.FailedCount++
				continue
			}

			result.SyncedCount++
		}
	}

	return result
}

func (s *ContestSyncService) shouldSyncPeriodic(system string, cached *regatta.Contest, now time.Time) bool {
	if cached == nil {
		return true
	}

	if system == systemCodeforces {
		return ShouldSyncCf(cached.Phase, cached.LastUpdated, now, s.syncInterval, s.intervalBeforeStart)
	}

	return ShouldSyncEjudge(cached.LastUpdated, now, s.syncInterval)
}

func (s *ContestSyncService) SyncContest(ctx context.Context, contestID int) error {
	registered, err := s.registry.GetContest(contestID)
	if err != nil {
		return err
	}
	system := registered.System

	adapter, ok := s.adapters[system]
	if !ok {
		return fmt.Errorf("adapter for system %s not found", system)
	}

	contest, err := adapter.FetchContest(ctx, contestID, integrations.FetchContestOptions{
		ScoringSettings: registered.ScoringSettings,
		TourSettings:    registered.TourSettings,
	})
	if err != nil {
		return fmt.Errorf("failed to sync contest %d (%s): %w", contestID, system, err)
	}

	applyRegisteredContestName(contest, registered)

	if err := s.contestRepo.Upsert(ctx, contest); err != nil {
		return fmt.Errorf("failed to save contest %d: %w", contestID, err)
	}

	logger.Infof(ctx, "successfully synced contest %d (%s)", contestID, system)
	return nil
}

func applyRegisteredContestName(contest *regatta.Contest, registered *repository.RegisteredContest) {
	if contest == nil || registered == nil {
		return
	}
	if name := strings.TrimSpace(registered.Name); name != "" {
		contest.ContestName = name
	}
}
