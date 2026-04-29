# ENSDF Intake Design: paceENSDF and Nuclei Deep Dive

Date: 2026-04-29
Repo: `/home/xel/git/sages-openclaw/workspace-riju/moscovium-statera-go`
Scope: Deep Dive 1 from `docs/superpowers/plans/2026-04-29-open-source-research-playbook.md`

This document is a design and research note only. It adds no new isotope records, half-lives, Q values, decay energies, branching ratios, or model coefficients. It defines a source-to-candidate-to-accepted workflow for future evaluated ENSDF intake in Moscovium Statera Go.

## Method

Startup and source-access method used:

1. Confirmed the Statera repository path and clean `main` branch state.
2. Pulled `origin main` with `git pull --ff-only origin main` before edits because the inner repository was clean.
3. Reviewed local project boundary files before task selection.
4. Queried public GitHub repository metadata for the two Deep Dive 1 projects.
5. Fetched direct raw README files from GitHub for both projects.
6. Inspected public repository root contents through the GitHub contents API.
7. Extracted architecture lessons only. No external code was copied.

Limitations:

- Review was limited to public metadata, README content, and root file listings.
- No ENSDF source archive was downloaded in this slice.
- No external parser output was imported into Statera.
- No scientific datum was promoted from README prose to Statera data.

## External projects reviewed

### `AaronMHurst/pace_ensdf`

Repository: `AaronMHurst/pace_ensdf`
URL: `https://github.com/AaronMHurst/pace_ensdf`
GitHub API URL reviewed: `https://api.github.com/repos/AaronMHurst/pace_ensdf`
Raw README URL reviewed: `https://raw.githubusercontent.com/AaronMHurst/pace_ensdf/main/README.md`
Root contents URL reviewed: `https://api.github.com/repos/AaronMHurst/pace_ensdf/contents`

Observed GitHub metadata from API response:

- HTTP status: `200`
- Stars: `10`
- Forks: `0`
- Primary language: `Python`
- Default branch: `main`
- License field: `NOASSERTION`
- Updated at: `2026-04-21T19:48:26Z`
- Description: `paceENSDF: Python Archive of Coincident Emissions from ENSDF. Package enabling interaction, manipulation, analysis, and visualization of radioactive-decay data from the ENSDF archive and corresponding coincidence gamma-gamma and gamma-X-ray emissions..`

Observed README evidence:

- HTTP status for `main` README: `200`
- README byte count: `94366`
- The README describes a Python package for access, manipulation, analysis, and visualization of radioactive-decay data from ENSDF.
- The README describes translation of ENSDF decay-scheme data into JSON structures.
- The README explicitly discusses quantity handling for asymmetric values, approximate values, and limit values.
- The README describes uncertainty-bearing coincidence gamma-gamma and gamma-X-ray structures.

Observed root contents:

- HTTP status: `200`
- Files/directories listed included `README.md`, `LICENSE`, `Hurst_EPJ_ND2022.pdf`, `installation.sh`, `requirements.txt`, `setup.py`, `tox.ini`, `notebook/`, `paceENSDF/`, and `tests/`.

Lessons for Statera:

1. Treat ENSDF parsing as a translation step into an explicit intermediate representation before any engine record is accepted.
2. Preserve the original ENSDF quantity qualifier, not only a normalized numeric value.
3. Store uncertainty policy alongside parsed values because approximate and limit values are not equivalent to ordinary measured central values.
4. Keep parsed candidate records separate from accepted Statera seed or workbook records.
5. Maintain tests around parser and validator behavior before UI or simulation use.

### `martukas/nuclei`

Repository: `martukas/nuclei`
URL: `https://github.com/martukas/nuclei`
GitHub API URL reviewed: `https://api.github.com/repos/martukas/nuclei`
Raw README URL reviewed: `https://raw.githubusercontent.com/martukas/nuclei/main/README.md`
Root contents URL reviewed: `https://api.github.com/repos/martukas/nuclei/contents`

Observed GitHub metadata from API response:

- HTTP status: `200`
- Stars: `26`
- Forks: `7`
- Primary language: `C++`
- Default branch: `main`
- License field: `GPL-3.0`
- Updated at: `2025-12-05T15:23:10Z`
- Description: `An Evaluated Nuclear Structure Data (ENSDF) parser, viewer and editor`

Observed README evidence:

- HTTP status for `main` README: `200`
- README byte count: `3333`
- The README describes Nuclei as an ENSDF parser, viewer, and editor.
- The README states a technical goal of complete and authentic parsing of ENSDF files, including information some projects neglect, with half-life uncertainties named as an example.
- The README describes a hierarchical object-oriented in-memory representation after ENSDF import.
- The README says parsing errors are printed in the terminal while files are loaded.
- The README cites `M. Nagl, et al., NIM A 726 (2013), 17-30, doi:10.1016/j.nima.2013.05.045` as a publication about the original Nuclei functionality and search method.

