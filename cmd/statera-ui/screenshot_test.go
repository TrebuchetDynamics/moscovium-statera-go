//go:build !cgo

package main

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/ui"
	"github.com/gogpu/ui/theme/material3"
	"github.com/gogpu/ui/widget"
)

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
