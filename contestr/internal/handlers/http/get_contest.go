package http

import (
	"contestr/internal/generated/server"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/togatoga/goforces"
	"net/http"
	"strconv"
)

type CodeforcesService interface {
	GetContest(ctx context.Context, contestID int) (*goforces.Standings, error)
}

type contestHandle struct {
	cfService CodeforcesService
}

func newContestHandle(
	cfService CodeforcesService,
) *contestHandle {
	return &contestHandle{
		cfService: cfService,
	}
}

func (s *contestHandle) GetContest(ectx echo.Context, contestId int) error {
	ctx := ectx.Request().Context()
	standings, err := s.cfService.GetContest(ctx, contestId)
	if err != nil {
		return err
	}

	response := server.GetContestResponse{
		Id:   strconv.FormatInt(standings.Contest.ID, 10),
		Name: standings.Contest.Name,
	}

	return ectx.JSON(http.StatusOK, response)
}
