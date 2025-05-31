package http

import (
	"contestr/internal/generated/server"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	hello   *helloHandle
	contest *contestHandle
}

func (h *Handlers) GetContest(ctx echo.Context, contestId int) error {
	return h.contest.GetContest(ctx, contestId)
}

func (h *Handlers) GetHello(ctx echo.Context, params server.GetHelloParams) error {
	return h.hello.GetHello(ctx, params)
}

func NewHandlers(
	cfService CodeforcesService,
) *Handlers {
	// TODO
	return &Handlers{
		hello:   newHelloHandle(),
		contest: newContestHandle(cfService),
	}
}
