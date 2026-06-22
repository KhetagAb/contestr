package configs

import (
	"os"
	"testing"
)

func TestResolveLogLevel(t *testing.T) {
	t.Setenv("APP_LOG_LEVEL", "debug")
	if got := ResolveLogLevel(AppConfig{LogLevel: "warn"}); got != "debug" {
		t.Fatalf("env override: got %q, want debug", got)
	}

	t.Setenv("APP_LOG_LEVEL", "")
	if got := ResolveLogLevel(AppConfig{LogLevel: "warn"}); got != "warn" {
		t.Fatalf("config fallback: got %q, want warn", got)
	}

	if got := ResolveLogLevel(AppConfig{}); got != "info" {
		t.Fatalf("default: got %q, want info", got)
	}
}

func TestResolveLogLevelIgnoresWhitespaceEnv(t *testing.T) {
	t.Setenv("APP_LOG_LEVEL", "  ")
	if got := ResolveLogLevel(AppConfig{LogLevel: "error"}); got != "error" {
		t.Fatalf("whitespace env: got %q, want error", got)
	}
	os.Unsetenv("APP_LOG_LEVEL")
}
