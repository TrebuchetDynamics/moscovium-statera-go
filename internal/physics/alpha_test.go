package physics

import (
	"math"
	"strings"
	"testing"
)

func TestRoyerModelMetadataIsAuditable(t *testing.T) {
	model := RoyerModel()

	if model.Name == "" {
		t.Fatal("Royer model missing Name")
	}
	if model.Reference == "" || !strings.Contains(model.Reference, "10.1103/PhysRevC.77.037602") {
		t.Fatalf("Royer model Reference = %q, want DOI 10.1103/PhysRevC.77.037602", model.Reference)
	}
	if model.EvidenceClass != EvidenceClassPeerReviewedModel {
		t.Fatalf("Royer model EvidenceClass = %q, want peer-reviewed-model", model.EvidenceClass)
	}
	for _, parity := range []ParityClass{ParityEvenEven, ParityOddZ, ParityOddN, ParityOddOdd} {
		if _, ok := model.Coefficients[parity]; !ok {
			t.Fatalf("Royer model missing coefficients for %q", parity)
		}
	}
}

func TestPredictReturnsLowerHalfLifeForHigherQAlpha(t *testing.T) {
	model := RoyerModel()

	low, err := model.Predict(115, 288, 10.0)
	if err != nil {
		t.Fatalf("Predict low Q: %v", err)
	}
	high, err := model.Predict(115, 288, 11.0)
	if err != nil {
		t.Fatalf("Predict high Q: %v", err)
	}
	if !(high.LogHalfLifeSeconds < low.LogHalfLifeSeconds) {
		t.Fatalf("higher Q_alpha should lower predicted log10(T_1/2): low=%v high=%v",
			low.LogHalfLifeSeconds, high.LogHalfLifeSeconds)
	}
}

func TestPredictRaisesHalfLifeForHigherZAtFixedQ(t *testing.T) {
	model := RoyerModel()

	lowZ, err := model.Predict(110, 286, 10.5)
	if err != nil {
		t.Fatalf("Predict lowZ: %v", err)
	}
	highZ, err := model.Predict(116, 286, 10.5)
	if err != nil {
		t.Fatalf("Predict highZ: %v", err)
	}
	if !(highZ.LogHalfLifeSeconds > lowZ.LogHalfLifeSeconds) {
		t.Fatalf("higher Z at fixed Q should raise predicted log10(T_1/2): lowZ=%v highZ=%v",
			lowZ.LogHalfLifeSeconds, highZ.LogHalfLifeSeconds)
	}
}

func TestPredictHindersOddANucleiRelativeToEvenEven(t *testing.T) {
	model := RoyerModel()

	even, err := model.Predict(114, 286, 10.5)
	if err != nil {
		t.Fatalf("Predict even-even: %v", err)
	}
	oddZ, err := model.Predict(115, 286, 10.5)
	if err != nil {
		t.Fatalf("Predict odd-Z: %v", err)
	}
	if !(oddZ.LogHalfLifeSeconds > even.LogHalfLifeSeconds) {
		t.Fatalf("odd-Z should hinder relative to even-even: even=%v oddZ=%v",
			even.LogHalfLifeSeconds, oddZ.LogHalfLifeSeconds)
	}
}

func TestPredictReturnsParityClassOnPrediction(t *testing.T) {
	model := RoyerModel()

	// 114 even, 116 even; 286 → N even, 287 → N odd, 288 → N even.
	cases := []struct {
		z, a int
		want ParityClass
	}{
		{114, 286, ParityEvenEven}, // Z=114 even, N=172 even
		{115, 287, ParityOddZ},     // Z=115 odd,  N=172 even
		{114, 287, ParityOddN},     // Z=114 even, N=173 odd
		{115, 288, ParityOddOdd},   // Z=115 odd,  N=173 odd
	}
	for _, c := range cases {
		got, err := model.Predict(c.z, c.a, 10.0)
		if err != nil {
			t.Fatalf("Predict(%d,%d): %v", c.z, c.a, err)
		}
		if got.ParityClass != c.want {
			t.Fatalf("Predict(%d,%d) parity = %q, want %q", c.z, c.a, got.ParityClass, c.want)
		}
	}
}

func TestPredictRejectsNonPositiveQAlpha(t *testing.T) {
	model := RoyerModel()
	if _, err := model.Predict(115, 288, 0); err == nil {
		t.Fatal("Predict accepted Q_alpha = 0")
	}
	if _, err := model.Predict(115, 288, -1); err == nil {
		t.Fatal("Predict accepted negative Q_alpha")
	}
}

func TestPredictRejectsInvalidZA(t *testing.T) {
	model := RoyerModel()
	if _, err := model.Predict(0, 288, 10.0); err == nil {
		t.Fatal("Predict accepted Z = 0")
	}
	if _, err := model.Predict(115, 0, 10.0); err == nil {
		t.Fatal("Predict accepted A = 0")
	}
	if _, err := model.Predict(115, 100, 10.0); err == nil {
		t.Fatal("Predict accepted A < Z")
	}
}

func TestPredictReturnsPositiveDuration(t *testing.T) {
	model := RoyerModel()
	got, err := model.Predict(115, 288, 10.75)
	if err != nil {
		t.Fatalf("Predict: %v", err)
	}
	if got.HalfLife <= 0 {
		t.Fatalf("HalfLife = %v, want positive", got.HalfLife)
	}
	if math.IsNaN(got.LogHalfLifeSeconds) || math.IsInf(got.LogHalfLifeSeconds, 0) {
		t.Fatalf("LogHalfLifeSeconds = %v, want finite", got.LogHalfLifeSeconds)
	}
}
