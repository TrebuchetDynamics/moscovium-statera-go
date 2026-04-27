//go:build !cgo

package main

import (
	"flag"
	"fmt"
	"log"
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

func main() {
	screenshotPath := flag.String("screenshot", "", "write an offscreen PNG screenshot and exit")
	flag.Parse()

	model := ui.DefaultModel()
	seed := widget.Hex(0x2F5D50)
	materialTheme := material3.New(seed)
	if *screenshotPath != "" {
		if err := saveScreenshot(*screenshotPath, model, materialTheme); err != nil {
			log.Fatal(err)
		}
		return
	}
	if !displayAvailable(os.Getenv) {
		path := headlessScreenshotPath(os.TempDir())
		if err := saveScreenshot(path, model, materialTheme); err != nil {
			log.Fatal(err)
		}
		log.Printf("no display detected; wrote offscreen screenshot: %s", path)
		return
	}

	gpuApp := gogpu.NewApp(gogpu.DefaultConfig().
		WithTitle(model.Spec.Title).
		WithSize(model.Spec.Width, model.Spec.Height).
		WithContinuousRender(false))

	appTheme := uitheme.DefaultLight()
	appTheme.Colors.Primary = seed
	appTheme.Colors.PrimaryDark = widget.Hex(0x1F463C)
	appTheme.Colors.PrimaryLight = widget.Hex(0x6F9C8D)
	app := uiapp.New(
		uiapp.WithWindowProvider(gpuApp),
		uiapp.WithPlatformProvider(gpuApp),
		uiapp.WithEventSource(gpuApp.EventSource()),
		uiapp.WithTheme(appTheme),
	)
	app.SetRoot(buildRoot(model, materialTheme))

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

func displayAvailable(getenv func(string) string) bool {
	return getenv("WAYLAND_DISPLAY") != "" || getenv("DISPLAY") != ""
}

func headlessScreenshotPath(tempDir string) string {
	return filepath.Join(tempDir, "statera-ui-headless.png")
}

func saveScreenshot(path string, model ui.AppModel, theme *material3.Theme) error {
	appTheme := uitheme.DefaultLight()
	app := uiapp.New(
		uiapp.WithWindowProvider(gpucontext.NullWindowProvider{W: model.Spec.Width, H: model.Spec.Height}),
		uiapp.WithTheme(appTheme),
	)
	app.SetRoot(buildRoot(model, theme))
	app.Frame()

	dc := gg.NewContext(model.Spec.Width, model.Spec.Height)
	dc.SetRGBA(0.96, 0.97, 0.96, 1)
	dc.DrawRectangle(0, 0, float64(model.Spec.Width), float64(model.Spec.Height))
	dc.Fill()
	app.Window().DrawTo(render.NewCanvas(dc, model.Spec.Width, model.Spec.Height))
	return dc.SavePNG(path)
}

func buildRoot(model ui.AppModel, theme *material3.Theme) widget.Widget {
	railItems := make([]widget.Widget, 0, len(model.Views)+1)
	railItems = append(railItems,
		primitives.Text("Views").FontSize(13).Bold().Color(widget.Hex(0x24483E)),
	)
	for _, view := range model.Views {
		railItems = append(railItems, railItem(view.Name, view.Description))
	}
	rail := primitives.Box(railItems...).
		Width(220).
		Padding(16).
		Gap(10).
		Background(widget.Hex(0xE7EFEA)).
		Rounded(8)

	content := primitives.Box(
		header(model),
		dashboardSection(model, theme),
		physicsSection(model, theme),
		sourcesSection(model, theme),
		contextSection(model, theme),
	).Padding(20).Gap(14)

	return primitives.HBox(
		rail,
		scrollview.New(content, scrollview.PainterOpt(material3.ScrollbarPainter{Theme: theme})),
	).Padding(18).Gap(16).Background(theme.Colors.Surface)
}

func header(model ui.AppModel) widget.Widget {
	return primitives.Box(
		primitives.Text(model.Spec.Title).FontSize(26).Bold().Color(widget.Hex(0x183D34)),
		primitives.Text("Evidence-first workspace for Moscovium isotope data, source records, and quarantined context.").FontSize(14).Color(widget.Hex(0x44504B)),
	).Gap(6)
}

func dashboardSection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	return section("Dashboard", []widget.Widget{
		primitives.HBox(
			metric("Verified isotopes", fmt.Sprintf("%d", model.Summary.VerifiedIsotopes)),
			metric("Citation records", fmt.Sprintf("%d", model.Summary.CitationRecords)),
			metric("Context records", fmt.Sprintf("%d", model.Summary.ContextRecords)),
		).Gap(10),
		labelValue("Decay chain", strings.Join(model.Summary.DecayChain, " -> ")),
		boundary(model.Summary.BoundaryNotice),
	}, theme)
}

func physicsSection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	children := []widget.Widget{primitives.Text("Track A isotope records").FontSize(14).Bold().Color(widget.Hex(0x24483E))}
	for _, isotope := range model.Isotopes {
		children = append(children, card(
			primitives.Text(fmt.Sprintf("%s  Z=%d  A=%d", isotope.ID, isotope.Z, isotope.A)).FontSize(14).Bold(),
			primitives.Text(fmt.Sprintf("half-life %s | daughter %s | citations %d", isotope.HalfLife, isotope.Daughter, isotope.CitationCount)).FontSize(12),
			primitives.Text(isotope.SourcePath).FontSize(11).Color(widget.Hex(0x5F6F68)),
		))
	}
	children = append(children, primitives.Text("Validator examples").FontSize(14).Bold().Color(widget.Hex(0x24483E)))
	for _, example := range model.ClaimExamples {
		children = append(children, card(
			primitives.Text(example.Label).FontSize(13).Bold(),
			primitives.Text(string(example.Result.Status)).FontSize(12).Color(statusColor(string(example.Result.Status))),
			primitives.Text(example.Result.Reason).FontSize(11).Color(widget.Hex(0x5F6F68)),
		))
	}
	return section("Physics", children, theme)
}

