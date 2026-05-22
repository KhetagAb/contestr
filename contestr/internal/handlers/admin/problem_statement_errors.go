package admin

import (
	"errors"
	"net/http"
	"strings"

	"contestr/internal/services/problem_statement"

	"github.com/aws/smithy-go"
)

func problemStatementHTTPError(err error) (status int, message string) {
	switch {
	case errors.Is(err, problem_statement.ErrNotConfigured):
		return http.StatusServiceUnavailable, err.Error()
	case errors.Is(err, problem_statement.ErrPDFTooLarge):
		return http.StatusRequestEntityTooLarge, err.Error()
	}

	msg := strings.TrimSpace(err.Error())
	if msg == "" {
		return http.StatusInternalServerError, "unknown error"
	}
	if strings.Contains(msg, "invalid problem_code") {
		return http.StatusBadRequest, msg
	}

	lower := strings.ToLower(msg)
	switch {
	case strings.Contains(lower, "accessdenied"),
		strings.Contains(lower, "access denied"),
		strings.Contains(lower, "invalidaccesskeyid"),
		strings.Contains(lower, "signaturedoesnotmatch"):
		return http.StatusBadGateway,
			"Object Storage: доступ запрещён — проверьте APP_OBJECT_STORAGE_ACCESS_KEY_ID / SECRET_ACCESS_KEY и роль storage.editor на бакет"
	case strings.Contains(lower, "nosuchbucket"):
		return http.StatusBadGateway,
			"Object Storage: бакет не найден — проверьте object_storage.bucket в config.yaml"
	case strings.Contains(lower, "put object:"):
		return http.StatusBadGateway, humanizeWrapped("загрузка PDF в Object Storage", msg)
	case strings.Contains(lower, "save metadata:"):
		return http.StatusInternalServerError, humanizeWrapped("сохранение метаданных в MongoDB", msg)
	}

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) && apiErr.ErrorMessage() != "" {
		return http.StatusBadGateway, "Object Storage: " + apiErr.ErrorMessage()
	}

	return http.StatusInternalServerError, msg
}

func humanizeWrapped(action, full string) string {
	if idx := strings.Index(full, ": "); idx >= 0 && idx < len(full)-2 {
		return action + ": " + strings.TrimSpace(full[idx+2:])
	}
	return action + ": " + full
}
