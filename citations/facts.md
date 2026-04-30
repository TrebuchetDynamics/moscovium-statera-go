# Citable Facts Database

**Last updated:** 2026-04-29
**Total verified facts:** 0
**Current modeling anchors from accepted seed records:** 2

This file is the cross-source index for facts used by Moscovium Statera Go. It starts conservative on purpose. Promote a fact to **verified** only after reading the source and recording the same fact in `citations/papers/{key}.md` or a source-review document with exact source location.

## Rules

- Preserve exact values and units from the source.
- Include source location such as page, section, table, figure, ENSDF line, or source-review excerpt.
- Include experimental or simulation conditions.
- Record evidence level: evaluated data, peer-reviewed, preprint, government document, institution statement, mainstream reporting, or other.
- Track B context facts must be marked `not-for-simulation`.
- Candidate browser or ML extraction must remain `candidate-for-review` until parsed, checked, and manually accepted.
- Model outputs must be labeled `theoretical/model`; never treat them as experimental facts.

## Index

1. [Moscovium Isotope Data](#moscovium-isotope-data)
2. [Trusted Modeling Anchors](#trusted-modeling-anchors)
3. [Discovery And Naming](#discovery-and-naming)
4. [Decay Chains](#decay-chains)
5. [Theoretical Stability](#theoretical-stability)
6. [Relativistic And Quantum Modeling](#relativistic-and-quantum-modeling)
7. [Contextual Events](#contextual-events)

## Moscovium Isotope Data

No cross-source verified facts recorded yet.

The current engine seed is not empty. It contains accepted seed records for `288Mc` and `290Mc` in `data/research.seed.json`, but this facts database still requires full source-location review before counting those rows as cross-source verified facts.

## Trusted Modeling Anchors

These anchors are allowed as current Statera modeling inputs because they are already present in the accepted seed record and pass the seed validation gate. They are not yet counted in `Total verified facts` until exact source-location review is finished.

| Anchor ID | Datum | Value | Unit | Isotope | Source path | Citation URLs | DOI trail | Allowed use | Status |
| --- | --- | ---: | --- | --- | --- | --- | --- | --- | --- |
| `anchor-288Mc-half-life-seed-v1` | Half-life | `0.17` | `s` | `288Mc` | `data/research.seed.json` | `https://www.nndc.bnl.gov/ensnds/288/Mc/adopted.pdf`; `https://www.nndc.bnl.gov/nudat3/getdecaydataset.jsp?dsid=288mc+a+decay+%280.17+s%29&nucleus=284NH` | `10.1103/PhysRevC.106.L031301`; `10.1016/j.nuclphysa.2003.11.001` | decay constant, mean life, Monte Carlo decay-time simulation, alpha residual target | accepted seed anchor; source-location review still required |
| `anchor-290Mc-half-life-seed-v1` | Half-life | `0.65` | `s` | `290Mc` | `data/research.seed.json` | `https://www.nndc.bnl.gov/ensnds/290/Mc/adopted.pdf`; `https://www.nndc.bnl.gov/ensnds/290/Mc/a_decay_51_ms.pdf` | `10.1103/PhysRevLett.104.142502`; `10.1103/PhysRevC.99.054306` | decay constant, mean life, Monte Carlo decay-time simulation, alpha residual target | accepted seed anchor; source-location review still required |

### Why these anchors matter

Half-lives are useful targets for nuclear decay models. A model can predict a half-life from inputs such as Q-alpha and a barrier-penetration formula, then Statera can compare the prediction to these accepted seed anchors.

Safe comparison shape:

```text
residual_log10 = log10(predicted_half_life_seconds / accepted_half_life_seconds)
```

If `residual_log10` is near `0`, the model is close for that isotope. If it is far from `0`, the model is not matching the accepted anchor.

Important boundary: these half-life anchors validate nuclear decay models. They do not validate electronic relativistic chemistry models.

## Discovery And Naming

No verified facts recorded yet.

Planned source review targets:

- IUPAC naming source for Moscovium.
- Primary synthesis/discovery papers.
- Institution reports only as context unless tied to peer-reviewed or standards-body sources.

## Decay Chains

No cross-source verified facts recorded yet.

Current candidate/source-review evidence exists for `288Mc` alpha decay in:

- `docs/research/nndc-288mc-alpha-decay-source-review-2026-04-29.md`

That review records browser-readable ENSDF text but keeps the values `candidate-only` until a parser and manual acceptance step preserve ENSDF notation, qualifiers, uncertainty, and source evidence.

## Theoretical Stability

No verified facts recorded yet.

Theory/model sources must be recorded as model evidence, not evaluated experimental data. Relevant queued source notes include:

- `citations/papers/giuliani2019superheavy.md`
- `citations/papers/erler2012-nuclear-landscape.md`

## Relativistic And Quantum Modeling

No verified numerical facts recorded yet.

Planning document:

- `docs/research/relativistic-quantum-simulation-roadmap.md`

Modeling boundary:

- Relativistic electronic/chemistry models can be compared to chemistry observables such as adsorption behavior.
- Nuclear half-life models can be compared to half-life anchors.
- A full Moscovium atom simulator must separate electron physics from nuclear physics and label every theoretical output.

## Contextual Events

No verified facts recorded yet.
