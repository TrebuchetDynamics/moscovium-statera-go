//go:build !cgo

package main

import (
	"log"

	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/ui"
	"github.com/gogpu/gg"
	_ "github.com/gogpu/gg/gpu"
	"github.com/gogpu/gg/integration/ggcanvas"
	"github.com/gogpu/gogpu"
	uiapp "github.com/gogpu/ui/app"
	"github.com/gogpu/ui/primitives"
	"github.com/gogpu/ui/render"
	uitheme "github.com/gogpu/ui/theme"
	"github.com/gogpu/ui/theme/material3"
	"github.com/gogpu/ui/widget"
)

func main() {
	spec := ui.DefaultSpec()

	gpuApp := gogpu.NewApp(gogpu.DefaultConfig().
		WithTitle(spec.Title).
		WithSize(spec.Width, spec.Height).
		WithContinuousRender(false))

	seed := widget.Hex(0x2F5D50)
	appTheme := uitheme.DefaultLight()
	appTheme.Colors.Primary = seed
	appTheme.Colors.PrimaryDark = widget.Hex(0x1F463C)
	appTheme.Colors.PrimaryLight = widget.Hex(0x6F9C8D)
	materialTheme := material3.New(seed)
	app := uiapp.New(
		uiapp.WithWindowProvider(gpuApp),
		uiapp.WithPlatformProvider(gpuApp),
		uiapp.WithEventSource(gpuApp.EventSource()),
		uiapp.WithTheme(appTheme),
	)
	app.SetRoot(buildRoot(spec, materialTheme))

	var canvas *ggcanvas.Canvas
	gpuApp.OnDraw(func(dc *gogpu.Context) {
		w, h := dc.Width(), dc.Height()
		if w <= 0 || h <= 0 {
			return
		}
		if canvas == nil {
			provider := gpuApp.GPUContextProvider()
			if provider == nil {
				return
			}
			var err error
			canvas, err = ggcanvas.New(provider, w, h)
			if err != nil {
				log.Printf("ggcanvas: %v", err)
				return
			}
		}

		app.Frame()
		cw, ch := canvas.Size()
		if cw != w || ch != h {
			if err := canvas.Resize(w, h); err != nil {
				log.Printf("resize: %v", err)
				return
			}
			cw, ch = w, h
		}

		if err := canvas.Draw(func(cc *gg.Context) {
			cc.SetRGBA(0.96, 0.97, 0.96, 1)
			cc.DrawRectangle(0, 0, float64(cw), float64(ch))
			cc.Fill()
			app.Window().DrawTo(render.NewCanvas(cc, cw, ch))
		}); err != nil {
			log.Printf("draw: %v", err)
			return
		}
		if err := canvas.Render(dc.RenderTarget()); err != nil {
			log.Printf("render: %v", err)
		}
	})
	gpuApp.OnClose(func() { gg.CloseAccelerator() })

	if err := gpuApp.Run(); err != nil {
		log.Fatal(err)
	}
}

func buildRoot(spec ui.Spec, theme *material3.Theme) widget.Widget {
	children := []widget.Widget{
		primitives.Text(spec.Title).FontSize(28).Bold(),
		primitives.Text("Evidence-first nuclear physics workspace for Moscovium isotope data.").FontSize(16),
	}
	for _, module := range spec.Modules {
		children = append(children,
			primitives.Box(
				primitives.Text(module.Name).FontSize(18).Bold(),
				primitives.Text(module.Description).FontSize(14),
			).Padding(16).Gap(8).Background(theme.Colors.SurfaceContainer),
		)
	}
	return primitives.Box(children...).Padding(28).Gap(14).Background(theme.Colors.Surface)
}
