package admin

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"contestr/internal/services/problem_statement"
	"contestr/pkg/logger"

	"github.com/labstack/echo/v4"
)

const maxPDFUploadBytes = 10 << 20

type ProblemStatementsHandle struct {
	svc *problem_statement.Service
}

func NewProblemStatementsHandle(svc *problem_statement.Service) *ProblemStatementsHandle {
	return &ProblemStatementsHandle{svc: svc}
}

func (h *ProblemStatementsHandle) GetAdminProblemStatements(ctx echo.Context, contestID int) error {
	if h.svc == nil {
		return ctx.JSON(http.StatusServiceUnavailable, map[string]string{
			"message": "object storage is not configured",
		})
	}
	resp, err := h.svc.GetAdminList(ctx.Request().Context(), contestID)
	if err != nil {
		return writeProblemStatementError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (h *ProblemStatementsHandle) PutAdminProblemStatement(ctx echo.Context, contestID int, problemCode string) error {
	if h.svc == nil {
		return ctx.JSON(http.StatusServiceUnavailable, map[string]string{
			"message": "object storage is not configured",
		})
	}
	contentType := strings.ToLower(strings.TrimSpace(ctx.Request().Header.Get("Content-Type")))
	if !strings.HasPrefix(contentType, "application/pdf") {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "Content-Type must be application/pdf",
		})
	}
	body := ctx.Request().Body
	defer body.Close()
	limited := io.LimitReader(body, maxPDFUploadBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "failed to read body"})
	}
	if len(data) == 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "empty pdf"})
	}
	if int64(len(data)) > maxPDFUploadBytes {
		return ctx.JSON(http.StatusRequestEntityTooLarge, map[string]string{"message": "pdf exceeds 10MB limit"})
	}

	err = h.svc.SaveOrReplace(
		ctx.Request().Context(),
		contestID,
		problemCode,
		bytes.NewReader(data),
		int64(len(data)),
	)
	if err != nil {
		return writeProblemStatementError(ctx, err)
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (h *ProblemStatementsHandle) DeleteAdminProblemStatement(ctx echo.Context, contestID int, problemCode string) error {
	if h.svc == nil {
		return ctx.JSON(http.StatusServiceUnavailable, map[string]string{
			"message": "object storage is not configured",
		})
	}
	if err := h.svc.Delete(ctx.Request().Context(), contestID, problemCode); err != nil {
		return writeProblemStatementError(ctx, err)
	}
	return ctx.NoContent(http.StatusNoContent)
}

func writeProblemStatementError(ctx echo.Context, err error) error {
	logger.Errorf(ctx.Request().Context(), "problem statement: %v", err)
	status, message := problemStatementHTTPError(err)
	return ctx.JSON(status, map[string]string{"message": message})
}
