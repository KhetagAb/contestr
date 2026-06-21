package codeforces

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/togatoga/goforces"
)

type Service interface {
	GetContestStandings(ctx context.Context, contestID int) (*goforces.Standings, error)
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
	standings, err := s.cfService.GetContestStandings(ctx, contestId)
	if err != nil {
		return err
	}

	return ectx.JSON(http.StatusOK, standings)
}
