# CLAUDE.md - moscovium-statera-go

This file applies to `/home/xel/git/sages-openclaw/workspace-riju/moscovium-statera-go`. Read local project docs and manifests before making changes.

If a sibling `AGENTS.md` exists, keep behavior consistent between the two instruction files.

<!-- karpathy-guidelines:start -->
## Karpathy-Inspired Agent Guardrails

Source: https://github.com/forrestchang/andrej-karpathy-skills at commit `2c60614`.

These guardrails supplement the local instructions above. Local project, safety, and user-specific rules win on conflict.

Tradeoff: they bias toward caution over speed for non-trivial work; use judgment for obvious one-line fixes.

### Think Before Coding

- State assumptions before implementing; ask when uncertainty would change the solution.
- Surface multiple interpretations and tradeoffs instead of silently picking one.
- Push back when a simpler approach meets the goal.

### Simplicity First

- Build the minimum code that solves the requested problem.
- Avoid speculative features, single-use abstractions, and unnecessary configurability.
- If the solution is growing large, stop and simplify before continuing.

### Surgical Changes

- Touch only files and lines required by the request.
- Preserve existing style, comments, and nearby code unless the task requires changing them.
- Clean up only dead code introduced by your own change; mention unrelated dead code instead of deleting it.

### Goal-Driven Execution

- Convert the request into verifiable success criteria before editing.
- For multi-step work, state a short plan with a verification check for each step.
- Loop until the relevant tests, builds, or manual checks prove the goal is met.
<!-- karpathy-guidelines:end -->

<!-- karpathy-project-adjustment:start -->
## Project-Specific Karpathy Adjustment

This section localizes the Karpathy guardrails for `workspace-riju/moscovium-statera-go`. Source inspiration: https://github.com/forrestchang/andrej-karpathy-skills at commit `2c60614`.

- Project family: Riju FeCIM / scientific simulation workspace.
- Local focus: ferroelectric compute-in-memory simulation, validation, citations, and reproducible Go/Fyne tooling.
- Stack cues: Go.
- Evidence to prefer: units, tolerances, cited sources, validation scripts, headless test output, numeric deltas, and before/after result tables.
- Surgical boundary: do not present simulation defaults as measured silicon facts; keep scientific claims cited, range-checked, and explicitly verified or marked unverified.
- Stop and ask when: a claim needs a citation, physical interpretation is ambiguous, or GUI-only checks would replace required headless validation.
<!-- karpathy-project-adjustment:end -->
