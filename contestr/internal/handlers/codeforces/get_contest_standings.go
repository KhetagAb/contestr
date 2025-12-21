package codeforces

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/togatoga/goforces"
	"net/http"
)

type Service interface {
	GetContest(ctx context.Context, contestID int) (*goforces.Standings, error)
	GetContestStatus(ctx context.Context, contestID int, options *goforces.ContestStatusOptions) ([]goforces.Submission, error)
}

type ContestHandle struct {
	cfService Service
}

func NewContestHandle(
	cfService Service,
) *ContestHandle {
	return &ContestHandle{
		cfService: cfService,
	}
}

func (s *ContestHandle) GetContest(ectx echo.Context, contestId int) error {
	ctx := ectx.Request().Context()
	standings, err := s.cfService.GetContest(ctx, contestId)
	if err != nil {
		return err
	}

	return ectx.JSON(http.StatusOK, standings)
}
