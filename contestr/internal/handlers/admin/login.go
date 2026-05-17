package admin

import (
	"errors"
	"net/http"

	"contestr/internal/auth"
	"contestr/internal/generated/server"

	"github.com/labstack/echo/v4"
)

type LoginHandle struct {
	auth *auth.Service
}

func NewLoginHandle(auth *auth.Service) *LoginHandle {
	return &LoginHandle{auth: auth}
}

func (h *LoginHandle) PostAdminAuthLogin(ctx echo.Context) error {
	if !h.auth.Enabled() {
		return ctx.JSON(http.StatusServiceUnavailable, server.Error{
			Message: "admin authentication is disabled",
		})
	}

	var req server.AdminLoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, server.Error{
			Message: "invalid request body",
		})
	}

	if err := h.auth.ValidateCredentials(req.Username, req.Password); err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return ctx.JSON(http.StatusUnauthorized, server.Error{
				Message: "invalid username or password",
			})
		}
		return ctx.JSON(http.StatusServiceUnavailable, server.Error{
			Message: err.Error(),
		})
	}

	token, expiresIn, err := h.auth.IssueToken(req.Username)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, server.Error{
			Message: "failed to issue token",
		})
	}

	return ctx.JSON(http.StatusOK, server.AdminLoginResponse{
		AccessToken: token,
		ExpiresIn:   int(expiresIn),
	})
}
