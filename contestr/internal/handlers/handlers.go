package handlers

import (
	"contestr/internal/handlers/codeforces"
	"contestr/internal/handlers/regatta"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	contestHandle        *codeforces.ContestHandle
	regattaContestHandle *regatta.ContestHandle
}

func (h *Handlers) GetRegattaContestStandings(ctx echo.Context, contestId int) error {
	return h.regattaContestHandle.GetContest(ctx, contestId)
}

func NewHandlers(
	contestHandle *codeforces.ContestHandle,
	regattaContestHandle *regatta.ContestHandle,
) *Handlers {
	return &Handlers{
		contestHandle:        contestHandle,
		regattaContestHandle: regattaContestHandle,
	}
}

func (h *Handlers) GetContest(ctx echo.Context, contestId int) error {
	return h.contestHandle.GetContest(ctx, contestId)
}
