# Moscovium Grand Questions: Research Plan and Evidence Boundaries

Date: 2026-04-29
Repo: `/home/xel/git/sages-openclaw/workspace-riju/moscovium-statera-go`

This document defines how Statera will answer Juan's large research questions about Moscovium without overclaiming. It separates current source-backed answers, candidate evidence, planned calculations, and explicitly blocked/speculative claims.

## Core rule

Moscovium Statera Go may explore big questions, but every answer must be labeled as one of:

- `evaluated-data`: derived from NNDC/ENSDF, NuDat, NUBASE, AME, IUPAC, or comparable evaluated/standards sources.
- `peer-reviewed-experimental`: derived from a read paper or source-backed experiment report.
- `peer-reviewed-model`: theoretical/model calculation with source, assumptions, coefficients, and limits visible.
- `candidate-for-review`: source text captured but not parsed/accepted yet.
- `blocked`: source inaccessible or evidence insufficient.
- `unsupported`: claim has no accepted path in current nuclear physics evidence.

No `candidate-for-review`, `blocked`, or `unsupported` item may seed simulations, screenshots, or UI defaults as if it were evaluated fact.

## Current accepted data baseline

Current accepted seed records in `data/research.seed.json`:

| Isotope | Z | A | N | Half-life seconds | Q_alpha MeV | Daughter | Status |
| --- | ---: | ---: | ---: | ---: | ---: | --- | --- |
| `288Mc` | 115 | 288 | 173 | 0.17 | 10.75 | `284Nh` | accepted seed record |
| `290Mc` | 115 | 290 | 175 | 0.65 | 10.45 | `286Nh` | accepted seed record |

Current browser source-review candidate:

- `docs/research/nndc-288mc-alpha-decay-source-review-2026-04-29.md`
- NuDat/ENSDF text for `288MC A DECAY (0.17 S)` was readable through browser DOM.
- Direct scripted fetch returned HTTP `429`; this must guide respectful, low-rate source access.
- Candidate values were documented but not promoted into accepted data.

## Question 1 — Is Moscovium possible? Has it been achieved?

### Current answer status

`peer-reviewed/evaluated answer in progress`

The project already treats Moscovium as a real element with accepted seed records for `288Mc` and `290Mc`, tied to evaluated-data URLs and DOI trails in `data/research.seed.json`.

However, the cross-source citable facts index at `citations/facts.md` still reports `Total verified facts: 0`; therefore the public-facing answer should cite the seed record and source-review docs, but the facts database needs to be populated before claiming a finalized research summary.

### Required final answer shape

A complete answer must include:

1. discovery/synthesis source references;
2. IUPAC naming recognition reference;
3. isotope records and observed decay chains;
4. distinction between existence of atoms and availability of macroscopic material.

### Planned work

- Add a `Discovery And Naming` section to `citations/facts.md` only after reviewing IUPAC and discovery sources.
- Add a source review for IUPAC 2016 naming source.
- Add a source review for Oganessian/JINR/DGFRS production papers.

## Question 2 — What is the half-life of this element?

### Current answer status

`evaluated-data seed answer exists for two isotopes`

Current accepted seed values:

- `288Mc`: `0.17 s`
- `290Mc`: `0.65 s`, interval fields `0.45 s` lower and `1.14 s` upper in the current seed record

### Important caveat

There is no single half-life for “the element” as a whole. Half-life belongs to each isotope. Moscovium isotopes have different neutron counts and different half-lives.

### Required final answer shape

The project should answer with an isotope table:

| Isotope | Half-life | Uncertainty/interval | Decay mode | Daughter | Evidence |
| --- | ---: | --- | --- | --- | --- |

### Planned calculations

1. Decay constant for each accepted isotope:

```text
lambda = ln(2) / T_1/2
```

2. Mean life:

```text
tau = T_1/2 / ln(2)
```

3. Monte Carlo decay-time simulation with fixed RNG seed:

```text
input: isotope, half-life, sample count, RNG seed
output: median, p05, p95 decay time
```

## Question 3 — How can we build/make Moscovium in a particle accelerator?

### Current answer status

`planned peer-reviewed-experimental answer`

