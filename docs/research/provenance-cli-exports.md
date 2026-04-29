# Provenance CLI Exports: TSV and JSON Examples

Date: 2026-04-29
Repo: `/home/xel/git/sages-openclaw/workspace-riju/moscovium-statera-go`
Scope: documentation for deterministic provenance node-table exports in `cmd/statera`

This document records how to use the current provenance node-table CLI exports. It adds no isotope records, half-lives, Q values, alpha energies, branching ratios, model coefficients, citation URLs, or DOI trails.

## Evidence boundary

The provenance table is an audit surface over already validated workbook records and documented blocked sources. It is not a new scientific source.

Allowed use:

- confirm which accepted isotope records are connected to source paths, citation URLs, and DOI trails;
- confirm which source-access blockers are documented;
- export deterministic TSV or JSON for accessibility review, scripts, or regression comparison.

Not allowed use:

- treating DOI metadata as article contents;
- treating blocked Royer 2008 or Wang 2015 nodes as validated coefficient/formula evidence;
- adding new scientific values without source content, units, uncertainties when available, and tests.

## Commands

Default TSV report with the existing decay-chain line followed by the provenance table:

```bash
CGO_ENABLED=0 go run ./cmd/statera
```

Filtered TSV report for current source-access blockers:

```bash
CGO_ENABLED=0 go run ./cmd/statera -provenance-node-type=blocked_source -provenance-status=blocked
```

Filtered JSON report for current source-access blockers:

```bash
CGO_ENABLED=0 go run ./cmd/statera -provenance-format=json -provenance-node-type=blocked_source -provenance-status=blocked
```

## TSV schema

The TSV provenance table starts with a summary line:

```text
provenance_node_table rows=<row_count> columns=8 [filters=node_type=<type>,status=<status>]
```

The table has exactly 8 tab-separated columns:

1. `node_id`
2. `node_type`
3. `status`
4. `source_path`
5. `doi_or_url`
6. `incoming_edges`
7. `outgoing_edges`
8. `orphan`

Rows are sorted deterministically by node ID.

## JSON schema

JSON mode emits only JSON so downstream tools can parse stdout directly. The top-level fields are:

- `report_type`: currently `provenance_node_table`.
- `rows`: number of records after filters.
- `columns`: currently `8`.
- `filters`: object with the active `node_type` and `status` filters when supplied.
- `records`: array of provenance node records.

Each JSON record has:

- `node_id`
- `node_type`
- `status`
- `source_path`
- `doi_or_url`
- `incoming_edges`
- `outgoing_edges`
- `orphan`

## Current expected counts

Current unfiltered command evidence from 2026-04-29:

```text
provenance_node_table rows=17 columns=8
```

Current filtered blocked-source evidence from 2026-04-29:

```text
provenance_node_table rows=2 columns=8 filters=node_type=blocked_source,status=blocked
```

Current blocked-source rows:

```text
blocked_source:royer2008alphaAnalytic	blocked_source	blocked	citations/papers/royer2008alpha-analytic.md	10.1103/PhysRevC.77.037602	0	2	false
blocked_source:wang2015alphaSystematics	blocked_source	blocked	citations/papers/wang2015alpha-systematics.md	10.1103/PhysRevC.92.064301	0	2	false
```

Interpretation:

- `blocked_source:royer2008alphaAnalytic` documents source-access blockage for DOI `10.1103/PhysRevC.77.037602` and local note `citations/papers/royer2008alpha-analytic.md`.
- `blocked_source:wang2015alphaSystematics` documents source-access blockage for DOI `10.1103/PhysRevC.92.064301` and local note `citations/papers/wang2015alpha-systematics.md`.
- Each blocked source has `outgoing_edges=2`: one edge to the DOI node and one edge to the source-path note.
- These blocked nodes are access/provenance audit nodes only. They do not validate model coefficients or formulas.

## Accessibility notes

- TSV mode is compact and screen-reader friendly when read line-by-line.
- JSON mode is better for scripts because it omits the decay-chain text line and emits parseable JSON only.
- Blank `status` fields on DOI, citation URL, and source-path nodes mean the underlying graph node has no accepted/blocked status; UI rendering may display those as `referenced` for readability without changing graph semantics.

## Verification command

Use this gate after changes touching CLI export behavior or this documentation:

```bash
CGO_ENABLED=0 go test ./...
```

For code changes, use the full repository gates from `AGENTS.md` instead of this docs-only minimum.
