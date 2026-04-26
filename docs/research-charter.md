# Research Charter

Moscovium Statera Go exists to make source-backed nuclear data easier to inspect, test, and teach.

## Scope

This project is limited to:

- evaluated nuclear data
- peer-reviewed nuclear physics
- Standard Model-compatible nuclear structure models
- shell-model, mean-field, fission-barrier, and alpha-decay calculations when explicitly sourced
- education and research visualization

This project excludes:

- unsupported claims about materials, energy generation, propulsion, or hidden technology
- claims without DOI, evaluated-data, or standards-body provenance
- speculative narratives presented as technical evidence
- UI elements that imply a result is more validated than its source supports

## Source Policy

The project follows a zero-trust source policy.

Every datum used by the engine must link to at least one of:

- a DOI from peer-reviewed literature
- an evaluated nuclear data record from NNDC/ENSDF, IAEA, IUPAC, GSI, JINR, or a comparable standards body
- a project-generated validation artifact with reproducible code and input data

When evaluated data and literature disagree, the disagreement must be recorded in documentation or data metadata before the value is used.

## Data Integrity

All isotope records must preserve:

- isotope identity
- atomic number `Z`
- mass number `A`
- half-life and uncertainty when available
- decay mode and daughter
- energy values and uncertainty when available
- source URL
- DOI trail
- evidence level

No code path may silently replace a missing citation with a default value.

## Modeling Ethics

The simulator may support theoretical exploration, but theoretical outputs must be labeled as theoretical. A model result is not an experimental claim.

The UI must maintain the same boundary: education views may explain, research views may cite, and design views may explore, but none may erase uncertainty.

## Contribution Rule

Pull requests that add data or physics behavior must include:

- a test
- source metadata
- a note on model limits
- a clear distinction between evaluated data and theoretical extrapolation

