# Alpha Systematics Lab Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add the first peer-reviewed-model calculation to Statera: an auditable Royer-style alpha-decay half-life predictor with explicit model metadata, plus a UI panel that compares evaluated half-lives to model predictions for the Track A seed isotopes.

**Architecture:** Create `internal/physics/alpha.go` with a `Model` struct holding named per-parity coefficients and a `Predict` function that returns log10 half-life and a derived `time.Duration`. The model carries a `peer-reviewed-model` evidence class so callers cannot accidentally treat its output as evaluated data. The UI gains an `AlphaSystematicsRecord` that pairs each evaluated isotope with its predicted half-life and the residual log10 difference; a new education lesson explains the evidence-class boundary. The renderer gets a new `Alpha Systematics` section.

**Tech Stack:** Go 1.25, pure Go, zero CGO, standard library `math` and `time`, existing `internal/physics` and `internal/ui` packages, `github.com/gogpu/ui` widgets, deterministic offscreen screenshots.

---

## File Structure

- Create `internal/physics/alpha.go`: `EvidenceClass`, `Parity`, `ParityClass`, `Coefficients`, `Model`, `Prediction`, `Predict`, and the `RoyerModel()` helper.
- Create `internal/physics/alpha_test.go`: monotonicity in Q_alpha, monotonicity in Z, parity classification, model-metadata invariants, and input validation.
- Modify `internal/ui/app.go`: add `AlphaSystematicsRecord`, `defaultAlphaSystematicsRecords`, append a lesson, wire records into `DefaultModel`.
- Modify `internal/ui/app_test.go`: assert at least two records, assert each carries evidence class `peer-reviewed-model`, assert lesson list grows.
- Modify `cmd/statera-ui/main.go`: add `alphaSystematicsSection`, slot it into `buildRoot`.
- Modify `screenshots/README.md`: document the new capture command and artifact name.
- Create `screenshots/statera-alpha-systematics-2026-04-27.png`: deterministic offscreen render.
- Modify `docs/academic-calculation-visualization-roadmap.md`: mark Alpha Systematics Lab as the active build slice and link this plan.

## Task 1: Physics package model and prediction tests

**Files:**
- Create: `internal/physics/alpha_test.go`

- [ ] **Step 1: Write the failing tests**

Write the full test file before any production code. The tests describe the surface and the algorithmic guarantees (monotonicity, parity hindrance, model metadata).

```go
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
```

- [ ] **Step 2: Run the tests to verify they fail**

Run:

```bash
CGO_ENABLED=0 go test ./internal/physics
```

Expected: FAIL with "undefined: RoyerModel" / "undefined: EvidenceClassPeerReviewedModel" etc.

- [ ] **Step 3: Commit the failing tests**

```bash
git add internal/physics/alpha_test.go
git commit -m "test: add alpha systematics model tests"
```

## Task 2: Physics package model and prediction implementation

**Files:**
- Create: `internal/physics/alpha.go`

- [ ] **Step 1: Write the implementation**

Implement the smallest code that makes Task 1's tests pass. The coefficients below are taken from the published Royer alpha-decay analytic-formula family (Royer 2000 Phys. G 26, 1149; refined in Royer & Zhang 2008, Phys. Rev. C 77, 037602). The values are the widely cited Royer parity-class fits. They are documented inline so an auditor can replace them with values transcribed directly from the Royer & Zhang 2008 PDF if any drift is found.

