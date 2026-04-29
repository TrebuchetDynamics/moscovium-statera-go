# Open-Source Research Playbook for Moscovium Statera Go

> **For Hermes:** Use this as the repository research notebook before implementing simulations, visualizations, or ML-assisted tooling. Do not copy code from external projects unless license-compatible and explicitly approved; extract architecture lessons, data-model patterns, test practices, and UI/reproducibility ideas.

**Goal:** Learn from mature open-source nuclear-data, simulation, visualization, provenance, and scientific-ML projects, then translate those lessons into small, source-backed Statera implementation slices.

**Repository:** `/home/xel/git/sages-openclaw/workspace-riju/moscovium-statera-go`

**Date:** 2026-04-29 CST

**Research method used:**

- GitHub public repository search API, first pass only.
- Direct raw README fetches from public GitHub repositories.
- GitHub Search API second pass was blocked by HTTP `403` rate limit. Treat this as a research limitation; repeat later with authenticated API or lower frequency.

**Host preflight at research time:**

- `git`: `/usr/bin/git`
- `go`: `/home/xel/.local/go-current/bin/go`
- `curl`: `/usr/bin/curl`
- `python3`: `/home/xel/.hermes/hermes-agent/venv/bin/python3`
- `rg`: missing; use `grep` or Hermes `search_files`

---

## Current Statera Baseline

The project is in Phase 1 bootstrap. Current source-backed seed data has only two Moscovium isotope records:

- `288Mc`
- `290Mc`

Therefore:

- deterministic validation and UI screenshots are ready now;
- real statistical ML on nuclear properties is not ready yet;
- ML should first target source triage, duplicate detection, metadata extraction candidates, and provenance hygiene;
- scientific prediction or surrogate modeling should wait until the dataset expands to at least tens of source-backed nuclide records.

---

## External Projects Reviewed

### 1. Geant4 — particle transport simulation architecture

**Repository:** `Geant4/geant4`

**URL:** `https://github.com/Geant4/geant4`

**Observed metadata from GitHub search:**

- Stars: `793`
- Forks: `373`
- Language: `C++`
- Description: `Geant4 toolkit for the simulation of the passage of particles through matter - NIM A 506 (2003) 250-303`

**README evidence fetched:** Geant4 is a toolkit for simulation of passage of particles through matter. Application domains include high-energy, nuclear, accelerator, medical, and space science. README cites three reference papers, including NIM A 506 (2003) 250-303.

**Lessons for Statera:**

1. Separate physics process definitions from the application shell.
2. Keep examples and reference papers visible near the implemented model.
3. Treat model domain as explicit; Geant4 is broad, but Statera must stay narrow.
4. Every simulation path should state its assumptions and references.

**Statera translation:**

- Add `internal/physics/simulation_metadata.go` with model name, evidence class, assumptions, and source references for each simulation type.
- Require UI simulation cards to display evidence class and source reference.
- Do not attempt a Geant4-like transport engine; Statera’s scope is nuclide-data workbench, not general particle transport.

---

### 2. OpenMC — Monte Carlo reproducibility and scientific-code packaging

**Repository:** `openmc-dev/openmc`

**URL:** `https://github.com/openmc-dev/openmc`

**Observed metadata from GitHub search:**

- Stars: `1014`
- Forks: `631`
- Language: `Python`
- Description: `OpenMC Monte Carlo Code`

**README evidence fetched:** OpenMC aims to provide a fully featured Monte Carlo particle transport code based on modern methods.

**Lessons for Statera:**

1. Monte Carlo research code needs deterministic seeds, reproducible outputs, and clear input decks.
2. Scientific CLI tools should emit machine-readable artifacts, not only console text.
3. Tests should verify statistical summary invariants without depending on exact random sequences unless the seed is fixed.

**Statera translation:**

- Decay-chain Monte Carlo must accept explicit RNG seed.
- Simulation summaries should be exportable as JSON.
- Tests should check `p05 <= median <= p95`, positive times, deterministic output with fixed seed, and rejection of invalid sample counts.

---

### 3. PyNE — nuclear engineering toolkit scope and data access

**Repository:** `pyne/pyne`

**URL:** `https://github.com/pyne/pyne`

**Observed metadata from GitHub search:**

