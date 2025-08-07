package util

import (
	"time"
)

const (
	PATTERN = "13.08.2002 12:00"
)

func ParseTimeOrPanic(t string) time.Time {
	result, err := time.Parse(PATTERN, t)

	if err != nil {
		panic("Parsed invalid time format: " + t)
	}
	return result
}
