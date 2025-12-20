package contest_sync

import (
	"contestr/internal/integrations"
	"contestr/internal/repository"
	"contestr/internal/services/contest_registry"
	"contestr/pkg/logger"
	"context"
	"time"
)

type ContestSyncService struct {
	registry     contest_registry.ContestRegistry
	adapters     map[string]integrations.ContestAdapter
	contestRepo  repository.ContestRepository
	syncInterval time.Duration
}

func NewContestSyncService(
	registry contest_registry.ContestRegistry,
	adapters map[string]integrations.ContestAdapter,
	contestRepo repository.ContestRepository,
	syncInterval time.Duration,
) *ContestSyncService {
	return &ContestSyncService{
		registry:     registry,
		adapters:     adapters,
		contestRepo:  contestRepo,
		syncInterval: syncInterval,
	}
}

func (s *ContestSyncService) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()

	s.syncAllContests(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			s.syncAllContests(ctx)
		}
	}
}

func (s *ContestSyncService) syncAllContests(ctx context.Context) {
	contestsBySystem := s.registry.GetAllContests()
	for system, contestIDs := range contestsBySystem {
		adapter, ok := s.adapters[system]
		if !ok {
			logger.Errorf(ctx, "adapter for system %s not found", system)
			continue
		}

		for _, contestID := range contestIDs {
			contest, err := adapter.FetchContest(ctx, contestID)
			if err != nil {
				logger.Errorf(ctx, "failed to sync contest %d (%s): %v", contestID, system, err)
				continue
			}

			if err := s.contestRepo.Upsert(ctx, contest); err != nil {
				logger.Errorf(ctx, "failed to save contest %d: %v", contestID, err)
			} else {
				logger.Infof(ctx, "successfully synced contest %d (%s)", contestID, system)
			}
		}
	}
}
