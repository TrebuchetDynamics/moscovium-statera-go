package physics

import (
	"math"
	"strings"
	"testing"
	"time"
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

func TestAlphaResidualSummaryAggregatesEvaluatedRecords(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, HalfLife: 170 * time.Millisecond, QAlphaMeV: 10.75, Daughter: "284Nh", CitationLink: "data/research.seed.json"},
		"290Mc": {Symbol: "Mc", Z: 115, A: 290, HalfLife: 650 * time.Millisecond, QAlphaMeV: 10.45, Daughter: "286Nh", CitationLink: "data/research.seed.json"},
	}

	summary, err := AlphaResidualSummary(catalog, []string{"288Mc", "290Mc"}, RoyerModel())
	if err != nil {
		t.Fatalf("AlphaResidualSummary: %v", err)
	}
	if summary.RecordCount != 2 {
		t.Fatalf("RecordCount = %d, want 2", summary.RecordCount)
	}
	if summary.SkippedCount != 0 {
		t.Fatalf("SkippedCount = %d, want 0", summary.SkippedCount)
	}
	if len(summary.Records) != 2 {
		t.Fatalf("len(Records) = %d, want 2", len(summary.Records))
	}
	for _, record := range summary.Records {
		if record.ModelName != RoyerModel().Name || record.EvidenceClass != EvidenceClassPeerReviewedModel {
			t.Fatalf("model metadata not preserved in residual record: %+v", record)
		}
		if record.LogResidual == 0 || math.IsNaN(record.LogResidual) || math.IsInf(record.LogResidual, 0) {
			t.Fatalf("LogResidual = %v, want finite non-zero residual", record.LogResidual)
		}
		if record.FactorError <= 0 || math.IsNaN(record.FactorError) || math.IsInf(record.FactorError, 0) {
			t.Fatalf("FactorError = %v, want finite positive factor", record.FactorError)
		}
	}
	if summary.MeanAbsoluteLogResidual <= 0 || summary.MaxAbsoluteLogResidual < summary.MeanAbsoluteLogResidual {
		t.Fatalf("summary residual magnitudes invalid: mean=%v max=%v", summary.MeanAbsoluteLogResidual, summary.MaxAbsoluteLogResidual)
	}
}

func TestAlphaResidualSummarySkipsZeroQRecords(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, HalfLife: 170 * time.Millisecond, QAlphaMeV: 10.75, Daughter: "284Nh", CitationLink: "data/research.seed.json"},
		"290Mc": {Symbol: "Mc", Z: 115, A: 290, HalfLife: 650 * time.Millisecond, QAlphaMeV: 0, Daughter: "286Nh", CitationLink: "data/research.seed.json"},
	}

	summary, err := AlphaResidualSummary(catalog, []string{"288Mc", "290Mc"}, RoyerModel())
	if err != nil {
		t.Fatalf("AlphaResidualSummary: %v", err)
	}
	if summary.RecordCount != 1 || summary.SkippedCount != 1 {
		t.Fatalf("counts = records %d skipped %d, want records 1 skipped 1", summary.RecordCount, summary.SkippedCount)
	}
	if !summary.Records[1].Skipped || !strings.Contains(summary.Records[1].SkipReason, "Q_alpha") {
		t.Fatalf("zero-Q record not skipped with Q_alpha reason: %+v", summary.Records[1])
	}
	if !math.IsNaN(summary.Records[1].LogResidual) || !math.IsNaN(summary.Records[1].FactorError) {
		t.Fatalf("skipped residuals = log %v factor %v, want NaN", summary.Records[1].LogResidual, summary.Records[1].FactorError)
	}
}

func TestAlphaResidualSummaryRejectsNaNResiduals(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, HalfLife: 170 * time.Millisecond, QAlphaMeV: math.NaN(), Daughter: "284Nh", CitationLink: "data/research.seed.json"},
	}

	if _, err := AlphaResidualSummary(catalog, []string{"288Mc"}, RoyerModel()); err == nil || !strings.Contains(err.Error(), "Q_alpha must be finite") {
		t.Fatalf("AlphaResidualSummary accepted NaN Q_alpha; err=%v", err)
	}
}
