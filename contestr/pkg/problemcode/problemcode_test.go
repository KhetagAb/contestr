package problemcode

import "testing"

func TestValidateAndRound(t *testing.T) {
	for _, code := range []string{"1A", "12B"} {
		if err := Validate(code); err != nil {
			t.Fatalf("%s: %v", code, err)
		}
	}
	for _, bad := range []string{"", "1a", "A1", "../1A"} {
		if err := Validate(bad); err == nil {
			t.Fatalf("expected error for %q", bad)
		}
	}
	r, err := Round("2A")
	if err != nil || r != 2 {
		t.Fatalf("Round(2A) = %d, %v", r, err)
	}
}