```go
package physics

import (
	"errors"
	"fmt"
	"math"
	"time"
)

type EvidenceClass string

const (
	EvidenceClassEvaluated         EvidenceClass = "evaluated"
	EvidenceClassPeerReviewedModel EvidenceClass = "peer-reviewed-model"
)

type ParityClass string

const (
	ParityEvenEven ParityClass = "even-even"
	ParityOddZ     ParityClass = "odd-Z"
	ParityOddN     ParityClass = "odd-N"
	ParityOddOdd   ParityClass = "odd-odd"
)

// Coefficients of the Royer log10(T_1/2 [s]) formula:
//
//	log10(T_1/2) = A + B * Aiso^(1/6) * sqrt(Z) + C * Z / sqrt(Q_alpha)
//
// where Aiso is the parent mass number, Z is the parent atomic number, and
// Q_alpha is in MeV. Coefficients vary by parent parity class.
type Coefficients struct {
	A float64
	B float64
	C float64
}

type Model struct {
	Name          string
	Reference     string
	EvidenceClass EvidenceClass
	Coefficients  map[ParityClass]Coefficients
	Notes         string
}

type Prediction struct {
	ParityClass        ParityClass
	LogHalfLifeSeconds float64
	HalfLife           time.Duration
	Coefficients       Coefficients
	EvidenceClass      EvidenceClass
}

// RoyerModel returns the Royer-family analytic alpha-decay half-life model.
// Coefficients are the parity-class fits from the Royer 2000 / Royer & Zhang 2008
// analytic-formula family (DOI 10.1103/PhysRevC.77.037602). Output is a
// peer-reviewed-model estimate, never evaluated data.
func RoyerModel() Model {
	return Model{
		Name:          "Royer 2008 analytic alpha-decay formula",
		Reference:     "DOI 10.1103/PhysRevC.77.037602",
		EvidenceClass: EvidenceClassPeerReviewedModel,
		Notes:         "log10(T_1/2 [s]) = A + B*Aiso^(1/6)*sqrt(Z) + C*Z/sqrt(Q_alpha). Parity-class coefficients from the Royer analytic-formula family; cross-check against the Royer & Zhang 2008 PDF before relying on absolute predictions.",
		Coefficients: map[ParityClass]Coefficients{
			ParityEvenEven: {A: -25.31, B: -1.1629, C: 1.5864},
			ParityOddZ:     {A: -26.65, B: -1.0859, C: 1.5848},
			ParityOddN:     {A: -25.68, B: -1.1423, C: 1.5920},
			ParityOddOdd:   {A: -29.48, B: -1.1130, C: 1.6971},
		},
	}
}

func (m Model) Predict(z, a int, qAlphaMeV float64) (Prediction, error) {
	if z <= 0 {
		return Prediction{}, fmt.Errorf("Predict: Z must be positive, got %d", z)
	}
	if a <= 0 {
		return Prediction{}, fmt.Errorf("Predict: A must be positive, got %d", a)
	}
	if a < z {
		return Prediction{}, fmt.Errorf("Predict: A=%d must be >= Z=%d", a, z)
	}
	if !(qAlphaMeV > 0) {
		return Prediction{}, fmt.Errorf("Predict: Q_alpha must be positive MeV, got %v", qAlphaMeV)
	}
	if math.IsNaN(qAlphaMeV) || math.IsInf(qAlphaMeV, 0) {
		return Prediction{}, errors.New("Predict: Q_alpha must be finite")
	}

	parity := classifyParity(z, a)
	coef, ok := m.Coefficients[parity]
	if !ok {
		return Prediction{}, fmt.Errorf("Predict: model has no coefficients for parity %q", parity)
	}

	logT := coef.A +
		coef.B*math.Pow(float64(a), 1.0/6.0)*math.Sqrt(float64(z)) +
		coef.C*float64(z)/math.Sqrt(qAlphaMeV)

	seconds := math.Pow(10, logT)
	if math.IsNaN(seconds) || math.IsInf(seconds, 0) {
		return Prediction{}, fmt.Errorf("Predict: predicted half-life is not finite (log10=%v)", logT)
	}
	half := durationFromSeconds(seconds)

	return Prediction{
		ParityClass:        parity,
		LogHalfLifeSeconds: logT,
		HalfLife:           half,
		Coefficients:       coef,
		EvidenceClass:      m.EvidenceClass,
	}, nil
}

func classifyParity(z, a int) ParityClass {
	zOdd := z%2 == 1
	nOdd := (a-z)%2 == 1
	switch {
	case !zOdd && !nOdd:
		return ParityEvenEven
	case zOdd && !nOdd:
		return ParityOddZ
	case !zOdd && nOdd:
		return ParityOddN
	default:
		return ParityOddOdd
	}
}

// durationFromSeconds clamps the predicted half-life into time.Duration without
// overflowing for very long-lived nuclei. Predictions above ~292 years saturate
// at math.MaxInt64 nanoseconds; the LogHalfLifeSeconds field stays exact.
func durationFromSeconds(seconds float64) time.Duration {
	const maxSeconds = float64(math.MaxInt64) / float64(time.Second)
	if seconds >= maxSeconds {
		return time.Duration(math.MaxInt64)
	}
	if seconds <= 0 {
		return 0
	}
	return time.Duration(seconds * float64(time.Second))
}
```

- [ ] **Step 2: Run the tests to verify they pass**

