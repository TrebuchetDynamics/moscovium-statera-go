package ui

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/physics"
	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/research"
)

func TestDefaultSpecUsesResearchTitleAndModules(t *testing.T) {
	spec := DefaultSpec()

	if spec.Title != "Moscovium Statera Go" {
		t.Fatalf("Title = %q", spec.Title)
	}
	if spec.Width < 900 || spec.Height < 600 {
		t.Fatalf("unexpected window size: %dx%d", spec.Width, spec.Height)
	}

	wantModules := []string{"Education", "Research", "Design", "Context"}
	if len(spec.Modules) != len(wantModules) {
		t.Fatalf("module count = %d, want %d", len(spec.Modules), len(wantModules))
	}
	for i, want := range wantModules {
		if spec.Modules[i].Name != want {
			t.Fatalf("module[%d] name = %q, want %q", i, spec.Modules[i].Name, want)
		}
		if spec.Modules[i].Description == "" {
			t.Fatalf("module[%d] missing description", i)
		}
	}
}

func TestDefaultModelExposesMVPViews(t *testing.T) {
	model := DefaultModel()

	got := make([]string, 0, len(model.Views))
	for _, view := range model.Views {
		got = append(got, view.Name)
	}

	want := []string{"Education", "Research", "Design", "Context"}
	if len(got) != len(want) {
		t.Fatalf("view count = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("view[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestDefaultModelExposesEducationResearchDesignRecords(t *testing.T) {
	model := DefaultModel()

	if got, want := len(model.EducationLessons), 4; got < want {
		t.Fatalf("EducationLessons count = %d, want at least %d", got, want)
	}
	if got, want := len(model.ResearchItems), 3; got < want {
		t.Fatalf("ResearchItems count = %d, want at least %d", got, want)
	}
	if got, want := len(model.DesignScenarios), 4; got < want {
		t.Fatalf("DesignScenarios count = %d, want at least %d", got, want)
	}
}

func TestDefaultModelEducationLessonsHaveProvenance(t *testing.T) {
	model := DefaultModel()

	for _, lesson := range model.EducationLessons {
		if lesson.Title == "" {
			t.Fatal("education lesson missing title")
		}
		if lesson.Objective == "" {
			t.Fatalf("%s missing objective", lesson.Title)
		}
		if lesson.SourcePath == "" {
			t.Fatalf("%s missing source path", lesson.Title)
		}
	}
}

func TestDefaultModelResearchItemsHaveTrackAndSource(t *testing.T) {
	model := DefaultModel()

	for _, item := range model.ResearchItems {
		if item.Track != "Track A" && item.Track != "Track B" {
			t.Fatalf("%s Track = %q, want Track A or Track B", item.Key, item.Track)
		}
		if item.SourcePath == "" {
			t.Fatalf("%s missing source path", item.Key)
		}
	}
}

func TestDefaultModelDesignScenariosGateSimulationUse(t *testing.T) {
	model := DefaultModel()

	for _, scenario := range model.DesignScenarios {
		if scenario.SimulationUse != "allowed" && scenario.SimulationUse != "blocked" {
			t.Fatalf("%s SimulationUse = %q, want allowed or blocked", scenario.Name, scenario.SimulationUse)
		}
		if scenario.Result.Status != physics.ClaimStatusSupportedByTrackA && scenario.SimulationUse != "blocked" {
			t.Fatalf("%s status %q SimulationUse = %q, want blocked", scenario.Name, scenario.Result.Status, scenario.SimulationUse)
		}
		if scenario.SourcePath == "" {
			t.Fatalf("%s missing source path", scenario.Name)
		}
	}
}

func TestDefaultModelSummaryKeepsTrackBBoundary(t *testing.T) {
	model := DefaultModel()

	if model.Summary.VerifiedIsotopes == 0 {
		t.Fatal("VerifiedIsotopes = 0, want nonzero")
	}
	if model.Summary.CitationRecords == 0 {
		t.Fatal("CitationRecords = 0, want nonzero")
	}
	if model.Summary.ContextRecords == 0 {
		t.Fatal("ContextRecords = 0, want nonzero")
	}
	if !strings.Contains(model.Summary.BoundaryNotice, "Track B context is excluded from simulation") {
		t.Fatalf("BoundaryNotice = %q, want Track B exclusion", model.Summary.BoundaryNotice)
	}
	if strings.Join(model.Summary.DecayChain, " -> ") != "288Mc -> 284Nh -> 280Rg -> 276Mt -> 272Bh -> 268Db -> 264Lr" {
		t.Fatalf("DecayChain = %v", model.Summary.DecayChain)
	}
}

func TestDefaultModelClaimExamplesCoverValidatorStatuses(t *testing.T) {
	model := DefaultModel()

	seen := map[physics.ClaimStatus]bool{}
	for _, example := range model.ClaimExamples {
		seen[example.Result.Status] = true
	}

	for _, status := range []physics.ClaimStatus{
		physics.ClaimStatusSupportedByTrackA,
		physics.ClaimStatusStabilityIncongruent,
		physics.ClaimStatusOutsideSupportedModel,
		physics.ClaimStatusInvalidClaim,
	} {
		if !seen[status] {
			t.Fatalf("claim examples missing status %q", status)
		}
	}
}

func TestDefaultModelContextRecordsAreNotForSimulation(t *testing.T) {
	model := DefaultModel()

	if len(model.ContextRecords) == 0 {
		t.Fatal("ContextRecords is empty")
	}
	for _, record := range model.ContextRecords {
		if record.SimulationUse != "prohibited" {
			t.Fatalf("%s SimulationUse = %q, want prohibited", record.Title, record.SimulationUse)
		}
		if !contains(record.Labels, "not-for-simulation") {
			t.Fatalf("%s missing not-for-simulation label: %v", record.Title, record.Labels)
		}
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func TestWorkbookUIRecordsPreserveCountsAndEvidence(t *testing.T) {
	records := WorkbookUIRecordsFromResearch([]research.WorkbookRecord{
		{
			ID:            "288Mc",
			Element:       "Moscovium",
			Symbol:        "Mc",
			Z:             115,
			A:             288,
			N:             173,
			Daughter:      "284Nh",
			EvidenceLevel: "ENSDF evaluated nuclear data plus peer-reviewed DOI trail",
			SourcePath:    "data/research.seed.json",
			CitationURLs:  []string{"https://example.invalid/citation-a", "https://example.invalid/citation-b"},
			DOIs:          []string{"10.1103/PhysRevC.106.L031301", "10.1016/j.nuclphysa.2003.11.001"},
		},
	})

	if got, want := len(records), 1; got != want {
		t.Fatalf("workbook UI record count = %d, want %d", got, want)
	}
	record := records[0]
	if record.ID != "288Mc" || record.Element != "Moscovium" || record.Symbol != "Mc" {
		t.Fatalf("identity fields not preserved: %+v", record)
	}
	if record.Z != 115 || record.A != 288 || record.N != 173 {
		t.Fatalf("nuclide counts not preserved: Z=%d A=%d N=%d", record.Z, record.A, record.N)
	}
	if record.Daughter != "284Nh" {
		t.Fatalf("Daughter = %q, want 284Nh", record.Daughter)
	}
	if record.CitationCount != 2 {
		t.Fatalf("CitationCount = %d, want 2", record.CitationCount)
	}
	if record.DOICount != 2 {
		t.Fatalf("DOICount = %d, want 2", record.DOICount)
	}
	if record.SourcePath != "data/research.seed.json" {
		t.Fatalf("SourcePath = %q, want data/research.seed.json", record.SourcePath)
	}
	if record.EvidenceClass != "ENSDF evaluated nuclear data plus peer-reviewed DOI trail" {
		t.Fatalf("EvidenceClass = %q", record.EvidenceClass)
	}
}

func TestDefaultModelExposesWorkbookRecords(t *testing.T) {
	model := DefaultModel()
	if got, want := len(model.Workbook), 2; got != want {
		t.Fatalf("Workbook count = %d, want %d", got, want)
	}
	for _, record := range model.Workbook {
		if record.ID == "" {
			t.Fatal("workbook record missing ID")
		}
		if record.CitationCount == 0 {
			t.Fatalf("%s CitationCount = 0, want nonzero", record.ID)
		}
		if record.DOICount == 0 {
			t.Fatalf("%s DOICount = 0, want nonzero", record.ID)
		}
		if record.SourcePath != "data/research.seed.json" {
			t.Fatalf("%s SourcePath = %q, want data/research.seed.json", record.ID, record.SourcePath)
		}
		if record.EvidenceClass == "" {
			t.Fatalf("%s missing EvidenceClass", record.ID)
		}
	}
}

func TestDefaultModelExposesProvenanceGraphTextSummary(t *testing.T) {
	model := DefaultModel()
	for _, want := range []string{
		"provenance_graph",
		"nodes=17",
		"edges=14",
		"isotopes=2",
		"doi_nodes=6",
		"citation_url_nodes=4",
		"blocked_source_nodes=2",
		"orphans=0",
		"blocked_sources=royer2008alphaAnalytic,wang2015alphaSystematics",
	} {
		if !strings.Contains(model.ProvenanceSummary, want) {
			t.Fatalf("ProvenanceSummary = %q, want substring %q", model.ProvenanceSummary, want)
		}
	}
}

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

func TestDefaultAlphaSystematicsRecordsSurfacesSkipReasons(t *testing.T) {
	catalog := physics.Catalog{
		// 288Mc absent: missing-catalog skip path.
		// 290Mc present but Q_alpha = 0: missing-Q-alpha skip path.
		"290Mc": {Symbol: "Mc", Z: 115, A: 290, HalfLife: 650 * time.Millisecond, QAlphaMeV: 0, Daughter: "286Nh", CitationLink: "data/research.seed.json"},
	}

	records := defaultAlphaSystematicsRecords(catalog)
	if got, want := len(records), 2; got != want {
		t.Fatalf("records count = %d, want %d", got, want)
	}

	byID := map[string]AlphaSystematicsRecord{}
	for _, record := range records {
		byID[record.IsotopeID] = record
	}

	missing, ok := byID["288Mc"]
	if !ok {
		t.Fatal("missing-catalog skip record for 288Mc not produced")
	}
	if !missing.Skipped {
		t.Fatal("288Mc record Skipped = false, want true")
	}
	if missing.SkipReason == "" {
		t.Fatal("288Mc SkipReason empty, want a missing-catalog message")
	}
	if !math.IsNaN(missing.LogResidual) {
		t.Fatalf("288Mc LogResidual = %v, want NaN", missing.LogResidual)
	}
	if missing.ModelName == "" || missing.ModelReference == "" || missing.EvidenceClass == "" {
		t.Fatalf("288Mc model identity fields not populated on skip: %+v", missing)
	}
	if missing.SourcePath != "internal/physics/alpha.go" {
		t.Fatalf("288Mc SourcePath = %q, want internal/physics/alpha.go", missing.SourcePath)
	}

	zeroQ, ok := byID["290Mc"]
	if !ok {
		t.Fatal("zero-Q skip record for 290Mc not produced")
	}
	if !zeroQ.Skipped {
		t.Fatal("290Mc record Skipped = false, want true")
	}
	if zeroQ.SkipReason == "" {
		t.Fatal("290Mc SkipReason empty, want a missing-Q-alpha message")
	}
	if !math.IsNaN(zeroQ.LogResidual) {
		t.Fatalf("290Mc LogResidual = %v, want NaN", zeroQ.LogResidual)
	}
}
