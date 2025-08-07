package util

import (
	"os"
	"time"
)

const (
	PATTERN = "13.08.2002 12:00"
)

func ParseTime(t string) time.Time {
	result, err := time.Parse(PATTERN, t)

	if err != nil {
		_, _ = os.Stderr.WriteString("Parsed invalid time format: " + t)
	}
	return result
}