Run:

```bash
CGO_ENABLED=0 go test ./internal/physics
```

Expected: PASS for all alpha tests and existing decay/claims tests.

- [ ] **Step 3: Commit**

```bash
git add internal/physics/alpha.go
git commit -m "feat: add royer alpha systematics predictor"
```

## Task 3: UI model record tests

**Files:**
- Modify: `internal/ui/app_test.go`

- [ ] **Step 1: Add failing tests at the bottom of the file**

Append the following tests after `TestDefaultModelContextRecordsAreNotForSimulation`. These tests will fail because the model fields and records do not exist yet.

```go
func TestDefaultModelExposesAlphaSystematicsRecords(t *testing.T) {
	model := DefaultModel()

	if got, want := len(model.AlphaSystematics), 2; got < want {
		t.Fatalf("AlphaSystematics count = %d, want at least %d", got, want)
	}
	for _, record := range model.AlphaSystematics {
		if record.IsotopeID == "" {
			t.Fatal("alpha systematics record missing IsotopeID")
		}
		if record.ModelName == "" {
			t.Fatalf("%s missing ModelName", record.IsotopeID)
		}
		if record.EvidenceClass != string(physics.EvidenceClassPeerReviewedModel) {
			t.Fatalf("%s EvidenceClass = %q, want peer-reviewed-model",
				record.IsotopeID, record.EvidenceClass)
		}
		if record.SourcePath == "" {
			t.Fatalf("%s missing SourcePath", record.IsotopeID)
		}
		if record.EvaluatedHalfLife <= 0 {
			t.Fatalf("%s EvaluatedHalfLife = %v, want positive", record.IsotopeID, record.EvaluatedHalfLife)
		}
		if record.PredictedHalfLife <= 0 {
			t.Fatalf("%s PredictedHalfLife = %v, want positive", record.IsotopeID, record.PredictedHalfLife)
		}
	}
}

func TestDefaultModelAlphaSystematicsCoversSeedIsotopes(t *testing.T) {
	model := DefaultModel()

	seen := map[string]bool{}
	for _, record := range model.AlphaSystematics {
		seen[record.IsotopeID] = true
	}
	for _, want := range []string{"288Mc", "290Mc"} {
		if !seen[want] {
			t.Fatalf("AlphaSystematics missing %s", want)
		}
	}
}

func TestDefaultModelEducationLessonsIncludePeerReviewedBoundary(t *testing.T) {
	model := DefaultModel()

	for _, lesson := range model.EducationLessons {
		if strings.Contains(strings.ToLower(lesson.Title), "peer-reviewed-model") ||
			strings.Contains(strings.ToLower(lesson.Concept), "peer-reviewed-model") {
			return
		}
	}
	t.Fatal("education lessons do not mention the peer-reviewed-model evidence class")
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run:

```bash
CGO_ENABLED=0 go test ./internal/ui
```

Expected: FAIL with "model.AlphaSystematics undefined" and "education lessons do not mention the peer-reviewed-model evidence class".

- [ ] **Step 3: Commit the failing tests**

```bash
git add internal/ui/app_test.go
git commit -m "test: cover alpha systematics ui records"
```

## Task 4: UI model implementation

**Files:**
- Modify: `internal/ui/app.go`

- [ ] **Step 1: Add the `AlphaSystematicsRecord` type below the existing record types**

Insert after the `ContextRecord` struct (around line 107):

```go
type AlphaSystematicsRecord struct {
	IsotopeID          string
	Z                  int
	A                  int
	ParityClass        string
	QAlphaMeV          float64
	EvaluatedHalfLife  time.Duration
	PredictedHalfLife  time.Duration
	LogResidual        float64
	ModelName          string
	ModelReference     string
	EvidenceClass      string
	SourcePath         string
}
```

- [ ] **Step 2: Add `AlphaSystematics` field to `AppModel`**

Replace the `AppModel` struct definition with:

```go
type AppModel struct {
	Spec             Spec
	Views            []ViewSpec
	Summary          Summary
	Isotopes         []IsotopeRecord
	ClaimExamples    []ClaimExample
	SourceRecords    []SourceRecord
	EducationLessons []LessonRecord
	ResearchItems    []ResearchItem
	DesignScenarios  []DesignScenario
	ContextRecords   []ContextRecord
	AlphaSystematics []AlphaSystematicsRecord
}
```

- [ ] **Step 3: Add the `defaultAlphaSystematicsRecords` helper near the bottom of the file**

Append after `defaultContextRecords`:

```go
func defaultAlphaSystematicsRecords(catalog physics.Catalog) []AlphaSystematicsRecord {
	model := physics.RoyerModel()
	ids := []string{"288Mc", "290Mc"}
	records := make([]AlphaSystematicsRecord, 0, len(ids))
	for _, id := range ids {
		isotope, ok := catalog[id]
		if !ok || isotope.QAlphaMeV <= 0 {
			continue
		}
		prediction, err := model.Predict(isotope.Z, isotope.A, isotope.QAlphaMeV)
		if err != nil {
			continue
		}
		residual := math.NaN()
		if isotope.HalfLife > 0 {
			evaluatedSeconds := isotope.HalfLife.Seconds()
			if evaluatedSeconds > 0 {
				residual = prediction.LogHalfLifeSeconds - math.Log10(evaluatedSeconds)
			}
		}
		records = append(records, AlphaSystematicsRecord{
			IsotopeID:         id,
			Z:                 isotope.Z,
			A:                 isotope.A,
			ParityClass:       string(prediction.ParityClass),
			QAlphaMeV:         isotope.QAlphaMeV,
			EvaluatedHalfLife: isotope.HalfLife,
			PredictedHalfLife: prediction.HalfLife,
			LogResidual:       residual,
			ModelName:         model.Name,
			ModelReference:    model.Reference,
			EvidenceClass:     string(prediction.EvidenceClass),
			SourcePath:        "internal/physics/alpha.go",
		})
	}
	return records
}
```

- [ ] **Step 4: Add `math` to the import block**

Replace the existing import block at the top of the file:

```go
import (
	"math"
	"time"

	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/physics"
)
```

- [ ] **Step 5: Wire records into `DefaultModel`**

Inside `DefaultModel`, after `designScenarios := defaultDesignScenarios(catalog)`, add:

```go
	alphaSystematics := defaultAlphaSystematicsRecords(catalog)
