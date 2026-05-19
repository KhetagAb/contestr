package codeforces

import "testing"

func TestNormalizeCFHandle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"tourist", "tourist"},
		{"g44870=Aleksandrova_Sofyya", "Aleksandrova_Sofyya"},
		{"g44870=timofeev-dmitry", "timofeev-dmitry"},
	}

	for _, tt := range tests {
		if got := normalizeCFHandle(tt.in); got != tt.want {
			t.Errorf("normalizeCFHandle(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
