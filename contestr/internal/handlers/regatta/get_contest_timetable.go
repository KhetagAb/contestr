package regatta

import (
	"context"
	"net/http"

	"contestr/internal/configs"
	regattasvc "contestr/internal/services/regatta"
	regattapkg "contestr/pkg/regatta"
	"contestr/pkg/logger"

	"github.com/labstack/echo/v4"
)

type TimetableViewer interface {
	GetTimetableView(ctx context.Context, contestID int, opts regattasvc.TimetableViewOptions) (*regattapkg.TimetableView, error)
}

type TimetableHandle struct {
	timetable TimetableViewer
	cfg       *configs.Config
}

func NewTimetableHandle(timetable TimetableViewer, cfg *configs.Config) *TimetableHandle {
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

func (h *TimetableHandle) GetContestTimetable(ctx echo.Context, contestId int) error {
	view, err := h.timetable.GetTimetableView(ctx.Request().Context(), contestId, h.viewOpts())
	if err != nil {
		logger.Errorf(ctx.Request().Context(), "error while getting contest timetable: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return ctx.JSON(http.StatusOK, view)
}