```

Then add `AlphaSystematics: alphaSystematics,` as the last field of the returned `AppModel` literal.

- [ ] **Step 6: Add a peer-reviewed-model lesson**

Inside `defaultEducationLessons`, append one more lesson to the returned slice:

```go
		{
			Title:      "Peer-Reviewed-Model Boundary",
			Objective:  "Read alpha half-life predictions as model output, not as evaluated data.",
			Concept:    "Royer-style analytic formulas live in the peer-reviewed-model evidence class. Their output never substitutes for evaluated half-lives and must always carry the model name and DOI.",
			SourcePath: "internal/physics/alpha.go",
		},
```

- [ ] **Step 7: Run the UI tests to verify they pass**

Run:

```bash
CGO_ENABLED=0 go test ./internal/ui
```

Expected: PASS for all alpha-systematics tests and the existing UI tests.

- [ ] **Step 8: Run the full test suite**

Run:

```bash
CGO_ENABLED=0 go test ./...
```

Expected: PASS.

- [ ] **Step 9: Commit**

```bash
git add internal/ui/app.go
git commit -m "feat: expose alpha systematics records in ui model"
```

## Task 5: Renderer integration

**Files:**
- Modify: `cmd/statera-ui/main.go`

- [ ] **Step 1: Add the `alphaSystematicsSection` function**

Insert after `designSection` (around line 242):

```go
func alphaSystematicsSection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	children := []widget.Widget{
		primitives.Text("Royer formula vs evaluated half-lives").FontSize(14).Bold().Color(widget.Hex(0x24483E)),
		boundary("Predictions are peer-reviewed-model output. They never replace evaluated half-lives."),
	}
	for _, record := range model.AlphaSystematics {
		residual := "n/a"
		if !math.IsNaN(record.LogResidual) {
			residual = fmt.Sprintf("%+.2f", record.LogResidual)
		}
		children = append(children, card(
			primitives.Text(record.IsotopeID).FontSize(13).Bold(),
			primitives.Text(fmt.Sprintf("Z=%d  A=%d  parity=%s  Q_alpha=%.2f MeV", record.Z, record.A, record.ParityClass, record.QAlphaMeV)).FontSize(11),
			primitives.Text(fmt.Sprintf("evaluated T_1/2 %s | predicted T_1/2 %s | log10 residual %s", record.EvaluatedHalfLife, record.PredictedHalfLife, residual)).FontSize(11).Color(widget.Hex(0x44504B)),
			primitives.Text(fmt.Sprintf("%s | %s | %s", record.ModelName, record.ModelReference, record.EvidenceClass)).FontSize(10).Color(widget.Hex(0x31574D)),
			primitives.Text(record.SourcePath).FontSize(10).Color(widget.Hex(0x5F6F68)),
		))
	}
	return section("Alpha Systematics", children, theme)
}
```

- [ ] **Step 2: Add `math` to the import block**

Replace the existing import block at the top of the file:

```go
import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/ui"
	"github.com/gogpu/gg"
	_ "github.com/gogpu/gg/gpu"
	"github.com/gogpu/gg/integration/ggcanvas"
	"github.com/gogpu/gogpu"
	"github.com/gogpu/gpucontext"
	uiapp "github.com/gogpu/ui/app"
	"github.com/gogpu/ui/core/scrollview"
	"github.com/gogpu/ui/primitives"
	"github.com/gogpu/ui/render"
	uitheme "github.com/gogpu/ui/theme"
	"github.com/gogpu/ui/theme/material3"
	"github.com/gogpu/ui/widget"
)
```

- [ ] **Step 3: Slot the section into `buildRoot`**

In `buildRoot`, replace the `content := primitives.Box(...)` call with one that includes the new section:

```go
	content := primitives.Box(
		header(model),
		moduleOverview(model),
		educationSection(model, theme),
		researchSection(model, theme),
		designSection(model, theme),
		alphaSystematicsSection(model, theme),
		contextSection(model, theme),
	).Padding(20).Gap(14)
