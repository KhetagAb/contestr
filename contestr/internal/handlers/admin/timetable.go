package admin

import (
	"context"
	"errors"
	"net/http"

	"contestr/internal/generated/server"
	regattasvc "contestr/internal/services/regatta"
	regattapkg "contestr/pkg/regatta"

	"github.com/labstack/echo/v4"
)

type TimetableService interface {
	CreateTimetable(ctx context.Context, timetable regattapkg.ToursTimetable) (*regattapkg.ToursTimetable, error)
	GetTimetable(ctx context.Context, contestID int) (*regattapkg.ToursTimetable, error)
	UpdateTimetable(ctx context.Context, timetable regattapkg.ToursTimetable) (*regattapkg.ToursTimetable, error)
	DeleteTimetable(ctx context.Context, contestID int) error
	MoveTimetableTour(ctx context.Context, contestID int, tourNumber int, newStartTime int) (*regattapkg.ToursTimetable, error)
	GetFirstNotStartedTimetableTour(ctx context.Context, contestID int) (*regattapkg.TourConfig, error)
}

type TimetableHandle struct {
	timetable TimetableService
}

type moveTimetableTourRequest struct {
	StartTime int `json:"start_time"`
}

func NewTimetableHandle(timetable TimetableService) *TimetableHandle {
	return &TimetableHandle{
		timetable: timetable,
	}
}

func (h *TimetableHandle) GetAdminTimetable(ctx echo.Context, contestId int) error {
	timetable, err := h.timetable.GetTimetable(ctx.Request().Context(), contestId)
	if err != nil {
		return writeTimetableError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, timetable)
}

func (h *TimetableHandle) PutAdminTimetable(ctx echo.Context, contestId int) error {
	var timetable regattapkg.ToursTimetable
	if err := ctx.Bind(&timetable); err != nil {
		return ctx.JSON(http.StatusBadRequest, server.Error{Message: "некорректное тело запроса"})
	}
	timetable.ContestId = contestId

	existing, err := h.timetable.GetTimetable(ctx.Request().Context(), contestId)
	if err != nil && !errors.Is(err, regattasvc.ErrTimetableNotFound) {
		return writeTimetableError(ctx, err)
	}

	var result *regattapkg.ToursTimetable
	if existing == nil {
		result, err = h.timetable.CreateTimetable(ctx.Request().Context(), timetable)
	} else {
		result, err = h.timetable.UpdateTimetable(ctx.Request().Context(), timetable)
	}
	if err != nil {
		return writeTimetableError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, result)
}

func (h *TimetableHandle) DeleteAdminTimetable(ctx echo.Context, contestId int) error {
	if err := h.timetable.DeleteTimetable(ctx.Request().Context(), contestId); err != nil {
		return writeTimetableError(ctx, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h *TimetableHandle) PatchAdminTimetableTourMove(ctx echo.Context, contestId int, tourNumber int) error {
	var req moveTimetableTourRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, server.Error{Message: "некорректное тело запроса"})
	}

	timetable, err := h.timetable.MoveTimetableTour(ctx.Request().Context(), contestId, tourNumber, req.StartTime)
	if err != nil {
		return writeTimetableError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, timetable)
}

func (h *TimetableHandle) GetAdminTimetableFirstNotStarted(ctx echo.Context, contestId int) error {
	tour, err := h.timetable.GetFirstNotStartedTimetableTour(ctx.Request().Context(), contestId)
	if err != nil {
		return writeTimetableError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, tour)
}

func writeTimetableError(ctx echo.Context, err error) error {
	switch {
	case errors.Is(err, regattasvc.ErrInvalidTimetable):
		return ctx.JSON(http.StatusBadRequest, server.Error{Message: err.Error()})
	case errors.Is(err, regattasvc.ErrTimetableAlreadyExists):
		return ctx.JSON(http.StatusConflict, server.Error{Message: err.Error()})
	case errors.Is(err, regattasvc.ErrTimetableNotFound), errors.Is(err, regattasvc.ErrTourNotFound):
		return ctx.JSON(http.StatusNotFound, server.Error{Message: err.Error()})
	default:
		return ctx.JSON(http.StatusInternalServerError, server.Error{Message: err.Error()})
	}
}