The roadmap already identifies hot-fusion synthesis campaigns and the reaction class involving `243Am + 48Ca` as a research target. The project must not provide a lab recipe as if it were an actionable production protocol. It can explain academically how experiments synthesize superheavy elements: target nucleus, projectile beam, compound nucleus, evaporation channels, separator/detector chain, observed alpha-decay chains, and cross sections.

### Required final answer shape

A safe academic answer should include:

1. target isotope;
2. projectile isotope;
3. beam energy / excitation-energy range from source;
4. evaporation channel/product isotope;
5. expected event count or observed chains;
6. cross section with units and uncertainty;
7. detector/separator context;
8. limits: atom-at-a-time production, short half-lives, no bulk material.

### Planned calculations

Synthesis yield/excitation-function explorer:

```text
expected_events = beam_particles * target_areal_density * cross_section * detection_efficiency
```

All variables require units and source provenance.

### Blocked until sourced

- Exact beam energies, cross sections, target thicknesses, and observed event counts cannot be added until the relevant papers/tables are read and recorded.

## Question 4 — Are there alternative ways to make Moscovium?

### Current answer status

`research question, not accepted engineering path`

The project can compare peer-reviewed production channels if sourced. It must not claim unsupported alternative production methods.

### Required evidence

For each proposed alternative:

- reaction equation;
- source DOI/evaluated reference;
- cross section or predicted yield;
- required beam/target feasibility;
- model vs experiment label;
- uncertainty and limitations.

### Statera implementation

Add a reaction-channel table with fields:

```text
target
projectile
compound nucleus
product isotope
evaporation channel
beam energy MeV
excitation energy MeV
cross section pb
uncertainty pb
observed chains
source DOI
status: experimental | theoretical-model | blocked
```

## Question 5 — Does Moscovium have “cool” nuclear, relativistic, or chemical properties?

### Current answer status

`partly planned; no accepted property claims yet beyond seed nuclear records`

Likely research categories:

1. nuclear decay properties: alpha decay, half-life, Q-alpha;
2. superheavy nuclear-structure effects: shell corrections, island-of-stability models;
3. relativistic chemistry: predicted or measured behavior of Mc/Nh chemistry;
4. fission/survival competition models.

### Critical boundary

Relativistic effects in superheavy elements are primarily electronic/chemical phenomena unless a source explicitly discusses nuclear-structure effects. They do not justify unsupported energy, propulsion, or clean power claims.

### Planned visualizations

- Alpha systematics residual dashboard.
- N-Z chart with `N=162` and `N=184` guide lines.
- Relativistic chemistry comparison cards only after source-backed chemistry papers are reviewed.
- Fission/survival maps only as model-labeled outputs.

## Question 6 — Can Moscovium be an energy source?

### Current answer status

`unsupported as practical energy source with current accepted evidence`

Current accepted Statera data indicates millisecond-to-subsecond half-lives for accepted seed isotopes. That is useful for nuclear-decay research, not evidence of a practical material energy source.

### What can be calculated safely

For a decay with energy `E` per atom, theoretical microscopic energy per atom can be estimated:

```text
E_joules = E_MeV * 1.602176634e-13 J/MeV
```

But this does not imply practical extractable energy because production is atom-scale, half-lives are short, and capture efficiency/availability are not established.

### Required before any energy discussion

1. accepted isotope inventory;
2. decay energies and branching ratios;
3. production yield/cross section;
4. number of atoms realistically producible;
5. energy capture mechanism, if any;
6. safety/radiation handling;
7. peer-reviewed source.

### Project stance

Statera may calculate decay energy bookkeeping as education, but must label any macroscopic energy-use claim as unsupported unless backed by peer-reviewed evidence and realistic production/yield calculations.

## Question 7 — Could it be used for propulsion?

### Current answer status

`unsupported`

The current validator already treats unsupported mechanism claims as outside supported model scope. Propulsion claims require specific, peer-reviewed mechanism, energy density, controllability, thrust/impulse calculation, material availability, and safety analysis.

### Required evidence for future reconsideration

- peer-reviewed propulsion mechanism using Moscovium or a specific isotope;
- isotope production feasibility;
- decay/control mechanism;
- thrust or energy calculation with units;
- comparison to known nuclear propulsion concepts;
- independent source verification.

