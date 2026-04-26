# Agent Instructions

## Project Boundary

Moscovium Statera Go is an academic nuclear-physics repository. Keep all prose sober, source-backed, and scoped to evaluated nuclear data or explicitly labeled theory.

Do not introduce unsupported claims, popular narratives, or non-standard physics. If a claim cannot be tied to evaluated data, a DOI, or a reproducible project artifact, do not add it.

## Development Rules

- Use test-driven development for behavior changes.
- Keep implementation pure Go and CGO-free.
- Do not add `gogpu/ui` as a dependency until the repository Go target is compatible with the library's documented requirements.
- Preserve citation links and DOI trails when editing data files.
- Keep seed data small and auditable.

## Verification

Before committing:

```bash
go test ./...
go run ./cmd/statera
go run ./cmd/statera-ui
git diff --check
```

## Data Rule

Every isotope record must include source provenance. Missing provenance is a blocking error, not a warning.

