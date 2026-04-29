//go:build !cgo

package main

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
		moduleOverview(model),
		educationSection(model, theme),
		researchSection(model, theme),
		designSection(model, theme),
		alphaSystematicsSection(model, theme),
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
		primitives.Text("Education, research, and constrained design for source-backed Moscovium data.").FontSize(14).Color(widget.Hex(0x44504B)),
		primitives.Text(fmt.Sprintf("Track A isotopes %d | research records %d | context records %d | %s", model.Summary.VerifiedIsotopes, model.Summary.CitationRecords, model.Summary.ContextRecords, strings.Join(model.Summary.DecayChain, " -> "))).FontSize(11).Color(widget.Hex(0x52645C)),
	).Gap(6)
}

type summaryCardSpec struct {
	Name        string
	Metric      string
	Description string
	Width       float32
}

func demoSummaryRows(model ui.AppModel) [][]summaryCardSpec {
	cards := []summaryCardSpec{
		{Name: "Education", Metric: fmt.Sprintf("%d lessons", len(model.EducationLessons)), Description: "Learn evaluated records and validator boundaries.", Width: 420},
		{Name: "Research", Metric: fmt.Sprintf("%d records", len(model.ResearchItems)), Description: "Inspect DOI, URL, queue, and provenance status.", Width: 420},
		{Name: "Design", Metric: fmt.Sprintf("%d scenarios", len(model.DesignScenarios)), Description: "Test constrained scenarios against Track A.", Width: 420},
		{Name: "Alpha", Metric: fmt.Sprintf("%d isotopes", len(model.AlphaSystematics)), Description: "Compare Royer model predictions to evaluated half-lives.", Width: 420},
	}
	return [][]summaryCardSpec{{cards[0], cards[1]}, {cards[2], cards[3]}}
}

func moduleOverview(model ui.AppModel) widget.Widget {
	rows := demoSummaryRows(model)
	rowWidgets := make([]widget.Widget, 0, len(rows))
	for _, row := range rows {
		cards := make([]widget.Widget, 0, len(row))
		for _, spec := range row {
			cards = append(cards, moduleSummary(spec))
		}
		rowWidgets = append(rowWidgets, primitives.HBox(cards...).Gap(12))
	}
	return primitives.Box(rowWidgets...).Gap(10)
}

func moduleSummary(spec summaryCardSpec) widget.Widget {
	return primitives.Box(
		primitives.Text(spec.Name).FontSize(14).Bold().Color(widget.Hex(0x183D34)),
		primitives.Text(spec.Metric).FontSize(12).Color(widget.Hex(0x246B45)),
		primitives.Text(spec.Description).FontSize(11).Color(widget.Hex(0x52645C)),
	).Width(spec.Width).Padding(10).Gap(4).Background(widget.Hex(0xFFFFFF)).Rounded(6).BorderStyle(1, widget.Hex(0xD8E1DC))
}

func educationSection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	children := []widget.Widget{
		primitives.Text("Guided learning path").FontSize(14).Bold().Color(widget.Hex(0x24483E)),
		boundary(model.Summary.BoundaryNotice),
	}
	for _, lesson := range model.EducationLessons {
		children = append(children, card(
			primitives.Text(lesson.Title).FontSize(13).Bold(),
			primitives.Text(lesson.Objective).FontSize(11).Color(widget.Hex(0x44504B)),
			primitives.Text(lesson.Concept).FontSize(10).Color(widget.Hex(0x5F6F68)),
			primitives.Text(lesson.SourcePath).FontSize(10).Color(widget.Hex(0x31574D)),
		))
	}
	return section("Education", children, theme)
}

func researchSection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	children := []widget.Widget{
		primitives.Text("Source inventory").FontSize(14).Bold().Color(widget.Hex(0x24483E)),
	}
	for _, source := range model.ResearchItems {
		children = append(children, card(
			primitives.Text(source.Key).FontSize(13).Bold(),
			primitives.Text(fmt.Sprintf("%s | %s | PDF: %s", source.Track, source.Status, source.PDF)).FontSize(11),
			primitives.Text(source.Identifier).FontSize(10).Color(widget.Hex(0x31574D)),
			primitives.Text(source.Relevance).FontSize(10).Color(widget.Hex(0x5F6F68)),
			primitives.Text(source.SourcePath).FontSize(10).Color(widget.Hex(0x5F6F68)),
		))
	}
	return section("Research", children, theme)
}

func designSection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	children := []widget.Widget{
		primitives.Text("Constrained scenarios").FontSize(14).Bold().Color(widget.Hex(0x24483E)),
	}
	for _, scenario := range model.DesignScenarios {
		children = append(children, card(
			primitives.Text(scenario.Name).FontSize(13).Bold(),
			primitives.Text(scenario.Goal).FontSize(11).Color(widget.Hex(0x44504B)),
			primitives.Text(fmt.Sprintf("inputs: %s", strings.Join(scenario.Inputs, ", "))).FontSize(10).Color(widget.Hex(0x5F6F68)),
			primitives.Text(fmt.Sprintf("%s | simulation use: %s", scenario.Result.Status, scenario.SimulationUse)).FontSize(11).Color(statusColor(string(scenario.Result.Status))),
			primitives.Text(scenario.Constraint).FontSize(10).Color(widget.Hex(0x5F6F68)),
			primitives.Text(scenario.SourcePath).FontSize(10).Color(widget.Hex(0x31574D)),
		))
	}
	return section("Design", children, theme)
}

func alphaSystematicsSection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	children := []widget.Widget{
		primitives.Text("Royer formula vs evaluated half-lives").FontSize(14).Bold().Color(widget.Hex(0x24483E)),
		boundary("Predictions are peer-reviewed-model output. They never replace evaluated half-lives."),
	}
	for _, record := range model.AlphaSystematics {
		if record.Skipped {
			children = append(children, card(
				primitives.Text(record.IsotopeID).FontSize(13).Bold(),
				skipNotice(fmt.Sprintf("skipped: %s", record.SkipReason)),
				primitives.Text(fmt.Sprintf("%s | %s | %s", record.ModelName, record.ModelReference, record.EvidenceClass)).FontSize(10).Color(widget.Hex(0x31574D)),
				primitives.Text(record.SourcePath).FontSize(10).Color(widget.Hex(0x5F6F68)),
			))
			continue
		}
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

func card(children ...widget.Widget) widget.Widget {
	return primitives.Box(children...).
		Padding(10).
		Gap(5).
		Background(widget.Hex(0xFFFFFF)).
		Rounded(6).
		BorderStyle(1, widget.Hex(0xD8E1DC))
}

func boundary(text string) widget.Widget {
	return primitives.Box(
		primitives.Text(text).FontSize(12).Color(widget.Hex(0x5C3B00)),
	).Padding(10).Background(widget.Hex(0xFFF4D8)).Rounded(6).BorderStyle(1, widget.Hex(0xE7C66A))
}

func skipNotice(text string) widget.Widget {
	return primitives.Box(
		primitives.Text(text).FontSize(11).Color(widget.Hex(0x4A4A4A)),
	).Padding(8).Background(widget.Hex(0xEFEFEF)).Rounded(6).BorderStyle(1, widget.Hex(0xCFCFCF))
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
