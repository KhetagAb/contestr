package handlers

import (
	"contestr/internal/generated/server"
	"contestr/internal/handlers/codeforces"
	"contestr/internal/handlers/regatta"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	helloHandle          *HelloHandle
	contestHandle        *codeforces.ContestHandle
	regattaContestHandle *regatta.ContestHandle
}

func (h *Handlers) GetRegattaContestStandings(ctx echo.Context, contestId int) error {
	return h.regattaContestHandle.GetContest(ctx, contestId)
}

func NewHandlers(
	helloHandle *HelloHandle,
	contestHandle *codeforces.ContestHandle,
	regattaContestHandle *regatta.ContestHandle,
) *Handlers {
	return &Handlers{
		helloHandle:          helloHandle,
		contestHandle:        contestHandle,
		regattaContestHandle: regattaContestHandle,
	}
}

func (h *Handlers) GetContest(ctx echo.Context, contestId int) error {
	return h.contestHandle.GetContest(ctx, contestId)
}

func (h *Handlers) GetHello(ctx echo.Context, params server.GetHelloParams) error {
	return h.helloHandle.GetHello(ctx, params)
}
