# Dual-Track Intake Design

## Goal

Moscovium Statera Go will preserve public context around Element 115 discourse without allowing unsupported narratives to become simulation inputs. The repository will keep evaluated nuclear data and DOI-backed physics in a strict calculation track, while quarantining media and public-claim context in a separate track that is searchable, cited, and explicitly excluded from physics defaults.

## Scope

This design covers the first repository-level architecture for:

- Track A: evaluated and peer-reviewed Moscovium data for calculations.
- Track B: contextual public discourse records for search and historical framing.
- A physics validator that compares testable claims against Track A.
- A fecim-style citation workflow for papers, source records, and disputed claims.

It does not implement a full RAG database, Hartree-Fock solver, or production UI in the first slice. Those remain later phases because each is large enough to require its own implementation plan.

## Track A: Simulation Substrate

Track A contains only source-backed nuclear physics records. Acceptable sources are evaluated nuclear data records, peer-reviewed papers, and standards or laboratory publications from sources such as ENSDF/NNDC, IUPAC/IUPAP, JINR, GSI, LBNL, LLNL, ORNL, or equivalent bodies.

Track A data may include isotope identity, `Z`, `A`, neutron count, half-life and uncertainty, decay mode, daughter isotope, alpha energy, `Q_alpha`, source URLs, DOI trails, and evidence level. Every isotope record remains blocked unless it has source provenance.

Track A is the only source for simulation defaults, decay traversal, stability checks, and any future numerical model.

## Track B: Contextual Archive

Track B records discourse about claims, not the truth of unsupported claims. It may include government oversight letters, mainstream reporting, institution statements, and claim metadata when those sources explain why public attention moved toward Element 115.

Track B records must be labeled with one or more of:

- `verified-event`: the event is supported by a reliable source, such as an official letter, institution notice, law-enforcement statement, or mainstream report.
- `media-claim`: the source reports or repeats a claim without independently proving it.
- `unsupported-linkage`: a source or community claim connects two facts without sufficient evidence.
- `physics-claim`: a claim can be expressed as a physics assertion that the validator may test.
- `not-for-simulation`: the record must never seed calculation defaults.

Examples:

- William Neil McCasland may be recorded as a missing-person event if sourced to law-enforcement reporting, congressional correspondence, or reliable news. Any Element 115 or UAP linkage remains `unsupported-linkage` unless a source proves it.
- Carl Grillmair may be recorded as a verified death or homicide when sourced to Caltech, the medical examiner, law enforcement, or reliable news. Claims that his death is evidence of an Element 115 pattern remain `unsupported-linkage`.
- Names or events submitted by users are intake candidates until independently sourced.

## Citation Workflow

The repository will adopt a small version of the `fecim-lattice-tools` citation structure:

```text
citations/
  README.md
  TEMPLATE.md
  refs.bib
  facts.md
  disputed.md
  papers/
  pdfs/
  queue/
  reports/
docs/lore/
  README.md
  intake-notes/
  records/
```

`citations/papers/` stores one Markdown record per paper or source. `citations/pdfs/` stores only open-access PDFs, accepted manuscripts, or other files that can be legally archived. Paywalled papers receive metadata-only records with DOI and source links.

`docs/lore/intake-notes/` may preserve user-submitted narrative text with a clear warning that it is not evidence. `docs/lore/records/` stores normalized Track B records after source review.

## Physics Validator

The validator receives a structured claim and returns a status, reason, and source trail.

Initial statuses:

- `supported-by-track-a`: the claim matches Track A within explicit uncertainty or categorical rules.
- `stability-incongruent`: the claim asserts long-lived or stable behavior for a known short-lived isotope.
- `outside-supported-model`: the claim cannot be represented in Standard Model-compatible nuclear physics used by this project.
- `insufficient-data`: the project lacks enough Track A data to test the claim.
- `invalid-claim`: the claim is malformed or mixes multiple assertions.

The first implementation slice should focus on decay traversal and stability checks for known isotope records. It should not assign fake numerical precision to unsupported concepts such as antigravity, propulsion, or frequency-based mechanisms.

## UI Direction

The gogpu/ui application should keep the same boundary as the data model:

- Education: source-backed explanations and decay-chain exploration.
- Research: citation browser, DOI trail, source metadata, and later model outputs.
- Design: clearly labeled theoretical exploration only.
- Context: Track B timeline and source records with warnings that lore records are not simulation data.

The Context view must not visually imply that unsupported linkage claims are validated.

## Legal And Source Handling

Automated research tooling may fetch metadata from Crossref, arXiv, DOI resolvers, and open repositories. It may download PDFs only when the source is open access, an accepted manuscript, a public-domain government document, or otherwise legally archivable. Otherwise, the project stores metadata and links only.

Any current-event record involving living people or recent deaths must use sober wording and cite reliable sources. If a source frames reporting as unconfirmed, the record must preserve that qualification.

## First Implementation Slice

The first implementation phase should avoid the full RAG system and UI expansion. It should instead add:

1. Stronger decay traversal behavior in `internal/physics/decay.go`.
2. A small Track A claim validator for isotope stability claims.
3. Citation/lore directory scaffolding modeled after `fecim-lattice-tools`.
4. Seed source records for core Moscovium papers and a small number of verified context sources.

This produces testable software while preserving the research boundary.

## Verification

Before claiming completion, run:

```bash
CGO_ENABLED=0 go test ./...
go run ./cmd/statera
CGO_ENABLED=0 go run ./cmd/statera-ui
git diff --check
```

If UI execution is blocked by local graphics constraints, record the exact error and run the non-UI tests and CLI path.
