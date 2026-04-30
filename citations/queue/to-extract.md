# Sources To Extract

Purpose: sources whose exact tables, values, equations, assumptions, or method details should be extracted into candidate records. Extraction remains candidate-only until reviewed.

## Candidate facts for evaluated isotope records

| Source | Candidate extraction targets | Required fields |
| --- | --- | --- |
| NNDC/NuDat/ENSDF `288Mc` | half-life, decay mode, daughter, alpha energy, Q-alpha, qualifiers | value, unit, uncertainty, source URL, dataset line/table |
| NNDC/NuDat/ENSDF `290Mc` | half-life, decay mode, daughter, alpha energy, Q-alpha, qualifiers | value, unit, uncertainty, source URL, dataset line/table |
| NUBASE 2020 | mass excess, half-life, decay modes, evaluated qualifiers for Mc/Nh chain | value, uncertainty, qualifier, table/page |
| AME 2020 | mass values or mass excess needed for Q-alpha worksheet | value, uncertainty, table/page |

## Candidate experimental facts

| Source | Candidate extraction targets | Allowed eventual use |
| --- | --- | --- |
| Oganessian 2022 Mc factory paper | reaction, beam/target, isotope assignment, cross section, event count, decay chain | production/evidence dashboard |
| Oganessian 2022 new `286Mc` paper | reaction, isotope assignment, decay chain, half-life context | isotope workbook expansion after review |
| IUPAC 2016 naming paper | official names/symbols for 113, 115, 117, 118 | standards/provenance page |

## Candidate model equations and coefficients

| Source | Candidate extraction targets | Acceptance blocker |
| --- | --- | --- |
| Royer and Zhang 2008 | alpha half-life formula, coefficients, domain, units | full source content needed; DOI metadata alone insufficient |
| Wang et al. 2015 | alpha model/systematics formula, coefficients, training domain | full source content needed; DOI metadata alone insufficient |
| Hosseini and Hassanabadi 2017 | comparison of alpha half-life approaches | exact formulas/tables needed |
| arXiv:1111.0505v1 | WKB approximation details for alpha half-life | full paper review needed |
| arXiv:1704.06334v1 | fission-barrier uncertainty method | full paper review needed |
| Ishizuka et al. 2023 | Langevin fission inputs/outputs/limitations | full paper review needed |

## Candidate electronic/chemistry facts and methods

| Source | Candidate extraction targets | Validation domain |
| --- | --- | --- |
| Yakushev et al. 2024 | relativistic chemistry effect statements, adsorption/gas chromatography observables if reported | chemistry/electronic only, not half-life |
| future Dirac-Hartree-Fock papers | Hamiltonian, basis, approximations, benchmark observables | electronic/atomic benchmarks |
| future relativistic DFT papers | functional, spin-orbit treatment, chemistry observables | chemistry/electronic benchmarks |

## Candidate numerical-method details

| Method source | Candidate extraction targets | Test target |
| --- | --- | --- |
| Monte Carlo references | RNG seed, convergence, confidence intervals | deterministic seed reproducibility |
| WKB/numerical quadrature references | integral bounds, potential model, quadrature tolerance | analytic integral tests plus half-life residuals |
| ODE/BVP solver references | method, tolerances, stability limits | manufactured solutions |
| Eigenvalue solver references | grid/basis, boundary conditions, convergence | square-well/hydrogen-like benchmarks |

## Non-extraction rule

Do not extract runtime facts from context-only sources. Track B sources can be summarized for narrative risk and misinformation boundaries only.
