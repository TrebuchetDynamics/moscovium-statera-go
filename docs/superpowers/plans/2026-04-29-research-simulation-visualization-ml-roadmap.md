# Research Simulation, Visualization, and ML Roadmap Implementation Plan

> **For Hermes:** Use `subagent-driven-development` only after Juan approves execution. Implement task-by-task with strict TDD for behavior changes.

**Goal:** Turn Moscovium Statera Go into a source-backed superheavy-nuclide research workbench with evaluated-data simulations, evidence-preserving visualizations, and safe ML-assisted research tooling.

**Architecture:** Keep the engine deterministic, pure Go, and source-first. Separate evaluated data, peer-reviewed model outputs, visual context, and ML suggestions so no UI or simulation path can blur their evidence class. Build the roadmap as small vertical slices: data schema + validation, physics calculation, UI model, screenshot, docs, and verification.

**Tech Stack:** Go 1.26.x on this host, pure-Go/CGO-free execution, `gogpu/ui` for UI shell, PNG screenshot generation through `cmd/statera-ui -screenshot`, JSON seed data in `data/`, citation markdown under `citations/`, tests with Go `testing`.

**Current verified context, 2026-04-29 07:50 CST:**

- Repo: `/home/xel/git/sages-openclaw/workspace-riju/moscovium-statera-go`
- Branch: `main`
- Current HEAD before this plan: `41c3cec15270`
- `go`: found at `/home/xel/.local/go-current/bin/go`
- `git`: found at `/usr/bin/git`
- `rg`: missing; use `grep` or Hermes `search_files` fallback
- Current seed records: `2` (`288Mc`, `290Mc`)
- Current tracked screenshots before this plan: `4`
- Current all-test event count from last run: `113` pass, `0` fail, `0` skip
- Existing hard boundary: only evaluated nuclear data, peer-reviewed nuclear physics, and explicitly labeled theory may enter calculations.

---

## Evidence Rules That Govern Every Slice

1. **Evaluated data** may come from ENSDF, NuDat, NUBASE, AME, NNDC, IAEA, IUPAC, GSI, JINR, or comparable standards bodies.
2. **Peer-reviewed model output** must name model, source, assumptions, coefficients, and evidence class.
3. **Visual context** may teach or explain but must not seed simulation defaults.
4. **ML output** is always suggestion, triage, extraction candidate, or anomaly flag until validated by source-backed tests.
5. Missing provenance is a blocking error.
6. No scientific value is accepted unless a test uses the value and a record preserves source URL / DOI trail / uncertainty when available.

---

## Milestone 1 — Evaluated Isotope Workbook v1

**Purpose:** Expand from two Moscovium seed records into a source-backed isotope workbook that can drive simulations, UI cards, and provenance graphs.

**Initial scope:** Add validated data structures and UI for the current chain, but only add daughter isotope scientific values after source URLs/DOIs are confirmed.

### Task 1.1: Add workbook data types without changing seed values

**Objective:** Define UI-neutral workbook records that can represent isotope identity, evaluated values, uncertainty, and evidence class.

**Files:**

- Create: `internal/research/workbook.go`
- Create: `internal/research/workbook_test.go`

**TDD RED:**

Add a test requiring `WorkbookFromSeed(seed, sourcePath)` to produce workbook records for `288Mc` and `290Mc` with:

- `ID`
- `Z`
- `A`
- `N`
- `HalfLifeSeconds`
- `QAlphaMeV`
- `Daughter`
- `EvidenceLevel`
- `SourcePath`
- at least one citation URL
- at least one DOI

Run:

```bash
CGO_ENABLED=0 go test ./internal/research -run TestWorkbookFromSeedPreservesProvenance -count=1 -v
```

Expected RED:

```text
undefined: WorkbookFromSeed
```

**GREEN implementation:**

Implement minimal conversion from existing validated `ResearchSeed` records. Do not add new isotope values.

**Verification:**

```bash
CGO_ENABLED=0 go test ./internal/research -count=1 -v
```

**Commit:**

```bash
git add internal/research/workbook.go internal/research/workbook_test.go
git commit -m "Add source-backed isotope workbook records"
```

### Task 1.2: Add workbook validation gate

**Objective:** Reject workbook records that lose provenance or identity consistency.

**Files:**