Observed root contents:

- HTTP status: `200`
- Files/directories listed included `.circleci/`, `.gitignore`, `CMakeLists.txt`, `LICENSE`, `README.md`, `TODO`, `cmake/`, `conanfile.txt`, `documentation/`, and `source/`.

Lessons for Statera:

1. Favor complete preservation of evaluated-data fields over task-specific partial extraction.
2. Separate parser storage, validation, and UI presentation so UI convenience cannot silently drop provenance or uncertainty.
3. Surface parse errors as first-class candidate status, not terminal-only logs.
4. Treat literature references as links to evidence that must remain attached to candidate records.
5. Do not reuse GPL-3.0 code in the MIT Statera repository without explicit legal review and approval; use only architecture lessons in this document.

## Statera ENSDF intake boundary

ENSDF intake should be a controlled pipeline with three states:

1. `source`: raw source evidence that can be fetched, hashed, and cited.
2. `candidate`: machine-parsed extraction that is not trusted as engine data.
3. `accepted`: human-reviewed, test-covered, provenance-complete record allowed to feed workbook, catalog, provenance graph, UI, and later simulation paths.

No direct path should exist from `source` or `candidate` to simulation defaults.

## Proposed source-to-candidate-to-accepted workflow

### Step 1: source registration

Create a source manifest entry before parsing any ENSDF-derived content.

Required source fields:

- `source_id`: stable local identifier.
- `source_type`: for example `ensdf_adopted_pdf`, `ensdf_dataset`, `nudat_html`, or `nndc_archive`.
- `source_url`: original HTTPS URL.
- `retrieved_at`: UTC timestamp.
- `retrieval_status`: HTTP status or local-access status.
- `content_sha256`: hash when content is stored locally.
- `content_bytes`: byte count when content is stored locally.
- `license_or_terms_note`: short note when known.
- `source_path`: local path when a fetched artifact is stored.

### Step 2: candidate extraction

Parse only into candidate records. Candidates may contain missing fields, parse warnings, qualifiers, and uncertain mappings.

Candidate output should be deterministic and network-free after source registration. The parser must not repair missing provenance or invent values.

### Step 3: candidate validation

Run validators that classify candidates into:

- `valid_candidate`: structurally coherent but not yet accepted.
- `requires_human_review`: structurally parseable but contains qualifiers, limits, conflicts, missing optional context, or low-confidence mapping.
- `blocked`: missing source identity, impossible nuclide identity, malformed units, unsupported decay mode, or broken provenance.

### Step 4: acceptance review

A candidate can become accepted only when all acceptance requirements are met:

- Source artifact and source URL are recorded.
- DOI trail is present when literature is used.
- Uncertainty and ENSDF qualifier are preserved when available.
- A Go test asserts the accepted value or acceptance behavior.
- The record names its evidence class.
- The record does not conflict with existing accepted records unless the conflict is documented.

### Step 5: accepted record generation

Accepted records should be generated into a Statera-owned data file or typed Go fixture, then validated through existing or new pure-Go validation gates. Candidate-only fields may remain in a candidate archive but must not become implicit engine defaults.

## Minimum candidate schema

The following schema is intentionally broader than the current `data/research.seed.json` record. It is for future ENSDF candidate intake, not current accepted engine data.

```json
{
  "schema": "moscovium-statera-go/ensdf-candidate/v1",
  "source": {
    "source_id": "string",
    "source_type": "string",
    "source_url": "https://example.invalid/source",
    "retrieved_at": "YYYY-MM-DDTHH:MM:SSZ",
    "retrieval_status": "string",
    "content_sha256": "hex string or empty when not stored",
    "content_bytes": 0,
    "source_path": "relative/path/or/empty"
  },
  "candidate": {
    "candidate_id": "string",
    "parent_nuclide_id": "288Mc-style identifier",
    "element": "element name from periodic-table helper",
    "symbol": "chemical symbol",
    "z": 0,
    "a": 0,
    "n": 0,
    "dataset_kind": "adopted | decay | other",
    "decay_mode": "alpha | beta_minus | electron_capture | spontaneous_fission | unknown",
    "daughter": "nuclide ID or empty",
    "half_life": {
      "value_seconds": null,
      "uncertainty_seconds": null,
      "lower_seconds": null,
      "upper_seconds": null,
      "ensdf_qualifier": "= | AP | LT | GT | LE | GE | empty",
      "raw_text": "original extracted text"
    },
    "q_alpha": {
      "value_mev": null,
      "uncertainty_mev": null,
      "ensdf_qualifier": "= | AP | LT | GT | LE | GE | empty",
      "raw_text": "original extracted text"
    },
    "alpha_energy": {
      "value_mev": null,
      "uncertainty_mev": null,
      "range_mev": [],
      "ensdf_qualifier": "= | AP | LT | GT | LE | GE | empty",
      "raw_text": "original extracted text"
    },
    "citation_urls": [],
    "dois": [],
    "references": [],
    "parse_warnings": [],
    "review_status": "candidate_only",
    "evidence_level": "string"
  }
}
```

