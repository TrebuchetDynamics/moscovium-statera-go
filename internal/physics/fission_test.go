package physics

import (
	"testing"
)

func TestSwiateckiSFHalfLife_290Mc(t *testing.T) {
	logT, err := SwiateckiSFHalfLife(115, 290)
	if err != nil {
		t.Fatalf("SF half-life: %v", err)
	}
	if logT < -20 || logT > 30 {
		t.Errorf("SF log10(T) = %.2f out of plausible range", logT)
	}
}

func TestSwiateckiSFHalfLife_SuperheavyRegion(t *testing.T) {
	tests := []struct {
		id  string
		z   int
		a   int
		min float64
		max float64
	}{
		{"288Mc", 115, 288, -20, 30},
		{"290Mc", 115, 290, -20, 30},
		{"284Nh", 113, 284, -20, 30},
		{"280Rg", 111, 280, -20, 30},
	}
	for _, tc := range tests {
		logT, err := SwiateckiSFHalfLife(tc.z, tc.a)
		if err != nil {
			t.Errorf("%s: %v", tc.id, err)
			continue
		}
		if logT < tc.min || logT > tc.max {
			t.Errorf("%s: logT=%.2f out of [%.0f, %.0f]", tc.id, logT, tc.min, tc.max)
		}
	}
}

func TestSwiateckiSFHalfLife_InvalidInputs(t *testing.T) {
	_, err := SwiateckiSFHalfLife(0, 290)
	if err == nil {
		t.Error("expected error for Z=0")
	}
	_, err = SwiateckiSFHalfLife(115, 100)
	if err == nil {
		t.Error("expected error for A < Z")
	}
}

func TestCompetitiveDecayMode_AlphaDominant(t *testing.T) {
	mode := CompetitiveDecayMode(115, 288, 10.75, -0.77, 15.0)
	if mode != "alpha" {
		t.Errorf("288Mc should be alpha-dominant, got %s", mode)
	}
}

func TestCompetitiveDecayMode_SFDominant(t *testing.T) {
	mode := CompetitiveDecayMode(111, 281, 0, 10.0, 1.0)
	if mode != "SF" {
		t.Errorf("should be SF when Q_alpha <= 0, got %s", mode)
	}
}
