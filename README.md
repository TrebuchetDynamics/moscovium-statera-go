# Moscovium Statera Go

Moscovium Statera Go is an evidence-first Go workspace for studying Moscovium isotope data, alpha-decay chains, and source-backed nuclear physics models.

The project is intentionally narrow. It is not a speculative materials platform, propulsion model, or popular-science claim engine. Its scope is evaluated nuclear data, peer-reviewed literature, and clearly labeled theoretical models.

## Goals

- Provide a small, auditable Go engine for Moscovium decay-chain traversal.
- Keep every isotope datum tied to evaluated records and DOI-backed source material.
- Build toward an education, research, and design UI without compromising source integrity.
- Keep the implementation pure Go with no CGO requirement.

## Current Status

This repository is in Phase 1 bootstrap:

- core decay traversal exists
- high-precision unit conversion exists
- Crossref DOI metadata retrieval exists
- seed data includes `288Mc` and `290Mc`
- gogpu/ui application shell exists for education, research, and design modules

## Quick Start

Prerequisite: Go 1.25 or newer.

```bash
CGO_ENABLED=0 go test ./...
go run ./cmd/statera
CGO_ENABLED=0 go run ./cmd/statera-ui
```

The gogpu/ui stack is zero-CGO and uses `github.com/go-webgpu/goffi`, which requires `CGO_ENABLED=0` when a local C compiler is installed.

## Research Boundary

Read [docs/research-charter.md](docs/research-charter.md) before adding models, data, or claims.

Short version:

- Standard nuclear physics only.
- No unverified claims.
- No non-standard forces or unsupported application claims.
- Every scientific datum must carry provenance.

## Layout

```text
cmd/statera/             CLI decay-chain smoke path
cmd/statera-ui/          gogpu/ui application shell
data/                    seed research records
docs/                    research charter and literature roadmap
internal/physics/        decay traversal and unit helpers
internal/research/       citation and DOI metadata utilities
```

## License

MIT License. See [LICENSE](LICENSE).
