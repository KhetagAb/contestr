package admin

import (
	"errors"
	"net/http"
	"strings"

	"contestr/internal/generated/server"
	"contestr/internal/repository"
	"contestr/internal/services/contest_admin"
	"contestr/pkg/regatta"

	"github.com/labstack/echo/v4"
)

type ContestsHandle struct {
	admin *contest_admin.Service
}

func NewContestsHandle(admin *contest_admin.Service) *ContestsHandle {
	return &ContestsHandle{admin: admin}
}

func (h *ContestsHandle) GetAdminContests(ctx echo.Context) error {
	return h.listContests(ctx)
}

func (h *ContestsHandle) PostAdminContest(ctx echo.Context) error {
	var req server.CreateContestRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, server.Error{Message: "invalid request body"})
	}

	name := ""
	if req.Name != nil {
		name = *req.Name
	}

	contest, err := h.admin.RegisterCodeforcesContest(
		ctx.Request().Context(),
		req.ContestId,
		name,
		fromAPIScoringSettings(req.ScoringSettings),
		fromAPITourSettings(req.TourSettings),
	)
	if err != nil {
		return writeContestError(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, toContestItem(*contest))
}

func (h *ContestsHandle) PatchAdminContestSettings(ctx echo.Context, contestId int) error {
	var req server.UpdateContestSettingsRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, server.Error{Message: "invalid request body"})
	}

	contest, err := h.admin.UpdateContestSettings(
		ctx.Request().Context(),
		contestId,
		fromAPIRequiredScoringSettings(req.ScoringSettings),
		fromAPIRequiredTourSettings(req.TourSettings),
	)
	if err != nil {
		return writeContestError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, toContestItem(*contest))
}

func (h *ContestsHandle) PostAdminContestRefresh(ctx echo.Context, contestId int) error {
	contest, err := h.admin.RefreshContest(ctx.Request().Context(), contestId)
	if err != nil {
		return writeContestError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, toContestItem(*contest))
}

func (h *ContestsHandle) DeleteAdminContest(ctx echo.Context, contestId int) error {
	if err := h.admin.DeleteContest(ctx.Request().Context(), contestId); err != nil {
		return writeContestError(ctx, err)
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (h *ContestsHandle) GetAdminContestHandles(ctx echo.Context, contestId int) error {
	mappings, err := h.admin.ListHandles(ctx.Request().Context(), contestId)
	if err != nil {
		return writeContestError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, toHandleItems(mappings))
}

func (h *ContestsHandle) PutAdminContestHandles(ctx echo.Context, contestId int) error {
	var req server.PutCodeforcesHandlesRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, server.Error{Message: "invalid request body"})
	}

	mappings := make([]repository.CodeforcesHandleMapping, 0, len(req.Handles))
	for _, item := range req.Handles {
		mappings = append(mappings, repository.CodeforcesHandleMapping{
			Handle: item.Handle,
			Name:   item.Name,
		})
	}

	if err := h.admin.UpsertHandles(ctx.Request().Context(), contestId, mappings); err != nil {
		return writeContestError(ctx, err)
	}

	updated, err := h.admin.ListHandles(ctx.Request().Context(), contestId)
	if err != nil {
		return writeContestError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, toHandleItems(updated))
}

func (h *ContestsHandle) DeleteAdminContestHandle(ctx echo.Context, contestId int, handle string) error {
	if err := h.admin.DeleteHandle(ctx.Request().Context(), contestId, handle); err != nil {
		return writeContestError(ctx, err)
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (h *ContestsHandle) listContests(ctx echo.Context) error {
	contests, err := h.admin.ListContests(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, server.Error{Message: "internal server error"})
	}

	items := make([]server.RegisteredContestItem, 0, len(contests))
	for _, c := range contests {
		items = append(items, toContestItem(c))
	}
	return ctx.JSON(http.StatusOK, items)
}

func toContestItem(c repository.RegisteredContest) server.RegisteredContestItem {
	settings := toAPIScoringSettings(regatta.NormalizeScoringSettings(c.ScoringSettings))
	tourSettings := toAPITourSettings(regatta.NormalizeTourSettings(c.TourSettings))
	return server.RegisteredContestItem{
		ContestId:       c.ContestID,
		Name:            c.Name,
		System:          c.System,
		ScoringSettings: settings,
		TourSettings:    tourSettings,
	}
}

func toHandleItems(mappings []repository.CodeforcesHandleMapping) []server.CodeforcesHandleItem {
	items := make([]server.CodeforcesHandleItem, 0, len(mappings))
	for _, m := range mappings {
		items = append(items, server.CodeforcesHandleItem{
			Handle: m.Handle,
			Name:   m.Name,
		})
	}
	return items
}

func fromAPIScoringSettings(settings *server.ScoringSettings) regatta.ScoringSettings {
	if settings == nil {
		return regatta.DefaultScoringSettings()
	}
	return fromAPIRequiredScoringSettings(*settings)
}

func fromAPIRequiredScoringSettings(settings server.ScoringSettings) regatta.ScoringSettings {
	return regatta.NormalizeScoringSettings(regatta.ScoringSettings{
		SolveInTimeBonus: settings.SolveInTimeBonus,
		OvertakeBonus:    settings.OvertakeBonus,
	})
}

func fromAPITourSettings(settings *server.TourSettings) regatta.TourSettings {
	if settings == nil {
		return regatta.DefaultTourSettings()
	}
	return fromAPIRequiredTourSettings(*settings)
}

func fromAPIRequiredTourSettings(settings server.TourSettings) regatta.TourSettings {
	return regatta.NormalizeTourSettings(regatta.TourSettings{
		GroupSize:           settings.GroupSize,
		ProblemsPerTour:     settings.ProblemsPerTour,
		GroupShufflePercent: settings.GroupShufflePercent,
	})
}

func toAPIScoringSettings(settings regatta.ScoringSettings) server.ScoringSettings {
	return server.ScoringSettings{
		SolveInTimeBonus: settings.SolveInTimeBonus,
		OvertakeBonus:    settings.OvertakeBonus,
	}
}

func toAPITourSettings(settings regatta.TourSettings) server.TourSettings {
	return server.TourSettings{
		GroupSize:           settings.GroupSize,
		ProblemsPerTour:     settings.ProblemsPerTour,
		GroupShufflePercent: settings.GroupShufflePercent,
	}
}

func writeContestError(ctx echo.Context, err error) error {
	if contest_admin.IsContestAlreadyRegistered(err) {
		return ctx.JSON(http.StatusConflict, server.Error{Message: "contest already registered"})
	}
	if contest_admin.IsContestNotRegistered(err) {
		return ctx.JSON(http.StatusNotFound, server.Error{Message: "contest not registered"})
	}
	if contest_admin.IsHandleNotFound(err) {
		return ctx.JSON(http.StatusNotFound, server.Error{Message: "handle mapping not found"})
	}

	msg := err.Error()
	if strings.Contains(msg, "codeforces contest not found") ||
		strings.Contains(msg, "failed to fetch codeforces contest") ||
		strings.Contains(msg, "codeforces contest.standings failed") ||
		strings.Contains(msg, "codeforces contest.status failed") {
		return ctx.JSON(http.StatusBadGateway, server.Error{Message: msg})
	}
	if errors.Is(err, repository.ErrContestAlreadyRegistered) {
		return ctx.JSON(http.StatusConflict, server.Error{Message: "contest already registered"})
	}

	return ctx.JSON(http.StatusBadRequest, server.Error{Message: msg})
}
