# Moscovium Statera Go

Moscovium Statera Go is an evidence-first Go workspace for studying Moscovium isotope data, alpha-decay chains, and source-backed nuclear physics models.

The project is intentionally narrow. It is not a speculative materials platform, propulsion model, or popular-science claim engine. Its scope is evaluated nuclear data, peer-reviewed literature, and clearly labeled theoretical models.

## Goals

- Provide a small, auditable Go engine for Moscovium decay-chain traversal.
- Keep every isotope datum tied to evaluated records and DOI-backed source material.
- Build toward an education, research, and design UI without compromising source integrity.
- Keep the implementation pure Go with no CGO requirement.

## Current Status

This repository is in Phase 2 (expanded platform):

- 5 Moscovium isotopes (287Mc–291Mc) with full ENSDF provenance
- 15 daughter isotopes (Nh, Rg, Mt, Bh, Db, Lr) across complete decay chains
- 6 alpha-decay half-life models: Royer, WKB, VSS, UNIV, Denisov, binding energy
- Spontaneous fission half-life model (Swiatecki)
- Monte Carlo decay simulation with fixed-seed determinism
- Liquid-drop binding energy with shell corrections (Bethe-Weizsäcker)
- Fusion-evaporation excitation function model
- ENSDF text parser and batch intake pipeline
- gogpu/ui dashboard: N-Z chart, binding energy, decay chain viewer, model comparison, alpha systematics, provenance node table
- All isotope records carry DOI trails and ENSDF citation URLs

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
cmd/statera/                  CLI decay-chain smoke path + JSON reports
cmd/statera-ui/               gogpu/ui dashboard application
data/                         seed research records (20 isotopes)
docs/                         research charter, literature roadmap, specs
internal/physics/             decay, alpha models (Royer, WKB, VSS, UNIV, Denisov),
                              binding energy, SF fission, excitation, claims, nuclide
internal/research/            ENSDF ingestion, provenance graph, workbook,
                              CrossRef DOI, source triage, seed validation
internal/ui/                  AppModel, view specs, data loading for all dashboard sections
```

## License

MIT License. See [LICENSE](LICENSE).
