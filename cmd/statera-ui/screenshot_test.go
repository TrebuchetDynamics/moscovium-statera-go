//go:build !cgo

package main

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/ui"
	"github.com/gogpu/ui/theme/material3"
	"github.com/gogpu/ui/widget"
)

func TestDemoSummaryRowsGroupCardsForScreenshotReadability(t *testing.T) {
	rows := demoSummaryRows(ui.DefaultModel())
	if got, want := len(rows), 2; got != want {
		t.Fatalf("summary row count = %d, want %d", got, want)
	}
	for i, row := range rows {
		if got, want := len(row), 2; got != want {
			t.Fatalf("summary row %d card count = %d, want %d", i, got, want)
		}
		for _, card := range row {
			if card.Width < 300 {
				t.Fatalf("summary card %q width = %.0f, want at least 300 px for demo screenshot readability", card.Name, card.Width)
			}
			if card.Description == "" {
				t.Fatalf("summary card %q missing description", card.Name)
			}
		}
	}
}

func TestWorkbookCardSpecsSummarizeProvenanceForScreenshot(t *testing.T) {
	cards := workbookCardSpecs(ui.DefaultModel())
	if got, minWant := len(cards), 2; got < minWant {
		t.Fatalf("workbook card count = %d, want at least %d", got, minWant)
	}
	for _, card := range cards {
		if card.ID == "" {
			t.Fatal("workbook card missing ID")
		}
		if card.Identity == "" {
			t.Fatalf("workbook card %q missing identity summary", card.ID)
		}
		if !strings.Contains(card.Provenance, "citations=") {
			t.Fatalf("workbook card %q provenance = %q, want citation count", card.ID, card.Provenance)
		}
		if !strings.Contains(card.Provenance, "DOIs=") {
			t.Fatalf("workbook card %q provenance = %q, want DOI count", card.ID, card.Provenance)
		}
		if !strings.Contains(card.Provenance, "source=data/research.seed.json") {
			t.Fatalf("workbook card %q provenance = %q, want source path", card.ID, card.Provenance)
		}
		if card.EvidenceClass == "" {
			t.Fatalf("workbook card %q missing evidence class", card.ID)
		}
	}
}

func TestProvenanceNodeTableSpecsRenderAllAuditRows(t *testing.T) {
	rows := provenanceNodeTableSpecs(ui.DefaultModel())
	if got, minWant := len(rows), 17; got < minWant {
		t.Fatalf("provenance node table row count = %d, want at least %d", got, minWant)
	}
	var isotope288 provenanceNodeTableSpec
	for _, row := range rows {
		if row.NodeID == "isotope:288Mc" {
			isotope288 = row
		}
		if row.NodeID == "" || row.NodeType == "" || row.Status == "" || row.EdgeSummary == "" || row.OrphanStatus == "" {
			t.Fatalf("provenance row has blank audit field: %+v", row)
		}
	}
	if isotope288.NodeID == "" {
		t.Fatal("missing isotope:288Mc provenance row")
	}
	if isotope288.Status != "accepted" {
		t.Fatalf("isotope:288Mc status = %q, want accepted", isotope288.Status)
	}
	if isotope288.SourcePath != "data/research.seed.json" {
		t.Fatalf("isotope:288Mc source = %q, want data/research.seed.json", isotope288.SourcePath)
	}
	if !strings.Contains(isotope288.EdgeSummary, "incoming=0 outgoing=5") {
		t.Fatalf("isotope:288Mc edge summary = %q, want incoming=0 outgoing=5", isotope288.EdgeSummary)
	}
	if isotope288.OrphanStatus != "orphan=false" {
		t.Fatalf("isotope:288Mc orphan status = %q, want orphan=false", isotope288.OrphanStatus)
	}
}

func TestScreenshotViewsCoverPlannedSimulationAndVisualizationArtifacts(t *testing.T) {
	views := screenshotViews(ui.DefaultModel())
	want := []string{"overview", "workbook", "provenance", "roadmap", "simulation", "alpha"}
	if got := len(views); got != len(want) {
		t.Fatalf("screenshot view count = %d, want %d", got, len(want))
	}
	for _, name := range want {
		view, ok := views[name]
		if !ok {
			t.Fatalf("missing screenshot view %q", name)
		}
		if view.Title == "" {
			t.Fatalf("view %q missing title", name)
		}
		if len(view.Cards) < 2 {
			t.Fatalf("view %q card count = %d, want at least 2", name, len(view.Cards))
		}
	}
	if _, err := screenshotRootForView(ui.DefaultModel(), material3.New(widget.Hex(0x2F5D50)), "not-a-view"); err == nil {
		t.Fatal("screenshotRootForView accepted unknown view")
	}
}

func TestSaveScreenshotWritesNonBlankPNG(t *testing.T) {
	path := filepath.Join(t.TempDir(), "statera.png")
	model := ui.DefaultModel()
	theme := material3.New(widget.Hex(0x2F5D50))

	if err := saveScreenshot(path, model, theme); err != nil {
		t.Fatalf("saveScreenshot: %v", err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open screenshot: %v", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		t.Fatalf("decode screenshot: %v", err)
	}
	if got, want := img.Bounds().Dx(), model.Spec.Width; got != want {
		t.Fatalf("width = %d, want %d", got, want)
	}
	if got, want := img.Bounds().Dy(), model.Spec.Height; got != want {
		t.Fatalf("height = %d, want %d", got, want)
	}
	if imageIsSingleColor(img) {
		t.Fatal("screenshot appears blank")
	}
}

func TestDisplayAvailableDetectsWaylandOrX11(t *testing.T) {
	cases := []struct {
		name string
		env  map[string]string
		want bool
	}{
		{name: "none", env: map[string]string{}, want: false},
		{name: "wayland", env: map[string]string{"WAYLAND_DISPLAY": "wayland-0"}, want: true},
		{name: "x11", env: map[string]string{"DISPLAY": ":0"}, want: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			getenv := func(key string) string { return tc.env[key] }
			if got := displayAvailable(getenv); got != tc.want {
				t.Fatalf("displayAvailable() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHeadlessScreenshotPathUsesTempDirectory(t *testing.T) {
	tempDir := t.TempDir()
	path := headlessScreenshotPath(tempDir)
	if !strings.HasPrefix(path, tempDir+string(os.PathSeparator)) {
		t.Fatalf("path = %q, want under %q", path, tempDir)
	}
	if filepath.Base(path) != "statera-ui-headless.png" {
		t.Fatalf("base = %q, want statera-ui-headless.png", filepath.Base(path))
	}
}

func imageIsSingleColor(img image.Image) bool {
	bounds := img.Bounds()
	first := img.At(bounds.Min.X, bounds.Min.Y)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img.At(x, y) != first {
				return false
			}
		}
	}
	return true
}