- Modify: `internal/research/workbook.go`
- Modify: `internal/research/workbook_test.go`

**TDD cases:**

- missing citation URLs rejects
- missing DOI rejects
- `N != A-Z` rejects
- blank `EvidenceLevel` rejects
- blank `SourcePath` rejects
- daughter ID malformed rejects when present

**Run:**

```bash
CGO_ENABLED=0 go test ./internal/research -run TestValidateWorkbookRejectsProvenanceLoss -count=1 -v
```

### Task 1.3: Add UI model for isotope workbook

**Objective:** Expose workbook records to `internal/ui.AppModel` without importing UI libraries into research code.

**Files:**

- Modify: `internal/ui/app.go`
- Modify: `internal/ui/app_test.go`

**TDD cases:**

- `DefaultModel()` exposes at least `2` workbook records.
- Every workbook UI record has ID, source path, citation count, DOI count, and evidence class.

**Run:**

```bash
CGO_ENABLED=0 go test ./internal/ui -run TestDefaultModelExposesIsotopeWorkbookRecords -count=1 -v
```

### Task 1.4: Render workbook section and screenshot

**Objective:** Add a visible `Isotope Workbook` section to the UI and generate a tested screenshot.

**Files:**

- Modify: `cmd/statera-ui/main.go`
- Modify: `cmd/statera-ui/screenshot_test.go`
- Create: `screenshots/statera-isotope-workbook-YYYY-MM-DD.png`

**TDD cases:**

- layout helper returns non-empty workbook cards
- cards expose ID, half-life seconds, Q-alpha MeV, DOI count, and source path
- screenshot remains `1180 x 760 px`, non-blank, and > `10 KiB`

**Run:**

```bash
CGO_ENABLED=0 go test ./cmd/statera-ui -run 'Workbook|Screenshot' -count=1 -v
CGO_ENABLED=0 go run ./cmd/statera-ui -screenshot screenshots/statera-isotope-workbook-YYYY-MM-DD.png
```

---

## Milestone 2 — Decay-Chain Simulation v1

**Purpose:** Convert the current deterministic chain into a research simulation with exact decay constants, mean life, and Monte Carlo timing, labeled as simulation derived from evaluated half-life inputs.

### Task 2.1: Add deterministic decay constants

**Files:**

- Modify: `internal/physics/decay.go`
- Modify: `internal/physics/decay_test.go`

**Calculations:**

For half-life `T_1/2` in seconds:

```text
lambda = ln(2) / T_1/2
mean_life = 1 / lambda = T_1/2 / ln(2)
```

**TDD cases:**

- `288Mc` half-life `0.17 s` gives positive decay constant.
- mean life is greater than half-life by factor `1 / ln(2)`.
- zero/negative half-life rejects.

**Run:**

```bash
CGO_ENABLED=0 go test ./internal/physics -run TestDecayConstantsFromHalfLife -count=1 -v
```

### Task 2.2: Add deterministic Monte Carlo API with fixed seed

**Files:**

- Create: `internal/physics/decay_montecarlo.go`
- Create: `internal/physics/decay_montecarlo_test.go`

**API sketch:**

```go
type DecaySample struct {
    IsotopeID string
    SampleIndex int
    TimeSeconds float64
}

type DecaySimulationSummary struct {
    IsotopeID string
    Samples int
    HalfLifeSeconds float64
    MedianSeconds float64
    P05Seconds float64
    P95Seconds float64
}
```

**TDD cases:**

- fixed RNG seed produces repeatable summaries
- `Samples <= 0` rejects
- median is positive
- `P05 <= Median <= P95`

**Run:**

```bash
CGO_ENABLED=0 go test ./internal/physics -run TestDecayMonteCarlo -count=1 -v
```

### Task 2.3: Add simulation UI records

**Files:**

- Modify: `internal/ui/app.go`
- Modify: `internal/ui/app_test.go`
- Modify: `cmd/statera-ui/main.go`

**UI fields:**

- isotope ID
- half-life seconds
- decay constant `s^-1`
- mean life seconds
- Monte Carlo sample count
- median / 5th / 95th percentile seconds
- evidence label: `simulation-derived-from-evaluated-half-life`

---

## Milestone 3 — Evidence and Provenance Graph v1

**Purpose:** Make source-to-datum trust visible and testable.

