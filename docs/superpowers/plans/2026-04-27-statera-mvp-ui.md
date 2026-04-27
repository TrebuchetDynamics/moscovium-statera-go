# Statera MVP UI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the first usable Statera desktop MVP with Dashboard, Physics, Sources, and Context views backed by a testable UI data model.

**Architecture:** Keep scientific/runtime behavior in `internal/ui` and `internal/physics`, and keep `cmd/statera-ui` as a renderer over a prepared `ui.AppModel`. The first slice uses source-backed static metadata tables for citations/context and the existing `physics.ValidateClaim` function for validator examples. `gogpu/ui` rendering follows the existing event-driven `gogpu.NewApp` plus `app.New` pattern.

**Tech Stack:** Go 1.25, pure Go, zero CGO, `github.com/gogpu/ui`, `github.com/gogpu/gogpu`, `github.com/gogpu/gg`, Markdown source records.

---

## File Structure

- Modify `internal/ui/app.go`: replace the minimal shell spec with a richer `AppModel`, deterministic source/context metadata, Track A catalog helpers, and validator examples.
- Modify `internal/ui/app_test.go`: test model views, summary counts, boundary text, and validator result statuses.
- Modify `cmd/statera-ui/main.go`: render the richer model using compact sections and cards while preserving the existing `gogpu/ui` draw loop.
- Create `screenshots/README.md`: document screenshot capture path and local display limitation.

## Task 1: UI App Model

**Files:**
- Modify: `internal/ui/app.go`
- Modify: `internal/ui/app_test.go`

- [ ] **Step 1: Write failing tests**

Add tests that assert:

- `DefaultModel()` exposes views named `Dashboard`, `Physics`, `Sources`, and `Context`.
- Summary counts are nonzero and the boundary text contains `Track B context is excluded from simulation`.
- Claim examples include statuses `supported-by-track-a`, `stability-incongruent`, `outside-supported-model`, and `invalid-claim`.
- Context records all have simulation use `prohibited`.

Run:

```bash
CGO_ENABLED=0 go test ./internal/ui
```

Expected: FAIL because `DefaultModel` and the new types do not exist.

- [ ] **Step 2: Implement the model**

Replace `internal/ui/app.go` with a deterministic UI model:

- Keep `Spec` with `Title`, `Width`, `Height`.
- Add `AppModel`, `ViewSpec`, `Summary`, `IsotopeRecord`, `ClaimExample`, `SourceRecord`, and `ContextRecord`.
- Add `DefaultModel() AppModel`.
- Build a Track A catalog with `288Mc`, daughter chain records for traversal, and `290Mc`.
- Use `physics.DecayChain` and `physics.ValidateClaim` for runtime-derived model fields.
- Use citation/context metadata pointing to existing repo files.

- [ ] **Step 3: Verify tests**

Run:

```bash
CGO_ENABLED=0 go test ./internal/ui
```

Expected: PASS.

## Task 2: MVP Renderer

**Files:**
- Modify: `cmd/statera-ui/main.go`

- [ ] **Step 1: Update renderer over `ui.AppModel`**

Modify `cmd/statera-ui/main.go` so:

- `main()` calls `ui.DefaultModel()`.
- Window title/size come from `model.Spec`.
- `buildRoot(model, materialTheme)` renders four sections: Dashboard, Physics, Sources, Context.
- Cards stay compact with radius <= 8.
- Text includes `not-for-simulation` on context records.

- [ ] **Step 2: Compile and test**

Run:

```bash
CGO_ENABLED=0 go test ./...
```

Expected: PASS.

Run:

```bash
go run ./cmd/statera
```

Expected output:

```text
288Mc -> 284Nh -> 280Rg -> 276Mt -> 272Bh -> 268Db -> 264Lr
```

## Task 3: Screenshot Documentation

**Files:**
- Create: `screenshots/README.md`

- [ ] **Step 1: Add screenshot instructions**

Create a README documenting:

- primary command: `CGO_ENABLED=0 go run ./cmd/statera-ui`
- expected local limitation if no graphical display: `wayland: WAYLAND_DISPLAY not set`
- screenshots should be stored in `screenshots/` when a graphical environment is available
- screenshots are visual verification artifacts, not simulation data

- [ ] **Step 2: Verify docs**

Run:

```bash
git diff --check
```

Expected: no output.

## Task 4: Full Verification And Screenshot Attempt

**Files:**
- No required tracked file changes.

- [ ] **Step 1: Run full tests**

Run:

```bash
CGO_ENABLED=0 go test ./...
```

Expected: PASS.

- [ ] **Step 2: Run CLI smoke path**

Run:

```bash
go run ./cmd/statera
```

Expected output:

```text
288Mc -> 284Nh -> 280Rg -> 276Mt -> 272Bh -> 268Db -> 264Lr
```

- [ ] **Step 3: Attempt UI run**

Run:

```bash
timeout 20s env CGO_ENABLED=0 go run ./cmd/statera-ui
```

Expected in this environment: either the app starts and times out, or it fails with `wayland: WAYLAND_DISPLAY not set`. Record exact output.

- [ ] **Step 4: Check git state**

Run:

```bash
git diff --check
git status --short
```

Expected: no tracked changes after final commit.

## Self-Review Notes

- Spec coverage: this plan implements the model, visible MVP sections, source/context boundaries, and screenshot documentation. It does not add RAG, parsing, or solvers.
- Placeholder scan: no unresolved placeholders remain.
- Type consistency: all new public names are defined in Task 1 before use in Task 2.
