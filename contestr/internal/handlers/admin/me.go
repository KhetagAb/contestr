package admin

import (
	"contestr/internal/auth"
	"contestr/internal/generated/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MeHandle struct{}

func NewMeHandle() *MeHandle {
	return &MeHandle{}
}

func (h *MeHandle) GetAdminMe(ctx echo.Context) error {
	username, ok := ctx.Get(auth.ContextUsernameKey).(string)
	if !ok || username == "" {
		return ctx.JSON(http.StatusUnauthorized, server.Error{
			Message: "unauthorized",
		})
	}

	return ctx.JSON(http.StatusOK, server.AdminMeResponse{
		Username: username,
		Ok:       true,
	})
}
