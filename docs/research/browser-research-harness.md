# Browser Research Harness for Academic Source Exploration

Repo: `/home/xel/git/sages-openclaw/workspace-riju/moscovium-statera-go`

Purpose: define a repeatable browser-assisted workflow for exploring academic nuclear-data websites, DOI landing pages, and source repositories without contaminating Statera with unsupported claims.

This is a harness protocol for Hermes/browser operation. It is not a scraper license, not a bypass mechanism, and not a substitute for reading source content.

## Research boundaries

Allowed targets:

- NNDC NuDat / ENSDF pages and PDFs
- DOI landing pages
- Crossref metadata pages/API
- AME2020 / NUBASE2020 official or journal pages
- IAEA / IUPAC / GSI / JINR source pages
- arXiv or journal pages for peer-reviewed or preprint papers
- GitHub repositories for open-source tooling research

Disallowed behavior:

- bypassing paywalls or access controls
- treating metadata as paper contents
- adding scientific values from snippets, abstracts, or memory
- importing data into simulation defaults before validation tests exist
- copying external code without license review

## Browser run checklist

For each target website or paper:

1. Record target identity:
   - URL
   - access date/time with timezone
   - reason for visit
   - expected source class: evaluated-data, experimental-paper, theoretical-model, review, tool, context

2. Open page with browser automation:
   - use `browser_navigate` for the target URL
   - use `browser_snapshot(full=true)` when text is central
   - use `browser_vision` only to describe layout or visual tables, never as sole scientific extraction evidence
   - use `browser_get_images` only when figures are relevant

3. Capture access evidence:
   - HTTP/visible status if available
   - title
   - DOI
   - venue/publisher
   - year
   - whether abstract, full text, PDF, data table, or only metadata is available

4. Classify the page:
   - `source-content-readable`
   - `metadata-only`
   - `blocked-access`
   - `not-relevant`
   - `candidate-for-manual-review`

5. If source content is readable:
   - record exact quote or table context in a repo doc
   - include units and uncertainty exactly as shown
   - preserve original text for audit
   - do not add to data defaults in the same step unless a TDD validator/test is added

6. If source content is blocked:
   - record the blocker, including exact visible error/status
   - do not infer coefficients or values from metadata
   - pivot to source-neutral work: provenance validation, docs, screenshots, or candidate schema

7. Save findings in repo:
   - source deep dive: `docs/research/<slug>-source-review.md`
   - intake design: `docs/research/<slug>-intake-design.md`
   - citation record: `citations/papers/<key>.md`
   - blocker note: roadmap blocker log or source review doc

## Source review template

```markdown
# Source Review: <short title>

Date: YYYY-MM-DD HH:MM TZ
Reviewer: Riju
Target URL: <url>
DOI: <doi or none>
Source class: evaluated-data | experimental-paper | theoretical-model | review | tool | context
Access status: source-content-readable | metadata-only | blocked-access | not-relevant

## Access evidence

- Browser URL opened:
- Page title:
- Visible status/error:
- Full text available: yes/no
- PDF available: yes/no
- Tables/figures available: yes/no

## Extracted evidence

Only include exact values if source content was readable.

| Datum | Value | Unit | Uncertainty | Source location | Notes |
| --- | --- | --- | --- | --- | --- |

## Candidate use in Statera

- Allowed use:
- Not allowed use:
- Required tests before acceptance:

## Blockers

- Access blocker:
- Data ambiguity:
- License/terms issue:

## Decision

accepted-for-candidate-intake | metadata-only | blocked | rejected
```

## Website-specific notes

### NNDC / ENSDF / NuDat

Use for evaluated-data intake candidates.

Required capture:

- nuclide ID
- dataset page or adopted PDF URL
- half-life text and uncertainty/qualifier
- decay mode
- daughter
- alpha energy or Q-alpha if shown
- retrieval URL
- whether PDF/table was actually readable

Never silently convert missing/approximate/limit values into ordinary central values.

### DOI landing pages

Use to confirm bibliographic identity only unless full article/PDF text is readable.

Metadata can confirm:

- title
- authors
- journal
- year
- DOI
- landing page URL

Metadata cannot confirm:

- formula coefficients
- data-table values
- figure numbers
- uncertainty values
- exact model assumptions

### APS pages

Prior Statera audits found APS DOI/abstract/PDF access returning HTTP `403` for Royer 2008 and Wang 2015. If the browser can now access content, record exact evidence and source location before changing any model code. If still blocked, do not infer values.

### GitHub repositories

Use to learn architecture and workflow patterns only.

Capture:

- repo full name
- URL
- stars/forks if API returns them
- language
- license
- README summary
- exact lesson for Statera
- whether code reuse is disallowed or needs review

## Screenshot handling

If a page or Statera UI screenshot is useful for Juan:

1. Capture screenshot using browser or Statera offscreen command.
2. Save under `screenshots/` only if curated and useful.
3. Verify PNG metrics:
   - width px
   - height px
   - bytes
   - unique colors
   - single color false
4. Add text caption in `screenshots/README.md` or relevant docs.
5. Do not rely on screenshot alone for scientific values.

## Acceptance criteria for a good browser research run

A successful run produces at least one of:

- a source review doc with exact access evidence,
- a blocker note with exact status/error,
- a citation markdown update with DOI/source identity,
- a candidate schema/design update,
- a screenshot with numeric metadata and text caption.

A successful run must not:

- add unsupported scientific values,
- promote candidate data to simulation defaults,
- claim coefficient correctness from metadata,
- omit blocked-source evidence.

## Recommended immediate browser targets

1. NNDC/NuDat pages for accepted `288Mc` and `290Mc` seed URLs.
2. `paceENSDF` documentation/examples for candidate schema details.
3. `Nuclei` documentation for uncertainty-preservation architecture.
4. DOI pages for AME2020 and NUBASE2020.
5. DOI pages and any legitimate accessible full text for Oganessian 2022 Mc production papers.

## Output discipline

Every browser research report to Juan must include:

- target URL count,
- pages readable count,
- metadata-only count,
- blocked count,
- repo docs changed,
- values extracted count,
- values accepted into engine count,
- tests run and pass/fail/skip counts,
- unresolved risks.
