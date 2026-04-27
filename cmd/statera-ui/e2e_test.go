//go:build !cgo

package main

import (
	"context"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestStateraUIBinaryWritesScreenshot(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess e2e test in short mode")
	}

	tempDir := t.TempDir()
	outPath := filepath.Join(tempDir, "e2e.png")

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", ".", "-screenshot", outPath)
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("go run -screenshot timed out after 90s\noutput:\n%s", output)
	}
	if err != nil {
		t.Fatalf("go run -screenshot: %v\noutput:\n%s", err, output)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("stat screenshot: %v", err)
	}
	if info.Size() < 10*1024 {
		t.Fatalf("screenshot size = %d bytes, want > 10KB (suggests blank/empty render)", info.Size())
	}

	file, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("open screenshot: %v", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		t.Fatalf("decode screenshot: %v", err)
	}
	if got, want := img.Bounds().Dx(), 1180; got != want {
		t.Fatalf("width = %d, want %d", got, want)
	}
	if got, want := img.Bounds().Dy(), 760; got != want {
		t.Fatalf("height = %d, want %d", got, want)
	}
	if imageIsSingleColor(img) {
		t.Fatal("binary-produced screenshot appears blank")
	}
}
