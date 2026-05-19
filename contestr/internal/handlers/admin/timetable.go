package admin

import (
	"context"
	"errors"
	"net/http"

	"contestr/internal/configs"
	"contestr/internal/generated/server"
	regattasvc "contestr/internal/services/regatta"
	regattapkg "contestr/pkg/regatta"

	"github.com/labstack/echo/v4"
)

type TimetableService interface {
	GetTimetableView(ctx context.Context, contestID int, opts regattasvc.TimetableViewOptions) (*regattapkg.TimetableView, error)
	SaveTimetable(ctx context.Context, contestID int, input regattapkg.SaveTimetableRequest, opts regattasvc.TimetableViewOptions) (*regattapkg.TimetableView, error)
	RemoveTimetable(ctx context.Context, contestID int) error
	AdvanceTimetable(ctx context.Context, contestID int, mode regattasvc.AdvanceMode, opts regattasvc.TimetableViewOptions) error
	UpdateActiveTourDuration(
		ctx context.Context,
		contestID int,
		durationSeconds int,
		opts regattasvc.TimetableViewOptions,
	) (*regattapkg.TimetableView, error)
}

type TimetableHandle struct {
	timetable TimetableService
	cfg       *configs.Config
}

func NewTimetableHandle(timetable TimetableService, cfg *configs.Config) *TimetableHandle {
	return &TimetableHandle{
		timetable: timetable,
		cfg:       cfg,
	}
}

func (h *TimetableHandle) viewOpts() regattasvc.TimetableViewOptions {
	return regattasvc.TimetableViewOptions{
		ServerAutoStartAvailable: h.cfg.TimetableSync.Interval > 0,
	}
}

func (h *TimetableHandle) GetAdminTimetable(ctx echo.Context, contestId int) error {
	view, err := h.timetable.GetTimetableView(ctx.Request().Context(), contestId, h.viewOpts())
	if err != nil {
		return writeTimetableError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, view)
}

func (h *TimetableHandle) PutAdminTimetable(ctx echo.Context, contestId int) error {
	var input regattapkg.SaveTimetableRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, server.Error{Message: "invalid request body"})
	}

	view, err := h.timetable.SaveTimetable(ctx.Request().Context(), contestId, input, h.viewOpts())
	if err != nil {
		return writeTimetableError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, view)
}

func (h *TimetableHandle) DeleteAdminTimetable(ctx echo.Context, contestId int) error {
	if err := h.timetable.RemoveTimetable(ctx.Request().Context(), contestId); err != nil {
		return writeTimetableError(ctx, err)
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (h *TimetableHandle) PatchAdminTimetableActiveTourDuration(ctx echo.Context, contestId int) error {
	var input regattapkg.UpdateActiveTourDurationRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, server.Error{Message: "invalid request body"})
	}

	view, err := h.timetable.UpdateActiveTourDuration(
		ctx.Request().Context(),
		contestId,
		input.Duration,
		h.viewOpts(),
	)
	if err != nil {
		return writeTimetableError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, view)
}

func (h *TimetableHandle) PostAdminTimetableAdvance(ctx echo.Context, contestId int) error {
	if err := h.timetable.AdvanceTimetable(ctx.Request().Context(), contestId, regattasvc.AdvanceManual, h.viewOpts()); err != nil {
		return writeTimetableError(ctx, err)
	}
	view, err := h.timetable.GetTimetableView(ctx.Request().Context(), contestId, h.viewOpts())
	if err != nil {
		return writeTimetableError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, view)
}

func writeTimetableError(ctx echo.Context, err error) error {
	switch {
	case errors.Is(err, regattasvc.ErrInvalidTimetable),
		errors.Is(err, regattasvc.ErrManualStartWithAutostart),
		errors.Is(err, regattasvc.ErrContestNotStarted),
		errors.Is(err, regattasvc.ErrNothingToAdvance):
		return ctx.JSON(http.StatusBadRequest, server.Error{Message: err.Error()})
	case errors.Is(err, regattasvc.ErrTimetableAlreadyExists):
		return ctx.JSON(http.StatusConflict, server.Error{Message: err.Error()})
	case errors.Is(err, regattasvc.ErrTimetableNotFound),
		errors.Is(err, regattasvc.ErrContestNotFound),
		errors.Is(err, regattasvc.ErrTourNotFound),
		errors.Is(err, regattasvc.ErrNoActiveTour):
		return ctx.JSON(http.StatusNotFound, server.Error{Message: err.Error()})
	default:
		return ctx.JSON(http.StatusInternalServerError, server.Error{Message: err.Error()})
	}
}