Without those, Statera should answer: `unsupported by current evidence`.

## Question 8 — Could it heat water or be cleaner than uranium?

### Current answer status

`unsupported as practical engineering claim`

Heating water is an energy conversion application. A valid comparison to uranium would require:

- obtainable mass or atom count;
- half-life and decay-chain energy release;
- controllability;
- radiation products;
- production cost/yield;
- waste/safety profile;
- source-backed reactor or heat-source design.

Current accepted seed records are atom-scale and short-lived. They do not support a practical heat-source claim.

### Safe future calculation

Statera can build a “material availability and energy bookkeeping” worksheet:

```text
atoms produced per experiment
energy per decay
estimated total energy if every decay were captured
capture efficiency assumption
comparison to production energy/cost
```

This must be marked educational/model-only unless sourced experimentally.

## Question 9 — What future use should we expect?

### Current answer status

`research/education outlook only`

Source-backed future expectation categories:

1. understanding superheavy nuclear stability;
2. testing nuclear models near island-of-stability regions;
3. refining synthesis techniques and detector methods;
4. probing relativistic chemistry trends for group 15 superheavy elements;
5. improving evaluated-data and provenance tooling.

Not currently supported:

- bulk material use;
- energy extraction technology;
- propulsion;
- clean power;
- hidden materials applications.

## Required Statera answer engine

To answer all questions responsibly, build these modules:

1. **Facts database:** source-reviewed, citable facts with units and source locations.
2. **Isotope workbook:** accepted isotope records, uncertainty, decay mode, daughters.
3. **Decay calculator:** decay constants, mean life, Monte Carlo decay timing.
4. **Production/yield calculator:** accelerator target/projectile/cross-section/event count model.
5. **Energy bookkeeping worksheet:** per-decay energy and impossible/practicality constraints.
6. **Chemistry/property review:** relativistic chemistry source cards.
7. **Claim validator:** marks claims supported, model-only, blocked, or unsupported.
8. **Screenshot gallery:** overview, workbook, provenance, roadmap, simulation preview, alpha model comparison.

## Immediate implementation plan

### Slice A — Screenshots for Juan

Implemented/underway:

- `overview`
- `workbook`
- `provenance`
- `roadmap`
- `simulation`
- `alpha`

These are UI screenshots, not final physics answers.

### Slice B — Grand Questions FAQ in UI

Add a new UI/data model section with question cards:

- “Is Moscovium possible?”
- “What are the half-lives?”
- “How is it synthesized?”
- “What can it be used for?”
- “What is unsupported?”

Each card must include an evidence status label.

### Slice C — Facts database population

Populate `citations/facts.md` only after source review docs exist for:

- IUPAC naming/recognition;
- Oganessian/DGFRS production papers;
- NNDC/ENSDF records for 288Mc and 290Mc;
- AME/NUBASE references for mass/half-life data if used;
- relativistic chemistry papers if discussing chemistry.

### Slice D — Calculators

Build in this order:

1. decay constant and mean life;
2. decay Monte Carlo;
3. alpha residual summary;
4. synthesis yield/event count;
5. energy bookkeeping worksheet;
6. chemistry/property comparison.

## Current answer summary for Juan

- Moscovium is treated by the project as a real synthesized superheavy element, but the public-facing facts database still needs source-reviewed entries.
- Current accepted Statera half-life records are `0.17 s` for `288Mc` and `0.65 s` for `290Mc`.
- Making Moscovium belongs to particle-accelerator hot-fusion research; exact production answers require source-reviewed reaction/cross-section data.
- Cool nuclear/chemical phenomena are plausible research topics, but claims must be source-labeled.
- Energy, propulsion, clean power, and water-heating applications are unsupported by current accepted evidence and must remain blocked unless future peer-reviewed evidence and feasibility calculations support them.

## Definition of success

This super-project is successful when Statera can answer each grand question with:

- evidence label;
- numeric values with units;
- source URLs and DOI trails;
- uncertainty or qualifier;
- calculations when applicable;
- screenshot-ready explanation;
- unsupported-claim warnings where needed.
