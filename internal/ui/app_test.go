package ui

import (
	"strings"
	"testing"

	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/physics"
)

func TestDefaultSpecUsesResearchTitleAndModules(t *testing.T) {
	spec := DefaultSpec()

	if spec.Title != "Moscovium Statera Go" {
		t.Fatalf("Title = %q", spec.Title)
	}
	if spec.Width < 900 || spec.Height < 600 {
		t.Fatalf("unexpected window size: %dx%d", spec.Width, spec.Height)
	}

	wantModules := []string{"Dashboard", "Physics", "Sources", "Context"}
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

	want := []string{"Dashboard", "Physics", "Sources", "Context"}
	if len(got) != len(want) {
		t.Fatalf("view count = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("view[%d] = %q, want %q", i, got[i], want[i])
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