- Stars: `313`
- Forks: `199`
- Language: `C++`
- Description: `PyNE: The Nuclear Engineering Toolkit`

**README status:** raw README was not fetched during this pass.

**Lessons for Statera:**

1. Nuclear toolkits grow around data access, material/isotope identity, and validated transformations.
2. Broad nuclear engineering scope can become large quickly; Statera should remain focused on superheavy nuclide research and provenance.
3. Caution: do not import PyNE’s broad scope into Statera.

**Statera translation:**

- Keep isotope identity helpers small and tested.
- Add workbook-style data access before adding advanced physics.
- Preserve the narrow research boundary in `docs/research-charter.md`.

---

### 4. paceENSDF — ENSDF parsing, JSON translation, and visualization

**Repository:** `AaronMHurst/pace_ensdf`

**URL:** `https://github.com/AaronMHurst/pace_ensdf`

**Observed metadata from GitHub search:**

- Stars: `10`
- Forks: `0`
- Language: `Python`
- Description: Python archive/toolkit for coincident emissions from ENSDF.

**README evidence fetched:** The project enables access, manipulation, analysis, and visualization of radioactive-decay data from ENSDF. README reports `3254` datasets covering alpha, beta-minus, and electron-capture/beta-plus decays, translated into JSON, with `92,264` deexcitation gamma rays and `41,094` levels.

**Lessons for Statera:**

1. ENSDF-to-JSON translation is a proven pattern for making evaluated data usable.
2. Bundled, versioned data artifacts improve reproducibility.
3. Visualizations should be backed by explicit translated records.

**Statera translation:**

- Add an `internal/research/ensdf_candidate.go` schema before importing more ENSDF data.
- Store raw intake candidates separate from simulation defaults.
- Add deterministic JSON validation tests for any ENSDF-derived records.

---

### 5. Nuclei — ENSDF parser, viewer, editor, full-data authenticity

**Repository:** `martukas/nuclei`

**URL:** `https://github.com/martukas/nuclei`

**Observed metadata from GitHub search:**

- Stars: `26`
- Forks: `7`
- Language: `C++`
- Description: ENSDF parser, viewer, and editor.

**README evidence fetched:** Nuclei aims at complete and authentic parsing of ENSDF files, explicitly noting that many projects neglect data deemed irrelevant, such as half-life uncertainties. It exports decay schemes as PDF or SVG and emphasizes comprehensive UI display.

**Lessons for Statera:**

1. Do not drop uncertainty fields because they are inconvenient.
2. Data authenticity matters more than quick plotting.
3. SVG/PDF export concepts are useful for future reports, but PNG screenshots are enough for the current phase.

**Statera translation:**

- Every isotope workbook record should preserve uncertainty when available.
- If a field is absent, represent absence explicitly; do not default it silently.
- UI cards should display whether uncertainty is present, absent, or model-derived.

---

### 6. nuclear-chart-plotter — publication-ready nuclear chart plots

**Repository:** `jonas-ka/nuclear-chart-plotter`

**URL:** `https://github.com/jonas-ka/nuclear-chart-plotter`

**Observed metadata from GitHub search:**

- Stars: `13`
- Forks: `7`
- Language: `HTML` per GitHub search result, README describes Python plotting workflow.

**README evidence fetched:** Project was updated to AME20 on `2024-10-03`, cites W.J. Huang, Chinese Physics C45, 030002, March 2021, and is for publication-ready nuclear chart plots based on Atomic Mass Evaluation 2020.

**Lessons for Statera:**

1. AME2020-backed nuclear charts are a direct fit for Statera’s Q-alpha/mass worksheet roadmap.
2. Plotting code should cite the data release.
3. Nuclear chart visualizations should be publication-aware but still data-labeled.

**Statera translation:**

- Build a simple `N-Z` table view before graphical heatmaps.
- When heatmaps arrive, each axis and cell must state source table and evidence class.
- Add screenshot captions with source and units.

---

### 7. NuChart — nuclear chart visualization workflow

**Repository:** `antoinebelley/NuChart`

**URL:** `https://github.com/antoinebelley/NuChart`

**Observed metadata from GitHub search:**

- Stars: `6`
- Forks: `1`
- Language: `Python`
- Description: package to automatically generate nuclear chart plots.