```

- [ ] **Step 4: Update the `moduleOverview` row to include alpha**

Replace `moduleOverview` with a four-card row so the metric strip mirrors the new section. The widths were 265 in a 3-card layout; reduce to 200 to keep the row inside the 1180-pixel window.

```go
func moduleOverview(model ui.AppModel) widget.Widget {
	return primitives.HBox(
		moduleSummary("Education", fmt.Sprintf("%d lessons", len(model.EducationLessons)), "Learn evaluated records and validator boundaries."),
		moduleSummary("Research", fmt.Sprintf("%d records", len(model.ResearchItems)), "Inspect DOI, URL, queue, and provenance status."),
		moduleSummary("Design", fmt.Sprintf("%d scenarios", len(model.DesignScenarios)), "Test constrained scenarios against Track A."),
		moduleSummary("Alpha", fmt.Sprintf("%d isotopes", len(model.AlphaSystematics)), "Compare Royer model predictions to evaluated half-lives."),
	).Gap(10)
}

func moduleSummary(name string, metric string, description string) widget.Widget {
	return primitives.Box(
		primitives.Text(name).FontSize(13).Bold().Color(widget.Hex(0x183D34)),
		primitives.Text(metric).FontSize(12).Color(widget.Hex(0x246B45)),
		primitives.Text(description).FontSize(10).Color(widget.Hex(0x52645C)),
	).Width(200).Padding(9).Gap(3).Background(widget.Hex(0xFFFFFF)).Rounded(6).BorderStyle(1, widget.Hex(0xD8E1DC))
}
```

- [ ] **Step 5: Run `go vet` and the screenshot test**

Run:

```bash
CGO_ENABLED=0 go vet ./...
CGO_ENABLED=0 go test ./cmd/statera-ui
```

Expected: PASS for `TestSaveScreenshotWritesNonBlankPNG` and the existing display tests.

- [ ] **Step 6: Commit**

```bash
git add cmd/statera-ui/main.go
git commit -m "feat: render alpha systematics section in statera-ui"
```

## Task 6: Deterministic offscreen screenshot

**Files:**
- Create: `screenshots/statera-alpha-systematics-2026-04-27.png`
- Modify: `screenshots/README.md`

- [ ] **Step 1: Capture the offscreen screenshot**

Run:

```bash
CGO_ENABLED=0 go run ./cmd/statera-ui -screenshot screenshots/statera-alpha-systematics-2026-04-27.png
```

Expected: command exits 0 and the PNG file exists.

- [ ] **Step 2: Verify the screenshot is non-blank**

Run:

```bash
file screenshots/statera-alpha-systematics-2026-04-27.png
```

Expected: `PNG image data, 1180 x 760, 8-bit/color RGBA, non-interlaced`.

- [ ] **Step 3: Document the new capture command in `screenshots/README.md`**

Append after the "Education/Research/Design capture command" block:

```
Alpha Systematics capture command:

