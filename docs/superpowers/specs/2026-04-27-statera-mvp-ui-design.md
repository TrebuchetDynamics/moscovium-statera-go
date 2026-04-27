# Statera MVP UI Design

## Goal

Build the first usable Moscovium Statera desktop app as a compact research workstation. The app should let a user inspect evaluated Moscovium data, test simple stability claims, browse source records, and view quarantined context records without weakening the repository's evidence boundary.

## Product Position

The MVP is not a landing page, speculative propulsion tool, or general conspiracy browser. It is a source-first nuclear-physics interface with a controlled context archive. The app should feel like a quiet laboratory dashboard: dense, readable, and built for repeated inspection.

## Users

- Students learning how evaluated nuclear data differs from unsupported public claims.
- Researchers or reviewers checking isotope provenance and decay-chain behavior.
- Project maintainers verifying that Track A simulation inputs and Track B context records remain separated.

## App Goals

1. Show what Track A currently knows about Moscovium isotopes.
2. Show the decay-chain path used by the CLI.
3. Let a user evaluate a small fixed set of example claims against the Track A validator.
4. List citation records and their read status without promoting `to-read` sources into facts.
5. List Track B context records with visible `not-for-simulation` boundaries.
6. Keep all text sober and source-scoped.

## MVP Views

### Dashboard

The first view summarizes the workspace state:

- verified isotope count
- citation record count
- context record count
- current decay-chain smoke path
- a short boundary notice that Track B context is excluded from simulation

This view should be the first screen because it tells the user what the app can currently prove.

### Physics

The Physics view focuses on runtime behavior:

- list Track A isotope records from the current seed data
- show `Z`, `A`, half-life, daughter, and citation count
- show the `288Mc` decay-chain traversal
- show fixed validator examples:
  - `288Mc minimum half-life >= 100ms` -> supported by Track A
  - `288Mc minimum half-life >= 1h` -> stability-incongruent
  - `antigravity propulsion` -> outside-supported-model
  - mixed half-life plus mechanism -> invalid-claim

The MVP should use fixed examples instead of a free-form claim editor so the first slice stays deterministic and testable.

### Sources

The Sources view lists citation records from `citations/papers/`:

- key
- track
- status
- DOI or URL
- PDF state
- one-line relevance

Records marked `to-read` must be presented as metadata records, not extracted facts.

### Context

The Context view lists normalized Track B records from `docs/lore/records/`:

- title
- labels
- event date when available
- source count
- simulation use status

The UI should make `not-for-simulation` visible on every context record. It should not imply that unsupported linkages are validated.

## Layout

Use a single-window `gogpu/ui` app with a left navigation rail and a main content area. The rail contains four module buttons: Dashboard, Physics, Sources, Context. The main area renders the selected view.

For the MVP, if interactive tab/button state is too expensive for the first pass, it is acceptable to render all four views as vertically stacked sections with stable headings. The visible content and testable data model matter more than navigation polish in the first slice.

## Visual Style

- Use Material 3 theme with the existing green seed `0x2F5D50`.
- Avoid decorative hero sections, oversized cards, or unsupported imagery.
- Use full-width bands and compact panels.
- Use cards only for repeated records.
- Keep cards at 8px radius or less.
- Use clear status labels, small tables, and source paths.

## Architecture

Create a UI-facing data model in `internal/ui` that is independent of `cmd/statera-ui` rendering:

- `AppModel`: holds summary metrics, views, physics examples, source records, and context records.
- `ViewSpec`: names and describes each view.
- `ClaimExample`: wraps `physics.Claim` and `physics.ClaimResult`.
- `SourceRecord`: metadata extracted from the current citation files.
- `ContextRecord`: metadata extracted from the current lore files.

The renderer in `cmd/statera-ui` should consume `internal/ui.AppModel` and build `gogpu/ui` widgets. This keeps most MVP behavior testable without a graphical environment.

## Data Sources

The MVP may use small hard-coded metadata derived from existing files, provided each record points back to the source path. It should not parse Markdown in the first slice unless that becomes simpler than keeping small source-backed metadata tables.

Track A runtime data should use the existing `physics.Catalog` shape and the existing validator. Lore/context data must not enter `physics.Catalog`.

## Screenshot Strategy

The primary screenshot path is:

```bash
CGO_ENABLED=0 go run ./cmd/statera-ui
```

If the local environment has no Wayland/X display, the app may fail at window creation. In that case the implementation should still provide:

- headless `internal/ui` tests proving all view text/data exists
- `CGO_ENABLED=0 go test ./...`
- a recorded failure message for the graphical run

If a graphical environment is available, capture screenshots under `screenshots/` and keep them out of runtime data.

## Non-Goals

- No free-form RAG search in this MVP.
- No SQLite-vec or WASM embeddings.
- No Hartree-Fock solver UI.
- No extraction of new facts from `to-read` papers.
- No claim that Track B events support Element 115 physics.

## Acceptance Criteria

- The app model exposes Dashboard, Physics, Sources, and Context views.
- Tests assert the four views exist and include the required boundary language.
- Tests assert validator examples produce `supported-by-track-a`, `stability-incongruent`, `outside-supported-model`, and `invalid-claim`.
- The UI renderer builds from the app model without duplicating scientific claims directly in `cmd/statera-ui`.
- `CGO_ENABLED=0 go test ./...` passes.
- `go run ./cmd/statera` still prints the existing decay-chain smoke path.
- `CGO_ENABLED=0 go run ./cmd/statera-ui` is attempted; screenshots are captured only if the local graphical environment supports it.
