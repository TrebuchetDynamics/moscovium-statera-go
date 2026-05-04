package research

import (
	"math"
	"testing"
)

func TestParseENSDFNuclideHeader(t *testing.T) {
	tests := []struct {
		line    string
		wantID  string
		wantErr bool
	}{
		{"288MC    ADOPTED LEVELS, GAMMAS                 202206", "288Mc", false},
		{"290MC    ADOPTED LEVELS", "290Mc", false},
		{"284NH    ADOPTED LEVELS, GAMMAS", "284Nh", false},
		{"", "", true},
		{"XYZ", "", true}, // no mass number
	}
	for _, tc := range tests {
		id, err := ParseENSDFNuclideHeader(tc.line)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseENSDFNuclideHeader(%q): expected error", tc.line)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseENSDFNuclideHeader(%q): unexpected error: %v", tc.line, err)
			continue
		}
		if id != tc.wantID {
			t.Errorf("ParseENSDFNuclideHeader(%q): got %q, want %q", tc.line, id, tc.wantID)
		}
	}
}

func TestParseENSDFHalfLifeField(t *testing.T) {
	tests := []struct {
		line        string
		wantSeconds float64
		wantErr     bool
	}{
		{"  T  1.700E-01 S", 0.17, false},
		{"  T  650E-3 S", 0.65, false},
		{"  T  100 MS", 0.1, false},
		{"  T  50 US", 5e-5, false},
		{"  T  10 M", 600.0, false},
		{"  T  2 H", 7200.0, false},
		{"  T  365 D", 31536000.0, false},
		{"  T  1 Y", 31556952.0, false},
		{"invalid", 0, true},
		{"", 0, true},
	}
	for _, tc := range tests {
		seconds, err := ParseENSDFHalfLifeField(tc.line)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseENSDFHalfLifeField(%q): expected error", tc.line)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseENSDFHalfLifeField(%q): unexpected error: %v", tc.line, err)
			continue
		}
		if math.Abs(seconds-tc.wantSeconds) > 1e-6 {
			t.Errorf("ParseENSDFHalfLifeField(%q): got %f, want %f", tc.line, seconds, tc.wantSeconds)
		}
	}
}

func TestParseENSDFQValueField(t *testing.T) {
	tests := []struct {
		line    string
		wantMeV float64
		wantErr bool
	}{
		{" Q  1.075E+04 KEV", 10.75, false},
		{" Q  10.45 MEV", 10.45, false},
		{" Q  5.0E+06 EV", 5.0, false},
		{"invalid", 0, true},
		{"", 0, true},
	}
	for _, tc := range tests {
		mev, err := ParseENSDFQValueField(tc.line)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseENSDFQValueField(%q): expected error", tc.line)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseENSDFQValueField(%q): unexpected error: %v", tc.line, err)
			continue
		}
		if math.Abs(mev-tc.wantMeV) > 1e-3 {
			t.Errorf("ParseENSDFQValueField(%q): got %f MeV, want %f MeV", tc.line, mev, tc.wantMeV)
		}
	}
}

func TestParseENSDFDaughterField(t *testing.T) {
	tests := []struct {
		line       string
		wantDaughter string
		wantErr    bool
	}{
		{" DA  284NH", "284Nh", false},
		{" DA  280RG", "280Rg", false},
		{" DA  276MT", "276Mt", false},
		{"invalid", "", true},
		{"", "", true},
	}
	for _, tc := range tests {
		daughter, err := ParseENSDFDaughterField(tc.line)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseENSDFDaughterField(%q): expected error", tc.line)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseENSDFDaughterField(%q): unexpected error: %v", tc.line, err)
			continue
		}
		if daughter != tc.wantDaughter {
			t.Errorf("ParseENSDFDaughterField(%q): got %q, want %q", tc.line, daughter, tc.wantDaughter)
		}
	}
}
