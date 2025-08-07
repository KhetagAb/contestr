package regatta

import (
	"contestr/pkg/logger"
	"contestr/pkg/regatta"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Regatta interface {
	GetContestResult(ctx context.Context, contestID int) (regatta.ContestStandings, error)
}

type ContestHandle struct {
	regatta Regatta
}

func NewContestHandle(
	regatta Regatta,
) *ContestHandle {
	return &ContestHandle{
		regatta: regatta,
	}
}

func (s *ContestHandle) GetContest(ectx echo.Context, contestId int) error {
	ctx := ectx.Request().Context()

	result, err := s.regatta.GetContestResult(ctx, contestId)
	if err != nil {
		logger.Errorf(ctx, "error while getting contest result: %v", err)
		return ectx.JSON(http.StatusInternalServerError, err.Error())
	}

	// TODO mapping

	return ectx.JSON(http.StatusOK, result)
}
