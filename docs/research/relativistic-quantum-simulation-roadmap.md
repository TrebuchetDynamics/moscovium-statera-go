# Relativistic Quantum Simulation Roadmap for Moscovium Statera Go

Date: 2026-04-29
Repo: `/home/xel/git/sages-openclaw/workspace-riju/moscovium-statera-go`

## Purpose

Juan asked for “more finished science” and a path toward complex simulations of the relativistic and quantum-mechanical behavior of a Moscovium atom. This document turns that into an evidence-first modeling plan.

The short version:

- A full Moscovium atom model is not one equation.
- The electrons need relativistic quantum chemistry.
- The nucleus needs nuclear-structure and decay models.
- Half-lives are calibration/validation targets for nuclear decay models, not for ordinary electronic-structure models.
- Every model output must be labeled `theoretical/model` unless it is directly evaluated experimental data.

## Separation of physics domains

| Domain | What it models | Useful outputs | Candidate validation data | Evidence boundary |
| --- | --- | --- | --- | --- |
| Relativistic electronic atom/chemistry | Electron orbitals around `Z=115` nucleus, spin-orbit effects, adsorption chemistry | orbital energies, ionization trends, bonding/adsorption comparisons | chemistry experiments such as Mc/Nh gas chromatography and peer-reviewed relativistic calculations | Does not predict nuclear half-life |
| Nuclear structure | Protons/neutrons, shells, deformation, binding energy, fission barrier | mass/binding estimates, shell corrections, deformation, fission barrier | AME/NUBASE/ENSDF, primary superheavy papers, model papers | Theoretical extrapolation is large for superheavy nuclei |
| Alpha decay | Parent/daughter Q-alpha and barrier tunneling | half-life prediction, residual versus evaluated half-life | accepted half-lives, Q-alpha values, alpha energies | Good first “model against measured data” target |
| Synthesis/reaction | Heavy-ion fusion, evaporation channels, survival probability | expected event counts, cross-section curves | peer-reviewed accelerator experiments | Needs beam energy, target thickness, efficiency, cross-section provenance |

## Why half-life is useful, and where it is not

Half-life is a useful target for nuclear-decay models. Example:

```text
model predicts alpha-decay half-life
compare to accepted half-life from ENSDF/NuDat/NUBASE
compute residual = log10(predicted_seconds / evaluated_seconds)
```

This is already partially implemented through the alpha residual workbench.

But half-life is not a validation target for an electronic orbital model. A relativistic electron-cloud calculation can explain chemistry trends and adsorption behavior, but it does not directly predict how fast the nucleus alpha-decays.

## Current accepted Statera anchors

Current source-backed seed records in `data/research.seed.json`:

| Isotope | Z | A | N | Half-life | Q-alpha | Daughter | Status |
| --- | ---: | ---: | ---: | ---: | ---: | --- | --- |
| `288Mc` | 115 | 288 | 173 | `0.17 s` | `10.75 MeV` | `284Nh` | accepted seed record |
| `290Mc` | 115 | 290 | 175 | `0.65 s`; interval fields `0.45 s` lower and `1.14 s` upper | `10.45 MeV` | `286Nh` | accepted seed record |

These are appropriate first anchors for simple decay-constant calculations and alpha-decay residual checks. They are not enough to validate a full electronic or nuclear-density-functional model.

## Similar work and source leads found or already queued

### Already in repo citation notes

| Key | Source | Relevance | Current use allowed |
| --- | --- | --- | --- |
| `yakushev2024mcnhchemistry` | Yakushev et al., Frontiers in Chemistry 2024, DOI `10.3389/fchem.2024.1474820` | relativistic chemistry and gas chromatography of Mc/Nh | chemistry education and future relativistic-chemistry source review |
| `giuliani2019superheavy` | Reviews of Modern Physics 2019, DOI `10.1103/RevModPhys.91.011001` | broad review of superheavy nuclear physics, fission, alpha decay, relativistic chemistry | roadmap and source discovery; not direct numeric defaults |
| `erler2012nuclearLandscape` | Nature 2012, DOI `10.1038/nature11188` | nuclear density functional theory / HFB landscape | model-labeled nuclear landscape planning |
| `nudat3ensdf` | NNDC/NuDat/ENSDF | evaluated nuclear data interface | source of candidate and accepted evaluated data with parser/review gate |
| `nubase2020evaluation` | Chinese Physics C 2021, DOI `10.1088/1674-1137/abddae` | evaluated nuclide properties | future source-backed half-life/mass table ingestion |

### arXiv search evidence from 2026-04-29

Some arXiv and Semantic Scholar calls were rate-limited. Usable arXiv results included:

| arXiv ID | Year | Title | Relevance |
| --- | ---: | --- | --- |
| `2409.04620v1` | 2024 | `Weak decays in superheavy nuclei` | superheavy decay modes and model context |
| `0901.0901v2` | 2009 | `Fission Barriers of Compound Superheavy Nuclei` | fission-barrier/survival modeling for synthesis |
| `2311.12011v2` | 2023 | `Multimodal fission from self-consistent calculations` | self-consistent fission calculations |
| `1902.10108v1` | 2019 | `Extension of nuclear landscape to hyperheavy nuclei` | covariant density functional theory near/higher than superheavy region |
| `1804.06395v2` | 2018 | `Hyperheavy nuclei: existence and stability` | limits of nuclear existence; model-labeled only |
| `2211.05671v1` | 2022 | `Structure and reaction study of Z=120 isotopes using non-relativistic and relativistic mean-field formalism` | example of comparing Skyrme-Hartree-Fock and relativistic mean-field approaches |

These are source leads, not accepted Statera values. Each requires a source-review document before becoming a model reference.

## Numerical-method ladder

