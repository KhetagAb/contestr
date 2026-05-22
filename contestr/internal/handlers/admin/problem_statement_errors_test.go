package admin

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"contestr/internal/services/problem_statement"
)

func TestProblemStatementHTTPError_putObject(t *testing.T) {
	err := fmt.Errorf("put object: %w", errors.New("AccessDenied: forbidden"))
	status, msg := problemStatementHTTPError(err)
	if status != http.StatusBadGateway {
		t.Fatalf("status=%d want %d", status, http.StatusBadGateway)
	}
	if msg == "" || msg == "Internal Server Error" {
		t.Fatalf("unexpected message: %q", msg)
	}
}

func TestProblemStatementHTTPError_notConfigured(t *testing.T) {
	status, msg := problemStatementHTTPError(problem_statement.ErrNotConfigured)
	if status != http.StatusServiceUnavailable {
		t.Fatalf("status=%d", status)
	}
	if msg != problem_statement.ErrNotConfigured.Error() {
		t.Fatalf("msg=%q", msg)
	}
}