```bash
CGO_ENABLED=0 go run ./cmd/statera-ui -screenshot screenshots/statera-alpha-systematics-2026-04-27.png
```
```

- [ ] **Step 4: Commit the screenshot and README update**

```bash
git add screenshots/statera-alpha-systematics-2026-04-27.png screenshots/README.md
git commit -m "docs: capture alpha systematics screenshot"
```

## Task 7: Roadmap status update

**Files:**
- Modify: `docs/academic-calculation-visualization-roadmap.md`

- [ ] **Step 1: Mark Alpha Systematics Lab as the active build slice**

Replace the "Recommended Next Build Slice" header line with:

```
## Active Build Slice: Alpha Systematics Lab
```

Then under that section, replace the `Minimum implementation:` block with:

```
Status: implemented in `internal/physics/alpha.go`, `internal/ui/app.go`, and `cmd/statera-ui/main.go` per `docs/superpowers/plans/2026-04-27-alpha-systematics-lab.md`.

Implementation receipts:

1. `internal/physics/alpha.go` exposes `RoyerModel()` and `Predict` for the Royer 2008 analytic alpha-decay formula family with parity-class coefficients and a `peer-reviewed-model` evidence-class label.
2. `internal/physics/alpha_test.go` covers Q_alpha monotonicity, Z monotonicity at fixed Q, parity hindrance, parity classification, and input validation.
3. `internal/ui/app.go` exposes `AlphaSystematicsRecord` and a peer-reviewed-model lesson; `cmd/statera-ui/main.go` renders an `Alpha Systematics` section.
4. The model output is labeled `peer-reviewed-model`, never `evaluated`, in both physics types and UI records.

Open audit follow-ups:

- Cross-check Royer 2008 PDF coefficients against the inline values in `internal/physics/alpha.go` and update if drift is found.
- Add a model-comparison overlay (Royer vs Wang or VSS) once a second analytic formula is intaken.
```

- [ ] **Step 2: Run the full test suite one more time**

Run:

```bash
CGO_ENABLED=0 go test ./...
```

Expected: PASS.

- [ ] **Step 3: Commit the roadmap update**

```bash
git add docs/academic-calculation-visualization-roadmap.md
git commit -m "docs: mark alpha systematics lab implemented"
```

---

## Self-Review

**Spec coverage** — Recommended Next Build Slice from `docs/academic-calculation-visualization-roadmap.md`:

1. "Add a `internal/physics/alpha.go` package with one published empirical formula and explicit model metadata." → Tasks 1-2.
2. "Add tests that compare monotonic behavior: increasing Q_alpha lowers predicted alpha half-life for fixed Z/A." → Task 1, `TestPredictReturnsLowerHalfLifeForHigherQAlpha`. Bonus monotonicity in Z and parity hindrance also covered.
3. "Add a UI section that compares evaluated half-life to calculated model half-life for seed isotopes." → Tasks 3-5.
4. "Label the result as `peer-reviewed-model`, not `evaluated`." → Task 1 metadata test, Task 3 record test, Task 5 renderer text.

**Placeholder scan** — No "TBD" / "implement later" / "appropriate error handling" / unreferenced types. Every code step shows the actual code. Coefficients are concrete numeric literals; the inline comment names the audit gap explicitly so it is not a hidden TODO.

**Type consistency** — `Model`, `Prediction`, `Coefficients`, `ParityClass`, `EvidenceClass`, `RoyerModel`, `Predict` referenced consistently across Tasks 1, 2, 3, 4. `AlphaSystematicsRecord` field names (`IsotopeID`, `Z`, `A`, `ParityClass`, `QAlphaMeV`, `EvaluatedHalfLife`, `PredictedHalfLife`, `LogResidual`, `ModelName`, `ModelReference`, `EvidenceClass`, `SourcePath`) are identical between Tasks 3, 4, and 5.

---

## Execution Handoff

Plan complete and saved to `docs/superpowers/plans/2026-04-27-alpha-systematics-lab.md`. Two execution options:

1. **Subagent-Driven (recommended)** - I dispatch a fresh subagent per task, review between tasks, fast iteration.
2. **Inline Execution** - Execute tasks in this session using executing-plans, batch execution with checkpoints.

Which approach?