### Stage 0 — accepted data and residuals

Goal: make the source-backed facts trustworthy before complex modeling.

Methods:

- deterministic half-life to decay-constant conversion;
- Monte Carlo exponential decay-time sampling with fixed RNG seed;
- alpha residuals versus accepted half-lives;
- provenance graph and facts database.

Already partly implemented:

- `CGO_ENABLED=0 go run ./cmd/statera -decay-simulation-format=json -decay-simulation-samples=64 -decay-simulation-seed=20260429`
- source-backed records: `288Mc`, `290Mc` only.

### Stage 1 — alpha-decay tunneling model

Goal: build a first nuclear model that can be validated against half-lives.

Candidate physics:

```text
alpha particle tunnels through Coulomb + nuclear potential barrier
half-life depends strongly on Q-alpha and barrier penetration probability
```

Candidate numerical methods:

- WKB integral through a one-dimensional effective potential;
- quadrature with explicit units;
- uncertainty propagation from Q-alpha and radius parameters;
- residual table against accepted half-lives.

Acceptance criteria:

1. DOI/source for the formula and coefficients is read, not just metadata.
2. Test fixtures include known isotope inputs and expected outputs from the source.
3. Output stores model name, coefficient source, version, input units, and residual.
4. Model results are labeled `theoretical/model`, never evaluated data.

### Stage 2 — relativistic electronic atom toy model

Goal: introduce relativistic quantum mechanics without pretending to solve all-electron Moscovium immediately.

Candidate physics:

- Dirac equation for hydrogen-like high-`Z` ion as an educational baseline;
- finite-nucleus correction only if sourced;
- compare nonrelativistic versus relativistic energy-level scaling.

Candidate numerical methods:

- radial Dirac equation integration;
- finite-difference or shooting methods for boundary-value problems;
- eigenvalue solvers with convergence tests.

Acceptance criteria:

1. Start with a validated hydrogen-like case where analytic or reference values exist.
2. Label result as an educational one-electron approximation, not neutral Moscovium chemistry.
3. Do not use it to claim material properties, energy production, or propulsion.

### Stage 3 — relativistic many-electron chemistry interface

Goal: connect Statera to real relativistic chemistry work instead of writing an unreliable all-electron code from scratch.

Candidate approaches to research:

- four-component Dirac-Coulomb calculations;
- Dirac-Coulomb-Breit corrections;
- relativistic DFT;
- relativistic coupled cluster;
- spin-orbit-coupled pseudopotentials / effective core potentials.

Statera implementation should likely be an adapter/metadata workbench first:

```text
source paper -> method card -> input summary -> output property -> uncertainty/limit -> comparison to experiment
```

Good first validation target:

- Mc/Nh gas chromatography and adsorption behavior from `yakushev2024mcnhchemistry`, once fully source-reviewed.

### Stage 4 — nuclear density functional / HFB planning

Goal: support model-labeled nuclear landscape views and fission/survival exploration.

Candidate methods:

- Hartree-Fock-Bogoliubov nuclear DFT;
- Skyrme energy-density functionals;
- covariant/relativistic mean-field theory;
- fission-barrier calculations;
- Langevin fission dynamics.

Statera should not implement a full HFB solver first. A safer sequence is:

1. source-review an existing calculation paper;
2. encode model metadata and outputs as read-only benchmark cases;
3. build plots/residual tables;
4. only then consider numerical solvers or external-code adapters.

## Differential equations and numerical methods we should actually use

| Method | Use in Statera | First safe target |
| --- | --- | --- |
| ODE initial/boundary solvers | radial Schrödinger/Dirac toy equations | hydrogen-like reference cases |
| Eigenvalue solvers | bound-state energy levels | finite-difference educational wells before high-Z atoms |
| Quadrature | WKB alpha barrier penetration | alpha half-life model |
| Monte Carlo | decay time, uncertainty propagation, reaction yield uncertainty | existing deterministic decay simulation |
| PDE/DFT solvers | full nuclear/electronic structure | not first; use source-reviewed benchmarks first |
| Langevin/stochastic differential equations | fission dynamics | later, after source-backed fission barrier data |

## Concrete next repo slices

1. **Trusted Facts v1**
   - Convert accepted seed half-lives into a facts/anchors table.
   - Keep separate columns for evaluated value, unit, source path, citation URL, DOI trail, and status.
   - Do not promote candidate browser text without parser/review.

2. **Alpha WKB model design**
   - Find and source-review a formula paper with coefficients and testable examples.
   - Add a model card before implementation.
   - Add tests before code.

3. **Relativistic one-electron educational solver**
   - Start with hydrogen-like Dirac equation or sourced analytic formula.
   - Validate against low-Z reference cases before setting `Z=115`.
   - Label as educational baseline only.

4. **Relativistic chemistry source review**
   - Fully review Yakushev 2024 and any method papers it cites.
   - Extract only source-backed method/property facts.
   - Build a chemistry comparison card, not a material claim.

5. **Nuclear DFT/HFB source-review notebook**
   - Review `giuliani2019superheavy`, `erler2012nuclearLandscape`, and arXiv source leads.
   - Record model class, equations solved, numerical method, observables, and uncertainty handling.

## Plain-English check

What we can safely do now:

- Use known half-lives as scoring targets for nuclear decay models.
- Build careful toy solvers for quantum/relativistic equations.
- Learn from professional DFT/HFB/Dirac/chemistry papers.
- Compare model output to accepted data when both have units and sources.

What we cannot honestly do yet:

- Claim a complete Moscovium atom simulator.
- Claim a model is correct because it makes a nice picture.
- Use electron-cloud simulations to prove nuclear half-life.
- Claim energy, propulsion, or bulk material behavior.
