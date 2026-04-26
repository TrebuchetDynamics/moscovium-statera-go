# Contributing

Contributions are welcome when they improve source-backed nuclear data, tests, documentation, or pure-Go implementation quality.

## Required For Pull Requests

- A test for behavior changes.
- Source metadata for data changes.
- Clear distinction between evaluated data and theoretical extrapolation.
- No CGO dependency unless the project charter changes.

## Data Contributions

Data changes must include:

- evaluated-data URL or standards-body source
- DOI trail when available
- uncertainty values when available
- evidence level
- note explaining why the source should be trusted

Do not add values from informal summaries or secondary websites as simulation defaults.

## Local Checks

```bash
go test ./...
go run ./cmd/statera
go run ./cmd/statera-ui
git diff --check
```

