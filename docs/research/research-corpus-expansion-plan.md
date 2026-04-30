# Research Corpus Expansion Plan

Date: 2026-04-29 19:50 CST -0600
Reviewer: Riju
Scope: Moscovium Statera Go research backlog and evidence architecture

## Purpose

Juan's requirement is correct: Statera must keep research in research, then turn only reviewed material into facts, model inputs, tests, and simulations. This document defines how the project fills the gaps with many papers, facts, techniques, and numerical methods without accidentally treating unreviewed text as physics truth.

## Current corpus state

Measured at repo path `/home/xel/git/sages-openclaw/workspace-riju/moscovium-statera-go`:

- Existing source-record markdown files in `citations/papers/`: 15.
- Existing detailed research docs in `docs/research/`: 8 before this document.
- Current fully verified facts in `citations/facts.md`: 0.
- Current accepted seed modeling anchors in `citations/facts.md`: 2.
- Current accepted seed isotope anchors: `288Mc` half-life `0.17 s`, `290Mc` half-life `0.65 s`.

These counts mean the project has a scaffold and a small trusted seed, but it does not yet have the large research library Juan wants.

## Evidence lanes

Every source must enter one of these lanes before code uses it.

| Lane | Source examples | Allowed use | Not allowed |
| --- | --- | --- | --- |
| Evaluated nuclear data | ENSDF, NuDat, NUBASE, AME | Candidate-to-accepted facts after exact extraction | Guessing missing values |
| Experimental paper | Oganessian synthesis papers, gas chromatography papers | Production, observed decay chains, measured chemistry context | Runtime defaults until values are extracted and reviewed |
| Theoretical model paper | Royer, Wang, WKB, DFT, HFB, CDFT, Langevin | Method description, model candidate, validation target | Treating formula coefficients as accepted from metadata alone |
| Numerical methods paper/tool | Monte Carlo, ODE/BVP, quadrature, eigenvalue, optimization methods | Implementation design and tests | Scientific data claims |
| Review paper | RMP, Progress in Particle and Nuclear Physics, Reports on Progress in Physics | Map of field and bibliography expansion | Replacing primary/evaluated data |
| Context-only source | News, government letters, popular articles | Track B narrative context | Any simulation input |

## Corpus targets

Near-term target size for a useful research workbench:

- Evaluated-data sources: at least 6 records.
- Experimental synthesis/decay sources: at least 12 records.
- Relativistic chemistry/electronic-structure sources: at least 10 records.
- Nuclear structure/decay/fission model sources: at least 20 records.
- Numerical methods / computational physics sources: at least 15 records.
- Open-source scientific-tool exemplars: at least 10 records.
- Extracted candidate facts: at least 50 rows before promotion.
- Fully verified facts: promote only after exact source location, units, uncertainty, and review status are recorded.

## Topic map for filling gaps

### 1. Moscovium existence and production

Need papers and facts for:

- `243Am + 48Ca` fusion-evaporation reactions.
- Beam energy, target, evaporation channel, decay-chain assignment.
- Cross sections with uncertainty.
- Reproducibility and discovery/naming decisions.
- IUPAC/IUPAP recognition and naming.

Current source leads:

- Oganessian et al. 2022, `10.1103/PhysRevC.106.L031301`.
- Oganessian et al. 2022, `10.1103/PhysRevC.106.064306`.
- IUPAC 2016 naming paper, `10.1515/pac-2016-0501`.

### 2. Evaluated isotope facts

Need exact source-location extraction for:

- Half-life values and qualifiers.
- Alpha decay mode and branching when available.
- Daughter assignments.
- Q-alpha and alpha-particle energies.
- Uncertainties and limits.

Current anchor values are not yet counted as full verified facts:

- `288Mc` half-life `0.17 s` from accepted seed record.
- `290Mc` half-life `0.65 s` from accepted seed record.

Required evaluated-data sources:

- NNDC ENSDF/NuDat pages.
- NUBASE 2020.
- AME 2020 mass tables.

### 3. Nuclear decay and structure models

Need sources for:

