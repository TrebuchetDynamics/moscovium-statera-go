# Statera Education Research Design Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the multi-module Statera app around Education, Research, and Design while preserving the Context quarantine boundary.

**Architecture:** Extend `internal/ui.AppModel` with module-specific data records and keep `cmd/statera-ui` as a pure renderer over that model. The Design module reuses `physics.ValidateClaim` so status labels are computed by existing physics code rather than copied into UI strings. Rendering stays deterministic for offscreen screenshots.

**Tech Stack:** Go 1.25, pure Go, zero CGO, `github.com/gogpu/ui`, Material3 theme, existing Markdown provenance files.

---

## File Structure

- Modify `internal/ui/app.go`: add `LessonRecord`, `ResearchItem`, and `DesignScenario`; update default modules and summaries; keep source-backed static records small.
- Modify `internal/ui/app_test.go`: add tests for Education, Research, and Design model requirements.
- Modify `cmd/statera-ui/main.go`: replace Dashboard/Physics/Sources sections with Education/Research/Design sections plus Context.
- Modify `screenshots/README.md`: document the new screenshot artifact.
- Create `screenshots/statera-education-research-design-2026-04-27.png`: deterministic offscreen screenshot.

## Task 1: Model Tests

**Files:**
- Modify: `internal/ui/app_test.go`

- [ ] **Step 1: Write failing tests**

Add tests that assert:

- `DefaultSpec()` exposes modules named `Education`, `Research`, `Design`, and `Context` in that order.
- `DefaultModel()` has at least four education lessons, at least three research items, and at least four design scenarios.
- Every education lesson has a title, objective, and source path.
- Every research item has a track of `Track A` or `Track B` and a source path.
- Every design scenario has `SimulationUse` set to `allowed` or `blocked`, and non-supported validator statuses are blocked.

Run:

```bash
CGO_ENABLED=0 go test ./internal/ui
```

Expected: FAIL because the new model fields and records do not exist yet.

## Task 2: Model Implementation

**Files:**
- Modify: `internal/ui/app.go`

- [ ] **Step 1: Add model types and default records**

Implement:

```go
type LessonRecord struct {
	Title      string
	Objective  string
	Concept    string
	SourcePath string
}

type ResearchItem struct {
	Key        string
	Track      string
	Status     string
	Identifier string
	PDF        string
	Relevance  string
	SourcePath string
}

type DesignScenario struct {
	Name          string
	Goal          string
	Inputs        []string
	Result        physics.ClaimResult
	SimulationUse string
	Constraint    string
	SourcePath    string
}
```

Update `AppModel` to include `EducationLessons`, `ResearchItems`, and `DesignScenarios`.

- [ ] **Step 2: Update default modules**

Update `DefaultSpec()` so `Spec.Modules` is:

```text
Education
Research
Design
Context
```

- [ ] **Step 3: Build default records**

Create deterministic helper functions:

- `defaultEducationLessons() []LessonRecord`
- `defaultResearchItems() []ResearchItem`
- `defaultDesignScenarios(catalog physics.Catalog) []DesignScenario`

Use existing source paths:

- `data/research.seed.json`
- `citations/papers/oganessian2022mcfactory.md`
- `citations/papers/iupac2016names.md`
- `citations/papers/houseoversight2026missing-scientists-letter.md`
- `docs/lore/records/mccasland-2026-missing-person.md`

For design scenarios, call `physics.ValidateClaim` for each scenario. Set `SimulationUse` to `allowed` only when the result status is `supported-by-track-a`; otherwise set it to `blocked`.

- [ ] **Step 4: Verify model tests**

Run:

```bash
CGO_ENABLED=0 go test ./internal/ui
```

Expected: PASS.

## Task 3: UI Renderer

**Files:**
- Modify: `cmd/statera-ui/main.go`

- [ ] **Step 1: Replace section composition**

Update `buildRoot()` to render:

```go
educationSection(model, theme)
researchSection(model, theme)
designSection(model, theme)
contextSection(model, theme)
```

Keep the rail generated from `model.Views`.

- [ ] **Step 2: Add section renderers**

Add:

- `educationSection(model ui.AppModel, theme *material3.Theme) widget.Widget`
- `researchSection(model ui.AppModel, theme *material3.Theme) widget.Widget`
- `designSection(model ui.AppModel, theme *material3.Theme) widget.Widget`

Education cards display title, objective, concept, and source path.

Research cards display key, track, status, identifier, PDF state, relevance, and source path.

Design cards display name, goal, inputs, validator status, simulation use, constraint, and source path.

- [ ] **Step 3: Verify full tests**

Run:

```bash
CGO_ENABLED=0 go test ./...
```

Expected: PASS.

## Task 4: Screenshot And Docs

**Files:**
- Modify: `screenshots/README.md`
- Create: `screenshots/statera-education-research-design-2026-04-27.png`

- [ ] **Step 1: Update screenshot docs**

Document:

```bash
CGO_ENABLED=0 go run ./cmd/statera-ui -screenshot screenshots/statera-education-research-design-2026-04-27.png
```

- [ ] **Step 2: Generate screenshot**

Run:

```bash
CGO_ENABLED=0 go run ./cmd/statera-ui -screenshot screenshots/statera-education-research-design-2026-04-27.png
```

Expected: PNG written at 1180x760.

## Task 5: Required Verification And Commit

**Files:**
- All changed files.

- [ ] **Step 1: Run required verification**

Run:

```bash
CGO_ENABLED=0 go test ./...
go run ./cmd/statera
CGO_ENABLED=0 go run ./cmd/statera-ui
CGO_ENABLED=0 go run ./cmd/statera-ui -screenshot screenshots/statera-education-research-design-2026-04-27.png
git diff --check
```

Expected:

- all tests pass
- CLI prints `288Mc -> 284Nh -> 280Rg -> 276Mt -> 272Bh -> 268Db -> 264Lr`
- UI command writes a headless screenshot in no-display environments
- screenshot command writes the tracked PNG
- diff check has no output

- [ ] **Step 2: Commit implementation**

Run:

```bash
git add internal/ui/app.go internal/ui/app_test.go cmd/statera-ui/main.go screenshots/README.md screenshots/statera-education-research-design-2026-04-27.png docs/superpowers/plans/2026-04-27-statera-education-research-design.md
git commit -m "feat: build education research design app"
```

## Self-Review Notes

- Spec coverage: every module in the spec maps to model tests and renderer work.
- Placeholder scan: no unresolved placeholders or hand-waved implementation steps remain.
- Type consistency: names in renderer tasks match the model types defined in Task 2.
