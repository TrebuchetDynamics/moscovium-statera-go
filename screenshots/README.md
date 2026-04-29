# Screenshots

This directory is for visual verification artifacts from the Statera desktop app.

Primary launch command:

```bash
CGO_ENABLED=0 go run ./cmd/statera-ui
```

Deterministic offscreen capture command:

```bash
CGO_ENABLED=0 go run ./cmd/statera-ui -screenshot screenshots/statera-mvp-offscreen-2026-04-27.png
```

Education/Research/Design capture command:

```bash
CGO_ENABLED=0 go run ./cmd/statera-ui -screenshot screenshots/statera-education-research-design-2026-04-27.png
```

Alpha Systematics capture command:

```bash
CGO_ENABLED=0 go run ./cmd/statera-ui -screenshot screenshots/statera-alpha-systematics-2026-04-27.png
```

Workbook capture command:

```bash
CGO_ENABLED=0 go run ./cmd/statera-ui -screenshot screenshots/statera-ui-workbook-2026-04-29.png
```

On a headless machine or a shell without a graphical session, the primary launch command writes an offscreen screenshot to the system temp directory and exits.

When a Wayland or X display is available, capture screenshots here and name them with the app view and date. Screenshots are review artifacts only. They are not simulation data, source provenance, or evaluated nuclear data.

## Manifest

| File | Width px | Height px | Pixels | Bytes | Unique colors | Single color | Text caption |
| --- | ---: | ---: | ---: | ---: | ---: | --- | --- |
| `screenshots/statera-ui-workbook-2026-04-29.png` | 1180 | 760 | 896800 | 95438 | 8290 | false | Statera UI screenshot with the evaluated isotope workbook section rendered from validated research seed rows; each workbook card summarizes isotope identity, neutron count, citation count, DOI count, source path, and evidence class without adding new scientific values. |
| `screenshots/statera-overview-2026-04-29.png` | 1180 | 760 | 896800 | 52872 | 4606 | false | Overview screenshot showing current Track A isotope count, research/context record counts, deterministic decay-chain text, boundary notice, and implemented workbench areas. |
| `screenshots/statera-workbook-2026-04-29.png` | 1180 | 760 | 896800 | 53767 | 4523 | false | Workbook screenshot showing accepted seed records as auditable isotope cards with Z, A, N, daughter, citation count, DOI count, source path, and no new scientific values. |
| `screenshots/statera-provenance-2026-04-29.png` | 1180 | 760 | 896800 | 89451 | 4943 | false | Provenance screenshot showing graph facts before drawings: graph summary plus source, DOI, citation, isotope, and blocker node-table rows. |
| `screenshots/statera-roadmap-2026-04-29.png` | 1180 | 760 | 896800 | 69415 | 4927 | false | Roadmap screenshot showing planned build order from evaluated workbook and provenance graph toward decay simulation, visual outputs, and safe ML source triage. |
| `screenshots/statera-simulation-2026-04-29.png` | 1180 | 760 | 896800 | 60473 | 4900 | false | Simulation-preview screenshot showing current deterministic chain and planned calculations: decay constant, mean life, fixed-seed Monte Carlo summaries, and current data limits. |
| `screenshots/statera-alpha-2026-04-29.png` | 1180 | 760 | 896800 | 57600 | 4650 | false | Alpha visualization screenshot showing Royer peer-reviewed-model boundary and current seed-isotope alpha comparison fields: Z, A, parity, Q_alpha, evaluated half-life, predicted half-life, and log residual. |
