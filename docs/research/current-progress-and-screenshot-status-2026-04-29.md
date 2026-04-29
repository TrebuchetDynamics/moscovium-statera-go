# Current Progress and Screenshot Status — 2026-04-29

Repo: `/home/xel/git/sages-openclaw/workspace-riju/moscovium-statera-go`

This note records concrete progress, visible artifacts, and current limitations for Juan-facing status reports. It is a status document, not a data source.

## Current git state at capture

- Branch: `main`
- HEAD: `46dbbf2f1a91`
- `origin/main`: `46dbbf2f1a91`
- Working tree before this note: clean
- `rg`: missing on host; use `grep` / file-search fallback

## Real progress made today

The repository has moved beyond planning-only work. Recent commits include:

```text
46dbbf2 Expose provenance graph node table
374e7b9 Add provenance graph data model
82bc4f1 Render evaluated workbook in UI screenshot
d123c8f Add workbook provenance UI model
385e1e1 Add isotope workbook validation gate
8975dbe Add source-backed isotope workbook records
a303840 Document ENSDF intake design
6048dbc Reject blank isotope citation links
5b25b57 docs: add open source research playbook
99e89a8 docs: plan research simulation visualization ml roadmap
41c3cec Improve Statera demo screenshot summary layout
1bcb84f Validate research seed notes provenance
```

## Current implemented research capabilities

1. Evaluated isotope seed validation for `288Mc` and `290Mc`.
2. Source-backed isotope workbook records derived from validated research seed data.
3. Workbook validation gate that rejects provenance loss and inconsistent identity.
4. Workbook UI model exposing identity and provenance counts for screenshots.
5. Evaluated workbook rendered in UI screenshot.
6. Provenance graph data model:
   - isotope nodes
   - DOI nodes
   - citation URL nodes
   - source path nodes
   - blocked source nodes
   - deterministic node table
7. Blocked-source representation for Royer/Wang source access blockers.
8. ENSDF intake design from `paceENSDF` and `Nuclei` research.
9. Open-source research playbook covering Geant4, OpenMC, PyNE, paceENSDF, Nuclei, nuclear chart projects, DeepXDE, SciML ModelingToolkit, and Nemo.
10. Tested offscreen screenshot generation.

## Current tracked screenshot artifacts

Tracked files:

```text
screenshots/statera-alpha-systematics-2026-04-27.png
screenshots/statera-education-research-design-2026-04-27.png
screenshots/statera-mvp-offscreen-2026-04-27.png
screenshots/statera-ui-demo-summary-2x2-2026-04-29.png
screenshots/statera-ui-workbook-2026-04-29.png
```

Current generated offscreen screenshot during this status pass:

```text
artifacts/screenshots/statera-current-2026-04-29-1209.png
bytes=95438
width_px=1180
height_px=760
pixels=896800
unique_colors_sampled_or_exact=8290
single_color=false
```

This `artifacts/` path is transient and ignored. Curated screenshots belong under `screenshots/`.

## Text description of current app screenshot

The app screenshot shows the `Education` view of `Moscovium Statera Go` with a left navigation rail and a main research/education dashboard.

Visible status line:

```text
Track A isotopes 2 | research records 3 | context records 2 | 288Mc -> 284Nh -> 280Rg -> 276Mt -> 272Bh -> 268Db -> 264Lr
```

Visible summary cards:

```text
Education: 5 lessons
Research: 3 records
Design: 4 scenarios
Alpha: 2 isotopes
```

Visible education lessons include:

1. `Evaluated Isotope Records` — explains auditable fields: `Z`, `A`, half-life, `Q_alpha`, daughter, provenance.
2. `Decay-Chain Traversal` — explains deterministic daughter traversal and catalog completeness.
3. `Half-Life Checks` — explains claim checking against catalog half-life.
4. `Supported Model Boundary` — partially visible; explains unsupported mechanism claims.

Visible boundary notice:

```text
Track B context is excluded from simulation and cannot seed physics defaults.
```

## What the screenshot proves

The current screenshot proves UI integration for:

- source-backed education path,
- 2x2 summary-card layout,
- current isotope/research/context counts,
- decay-chain text summary,
- strict Track B exclusion notice,
- screenshot generation in headless/offscreen mode.

It does not yet show a plotted graph, heatmap, Monte Carlo trace, or interactive paper browser.

## Current limitations

1. Only `2` accepted Track A isotope seed records exist.
2. The visible screenshot is primarily text/card UI, not a numeric plot.
3. Provenance graph exists as data and UI summary/table model, but not a visual graph drawing yet.
4. Simulations are still deterministic validation/model summaries; Monte Carlo decay-chain simulation is planned but not yet implemented.
5. ML is still planned as source triage / duplicate detection; no scientific ML prediction is implemented or justified yet.
6. Royer/Wang coefficient/model work remains blocked by source access; APS returned HTTP `403` in prior audits.

## Next concrete build slices

Recommended order:

1. Render provenance node table in UI screenshot so Juan can see source graph facts, not only education cards.
2. Add screenshot manifest tests for captions, PNG dimensions, byte size, and nonblank status.
3. Add deterministic CLI command for provenance graph/node-table output.
4. Add decay constants and mean-life calculations for workbook isotopes.
5. Add Monte Carlo decay timing only after deterministic decay constants are covered.

## Reporting rule

When showing this project to Juan, use text and numbers:

- exact screenshot path,
- width and height in pixels,
- byte size,
- unique-color count,
- `single_color=false`,
- test pass/fail/skip counts,
- commit hash,
- explicit limitations.
