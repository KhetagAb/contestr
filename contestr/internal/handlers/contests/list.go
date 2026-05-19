package contests

import (
	"net/http"

	"contestr/internal/generated/server"
	"contestr/internal/services/contest_admin"

	"github.com/labstack/echo/v4"
)

type ListHandle struct {
	admin *contest_admin.Service
}

func NewListHandle(admin *contest_admin.Service) *ListHandle {
	return &ListHandle{admin: admin}
}

func (h *ListHandle) GetContests(ctx echo.Context) error {
	contests, err := h.admin.ListContests(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, server.Error{Message: "internal server error"})
	}

	items := make([]server.RegisteredContestItem, 0, len(contests))
	for _, c := range contests {
		items = append(items, server.RegisteredContestItem{
			ContestId: c.ContestID,
			Name:      c.Name,
			System:    c.System,
		})
	}
	return ctx.JSON(http.StatusOK, items)
}