Schema notes:

- JSON `null` is preferred over numeric zero for unknown candidate measurements so missing values cannot be confused with real values.
- `raw_text` must be preserved for auditability.
- `ensdf_qualifier` must be preserved because approximate and limit values need different acceptance rules.
- The current accepted seed structure may stay narrower until an accepted workbook schema exists.

## Validation rules for future implementation

Minimum validation rules before a candidate can be accepted:

1. `schema` must equal `moscovium-statera-go/ensdf-candidate/v1` for candidate files.
2. `source.source_url` must be non-blank HTTPS, and `source.retrieval_status` must be non-blank.
3. If `source.source_path` is present, `content_sha256` must be a lowercase SHA-256 hex string and `content_bytes` must be positive.
4. `candidate.candidate_id` must be unique within the candidate file.
5. `parent_nuclide_id` must parse with the shared nuclide parser.
6. `symbol`, `z`, `a`, and `n` must match `parent_nuclide_id` and `n == a - z`.
7. `element` must match the periodic-table helper for `symbol`.
8. For alpha-decay candidates with a daughter, `daughter` must satisfy shared alpha-daughter validation.
9. `citation_urls` must be non-empty before acceptance; entries must be non-blank HTTPS URLs and duplicates must be rejected.
10. `dois` must be non-empty before acceptance when literature is part of the evidence trail; entries must start with `10.`, contain no whitespace, and duplicates must be rejected.
11. Half-life numeric values, if present, must be positive and must preserve qualifier and raw text.
12. Half-life intervals, if present, must include both lower and upper bounds and must bracket the nominal value.
13. Q-alpha numeric values, if present, must be positive; Q-alpha uncertainty, if present, must be non-negative.
14. Alpha-energy numeric values, if present, must be positive; uncertainty must be non-negative; ranges must have exactly 2 positive ascending bounds.
15. Any `ensdf_qualifier` outside the accepted vocabulary must block acceptance.
16. Candidates with qualifier `AP`, `LT`, `GT`, `LE`, or `GE` must require human review before acceptance even if numeric normalization is possible.
17. `parse_warnings` must be preserved; non-empty warnings should default to `requires_human_review` unless a test documents why acceptance is safe.
18. Accepted records must include a Go test that uses the accepted datum or validation behavior.

## Candidate status model

Recommended status values:

- `candidate_only`: parsed but not reviewed.
- `requires_human_review`: structurally parseable but acceptance cannot be automatic.
- `blocked_source`: source could not be fetched, hashed, cited, or legally stored.
- `blocked_identity`: nuclide identity or daughter relationship is incoherent.
- `blocked_provenance`: citation URL, DOI trail, source path, or evidence level is missing or malformed.
- `accepted`: validated, reviewed, source-backed, and test-covered.

## Statera implementation translation

Recommended next implementation slices:

1. Add duplicate URL and duplicate DOI validation to the existing research seed validator. This is source-neutral and aligns with the candidate validation rules above.
2. Add a new pure-Go `internal/research/ensdf_candidate.go` type with schema constant and validation only; no parser and no scientific values.
3. Add tests for candidate schema rejection, source provenance fields, candidate ID uniqueness, nuclide identity, alpha daughter validation, qualifier vocabulary, and candidate status transitions.
4. Add a candidate fixture with placeholder structural values only if the fixture contains no accepted scientific measurements.
5. Add documentation connecting accepted candidate records to the future Evaluated Isotope Workbook v1.

## Blockers and non-goals

Blocked or intentionally deferred:

- No external ENSDF archive was downloaded in this slice.
- No source PDF or ENSDF record was read for a new isotope value.
- No accepted daughter isotope records were added.
- No model coefficients were audited or changed.
- No parser implementation was added.
- No UI or screenshot changes were made.

Non-goals:

- Do not make Statera a general ENSDF editor.
- Do not import GPL-3.0 source code from Nuclei into this MIT repository.
- Do not normalize approximate or limit ENSDF values into default simulation values without explicit review and tests.
- Do not use ML to accept scientific values; ML may only triage sources or candidates.

## Acceptance checklist for this design slice

- Exact external repositories and URLs reviewed: complete.
- Source-to-candidate-to-accepted workflow documented: complete.
- Minimum Statera ENSDF candidate schema documented: complete.
- At least 5 validation rules documented: complete; 18 rules are listed.
- No new scientific values added to Statera data: complete.
- Expected verification gate: `CGO_ENABLED=0 go test ./...`.