**README evidence fetched:** NuChart generates nuclear chart plots and references ab initio reach plots. README says some data lack references and asks contributors to help fill references.

**Lessons for Statera:**

1. Visualization datasets often have missing references; Statera must reject or quarantine those for simulation.
2. Contribution workflows should make missing references visible.
3. A chart can be useful even while some regions are marked source-needed.

**Statera translation:**

- Add `source_needed` status in UI rather than hiding missing data.
- Provenance graph should include missing-reference warnings.
- Screenshot captions should identify source-needed regions.

---

### 8. DeepXDE — scientific ML / physics-informed learning

**Repository:** `lululxvi/deepxde`

**URL:** `https://github.com/lululxvi/deepxde`

**Observed metadata from GitHub search:**

- Stars: `4119`
- Forks: `959`
- Language: `Python`
- Description: scientific machine learning and physics-informed learning.

**README evidence fetched:** DeepXDE is a library for scientific machine learning and physics-informed learning.

**Lessons for Statera:**

1. Scientific ML frameworks separate model definition, training data, losses, and validation.
2. ML outputs require validation datasets and explicit error metrics.
3. Physics-informed learning is not a shortcut around source-backed data.

**Statera translation:**

- Do not build physics ML until dataset is large enough.
- Start with ML-candidate schemas and deterministic source triage.
- Later, if modeling residuals, require train/test split, error metrics, and model cards.

---

### 9. ModelingToolkit.jl — symbolic/scientific model composition

**Repository:** `SciML/ModelingToolkit.jl`

**URL:** `https://github.com/SciML/ModelingToolkit.jl`

**Observed metadata from GitHub search:**

- Stars: `1627`
- Forks: `252`
- Language: `Julia`
- Description: acausal modeling framework for automatically parallelized scientific machine learning and symbolics.

**README evidence fetched:** README identifies ModelingToolkit.jl as part of SciML with documentation, CI, and code style badges.

**Lessons for Statera:**

1. Model composition benefits from explicit symbolic variables and units.
2. Documentation and CI visibility are part of scientific trust.
3. Model interfaces should be inspectable, not opaque.

**Statera translation:**

- For each physics model, define explicit input/output structs with units in field names or comments.
- Add model metadata: inputs, outputs, evidence class, reference, assumptions.
- Avoid opaque function signatures like `Predict(float64, float64)` without named fields.

---

### 10. Nemo — provenance graph debugging in Go

**Repository:** `numbleroot/nemo`

**URL:** `https://github.com/numbleroot/nemo`

**Observed metadata from GitHub search:**

- Stars: `19`
- Forks: `3`
- Language: `Go`
- Description: analyzes provenance graphs from distributed-system fault injection.

**README evidence fetched:** Nemo debugs distributed systems by analyzing provenance graphs from fault injection and is implemented with Go and Docker tooling.

**Lessons for Statera:**

1. Provenance graphs are useful beyond data science; they can explain why a result exists.
2. Go is suitable for provenance graph tooling.
3. Graph analysis should be queryable and testable, not just drawn.

**Statera translation:**

- Build provenance graph as a pure Go data model first.
- Add tests for node counts, edge counts, orphan detection, and blocked-source labeling.
- Render text summaries before graphical graph layouts.

---

## GitHub Search Results Captured

First-pass GitHub repository search returned these relevant examples:

| Query | Repository | Stars | Forks | Language | Lesson |
| --- | --- | ---: | ---: | --- | --- |
| `geant4` | `Geant4/geant4` | 793 | 373 | C++ | simulation scope, reference-paper discipline |
| `pyne nuclear` | `pyne/pyne` | 313 | 199 | C++ | nuclear toolkit data/model boundaries |
| `openmc` | `openmc-dev/openmc` | 1014 | 631 | Python | Monte Carlo reproducibility |
| `ENSDF` | `martukas/nuclei` | 26 | 7 | C++ | comprehensive ENSDF parsing and UI |
| `ENSDF` | `AaronMHurst/pace_ensdf` | 10 | 0 | Python | ENSDF to JSON and analysis/visualization toolkit |
| `nuclear chart` | `jonas-ka/nuclear-chart-plotter` | 13 | 7 | HTML/Python docs | AME2020-backed publication plots |
| `nuclear chart` | `antoinebelley/NuChart` | 6 | 1 | Python | nuclear chart generation and missing-reference visibility |
| `provenance graph` | `numbleroot/nemo` | 19 | 3 | Go | Go provenance graph analysis |
| `physics machine learning` | `lululxvi/deepxde` | 4119 | 959 | Python | scientific ML discipline and validation |
| `physics machine learning` | `SciML/ModelingToolkit.jl` | 1627 | 252 | Julia | symbolic scientific model composition |