### Task 3.1: Add graph data model

**Files:**

- Create: `internal/research/provenance_graph.go`
- Create: `internal/research/provenance_graph_test.go`

**Node types:**

```text
isotope-record
citation-url
doi
evidence-class
model
blocked-source
ui-output
```

**Edge types:**

```text
record-cites-url
record-cites-doi
record-has-evidence-class
model-uses-source
ui-renders-record
source-blocked-by-access
```

**TDD cases:**

- graph generated from seed has at least one isotope node per seed record
- each isotope node has DOI and citation URL edges
- no orphan DOI nodes
- duplicate DOI edges are rejected or deduplicated deterministically

### Task 3.2: Add blocked-source nodes for Royer/Wang access blockers

**Files:**

- Modify: `internal/research/provenance_graph.go`
- Modify: `internal/research/provenance_graph_test.go`
- Read-only evidence source: `docs/academic-calculation-visualization-roadmap.md`

**TDD cases:**

- Royer 2008 blocker appears as blocked source with DOI `10.1103/PhysRevC.77.037602`
- Wang 2015 blocker appears as blocked source with DOI `10.1103/PhysRevC.92.064301`
- both are marked `not-simulation-default`

### Task 3.3: Render graph summary in UI

**Files:**

- Modify: `internal/ui/app.go`
- Modify: `cmd/statera-ui/main.go`
- Create screenshot: `screenshots/statera-provenance-graph-YYYY-MM-DD.png`

**UI requirements:**

- count isotope nodes
- count DOI nodes
- count citation URL nodes
- count blocked-source nodes
- list blocked sources textually
- no visual-only graph; every graph fact must be in readable text

---

## Milestone 4 — Alpha Systematics Residual Lab v2

**Purpose:** Improve the existing Royer model lab without adding unsupported coefficients.

### Task 4.1: Add residual summary stats

**Files:**

- Modify: `internal/physics/alpha.go`
- Modify: `internal/physics/alpha_test.go`

**Calculations:**

For each isotope with evaluated half-life and prediction:

```text
log_residual = log10(predicted_seconds / evaluated_seconds)
factor_error = predicted_seconds / evaluated_seconds
```

Aggregate:

```text
mean absolute log residual
max absolute log residual
record count
skipped count
```

**TDD cases:**

- aggregate over `288Mc`, `290Mc` has count `2`
- skipped zero-Q records are excluded from residual stats
- NaN residuals reject or skip deterministically

### Task 4.2: Add residual table screenshot

**Files:**

- Modify: `cmd/statera-ui/main.go`
- Modify: `cmd/statera-ui/screenshot_test.go`
- Create: `screenshots/statera-alpha-residuals-YYYY-MM-DD.png`

**UI requirements:**

- isotope ID
- evaluated half-life
- predicted half-life
- log residual
- factor error
- model name and evidence class

**Boundary text:**

```text
Royer output is peer-reviewed-model output and never replaces evaluated half-lives.
```

---

## Milestone 5 — Q-alpha and Mass-Energy Worksheet v1

**Purpose:** Prepare for AME2020-backed mass-energy calculations without adding unsupported data.

### Task 5.1: Add schema for mass table records

**Files:**

- Create: `data/mass_records.example.json`
- Create: `internal/research/mass_records.go`
- Create: `internal/research/mass_records_test.go`

**Fields:**

- isotope ID
- Z
- A
- mass excess MeV
- mass excess uncertainty MeV
- source URL
- DOI
- evidence level
- source table/line note

**TDD cases:**

- missing uncertainty rejects
- missing DOI rejects
- malformed isotope ID rejects
- `Z/A` mismatch rejects

### Task 5.2: Add Q-alpha calculation function

**Files:**

- Create: `internal/physics/qalpha.go`
- Create: `internal/physics/qalpha_test.go`

**Calculation:**

Use mass excess values only when parent, daughter, and alpha particle records are all available and source-backed.

**TDD cases:**

- rejects missing daughter mass
- rejects missing alpha mass
- propagates uncertainty by quadrature
- preserves source labels in result

---

## Milestone 6 — Synthesis Yield and Excitation-Function Explorer v1

**Purpose:** Model event-count calculations from source-backed reaction/cross-section data.

### Task 6.1: Add reaction campaign schema

