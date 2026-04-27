# Literature Roadmap

This file records intake targets. It is not evidence until a source is read, cited, and added to `data/` or future citation records.

## Priority Topics

1. Evaluated-data ingestion from ENSDF/NuDat, NUBASE2020, and AME2020.
2. Alpha-decay systematics for superheavy nuclei, including model residuals against evaluated half-lives.
3. Recent decay spectroscopy and production channels for products of `243Am + 48Ca`.
4. Reaction-yield and excitation-function analysis for hot-fusion synthesis campaigns.
5. Relativistic chemistry of Moscovium and Nihonium on silicon oxide and gold surfaces.
6. Theoretical work around `N=162`, `N=184`, shell corrections, fission barriers, and survival probabilities.
7. Evidence and uncertainty visualization for source-to-datum provenance.

## Feature Roadmap

See `docs/academic-calculation-visualization-roadmap.md` for the current mapping from academic literature to possible Statera calculations and UI views.

## Intake Standard

For each paper:

- capture DOI and bibliographic metadata
- identify whether it is experimental, evaluated, theoretical, or review work
- extract only exact values with context
- record uncertainty and model assumptions
- avoid adding values to simulation defaults until tests use them
