package transport

import (
	"contestr/internal/handlers"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Дополнительные публичные маршруты regatta (пока не в сгенерированном api.gen.go).
func registerPublicRegattaRoutes(e *echo.Echo, h *handlers.Handlers) {
	e.GET("/api/regatta/contests/:contest_id/participants", func(c echo.Context) error {
		contestID, err := strconv.Atoi(c.Param("contest_id"))
		if err != nil || contestID <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid contest_id")
		}
		return h.GetRegattaContestParticipants(c, contestID)
	})
	e.GET("/api/regatta/contests/:contest_id/timetable", func(c echo.Context) error {
		contestID, err := strconv.Atoi(c.Param("contest_id"))
		if err != nil || contestID <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid contest_id")
		}
		return h.GetRegattaContestTimetable(c, contestID)
	})
}