**Files:**

- Create: `data/reaction_campaigns.example.json`
- Create: `internal/research/reaction_campaign.go`
- Create: `internal/research/reaction_campaign_test.go`

**Fields:**

- projectile
- target
- compound nucleus
- product isotope
- beam energy MeV
- excitation energy MeV
- cross section pb
- cross section uncertainty pb
- observed chains
- source URL
- DOI

**TDD cases:**

- negative cross section rejects
- missing uncertainty rejects
- missing DOI rejects
- unsupported reaction labels reject

### Task 6.2: Add expected event calculator

**Files:**

- Create: `internal/physics/yield.go`
- Create: `internal/physics/yield_test.go`

**Inputs:**

- beam particles
- target areal density
- cross section
- detection efficiency

**Output:**

- expected events
- uncertainty interval

**Boundary:**

This is a planning/education calculator, not an experimental claim.

---

## Milestone 7 — ML-Assisted Research Tools v1

**Purpose:** Use ML-like classification and duplicate detection for research operations without letting ML generate scientific defaults.

### Task 7.1: Add deterministic source triage baseline first

**Files:**

- Create: `internal/research/source_triage.go`
- Create: `internal/research/source_triage_test.go`

**Classes:**

```text
evaluated-data
experimental-paper
theoretical-model
review
context-not-for-simulation
blocked-metadata-only
unknown-needs-human-review
```

**Input:** citation markdown metadata and text snippets.

**TDD cases:**

- `nudat3ensdf` classifies as `evaluated-data`
- `royer2008alpha-analytic` classifies as `theoretical-model` or `blocked-metadata-only` depending access status
- House Oversight context record classifies as `context-not-for-simulation`

### Task 7.2: Add duplicate DOI and duplicate URL detector

**Files:**

- Modify: `internal/research/source_triage.go`
- Modify: `internal/research/source_triage_test.go`

**TDD cases:**

- exact duplicate DOI detected
- DOI normalization trims whitespace and lowercases DOI prefix
- duplicate URL detected after trimming
- duplicates do not count as independent evidence breadth

### Task 7.3: Add optional ML interface behind explicit candidate label

**Files:**

- Create: `internal/research/ml_candidates.go`
- Create: `internal/research/ml_candidates_test.go`

**Rules:**

- ML suggestions must include `CandidateOnly: true`
- ML suggestions must include `RequiresHumanReview: true`
- ML suggestions cannot be converted to simulation defaults
- any function that attempts to promote candidate-only data without source validation must reject

**No network or model runtime required for v1.** This is a safe interface and validation layer before any ML model is integrated.

---

## Milestone 8 — Screenshot Gallery and Demo Docs

**Purpose:** Make research progress visible and auditable for demos.

### Task 8.1: Add screenshot manifest

**Files:**

- Create: `screenshots/manifest.json`
- Create: `screenshots/README.md` if it needs update
- Create: `cmd/statera-ui/screenshot_manifest_test.go`

**Fields per screenshot:**

- path
- width px
- height px
- byte size minimum
- view name
- caption for blind users
- source command
- commit hash produced from

**TDD cases:**

- every manifest screenshot exists
- every PNG decodes
- each PNG is `1180 x 760 px`
- each PNG is non-blank
- every screenshot has a text caption

### Task 8.2: Generate four research demo screenshots

**Files:**

- Create: `screenshots/statera-education-YYYY-MM-DD.png`
- Create: `screenshots/statera-research-YYYY-MM-DD.png`
- Create: `screenshots/statera-design-YYYY-MM-DD.png`
- Create: `screenshots/statera-alpha-systematics-YYYY-MM-DD.png`
- Create or modify: `docs/demo-screenshots.md`

**Commands:**

```bash
CGO_ENABLED=0 go run ./cmd/statera-ui -screenshot screenshots/statera-education-YYYY-MM-DD.png
```

If per-view screenshot selection does not exist yet, first add `-view education|research|design|alpha` using TDD.

---

## Milestone 9 — Documentation and Research Status Page

**Purpose:** Keep planning visible in the repo and prevent future scope drift.

### Task 9.1: Add research workbench status document

**Files:**

- Create: `docs/research-workbench-status.md`

**Content sections:**

- implemented
- verified test counts
- screenshots
- source-backed datasets
- blocked source access
- not-for-simulation context
- next accepted slices

