package transport

import (
	"contestr/internal/handlers"
	"contestr/pkg/problemcode"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func registerProblemStatementRoutes(e *echo.Echo, h *handlers.Handlers) {
	e.PUT("/api/admin/contests/:contest_id/problem-statements/:problem_code", func(c echo.Context) error {
		contestID, err := strconv.Atoi(c.Param("contest_id"))
		if err != nil || contestID <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid contest_id")
		}
		code := c.Param("problem_code")
		if err := problemcode.Validate(code); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid problem_code")
		}
		return h.PutAdminContestProblemStatement(c, contestID, code)
	})
}
