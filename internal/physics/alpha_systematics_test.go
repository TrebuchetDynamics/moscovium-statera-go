package physics

import (
	"math"
	"testing"
)

func TestPredictVSSHalfLife_288Mc(t *testing.T) {
	logT, err := predictVSSHalfLife(115, 288, 10.75)
	if err != nil {
		t.Fatalf("VSS: %v", err)
	}
	if math.IsNaN(logT) || math.IsInf(logT, 0) {
		t.Errorf("VSS produced non-finite logT for 288Mc")
	}
}

func TestPredictUNIVHalfLife_290Mc(t *testing.T) {
	logT, err := predictUNIVHalfLife(115, 290, 10.45)
	if err != nil {
		t.Fatalf("UNIV: %v", err)
	}
	if math.IsNaN(logT) || math.IsInf(logT, 0) {
		t.Errorf("UNIV produced non-finite logT for 290Mc")
	}
}

func TestPredictDenisovHalfLife_SuperheavyRegion(t *testing.T) {
	for _, id := range []string{"288Mc", "290Mc", "284Nh", "280Rg"} {
		var z, a int
		var q float64
		switch id {
		case "288Mc":
			z, a, q = 115, 288, 10.75
		case "290Mc":
			z, a, q = 115, 290, 10.45
		case "284Nh":
			z, a, q = 113, 284, 10.03
		case "280Rg":
			z, a, q = 111, 280, 10.00
		}
		logT, err := predictDenisovHalfLife(z, a, q)
		if err != nil {
			t.Errorf("%s Denisov: %v", id, err)
			continue
		}
		if math.IsNaN(logT) || math.IsInf(logT, 0) {
			t.Errorf("%s Denisov produced non-finite logT", id)
		}
		seconds := math.Pow(10, logT)
		if seconds <= 0 {
			t.Errorf("%s Denisov predicted negative/zero half-life", id)
		}
	}
}

func TestAlphaFormulas_ConsistentParity(t *testing.T) {
	for _, tc := range []struct {
		name string
		fn   func(z, a int, q float64) (float64, error)
	}{
		{"VSS", predictVSSHalfLife},
		{"UNIV", predictUNIVHalfLife},
		{"Denisov", predictDenisovHalfLife},
	} {
		logT, err := tc.fn(115, 288, 10.75)
		if err != nil {
			t.Errorf("%s: %v", tc.name, err)
			continue
		}
		if math.IsNaN(logT) || math.IsInf(logT, 0) {
			t.Errorf("%s produced non-finite value for 288Mc", tc.name)
		}
	}
}

func TestAlphaFormulas_InvalidInputs(t *testing.T) {
	fns := []struct {
		name string
		fn   func(z, a int, q float64) (float64, error)
	}{
		{"VSS", predictVSSHalfLife},
		{"UNIV", predictUNIVHalfLife},
		{"Denisov", predictDenisovHalfLife},
	}
	for _, f := range fns {
		_, err := f.fn(0, 288, 10.0)
		if err == nil {
			t.Errorf("%s: expected error for Z=0", f.name)
		}
		_, err = f.fn(115, 288, -1.0)
		if err == nil {
			t.Errorf("%s: expected error for negative Q", f.name)
		}
	}
}
