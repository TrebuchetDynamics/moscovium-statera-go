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
	if got, want := len(cards), 2; got != want {
		t.Fatalf("workbook card count = %d, want %d", got, want)
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
