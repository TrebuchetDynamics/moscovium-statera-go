# Agent Instructions

## Project Boundary

Moscovium Statera Go is an academic nuclear-physics repository. Keep all prose sober, source-backed, and scoped to evaluated nuclear data or explicitly labeled theory.

Do not introduce unsupported claims, popular narratives, or non-standard physics. If a claim cannot be tied to evaluated data, a DOI, or a reproducible project artifact, do not add it.

## Development Rules

- Use test-driven development for behavior changes.
- Keep implementation pure Go and CGO-free.
- Keep gogpu/ui integration compatible with the repository Go target.
- Preserve citation links and DOI trails when editing data files.
- Keep seed data small and auditable.

## Verification

Before committing:

```bash
CGO_ENABLED=0 go test ./...
go run ./cmd/statera
CGO_ENABLED=0 go test ./cmd/statera-ui
git diff --check
```

For a desktop smoke test, run:

```bash
CGO_ENABLED=0 go run ./cmd/statera-ui
```

## Data Rule

Every isotope record must include source provenance. Missing provenance is a blocking error, not a warning.
