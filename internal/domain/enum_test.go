package domain

import "testing"

func TestParseEnum(t *testing.T) {
	tests := []struct {
		in      string
		wantTot int
		wantOk  bool
	}{
		{"3", 3, true},
		{"3,5", 8, true},
		{"4-4", 8, true},
		{"3,4-5,2", 14, true},
		{" 3 , 5 ", 8, true},
		{"", 0, false},
		{"0", 0, false},
		{"3,,5", 0, false},
		{"a", 0, false},
	}

	for _, tt := range tests {
		e, err := ParseEnum(tt.in)
		if tt.wantOk && err != nil {
			t.Fatalf("ParseEnum(%q) unexpected error: %v", tt.in, err)
		}
		if !tt.wantOk && err == nil {
			t.Fatalf("ParseEnum(%q) expected error, got none: %+v", tt.in, e)
		}
		if tt.wantOk && e.Total != tt.wantTot {
			t.Fatalf("ParseEnum(%q) total=%d want=%d", tt.in, e.Total, tt.wantTot)
		}
	}
}

func TestNormalizedAnswerLen(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"CAT", 3},
		{"ICE-CREAM", 8},
		{"NEW YORK", 7},
		{"O'NEIL", 5},
		{"", 0},
	}

	for _, tt := range tests {
		if got := NormalizedAnswerLen(tt.in); got != tt.want {
			t.Fatalf("NormalizedAnswerLen(%q)=%d want=%d", tt.in, got, tt.want)
		}
	}
}
