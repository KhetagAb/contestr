package contest_sync

import (
	"contestr/internal/integrations"
	"contestr/internal/repository"
	"contestr/internal/services/contest_registry"
	"contestr/pkg/logger"
	"context"
	"fmt"
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

	s.SyncAllContests(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			s.SyncAllContests(ctx)
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

func (s *ContestSyncService) SyncAllContests(ctx context.Context) *SyncResult {
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

	for system, contestIDs := range contestsBySystem {
		if len(contestIDs) == 0 {
			continue
		}

		adapter, ok := s.adapters[system]
		if !ok {
			errMsg := fmt.Sprintf("adapter for system %s not found", system)
			logger.Errorf(ctx, errMsg)
			result.ErrorMessages = append(result.ErrorMessages, errMsg)
			result.FailedCount += len(contestIDs)
			continue
		}

		for _, contestID := range contestIDs {
			contest, err := adapter.FetchContest(ctx, contestID)
			if err != nil {
				errMsg := fmt.Sprintf("failed to sync contest %d (%s): %v", contestID, system, err)
				logger.Errorf(ctx, errMsg)
				result.ErrorMessages = append(result.ErrorMessages, errMsg)
				result.FailedCount++
				continue
			}

			if err := s.contestRepo.Upsert(ctx, contest); err != nil {
				errMsg := fmt.Sprintf("failed to save contest %d: %v", contestID, err)
				logger.Errorf(ctx, errMsg)
				result.ErrorMessages = append(result.ErrorMessages, errMsg)
				result.FailedCount++
			} else {
				logger.Infof(ctx, "successfully synced contest %d (%s)", contestID, system)
				result.SyncedCount++
			}
		}
	}

	return result
}

func (s *ContestSyncService) SyncContest(ctx context.Context, contestID int) error {
	system, err := s.registry.GetSystem(contestID)
	if err != nil {
		return err
	}

	adapter, ok := s.adapters[system]
	if !ok {
		return fmt.Errorf("adapter for system %s not found", system)
	}

	contest, err := adapter.FetchContest(ctx, contestID)
	if err != nil {
		return fmt.Errorf("failed to sync contest %d (%s): %w", contestID, system, err)
	}

	if err := s.contestRepo.Upsert(ctx, contest); err != nil {
		return fmt.Errorf("failed to save contest %d: %w", contestID, err)
	}

	logger.Infof(ctx, "successfully synced contest %d (%s)", contestID, system)
	return nil
}
