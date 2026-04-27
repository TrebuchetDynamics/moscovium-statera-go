# Statera Education Research Design App Design

## Goal

Turn the MVP dashboard into a multi-module learning workstation with three primary goals: Education, Research, and Design. The app should teach evaluated Moscovium physics, expose source provenance, and let users inspect constrained design scenarios without allowing unsupported claims or Track B context to become simulation data.

## Position

The app is an academic nuclear-physics tool, not a speculative propulsion simulator. It should make the evidence boundary useful and visible: Track A records support calculations and validator examples; Track B records explain public context but remain not-for-simulation.

## Modules

### Education

Education is the first user-facing goal. It presents a short learning path that explains:

- what evaluated isotope records contain
- how the `288Mc` decay-chain traversal works
- why half-life claims can be checked against Track A data
- why non-standard mechanisms are outside the supported model

The module should use concise lessons and deterministic validator examples. It should avoid free-form prompts in this slice so the learning path remains testable.

### Research

Research is the provenance module. It presents:

- citation records with track, status, identifier, PDF state, and source path
- read queues or source-status summaries
- a clear warning that `to-read` metadata is not extracted fact

The module should help a reviewer answer: which records exist, where they live, and whether they are safe to use as simulation substrate.

### Design

Design is a constrained sandbox, not an open-ended generator. It presents named design scenarios that show:

- the scenario goal
- required Track A inputs
- validator result
- whether the scenario can affect simulation defaults
- a source-backed constraint note

The initial scenarios should be fixed examples derived from existing validator behavior: a supported short half-life check, a long half-life stability failure, an unsupported mechanism failure, and a mixed-claim rejection.

### Context

Context remains a boundary module. It lists Track B context records with event date, labels, source count, source path, and `simulation use: prohibited`. It must state that context records do not validate unsupported linkages.

## Information Architecture

The left rail lists the primary modules in this order:

1. Education
2. Research
3. Design
4. Context

The main content renders vertically stacked sections in the same order. The app can later add interactive navigation, but this implementation should keep rendering deterministic for screenshots and tests.

## Data Model

Extend `internal/ui.AppModel` with module-specific records:

- `EducationLessons []LessonRecord`
- `ResearchItems []ResearchItem`
- `DesignScenarios []DesignScenario`
- existing `ContextRecords []ContextRecord`

Each record must include source or constraint provenance. Design scenarios must include `physics.ClaimResult` values from the existing validator rather than duplicating outcome strings by hand.

## Rendering

Keep `cmd/statera-ui` as a renderer over `internal/ui.AppModel`.

- Education uses compact lesson cards.
- Research uses source cards with status and provenance.
- Design uses scenario cards with status coloring from validator results.
- Context keeps the existing quarantine/boundary treatment.

The visual style remains restrained: compact panels, no landing page, no speculative imagery, no nested cards, and card radii no greater than 8px.

## Non-Goals

- No RAG search interface in this slice.
- No free-form claim editor.
- No solver implementation.
- No import of Track B data into `physics.Catalog`.
- No new factual assertions from `to-read` sources.

## Acceptance Criteria

- `DefaultSpec()` exposes modules named Education, Research, Design, and Context in that order.
- `DefaultModel()` includes at least four education lessons, at least three research items, and at least four design scenarios.
- Tests assert every education lesson has a source path.
- Tests assert every design scenario has `SimulationUse` set to `allowed` or `blocked`, and unsupported examples are blocked.
- Tests assert research items preserve Track A or Track B labels and source paths.
- The UI screenshot shows Education, Research, and Design sections in the first-pass multi-module app.
- Required verification commands pass:

```bash
CGO_ENABLED=0 go test ./...
go run ./cmd/statera
CGO_ENABLED=0 go run ./cmd/statera-ui
CGO_ENABLED=0 go run ./cmd/statera-ui -screenshot screenshots/statera-education-research-design-2026-04-27.png
git diff --check
```
