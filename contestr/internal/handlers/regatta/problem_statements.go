package regatta

import (
	"net/http"

	"contestr/internal/services/problem_statement"

	"github.com/labstack/echo/v4"
)

type ProblemStatementsHandle struct {
	svc *problem_statement.Service
}

func NewProblemStatementsHandle(svc *problem_statement.Service) *ProblemStatementsHandle {
	return &ProblemStatementsHandle{svc: svc}
}

type publicStatementsResponse struct {
	Statements map[string]string `json:"statements"`
}

func (h *ProblemStatementsHandle) GetRegattaProblemStatements(ctx echo.Context, contestID int) error {
	if h.svc == nil {
		return ctx.JSON(http.StatusOK, publicStatementsResponse{Statements: map[string]string{}})
	}
	statements, err := h.svc.ListPublicStatements(ctx.Request().Context(), contestID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	if statements == nil {
		statements = map[string]string{}
	}
	return ctx.JSON(http.StatusOK, publicStatementsResponse{Statements: statements})
}
