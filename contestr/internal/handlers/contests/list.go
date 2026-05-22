package contests

import (
	"net/http"

	"contestr/internal/generated/server"
	"contestr/internal/services/contest_admin"
	"contestr/pkg/regatta"

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
		settings := regatta.NormalizeScoringSettings(c.ScoringSettings)
		tourSettings := regatta.NormalizeTourSettings(c.TourSettings)
		items = append(items, server.RegisteredContestItem{
			ContestId: c.ContestID,
			Name:      c.Name,
			System:    c.System,
			ScoringSettings: server.ScoringSettings{
				SolveInTimeBonus: settings.SolveInTimeBonus,
				OvertakeBonus:    settings.OvertakeBonus,
			},
			TourSettings: server.TourSettings{
				GroupSize:           tourSettings.GroupSize,
				ProblemsPerTour:     tourSettings.ProblemsPerTour,
				GroupShufflePercent: tourSettings.GroupShufflePercent,
			},
		})
	}
	return ctx.JSON(http.StatusOK, items)
}
