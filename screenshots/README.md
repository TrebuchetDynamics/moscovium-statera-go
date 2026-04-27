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

On a headless machine or a shell without a graphical session, the primary launch command writes an offscreen screenshot to the system temp directory and exits.

When a Wayland or X display is available, capture screenshots here and name them with the app view and date. Screenshots are review artifacts only. They are not simulation data, source provenance, or evaluated nuclear data.
