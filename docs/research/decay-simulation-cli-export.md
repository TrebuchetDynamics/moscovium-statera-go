# Decay Simulation CLI Export

Date: 2026-04-29
Repo: `/home/xel/git/sages-openclaw/workspace-riju/moscovium-statera-go`

This note documents the deterministic JSON export added for Decay-Chain Simulation v1. It is an export and audit surface only. It adds no isotope records, daughter half-lives, Q-values, citation URLs, DOI trails, or model coefficients.

## Command

```bash
CGO_ENABLED=0 go run ./cmd/statera -decay-simulation-format=json -decay-simulation-samples=64 -decay-simulation-seed=20260429
```

## Evidence boundary

- The export reads `data/research.seed.json` through the validated research workbook path.
- Current output covers only accepted seed records: `288Mc` and `290Mc`.
- Daughter IDs present in the seed are not expanded into simulated daughter isotope records, because daughter half-lives are not accepted in the current seed catalog.
- The source path in each simulation record is `data/research.seed.json`, preserving that the values came through the seed validation gate.
- Fixed-seed Monte Carlo samples are deterministic for the same sample count and seed.

## JSON top-level fields

| Field | Meaning |
| --- | --- |
| `report_type` | Always `decay_simulation_seed_summary` for this export. |
| `source_path` | Validated source JSON path: `data/research.seed.json`. |
| `sample_count` | Monte Carlo sample count supplied by `-decay-simulation-samples`; must be positive. |
| `seed` | Deterministic RNG seed supplied by `-decay-simulation-seed`; must be non-zero. |
| `records` | One simulation record per accepted seed isotope, sorted by isotope ID. |

## Per-record fields

| Field | Meaning |
| --- | --- |
| `id` | Accepted seed isotope ID. |
| `source_path` | Source path used by the simulation record. |
| `half_life_seconds` | Source-backed half-life from the accepted seed record. |
| `decay_constant_per_second` | `ln(2) / half_life_seconds`. |
| `mean_life_seconds` | `1 / decay_constant_per_second`. |
| `monte_carlo_p05_seconds` | Fixed-seed exponential-decay 5th percentile. |
| `monte_carlo_median_seconds` | Fixed-seed exponential-decay median. |
| `monte_carlo_p95_seconds` | Fixed-seed exponential-decay 95th percentile. |

## Current numeric smoke output

With `samples=64` and `seed=20260429`:

| ID | half_life_seconds | decay_constant_per_second | mean_life_seconds | p05 seconds | median seconds | p95 seconds |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `288Mc` | 0.17 | 4.077336356234972 | 0.2452581569511238 | 0.013032510769998121 | 0.16061436121022374 | 0.9516687441729552 |
| `290Mc` | 0.65 | 1.0663802777845313 | 0.9377517765778262 | 0.04983018823822811 | 0.6141137340390908 | 3.6387334336024755 |

## Not allowed from this export

- Do not treat Monte Carlo quantiles as evaluated nuclear data.
- Do not infer daughter isotope half-lives from parent seed rows.
- Do not use this export to validate Royer/Wang coefficients or formulas; those remain source-access blockers until source content is read.