- Deterministic decay constants: `lambda = ln(2) / T1/2`.
- Monte Carlo decay-time sampling from exponential distributions.
- WKB alpha-decay barrier penetration.
- Semi-empirical alpha half-life formulas.
- Fission barriers and spontaneous fission competition.
- Nuclear shell effects and island-of-stability predictions.
- HFB/DFT/CDFT/RMF nuclear structure calculations.

Live arXiv metadata sampled on 2026-04-29 found these leads:

- `0803.4151v2`, "Superheavy Elements in the Magic Islands".
- `1111.0505v1`, "New approach for alpha decay half-lives of superheavy nuclei and applicability of WKB approximation".
- `0811.1619v1`, "Alpha decay chains from superheavy nuclei".
- `1810.04421v1`, "Study on alpha decay chains of Z = 122 superheavy nuclei with deformation effects and Langer modification".
- `1704.06334v1`, "Assessing theoretical uncertainties in fission barriers of superheavy nuclei".

These are method leads only until full source review.

### 4. Relativistic atom and chemistry models

Need sources for:

- Dirac equation and relativistic quantum mechanics basics.
- Dirac-Hartree-Fock and relativistic coupled cluster.
- Relativistic DFT for superheavy elements.
- Spin-orbit splitting and 7p/8s orbital effects.
- Moscovium/nihonium chemistry, adsorption, volatility, gas chromatography.

Important boundary:

- Nuclear half-life can validate nuclear decay models.
- Nuclear half-life does not validate electron-cloud or chemistry models.

Current source lead:

- Yakushev et al. 2024, `10.3389/fchem.2024.1474820`.

### 5. Numerical and computational methods

Need sources and tests for:

- Monte Carlo reproducibility and random-seed discipline.
- Quadrature for WKB action integrals.
- ODE initial-value and boundary-value solvers.
- Eigenvalue solvers for Schrödinger/Dirac educational toy problems.
- Finite difference / finite element grids.
- Uncertainty propagation and residual metrics.
- Optimization / parameter fitting with held-out validation.
- Model cards and provenance for ML-assisted extraction.

These methods are allowed before all physics facts are complete if outputs are clearly labeled toy, candidate, or simulation-only.

## Source-to-fact pipeline

1. Discover source.
2. Record bibliographic identity under `citations/papers/` or queue.
3. Classify source lane.
4. Extract candidate facts with exact source location.
5. Preserve units, uncertainty, qualifiers, and caveats.
6. Store candidate facts as candidate-only.
7. Add tests for parsing and validation.
8. Promote only reviewed facts into `citations/facts.md` and runtime seed files.
9. Record residual/model validation outputs separately from experimental data.

## Candidate fact row schema

Use this schema for future extraction tables:

| Field | Required | Example |
| --- | --- | --- |
| nuclide_id | yes | `288Mc` |
| datum_type | yes | `half_life` |
| value | yes unless qualitative | `0.17` |
| unit | yes when numeric | `s` |
| uncertainty | yes when available | `+0.14/-0.05 s` |
| qualifier | yes when available | `approximate`, `limit`, `adopted` |
| source_id | yes | `nudat3ensdf` |
| source_location | yes | table, line, section, page, URL fragment |
| evidence_lane | yes | `evaluated-data` |
| review_status | yes | `candidate`, `accepted`, `blocked` |
| allowed_use | yes | `model-validation-anchor` |

## How this supports building Statera

- More papers create the corpus.
- Candidate facts create extraction pressure without contaminating physics code.
- Accepted facts feed deterministic calculators and tests.
- Techniques become implementation tasks.
- Models get residuals against accepted anchors.
- Screenshots and CLI exports show progress without pretending the science is finished.

## Immediate next slices

1. Create a candidate fact table for NNDC/NuDat `288Mc` and `290Mc` half-lives without promoting new facts.
2. Add source records for the arXiv WKB/fission/HFB leads.
3. Add a numerical-methods registry that maps method to equations, inputs, outputs, validation target, and source status.
4. Build an alpha-decay WKB worksheet as a toy/candidate model, then compare residuals to accepted half-life anchors.
5. Expand relativistic chemistry bibliography, especially Dirac-Hartree-Fock/relativistic DFT sources for group 15 superheavy elements.

## Safety rule

A paper in the repo is not automatically a fact. A method in the repo is not automatically a validated model. A model output is not experimental data.
