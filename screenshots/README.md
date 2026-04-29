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
