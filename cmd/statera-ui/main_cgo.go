//go:build cgo

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "statera-ui uses gogpu/ui through the zero-CGO WebGPU stack.")
	fmt.Fprintln(os.Stderr, "Run with: CGO_ENABLED=0 go run ./cmd/statera-ui")
	os.Exit(2)
}
