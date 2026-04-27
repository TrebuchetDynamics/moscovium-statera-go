# Citation System

This directory tracks source records for Moscovium Statera Go. Citations are treated like code: versioned, reviewable, searchable with standard command-line tools, and conservative by default.

The goal is to make every scientific or contextual claim traceable to one of three statuses:

- **Cited:** supported by an external source recorded in `citations/papers/`.
- **Demonstrated:** supported by a reproducible project artifact.
- **Quarantined:** preserved as Track B context and marked `not-for-simulation`.

Claims that cannot fit one of those categories should not be used in repository prose or runtime behavior.

## Directory Layout

| Path | Purpose |
|---|---|
| `TEMPLATE.md` | Template for one source record |
| `papers/` | One Markdown file per paper or source |
| `facts.md` | Cross-source index of verified facts |
| `disputed.md` | Conflicting, weak, or contested claims |
| `refs.bib` | BibTeX bibliography for source records |
| `pdfs/` | Optional local open-access PDFs |
| `queue/` | Intake queues |
| `reports/` | Citation audit and coverage reports |

## Rules

- Track A facts require evaluated data, DOI-backed literature, or standards-body provenance.
- Track B records describe discourse about claims and must not seed simulations.
- Local PDFs are allowed only for open-access, accepted-manuscript, public-domain, or otherwise legally archivable sources.
- Missing provenance is a blocking error for isotope data.