func sourcesSection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	children := make([]widget.Widget, 0, len(model.SourceRecords))
	for _, source := range model.SourceRecords {
		children = append(children, card(
			primitives.Text(source.Key).FontSize(13).Bold(),
			primitives.Text(fmt.Sprintf("%s | %s | PDF: %s", source.Track, source.Status, source.PDF)).FontSize(12),
			primitives.Text(source.Identifier).FontSize(11).Color(widget.Hex(0x31574D)),
			primitives.Text(source.Relevance).FontSize(11).Color(widget.Hex(0x5F6F68)),
			primitives.Text(source.SourcePath).FontSize(11).Color(widget.Hex(0x5F6F68)),
		))
	}
	return section("Sources", children, theme)
}

func contextSection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	children := make([]widget.Widget, 0, len(model.ContextRecords))
	for _, record := range model.ContextRecords {
		children = append(children, card(
			primitives.Text(record.Title).FontSize(13).Bold(),
			primitives.Text(fmt.Sprintf("event %s | sources %d | simulation use: %s", record.EventDate, record.SourceCount, record.SimulationUse)).FontSize(12),
			primitives.Text(strings.Join(record.Labels, ", ")).FontSize(11).Color(widget.Hex(0x7A3A00)),
			primitives.Text(record.SourcePath).FontSize(11).Color(widget.Hex(0x5F6F68)),
		))
	}
	children = append(children, boundary("Context records are not-for-simulation and do not validate unsupported linkages."))
	return section("Context", children, theme)
}

func section(title string, children []widget.Widget, theme *material3.Theme) widget.Widget {
	items := []widget.Widget{primitives.Text(title).FontSize(18).Bold().Color(widget.Hex(0x183D34))}
	items = append(items, children...)
	return primitives.Box(items...).
		Padding(14).
		Gap(9).
		Background(theme.Colors.SurfaceContainer).
		Rounded(8).
		BorderStyle(1, widget.Hex(0xD4DED8))
}

func railItem(name string, _ string) widget.Widget {
	return primitives.Box(
		primitives.Text(name).FontSize(14).Bold().Color(widget.Hex(0x183D34)),
	).Padding(10).Gap(4).Background(widget.Hex(0xF6FAF7)).Rounded(6)
}

func metric(label string, value string) widget.Widget {
	return primitives.Box(
		primitives.Text(value).FontSize(22).Bold().Color(widget.Hex(0x183D34)),
		primitives.Text(label).FontSize(11).Color(widget.Hex(0x52645C)),
	).Width(160).Padding(10).Gap(3).Background(widget.Hex(0xF6FAF7)).Rounded(6)
}

func card(children ...widget.Widget) widget.Widget {
	return primitives.Box(children...).
		Padding(10).
		Gap(5).
		Background(widget.Hex(0xFFFFFF)).
		Rounded(6).
		BorderStyle(1, widget.Hex(0xD8E1DC))
}

func labelValue(label string, value string) widget.Widget {
	return card(
		primitives.Text(label).FontSize(11).Bold().Color(widget.Hex(0x52645C)),
		primitives.Text(value).FontSize(13).Color(widget.Hex(0x183D34)),
	)
}

func boundary(text string) widget.Widget {
	return primitives.Box(
		primitives.Text(text).FontSize(12).Color(widget.Hex(0x5C3B00)),
	).Padding(10).Background(widget.Hex(0xFFF4D8)).Rounded(6).BorderStyle(1, widget.Hex(0xE7C66A))
}

func statusColor(status string) widget.Color {
	switch status {
	case "supported-by-track-a":
		return widget.Hex(0x246B45)
	case "stability-incongruent", "outside-supported-model", "invalid-claim":
		return widget.Hex(0x8A3D00)
	default:
		return widget.Hex(0x44504B)
	}
}