Second-pass GitHub repository search was blocked:

```text
HTTP Error 403: rate limit exceeded
```

Repeat later with authentication or lower frequency for additional projects such as `radioactivedecay`, `pynucastro`, AME2020 utilities, NUBASE utilities, and isotope visualization packages.

---

## How To Research Each External Project Before Borrowing Ideas

For each external repository, perform this checklist:

1. **Identity check**
   - repo full name
   - license
   - latest commit activity
   - primary language
   - documentation status

2. **Scope check**
   - what problem it solves
   - what it explicitly does not solve
   - whether scope overlaps with Statera

3. **Data model check**
   - how it names nuclides/isotopes
   - how it stores uncertainty
   - how it links sources
   - whether it supports evaluated vs model labels

4. **Validation check**
   - unit tests
   - golden files
   - CI
   - examples with expected outputs

5. **Visualization check**
   - chart types
   - screenshot/export support
   - accessibility/text alternatives
   - source labels in plots

6. **ML check, if applicable**
   - training data source
   - test split
   - metrics
   - uncertainty/error reporting
   - whether outputs are safely labeled as predictions

7. **Statera adoption decision**
   - adopt pattern now
   - defer until dataset expands
   - reject as out-of-scope
   - blocked by license/source access

Record results in `docs/research-workbench-status.md` or a focused `docs/superpowers/plans/YYYY-MM-DD-*.md` file.

---

## Architecture Patterns To Adopt

### Pattern A — Data-first, model-second

From ENSDF/Nuclei/paceENSDF.

Statera implementation:

1. validated raw source record
2. candidate record
3. accepted evaluated record
4. simulation input record
5. UI rendering record

Never skip directly from paper text or ML extraction to simulation default.

### Pattern B — Model metadata is mandatory

From Geant4, DeepXDE, ModelingToolkit.jl.

Every Statera model should carry:

```go
type ModelMetadata struct {
    Name string
    EvidenceClass string
    Reference string
    DOI string
    Assumptions []string
    Inputs []string
    Outputs []string
    Limitations []string
}
```

### Pattern C — Graph facts before graph drawings

From Nemo and provenance graph projects.

Build first:

```go
type ProvenanceNode struct { ... }
type ProvenanceEdge struct { ... }
```

Then render:

- node counts
- edge counts
- orphan warnings
- blocked-source list
- eventually a graph view

### Pattern D — Screenshot artifacts are tests, not decorations

From scientific visualization lessons.

Each screenshot should have:

- source command
- dimensions
- byte size
- non-blank check
- text caption
- commit hash
- evidence class notes

### Pattern E — ML starts as triage, not physics

From DeepXDE/SciML plus current Statera data limits.

Safe ML order:

1. duplicate DOI/URL detection
2. source type classifier
3. extraction candidate schema
4. human/test validation
5. only later residual modeling

---

## Concrete Implementation Additions To Existing Roadmap

Add these tasks to the previously committed roadmap `docs/superpowers/plans/2026-04-29-research-simulation-visualization-ml-roadmap.md` when execution begins.

### Add to Milestone 1 — Workbook

- Include explicit `UncertaintyStatus` field:
  - `present`
  - `absent-in-source`
  - `not-yet-intaken`
  - `model-derived`

Reason: Nuclei warns that uncertainty is often dropped; Statera must not drop it silently.

### Add to Milestone 2 — Decay Simulation

- Add `SimulationMetadata` to every Monte Carlo result.
- Include RNG seed in JSON output.

Reason: OpenMC-style reproducibility.

### Add to Milestone 3 — Provenance Graph

- Add orphan-source and orphan-DOI tests.
- Add blocked-source nodes for Royer and Wang access blockers.

Reason: Nemo/provenance pattern.

### Add to Milestone 5 — Q-alpha/Mass Worksheet

