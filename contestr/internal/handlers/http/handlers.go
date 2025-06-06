package http

import (
	"contestr/internal/generated/server"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	helloHandle   *HelloHandle
	contestHandle *ContestHandle
}

func NewHandlers(
	helloHandle *HelloHandle,
	contestHandle *ContestHandle,
) *Handlers {
	return &Handlers{
		helloHandle:   helloHandle,
		contestHandle: contestHandle,
	}
}

func (h *Handlers) GetContest(ctx echo.Context, contestId int) error {
	return h.contestHandle.GetContest(ctx, contestId)
}

func (h *Handlers) GetHello(ctx echo.Context, params server.GetHelloParams) error {
	return h.helloHandle.GetHello(ctx, params)
}