### Task 9.2: Link roadmap files together

**Files:**

- Modify: `README.md`
- Modify: `docs/literature-roadmap.md`
- Modify: `docs/academic-calculation-visualization-roadmap.md`

**Rules:**

- Do not claim planned features are implemented.
- Use `planned`, `implemented`, `blocked`, and `source-needed` status labels.

---

## Required Verification Before Every Code Commit

Run from repo root:

```bash
CGO_ENABLED=0 go test ./...
go run ./cmd/statera
CGO_ENABLED=0 go test ./cmd/statera-ui
git diff --check
CGO_ENABLED=0 go test -json ./... > /tmp/moscovium-go-test.json
```

Extract test counts:

```bash
python3 - <<'PY'
import json
from collections import Counter
path='/tmp/moscovium-go-test.json'
counts=Counter(); packages=set(); package_results={}
for line in open(path):
    if not line.strip():
        continue
    ev=json.loads(line)
    pkg=ev.get('Package')
    action=ev.get('Action')
    test=ev.get('Test')
    if pkg:
        packages.add(pkg)
    if action in ('pass','fail','skip') and test:
        counts[action]+=1
    elif action in ('pass','fail','skip') and pkg and not test:
        package_results[pkg]=action
print('packages_total='+str(len(packages)))
for k in ('pass','fail','skip'):
    print(f'tests_{k}={counts[k]}')
print('tests_total='+str(sum(counts.values())))
print('package_results='+','.join(f'{p}:{package_results.get(p,"none")}' for p in sorted(packages)))
PY
```

For screenshot-producing tasks, also verify PNG metadata using Go stdlib or an existing test, not visual inspection alone:

```text
width_px = 1180
height_px = 760
single_color = false
bytes > 10240
```

---

## Commit Strategy

Commit each independent slice separately.

Examples:

```bash
git commit -m "Add source-backed isotope workbook records"
git commit -m "Add decay-chain Monte Carlo summaries"
git commit -m "Render provenance graph summary"
git commit -m "Add source triage duplicate DOI detector"
git commit -m "Document research screenshot gallery"
```

Push after every verified slice:

```bash
git push origin main
git rev-parse HEAD
git rev-parse origin/main
```

The local HEAD must equal `origin/main` before reporting completion.

---

## Risk Register

| Risk | Impact | Mitigation |
| --- | --- | --- |
| APS HTTP `403` blocks Royer/Wang source details | Cannot validate coefficients or add second formula | Keep blockers documented; do not infer formulas from metadata |
| Only 2 seed isotope records currently exist | ML and residual statistics are weak | Prioritize evaluated isotope workbook expansion |
| UI screenshots can become visual-only evidence | Bad for Juan and weak auditability | Every screenshot needs numeric metadata and text caption |
| ML suggestions may be mistaken for facts | Scientific integrity failure | Candidate-only flags, requires-human-review flags, no default promotion path |
| Context/lore records can contaminate simulation | Unsupported claims could leak into engine | Track B stays `not-for-simulation`, tested in validators |

---

## Recommended Execution Order

1. Milestone 1: Evaluated Isotope Workbook v1
2. Milestone 3: Evidence and Provenance Graph v1
3. Milestone 8: Screenshot Gallery and Demo Docs
4. Milestone 2: Decay-Chain Simulation v1
5. Milestone 4: Alpha Systematics Residual Lab v2
6. Milestone 7: ML-Assisted Research Tools v1
7. Milestone 5: Q-alpha and Mass-Energy Worksheet v1
8. Milestone 6: Synthesis Yield and Excitation-Function Explorer v1

Reason: this order improves source integrity and demos first, then expands simulations, then adds safe ML and more advanced physics worksheets.

---

## Definition of Done for the Full Roadmap

The roadmap is complete when the repo has:

- source-backed isotope workbook records and validation
- decay-chain deterministic and Monte Carlo summaries
- provenance graph data and UI summary
- alpha residual dashboard with model/evidence labels
- screenshot gallery with text captions and PNG tests
- deterministic source triage and duplicate detection
- safe ML candidate interfaces that cannot seed simulation defaults
- required gates passing with `0` failures
- all implemented claims linked to source files, DOI trails, or reproducible artifacts
