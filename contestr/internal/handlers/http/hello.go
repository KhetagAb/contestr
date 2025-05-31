package http

import (
	"contestr/internal/generated/server"
	"contestr/pkg/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type helloHandle struct {
}

func newHelloHandle() *helloHandle {
	return &helloHandle{}
}

// GetHello реализует обработчик для ручки /hello
func (s *helloHandle) GetHello(ctx echo.Context, params server.GetHelloParams) error {
	name := "World"
	if params.Name != nil && *params.Name != "" {
		name = *params.Name
	}

	timestamp := time.Now()
	response := server.Hello{
		Message:   "Hello, " + name + "!",
		Timestamp: &timestamp,
	}

	logger.Info(ctx.Request().Context(), "hello request processed")
	return ctx.JSON(http.StatusOK, response)
}
