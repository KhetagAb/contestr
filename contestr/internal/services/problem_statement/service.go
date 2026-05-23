package problem_statement

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"contestr/internal/repository"
	"contestr/internal/storage/objectstorage"
	"contestr/pkg/problemcode"
	"contestr/pkg/regatta"

	"github.com/google/uuid"
)

const maxPDFSize = 10 << 20 // 10 MiB

var (
	ErrNotConfigured = errors.New("object storage is not configured")
	ErrPDFTooLarge   = errors.New("pdf exceeds size limit")
)

type TourLister interface {
	FindByContestID(ctx context.Context, contestID int) ([]regatta.Tour, error)
}

type BlobStore interface {
	PutObject(ctx context.Context, objectKey string, body io.Reader, size int64) error
	DeleteObject(ctx context.Context, objectKey string) error
	PublicURL(objectKey string) string
}

type Service struct {
	storage BlobStore
	repo    repository.ProblemStatementRepository
	tours   TourLister
}

func NewService(
	storage BlobStore,
	repo repository.ProblemStatementRepository,
	tours TourLister,
) *Service {
	return &Service{storage: storage, repo: repo, tours: tours}
}

func (s *Service) configured() bool {
	return s != nil && s.storage != nil
}

type AdminItem struct {
	ProblemCode string    `json:"problem_code"`
	Status      string    `json:"status"`
	PublicURL   string    `json:"public_url,omitempty"`
	UploadedAt  string    `json:"uploaded_at,omitempty"`
	SizeBytes   int64     `json:"size_bytes,omitempty"`
}

type AdminListResponse struct {
	Items []AdminItem `json:"items"`
}

func (s *Service) SaveOrReplace(ctx context.Context, contestID int, problemCode string, r io.Reader, size int64) error {
	if !s.configured() {
		return ErrNotConfigured
	}
	if err := problemcode.Validate(problemCode); err != nil {
		return err
	}
	if size <= 0 || size > maxPDFSize {
		return ErrPDFTooLarge
	}

	existing, err := s.repo.Get(ctx, contestID, problemCode)
	if err != nil {
		return err
	}

	objectKey := objectstorage.ObjectKey(contestID, uuid.NewString())
	if err := s.storage.PutObject(ctx, objectKey, r, size); err != nil {
		return fmt.Errorf("put object: %w", err)
	}

	doc := repository.ProblemStatement{
		ContestID:   contestID,
		ProblemCode: problemCode,
		ObjectKey:   objectKey,
		UploadedAt:  nowUTC(),
		SizeBytes:   size,
	}
	if err := s.repo.Upsert(ctx, doc); err != nil {
		_ = s.storage.DeleteObject(ctx, objectKey)
		return fmt.Errorf("save metadata: %w", err)
	}

	if existing != nil && existing.ObjectKey != "" && existing.ObjectKey != objectKey {
		_ = s.storage.DeleteObject(ctx, existing.ObjectKey)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, contestID int, problemCode string) error {
	if !s.configured() {
		return ErrNotConfigured
	}
	if err := problemcode.Validate(problemCode); err != nil {
		return err
	}

	existing, err := s.repo.Get(ctx, contestID, problemCode)
	if err != nil {
		return err
	}
	if existing == nil {
		return nil
	}

	if existing.ObjectKey != "" {
		_ = s.storage.DeleteObject(ctx, existing.ObjectKey)
	}
	return s.repo.Delete(ctx, contestID, problemCode)
}

func (s *Service) DeleteByContest(ctx context.Context, contestID int) error {
	if !s.configured() {
		return nil
	}
	docs, err := s.repo.ListByContest(ctx, contestID)
	if err != nil {
		return err
	}
	for _, doc := range docs {
		if doc.ObjectKey != "" {
			_ = s.storage.DeleteObject(ctx, doc.ObjectKey)
		}
	}
	return s.repo.DeleteByContest(ctx, contestID)
}

func (s *Service) GetAdminList(ctx context.Context, contestID int) (AdminListResponse, error) {
	if !s.configured() {
		return AdminListResponse{}, ErrNotConfigured
	}
	docs, err := s.repo.ListByContest(ctx, contestID)
	if err != nil {
		return AdminListResponse{}, err
	}

	items := make([]AdminItem, 0, len(docs))
	for _, doc := range docs {
		items = append(items, AdminItem{
			ProblemCode: doc.ProblemCode,
			Status:      "uploaded",
			PublicURL:   s.storage.PublicURL(doc.ObjectKey),
			UploadedAt:  doc.UploadedAt.UTC().Format(timeRFC3339),
			SizeBytes:   doc.SizeBytes,
		})
	}
	return AdminListResponse{Items: items}, nil
}

func (s *Service) ListPublicStatements(ctx context.Context, contestID int) (map[string]string, error) {
	if !s.configured() {
		return map[string]string{}, nil
	}

	maxRound, err := s.maxStartedRound(ctx, contestID)
	if err != nil {
		return nil, err
	}
	if maxRound == 0 {
		return map[string]string{}, nil
	}

	docs, err := s.repo.ListByContest(ctx, contestID)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string)
	for _, doc := range docs {
		if strings.TrimSpace(doc.ObjectKey) == "" {
			continue
		}
		round, err := problemcode.Round(doc.ProblemCode)
		if err != nil || round > maxRound {
			continue
		}
		url := s.storage.PublicURL(doc.ObjectKey)
		if !isOurPublicURL(url) {
			continue
		}
		out[doc.ProblemCode] = url
	}
	return out, nil
}

func (s *Service) maxStartedRound(ctx context.Context, contestID int) (int, error) {
	tours, err := s.tours.FindByContestID(ctx, contestID)
	if err != nil {
		return 0, err
	}
	maxRound := 0
	for _, tour := range tours {
		if tour.IsPause {
			continue
		}
		if tour.Round > maxRound {
			maxRound = tour.Round
		}
	}
	return maxRound, nil
}

func isOurPublicURL(url string) bool {
	return strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://")
}

const timeRFC3339 = "2006-01-02T15:04:05Z07:00"

func nowUTC() time.Time {
	return time.Now().UTC()
}
