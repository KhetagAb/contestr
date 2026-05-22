package problemcode

import (
	"fmt"
	"regexp"
	"strconv"
)

var validCode = regexp.MustCompile(`^\d+[A-Z]$`)

// Validate checks regatta problem codes like 1A, 2B.
func Validate(code string) error {
	if !validCode.MatchString(code) {
		return fmt.Errorf("invalid problem_code %q", code)
	}
	return nil
}

// Round extracts tour round number from problem code (e.g. 2A -> 2).
func Round(code string) (int, error) {
	if err := Validate(code); err != nil {
		return 0, err
	}
	i := 0
	for i < len(code) && code[i] >= '0' && code[i] <= '9' {
		i++
	}
	n, err := strconv.Atoi(code[:i])
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("invalid problem_code %q", code)
	}
	return n, nil
}
