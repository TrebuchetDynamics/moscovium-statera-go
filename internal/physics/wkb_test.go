package physics

import (
	"math"
	"testing"
)

func TestPredictWKBHalfLife_288Mc(t *testing.T) {
	logT, err := predictWKBHalfLife(115, 288, 10.75)
	if err != nil {
		t.Fatalf("WKB Predict: %v", err)
	}
	if logT < -3 || logT > 3 {
		t.Errorf("WKB log10(T) = %.2f, expected near -0.77 (0.17 s)", logT)
	}

	seconds := math.Pow(10, logT)
	if math.IsNaN(seconds) || math.IsInf(seconds, 0) {
		t.Errorf("WKB seconds not finite for logT=%.2f", logT)
	}
}

func TestPredictWKBHalfLife_InvalidInputs(t *testing.T) {
	_, err := predictWKBHalfLife(0, 288, 10.0)
	if err == nil {
		t.Error("expected error for Z=0")
	}
	_, err = predictWKBHalfLife(115, 288, -1.0)
	if err == nil {
		t.Error("expected error for negative Q")
	}
	_, err = predictWKBHalfLife(115, 100, 10.0)
	if err == nil {
		t.Error("expected error for A < Z (100 < 115)")
	}
}

func TestPredictWKBHalfLife_ParityIndependence(t *testing.T) {
	logTEvenEven, _ := predictWKBHalfLife(114, 286, 10.0)
	logTOddOdd, _ := predictWKBHalfLife(115, 288, 10.75)
	logTOddZ, _ := predictWKBHalfLife(115, 285, 10.5)
	logTOddN, _ := predictWKBHalfLife(114, 287, 10.0)

	for _, logT := range []float64{logTEvenEven, logTOddOdd, logTOddZ, logTOddN} {
		if math.IsNaN(logT) || math.IsInf(logT, 0) {
			t.Error("WKB produced NaN/Inf for parity class")
		}
	}
}