- Start with text/table view before heatmap.
- Require AME2020 source-table references.

Reason: nuclear-chart-plotter and NuChart rely on AME-backed plotting.

### Add to Milestone 7 — ML Research Tools

- Add `CandidateOnly` and `RequiresHumanReview` flags before any ML extractor.
- Add tests proving candidate data cannot enter `ResearchSeed.Catalog`.

Reason: prevent ML output from becoming scientific default.

### Add to Milestone 8 — Screenshot Gallery

- Add `screenshots/manifest.json` and tests that every screenshot has a caption.

Reason: Juan is blind and screenshots must have text equivalents.

---

## Recommended Open-Source Deep Dives, In Order

### Deep Dive 1: paceENSDF and Nuclei

**Why first:** closest to evaluated decay data and ENSDF workflows.

**Questions to answer:**

- How do they represent decay datasets?
- How do they preserve uncertainty?
- How do they handle missing fields?
- What can Statera mimic in a much smaller Go schema?

**Deliverable:** `docs/research/ensdf-intake-design.md`

### Deep Dive 2: nuclear-chart-plotter and NuChart

**Why second:** informs screenshot-friendly isotope workbook and N-Z views.

**Questions to answer:**

- How are AME/NUBASE data shaped for plotting?
- What labels and legends are necessary?
- How can Statera provide text alternatives?

**Deliverable:** `docs/research/nuclear-chart-visualization-design.md`

### Deep Dive 3: OpenMC

**Why third:** informs Monte Carlo reproducibility, fixed seeds, input decks, and outputs.

**Questions to answer:**

- How does it document stochastic runs?
- How are input files and outputs separated?
- What test patterns are useful for stochastic code?

**Deliverable:** `docs/research/decay-monte-carlo-design.md`

### Deep Dive 4: Nemo / provenance graph projects

**Why fourth:** Statera’s differentiator is evidence integrity.

**Questions to answer:**

- How should graph nodes/edges be represented in Go?
- How do we test graph consistency?
- What graph queries matter before visualization?

**Deliverable:** `docs/research/provenance-graph-design.md`

### Deep Dive 5: DeepXDE / SciML ModelingToolkit

**Why fifth:** ML is useful only after data/provenance foundations exist.

**Questions to answer:**

- How are models documented?
- How are metrics reported?
- What minimal model card does Statera need?

**Deliverable:** `docs/research/ml-candidate-safety-design.md`

---

## Next Research Task To Execute

Start with:

```text
Deep Dive 1: paceENSDF and Nuclei
```

Reason:

- closest to current data expansion need;
- directly informs Evaluated Isotope Workbook v1;
- helps avoid dropping uncertainty or source details;
- does not require ML or blocked APS article access.

Acceptance criteria for Deep Dive 1:

- create `docs/research/ensdf-intake-design.md`
- cite exact external repos reviewed
- describe minimum Statera ENSDF candidate schema
- define source-to-candidate-to-accepted-data workflow
- add at least 5 validation rules
- no scientific values added
- `CGO_ENABLED=0 go test ./...` passes

---

## Blockers and Limits

1. GitHub search API second pass hit HTTP `403` rate limit.
   - Impact: incomplete broad discovery.
   - Workaround: direct README fetches for known projects; retry later authenticated.

2. Royer/Wang paper content remains access-blocked by APS HTTP `403` from prior audit.
   - Impact: no coefficient cross-check or second alpha model intake.
   - Workaround: focus on source-neutral architecture, provenance graph, and evaluated-data intake.

3. Current dataset has only `2` seed isotope records.
   - Impact: no meaningful statistical ML on nuclear properties.
   - Workaround: ML starts with source triage and duplicate detection only.

---

## Bottom Line

Open-source projects confirm the direction:

- use ENSDF/AME/NUBASE-style data-first architecture;
- preserve uncertainty and source identity;
- make simulation reproducible with explicit inputs and seeds;
- build provenance graph facts before visual graph drawings;
- treat screenshots as tested artifacts with text captions;
- keep ML in candidate/triage mode until the dataset is much larger.

The next concrete move is not more ML. It is the evaluated isotope workbook plus ENSDF intake design, because every later simulation, visualization, and ML feature depends on source-backed, validated records.
