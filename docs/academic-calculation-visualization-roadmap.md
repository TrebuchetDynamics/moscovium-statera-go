# Academic Calculation And Visualization Roadmap

This roadmap records academically grounded features that Statera can calculate or visualize next. It is a planning artifact, not a data file. A calculation may enter the engine only after its input records carry evaluated-data provenance, DOI trails, units, uncertainties, and tests.

## Design Rule

Every feature belongs to one of three evidence classes:

- `evaluated`: derived from ENSDF, NuDat, AME, NUBASE, or other evaluated data.
- `peer-reviewed-model`: computed from a named published model with source assumptions visible.
- `visual-context`: useful for education or review, but not permitted to seed simulation defaults.

Track B context remains excluded from all calculations.

## High-Value Calculation Targets

| Priority | Candidate | Evidence class | Inputs | Calculations | Visualizations | Implementation risk |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | Evaluated isotope workbook | evaluated | ENSDF/NuDat, NUBASE2020, `data/research.seed.json` | neutron number, decay constant, mean life, branching display, daughter traversal | nuclide cards, decay-chain timeline, log half-life strip | Low |
| 2 | Q_alpha and mass-energy worksheet | evaluated / peer-reviewed-model | AME2020 mass excess tables; optional model mass tables for extrapolation | `Q_alpha`, separation energies, mass-excess residuals, uncertainty propagation | N-Z heatmaps, isotope detail panel, residual plots | Medium |
| 3 | Alpha half-life systematics lab | peer-reviewed-model | evaluated `Q_alpha`, Z, A, parity class; Royer/VSS/UNIV/SemFIS formulas | predicted `log10(T_1/2)`, model residuals against evaluated half-lives, sensitivity to Q_alpha | model-comparison bands, Geiger-Nuttall-style plot, residual histogram | Medium |
| 4 | Synthesis yield and excitation-function explorer | peer-reviewed experimental | beam energy, excitation energy, observed chains, cross sections, target thickness, beam dose | expected event counts, chains per beam dose, cross-section curves, channel comparison | excitation-function chart, event-rate calculator, reaction-channel table | Medium |
| 5 | Fission barrier and survival map | peer-reviewed-model | fission barrier grids, excitation energy, spontaneous-fission half-lives, model label | barrier trends, survival probability proxies, alpha-vs-fission competition indicators | barrier heatmap, decay-mode competition map, uncertainty badges | High |
| 6 | Shell-structure / island-of-stability atlas | peer-reviewed-model | HFB/DFT or macroscopic-microscopic model tables | shell gaps, shell-correction energy, beta deformation, predicted bound-region masks | N-Z shell maps, N=162/N=184 guide lines, model-disagreement overlay | High |
| 7 | Relativistic chemistry module | peer-reviewed experimental/model | adsorption enthalpy, detector deposition counts, relativistic DFT references | adsorption-energy comparison, relative reactivity ranking, homolog comparison | quartz/gold deposition strip, adsorption enthalpy intervals, periodic group comparison | Medium |
| 8 | Evidence and uncertainty graph | project-generated from sources | citation records, evaluated data records, model outputs, claims | source-to-datum graph, evidence class propagation, blocked/default-safe states | provenance graph, claim audit trail, uncertainty stack | Low-Medium |

## Source-Backed Observations

- NuDat exposes half-life, decay mode, Q-values, separation energies, excited states, cross sections, and fission-yield datasets, so it can support broad education visualizations before Statera imports any new model code.
- NUBASE2020 evaluates mass excess, isomer excitation, half-life, spin/parity, decay modes, discovery year, and bibliographic information. It also flags non-experimental estimates, which is important for UI labeling.
- AME2020 publishes mass tables and reaction/decay-energy tables, including `Q_alpha`, neutron/proton separation energies, and covariance data. This is the cleanest path to a mass-energy worksheet.
- The 2022 DGFRS-2 `243Am + 48Ca` Mc experiments report production channels, excitation-energy intervals, observed decay chains, cross sections, alpha energies, and half-lives. These support an experimental yield/excitation-function module.
- Alpha-decay systematics papers explicitly compare empirical formulas and mass models. This supports a model-comparison lab, but outputs must be labeled theoretical unless reproduced from evaluated inputs.
- Recent superheavy chemistry work reports Mc/Nh gas-solid chromatography on silicon oxide and gold surfaces and ties trends to relativistic valence-orbital effects. This supports an education visualization, not nuclear stability claims.
- Superheavy fission and shell-structure papers are model-heavy. They are valuable for design exploration, but should be quarantined behind model labels, confidence notes, and no-default warnings.

## Recommended Next Build Slice

Build the `Alpha Systematics Lab` first.

Reasoning:

- It extends the existing Track A validator without importing Track B context.
- It uses small auditable inputs already present in the seed records: Z, A, half-life, Q_alpha, and daughter.
- It gives Education a strong interactive explanation: Q_alpha controls alpha-decay half-life estimates through tunneling/systematics.
- It gives Research a clear DOI-backed queue: Royer 2008, Wang et al. 2015, Hosseini and Hassanabadi 2017, AME2020, and NUBASE2020.
- It gives Design a constrained model-comparison scenario without implying a stable isotope.

Minimum implementation:

1. Add a `internal/physics/alpha.go` package with one published empirical formula and explicit model metadata.
2. Add tests that compare monotonic behavior: increasing Q_alpha lowers predicted alpha half-life for fixed Z/A.
3. Add a UI section that compares evaluated half-life to calculated model half-life for seed isotopes.
4. Label the result as `peer-reviewed-model`, not `evaluated`.

## References To Intake

- `giuliani2019superheavy` - Review framing for superheavy nuclear/atomic/chemistry questions.
- `nudat3ensdf` - Evaluated decay and nuclear-property data surface.
- `nubase2020evaluation` - Evaluated nuclide properties and estimate flags.
- `ame2020tables` - Mass and reaction-energy tables.
- `oganessian2022new286mc` - Additional Mc reaction/channel data.
- `wang2015alphaSystematics` - Q_alpha and half-life model comparison for superheavy nuclei.
- `royer2008alphaAnalytic` - Analytic alpha half-life formulas including superheavy predictions.
- `hosseini2017alphaApproaches` - Superheavy alpha half-life systematics and preformation-factor discussion.
- `ishizuka2023fissionLangevin` - Model fission-yield/TKE visual targets.
- `yakushev2024mcnhchemistry` - Mc/Nh gas chromatography and adsorption enthalpy visualization target.
