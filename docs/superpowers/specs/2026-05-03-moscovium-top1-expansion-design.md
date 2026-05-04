# Moscovium Statera Go — Top 1 Expansion Design

**Date**: 2026-05-03
**Status**: Approved — proceed to writing-plans for implementation

## Overview

Expand moscovium-statera-go from bootstrap infrastructure (2 seed isotopes, Royer model, basic UI) into a comprehensive Element 115 research, simulation, and visualization platform. Execute in parallel waves.

## Design Principles

1. **Provenance-first**: Every new isotope, coefficient, and model output must carry source provenance. Missing provenance is a blocking error.
2. **Pure Go, CGO-free**: All code must compile with `CGO_ENABLED=0 go build ./...`
3. **Test-driven**: New functionality requires tests before or alongside implementation
4. **Model labels**: Simulation outputs are labeled `peer-reviewed-model` or `evaluated` — never ambiguous
5. **Existing architecture preserved**: Expand `internal/physics/`, `internal/research/`, `internal/ui/` — no breaking renames

## Wave Structure

### Wave 1 — Data Breadth
**Goal**: Expand from 2 to ~15 isotope records with full provenance

- [ ] Intake ENSDF/NuDat data for 287Mc, 288Mc, 289Mc, 290Mc, 291Mc (all known Mc isotopes)
- [ ] Add evaluated half-lives, Q-alpha, decay modes for all Mc isotopes from ENSDF
- [ ] Add daughter chain isotopes: 284Nh, 286Nh, 280Rg, 276Mt, 272Bh, 268Db with evaluated data
- [ ] Add Spontaneous Fission branching ratios where known
- [ ] Add beta-decay and electron capture data where applicable
- [ ] Populate citations/facts.md with cross-source verified facts
- [ ] Expand `data/research.seed.json` with source-reviewed records
- **Files**: `data/research.seed.json`, `internal/research/seed.go`, `internal/research/ensdf_candidate.go`
- **Verification**: `CGO_ENABLED=0 go test ./...` + provenance validation for every new record

### Wave 2 — Simulation Depth
**Goal**: Add physics models beyond Royer analytic formula

- [ ] Alpha WKB barrier penetration model (Stage 1 from relativistic roadmap)
- [ ] Liquid-drop binding energy (Bethe-Weizsäcker mass formula)
- [ ] Shell-model corrections (Strutinsky or equivalent)
- [ ] Spontaneous fission half-life model (Swiatecki or equivalent)
- [ ] Beta-decay / electron capture half-life estimates
- [ ] Additional alpha systematics: VSS formula, UNIV formula, Denisov formula
- [ ] Excitation function model for synthesis cross-sections (e.g., 243Am+48Ca → 288Mc+3n)
- [ ] Q-value systematics (alpha, beta, EC)
- **Files**: New files in `internal/physics/` (e.g., `wkb.go`, `binding.go`, `fission.go`, `beta.go`, `excitation.go`)
- **Verification**: Unit tests comparing model predictions against evaluated data from Wave 1 seeds

### Wave 3 — Visualization & UI Excellence
**Goal**: Transform gogpu/ui app from card-based report into interactive scientific dashboard

- [ ] 3D nucleus hull rendering (surface deformation, proton/neutron distributions)
- [ ] Interactive N-Z chart with all superheavy isotopes plotted
- [ ] Alpha energy spectra plots (bar chart of decay energies)
- [ ] Animated decay-chain viewer (step through chain with timing)
- [ ] Binding energy per nucleon chart
- [ ] Q-value systematics dashboard
- [ ] Education module: Guided walkthroughs with interactive elements
- [ ] Research dashboard: Provenance graph as interactive node-edge diagram
- [ ] Design module: Real-time claim validator with visual feedback
- [ ] Excitation function viewer (cross-section vs energy)
- **Files**: `internal/ui/app.go`, new `internal/ui/` components, `cmd/statera-ui/main.go`
- **Framework**: gogpu/ui + ggcanvas (existing stack, no new dependencies)
- **Screenshots**: Every new view must have a headless screenshot test

### Wave 4 — Integration & Polish
- [ ] Comprehensive test coverage (>80% for new code)
- [ ] Performance profiling and optimization
- [ ] Documentation updates (README, research charter)
- [ ] Citation/papers directory cleanup and expansion
- [ ] Full end-to-end screenshot verification

## Technical Boundaries

- No CGO dependencies
- No unsupported physics claims
- All model outputs explicitly labeled `peer-reviewed-model` or `evaluated`
- Do not modify existing function signatures without backward compatibility
- Track B data stays quarantined

## Success Criteria

1. **Data**: 15+ isotope records with complete ENSDF provenance
2. **Simulation**: 5+ physics models producing validated predictions
3. **Visualization**: 8+ interactive views in gogpu/ui dashboard
4. **Tests**: All new code covered, existing tests still pass
5. **Build**: `CGO_ENABLED=0 go test ./...` clean
