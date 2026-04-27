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

On a headless machine or a shell without a graphical session, the run can fail before opening a window with:

```text
wayland: WAYLAND_DISPLAY not set
```

When a Wayland or X display is available, capture screenshots here and name them with the app view and date. Screenshots are review artifacts only. They are not simulation data, source provenance, or evaluated nuclear data.
