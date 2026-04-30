# Numerical Methods Registry

Date: 2026-04-29 19:50 CST -0600
Reviewer: Riju
Scope: numerical and computational methods needed for Moscovium Statera Go

## Purpose

This registry keeps numerical methods in research until each method has sources, implementation tests, units, input/output definitions, and validation targets. It prevents a toy solver from being mistaken for finished science.

## Method readiness levels

| Level | Meaning | Runtime status |
| --- | --- | --- |
| M0 | Name only | Do not implement as science |
| M1 | Source lead recorded | Documentation only |
| M2 | Equations and assumptions extracted | Candidate implementation allowed |
| M3 | Unit tests and deterministic examples | Toy or educational output allowed |
| M4 | Compared to accepted facts with residuals | Model-validation output allowed |
| M5 | Reviewed against multiple sources/data sets | Candidate production science module |

## Registry

| Method | Domain | Current level | Inputs | Outputs | Validation target | Evidence boundary |
| --- | --- | ---: | --- | --- | --- | --- |
| Decay constant from half-life | Nuclear decay | M4 | accepted half-life in seconds | decay constant `s^-1`, mean life in seconds | exact formula check and accepted seed anchors | implemented deterministic calculation, not new experimental data |
| Exponential Monte Carlo decay sampling | Nuclear decay | M4 | half-life seconds, sample count, RNG seed | reproducible quantiles in seconds | fixed-seed tests and seed anchors | simulation summary only |
| Alpha residual metric | Nuclear decay | M3 | predicted half-life, evaluated half-life | `log10(predicted/evaluated)`, factor error | accepted half-life anchors | model comparison only |
| WKB alpha barrier penetration | Nuclear decay | M1 | Q-alpha, daughter charge, radius/barrier model | tunneling probability, half-life estimate | half-life anchors and source formula tests | not implemented; source review required |
| Semi-empirical alpha formulas | Nuclear decay | M1 | Z, A, Q-alpha, coefficients | half-life prediction | residuals against accepted isotope facts | coefficients blocked until source content reviewed |
| Fission barrier uncertainty estimation | Nuclear structure | M1 | deformation coordinates, functional/model parameters | barrier height/uncertainty | literature benchmark tables | not implemented |
| Langevin fission dynamics | Nuclear structure | M1 | potential energy surface, inertia, friction, random force | fission probability/path statistics | fission-paper benchmarks | not implemented |
| HFB / nuclear DFT | Nuclear structure | M1 | energy density functional, pairing, deformation basis | masses, shell effects, barriers | AME/NUBASE and literature benchmark nuclei | likely external-tool design first |
| CDFT / RMF nuclear models | Nuclear structure | M1 | covariant functional parameters | nuclear landscape predictions | known superheavy benchmark papers | not implemented |
| Relativistic DFT chemistry | Electronic/chemistry | M1 | nuclear charge, basis, functional, relativistic Hamiltonian | orbital/adsorption/chemical observables | chemistry papers, not half-life | separate from nuclear validation |
| Dirac-Hartree-Fock toy solver | Electronic/education | M0 | Coulomb-like potential, grid/basis | bound-state energies/orbitals | analytic/simple benchmarks | educational only unless source-reviewed |
| Schrödinger finite-difference toy solver | Education/numerics | M0 | potential, grid spacing, boundary conditions | eigenvalues/eigenvectors | textbook benchmarks | not Moscovium physics |
| ODE IVP solver | Numerical infrastructure | M0 | derivative function, step size, tolerances | trajectory with error estimates | manufactured solutions | infrastructure only |
| Boundary-value/eigenvalue solver | Numerical infrastructure | M0 | differential operator, boundary conditions | eigenvalues/eigenfunctions | analytic wells/hydrogen-like tests | infrastructure only |
| Quadrature/integration | Numerical infrastructure | M0 | integrand, bounds, tolerance | integral with error estimate | analytic integrals | needed for WKB |
| Uncertainty propagation | Model validation | M0 | value distributions/covariance | propagated interval | synthetic and source examples | required before claiming precision |
| ML source triage | Research tooling | M3 | title/abstract/metadata/text snippets | candidate lane/review flags | deterministic tests and human review | never promotes facts alone |
| Duplicate DOI/URL detection | Research tooling | M3 | source records | duplicate warnings | deterministic string/DOI tests | metadata hygiene only |

## Minimum implementation contract

Before a method can produce project-facing output, it needs:

1. Method source record or explicit toy/educational label.
2. Input units.
3. Output units.
4. Deterministic test cases.
5. Failure-mode tests.
6. Provenance label in CLI/UI output.
7. Statement of what the method cannot validate.

## Validation domains

### Nuclear decay validation

Allowed anchors:

- Half-life seconds.
- Decay mode.
- Daughter identity.
- Q-alpha when source-reviewed.

Useful metric:

```text
residual_log10 = log10(predicted_half_life_seconds / accepted_half_life_seconds)
```

Interpretation:

- `0`: predicted half-life equals accepted half-life.
- `+1`: prediction is 10 times too long.
- `-1`: prediction is 10 times too short.

### Electronic/chemistry validation

Do not use nuclear half-life. Use sources for:

- Adsorption enthalpy.
- Volatility/gas chromatography behavior.
- Relativistic orbital ordering.
- Ionization/chemical trends where source-backed.

### Numerical solver validation

Use manufactured or textbook problems first:

- Harmonic oscillator.
- Square well.
- Hydrogen-like Coulomb potential at non-Moscovium Z for stability checks.
- Known ODE/integral solutions.

Only after numerical correctness should Moscovium-specific source-backed parameters be introduced.

## Next registry actions

1. Create paper records for each M1 method lead.
2. Extract equations and assumptions for WKB alpha decay.
3. Build a candidate WKB worksheet with explicit source labels.
4. Add model-card output for every nontrivial model.
5. Expand the registry with exact source IDs as citations are reviewed.
