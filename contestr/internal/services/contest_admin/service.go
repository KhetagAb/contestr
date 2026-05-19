package contest_admin

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"contestr/internal/integrations/codeforces"
	"contestr/internal/repository"
	contestsync "contestr/internal/services/contest_sync"
)

const systemCodeforces = "codeforces"

type Service struct {
	registeredRepo repository.RegisteredContestRepository
	handleRepo     repository.CodeforcesHandleRepository
	cfService      *codeforces.Service
	syncService    *contestsync.ContestSyncService
}

func NewService(
	registeredRepo repository.RegisteredContestRepository,
	handleRepo repository.CodeforcesHandleRepository,
	cfService *codeforces.Service,
	syncService *contestsync.ContestSyncService,
) *Service {
	return &Service{
		registeredRepo: registeredRepo,
		handleRepo:     handleRepo,
		cfService:      cfService,
		syncService:    syncService,
	}
}

func (s *Service) ListContests(ctx context.Context) ([]repository.RegisteredContest, error) {
	return s.registeredRepo.List(ctx)
}

func (s *Service) RegisterCodeforcesContest(ctx context.Context, contestID int, nameOverride string) (*repository.RegisteredContest, error) {
	if contestID <= 0 {
		return nil, fmt.Errorf("invalid contest id")
	}

	existing, err := s.registeredRepo.GetByContestID(ctx, contestID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, repository.ErrContestAlreadyRegistered
	}

	standings, err := s.cfService.GetContest(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("codeforces contest not found or unavailable: %w", err)
	}

	name := strings.TrimSpace(nameOverride)
	if name == "" {
		name = standings.Contest.Name
	}

	contest := repository.RegisteredContest{
		ContestID: contestID,
		System:    systemCodeforces,
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.registeredRepo.Create(ctx, contest); err != nil {
		return nil, err
	}

	if err := s.syncService.SyncContest(ctx, contestID); err != nil {
		return &contest, fmt.Errorf("contest registered but sync failed: %w", err)
	}

	return &contest, nil
}

func (s *Service) DeleteContest(ctx context.Context, contestID int) error {
	return s.registeredRepo.Delete(ctx, contestID)
}

func (s *Service) ListHandles(ctx context.Context, contestID int) ([]repository.CodeforcesHandleMapping, error) {
	if err := s.ensureContestExists(ctx, contestID); err != nil {
		return nil, err
	}
	return s.handleRepo.ListByContestID(ctx, contestID)
}

func (s *Service) UpsertHandles(ctx context.Context, contestID int, mappings []repository.CodeforcesHandleMapping) error {
	if err := s.ensureContestExists(ctx, contestID); err != nil {
		return err
	}
	return s.handleRepo.UpsertMany(ctx, contestID, mappings)
}

func (s *Service) DeleteHandle(ctx context.Context, contestID int, handle string) error {
	if err := s.ensureContestExists(ctx, contestID); err != nil {
		return err
	}
	return s.handleRepo.DeleteOne(ctx, contestID, handle)
}

func (s *Service) ensureContestExists(ctx context.Context, contestID int) error {
	contest, err := s.registeredRepo.GetByContestID(ctx, contestID)
	if err != nil {
		return err
	}
	if contest == nil {
		return repository.ErrContestNotRegistered
	}
	return nil
}

func IsContestAlreadyRegistered(err error) bool {
	return errors.Is(err, repository.ErrContestAlreadyRegistered)
}

func IsContestNotRegistered(err error) bool {
	return errors.Is(err, repository.ErrContestNotRegistered)
}

func IsHandleNotFound(err error) bool {
	return errors.Is(err, repository.ErrHandleNotFound)
}
