# Decay Validator Intake Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the first testable dual-track slice: stronger decay traversal, a Track A stability validator, and fecim-style citation/lore scaffolding.

**Architecture:** Keep simulation behavior in `internal/physics` and keep contextual/lore material out of runtime defaults. The validator consumes the existing `physics.Catalog` and returns explicit statuses instead of blending unsupported discourse into numeric model outputs. Citation and lore records are plain Markdown so every claim is reviewable with normal git diffs.

**Tech Stack:** Pure Go, standard library only, Markdown citation records, no CGO, no RAG database in this slice.

---

## Scope Boundary

This plan implements the first slice from `docs/superpowers/specs/2026-04-27-dual-track-intake-design.md`.

It deliberately does not implement SQLite-vec, WASM embeddings, Hartree-Fock solvers, or a new gogpu/ui Context tab. Those are separate subsystems and should get their own specs and plans after this slice is merged.

## File Structure

- Modify `internal/physics/decay.go`: keep the public `DecayChain` API, add key/ID validation, path-aware missing-daughter errors, and path-aware cycle errors.
- Modify `internal/physics/decay_test.go`: add regression tests for catalog key mismatches, missing daughter context, and cycle path wording.
- Create `internal/physics/claims.go`: define Track A claim types and validator statuses.
- Create `internal/physics/claims_test.go`: test stability and outside-model claim behavior.
- Create `citations/`: fecim-style source record scaffold for peer-reviewed papers, evaluated sources, and government/media context sources.
- Create `docs/lore/`: quarantined Track B archive for user-submitted discourse and normalized context records.

## Task 1: Harden Decay Traversal

**Files:**
- Modify: `internal/physics/decay.go`
- Modify: `internal/physics/decay_test.go`

- [ ] **Step 1: Add failing traversal regression tests**

Modify the import block in `internal/physics/decay_test.go` to include `strings`:

```go
import (
	"reflect"
	"strings"
	"testing"
	"time"
)
```

Append these tests to `internal/physics/decay_test.go`:

```go
func TestDecayChainRejectsCatalogKeyMismatch(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 290, CitationLink: "seed"},
	}

	_, err := DecayChain("288Mc", catalog)
	if err == nil {
		t.Fatal("DecayChain returned nil error for catalog key mismatch")
	}
	if !strings.Contains(err.Error(), "catalog key 288Mc does not match isotope ID 290Mc") {
		t.Fatalf("error = %q, want catalog key mismatch", err)
	}
}

func TestDecayChainReportsMissingDaughterParent(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, Daughter: "284Nh", CitationLink: "seed"},
	}

	_, err := DecayChain("288Mc", catalog)
	if err == nil {
		t.Fatal("DecayChain returned nil error for missing daughter")
	}
	if !strings.Contains(err.Error(), "daughter 284Nh referenced by 288Mc not found in catalog") {
		t.Fatalf("error = %q, want missing daughter with parent context", err)
	}
}

func TestDecayChainReportsCyclePath(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, Daughter: "284Nh", CitationLink: "seed"},
		"284Nh": {Symbol: "Nh", Z: 113, A: 284, Daughter: "288Mc", CitationLink: "seed"},
	}

	_, err := DecayChain("288Mc", catalog)
	if err == nil {
		t.Fatal("DecayChain returned nil error for cyclic decay graph")
	}
	if !strings.Contains(err.Error(), "cycle detected: 288Mc -> 284Nh -> 288Mc") {
		t.Fatalf("error = %q, want cycle path", err)
	}
}
```

- [ ] **Step 2: Run the failing tests**

Run:

```bash
CGO_ENABLED=0 go test ./internal/physics
```

Expected: FAIL. The key-mismatch and missing-daughter tests fail because current traversal does not validate catalog keys or include parent context. The new cycle path test fails because the current error only says `cycle detected at 288Mc`.

- [ ] **Step 3: Implement path-aware traversal**

Modify `internal/physics/decay.go` so the import block includes `strings`:

```go
import (
	"errors"
	"fmt"
	"strings"
	"time"
)
```

Replace `DecayChain`, `traverse`, and add `formatCycle`:

```go
func DecayChain(start string, catalog Catalog) ([]Isotope, error) {
	if start == "" {
		return nil, errors.New("start isotope is required")
	}

	visited := make(map[string]bool)
	chain := make([]Isotope, 0, 8)
	if err := traverse(start, "", catalog, visited, nil, &chain); err != nil {
		return nil, err
	}
	return chain, nil
}

func traverse(id string, parent string, catalog Catalog, visited map[string]bool, path []string, chain *[]Isotope) error {
	if visited[id] {
		return fmt.Errorf("cycle detected: %s", formatCycle(path, id))
	}

	isotope, ok := catalog[id]
	if !ok {
		if parent != "" {
			return fmt.Errorf("daughter %s referenced by %s not found in catalog", id, parent)
		}
		return fmt.Errorf("isotope %s not found in catalog", id)
	}
	if isotope.ID() != id {
		return fmt.Errorf("catalog key %s does not match isotope ID %s", id, isotope.ID())
	}
	if err := isotope.Validate(); err != nil {
		return err
	}

	visited[id] = true
	path = append(path, id)
	*chain = append(*chain, isotope)
	if isotope.Daughter == "" {
		return nil
	}
	return traverse(isotope.Daughter, id, catalog, visited, path, chain)
}

func formatCycle(path []string, id string) string {
	start := 0
	for index, pathID := range path {
		if pathID == id {
			start = index
			break
		}
	}

	cycle := append([]string{}, path[start:]...)
	cycle = append(cycle, id)
	return strings.Join(cycle, " -> ")
}
```

- [ ] **Step 4: Verify traversal tests pass**

Run:

```bash
CGO_ENABLED=0 go test ./internal/physics
```

Expected: PASS.

- [ ] **Step 5: Commit traversal hardening**

Run:

```bash
git add internal/physics/decay.go internal/physics/decay_test.go
git commit -m "fix: harden decay chain traversal"
```

## Task 2: Add Track A Stability Claim Validator

**Files:**
- Create: `internal/physics/claims.go`
- Create: `internal/physics/claims_test.go`

- [ ] **Step 1: Write failing claim validator tests**

Create `internal/physics/claims_test.go`:

```go
package physics

import (
	"strings"
	"testing"
	"time"
)

func TestValidateClaimRejectsLongLivedKnownShortHalfLife(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, HalfLife: 170 * time.Millisecond, CitationLink: "https://www.nndc.bnl.gov/ensnds/288/Mc/adopted.pdf"},
	}

	result := ValidateClaim(Claim{
		Kind:            ClaimKindMinimumHalfLife,
		IsotopeID:       "288Mc",
		MinimumHalfLife: time.Hour,
	}, catalog)

	if result.Status != ClaimStatusStabilityIncongruent {
		t.Fatalf("status = %q, want %q", result.Status, ClaimStatusStabilityIncongruent)
	}
	if !strings.Contains(result.Reason, "known half-life 170ms is shorter than claimed minimum 1h0m0s") {
		t.Fatalf("reason = %q, want half-life comparison", result.Reason)
	}
	if len(result.Evidence) != 1 || result.Evidence[0] == "" {
		t.Fatalf("evidence = %v, want citation link", result.Evidence)
	}
}

func TestValidateClaimSupportsMinimumHalfLifeWithinKnownData(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, HalfLife: 170 * time.Millisecond, CitationLink: "https://www.nndc.bnl.gov/ensnds/288/Mc/adopted.pdf"},
	}

	result := ValidateClaim(Claim{
		Kind:            ClaimKindMinimumHalfLife,
		IsotopeID:       "288Mc",
		MinimumHalfLife: 100 * time.Millisecond,
	}, catalog)

	if result.Status != ClaimStatusSupportedByTrackA {
		t.Fatalf("status = %q, want %q", result.Status, ClaimStatusSupportedByTrackA)
	}
}

func TestValidateClaimReturnsInsufficientDataForUnknownIsotope(t *testing.T) {
	result := ValidateClaim(Claim{
		Kind:            ClaimKindMinimumHalfLife,
		IsotopeID:       "299Mc",
		MinimumHalfLife: time.Second,
	}, Catalog{})

	if result.Status != ClaimStatusInsufficientData {
		t.Fatalf("status = %q, want %q", result.Status, ClaimStatusInsufficientData)
	}
}

func TestValidateClaimRejectsUnsupportedMechanism(t *testing.T) {
	result := ValidateClaim(Claim{
		Kind:      ClaimKindMechanism,
		Mechanism: "antigravity propulsion",
	}, Catalog{})

	if result.Status != ClaimStatusOutsideSupportedModel {
		t.Fatalf("status = %q, want %q", result.Status, ClaimStatusOutsideSupportedModel)
	}
	if !strings.Contains(result.Reason, "outside Standard Model-compatible nuclear physics") {
		t.Fatalf("reason = %q, want model-boundary explanation", result.Reason)
	}
}

func TestValidateClaimRejectsMalformedMinimumHalfLifeClaim(t *testing.T) {
	result := ValidateClaim(Claim{
		Kind:            ClaimKindMinimumHalfLife,
		IsotopeID:       "",
		MinimumHalfLife: time.Second,
	}, Catalog{})

	if result.Status != ClaimStatusInvalidClaim {
		t.Fatalf("status = %q, want %q", result.Status, ClaimStatusInvalidClaim)
	}
}
```

- [ ] **Step 2: Run the failing tests**

Run:

```bash
CGO_ENABLED=0 go test ./internal/physics
```

Expected: FAIL with undefined `ValidateClaim`, `Claim`, and status constants.

- [ ] **Step 3: Implement the claim validator**

Create `internal/physics/claims.go`:

```go
package physics

import (
	"fmt"
	"time"
)

type ClaimKind string

const (
	ClaimKindMinimumHalfLife ClaimKind = "minimum-half-life"
	ClaimKindMechanism       ClaimKind = "mechanism"
)

type ClaimStatus string

const (
	ClaimStatusSupportedByTrackA     ClaimStatus = "supported-by-track-a"
	ClaimStatusStabilityIncongruent  ClaimStatus = "stability-incongruent"
	ClaimStatusOutsideSupportedModel ClaimStatus = "outside-supported-model"
	ClaimStatusInsufficientData      ClaimStatus = "insufficient-data"
	ClaimStatusInvalidClaim          ClaimStatus = "invalid-claim"
)

type Claim struct {
	Kind            ClaimKind
	IsotopeID       string
	MinimumHalfLife time.Duration
	Mechanism       string
}

type ClaimResult struct {
	Status   ClaimStatus
	Reason   string
	Evidence []string
}

func ValidateClaim(claim Claim, catalog Catalog) ClaimResult {
	switch claim.Kind {
	case ClaimKindMinimumHalfLife:
		return validateMinimumHalfLifeClaim(claim, catalog)
	case ClaimKindMechanism:
		return validateMechanismClaim(claim)
	default:
		return ClaimResult{
			Status: ClaimStatusInvalidClaim,
			Reason: fmt.Sprintf("unknown claim kind %q", claim.Kind),
		}
	}
}

func validateMinimumHalfLifeClaim(claim Claim, catalog Catalog) ClaimResult {
	if claim.IsotopeID == "" {
		return ClaimResult{Status: ClaimStatusInvalidClaim, Reason: "minimum half-life claim requires isotope ID"}
	}
	if claim.MinimumHalfLife <= 0 {
		return ClaimResult{Status: ClaimStatusInvalidClaim, Reason: "minimum half-life claim requires positive duration"}
	}

	isotope, ok := catalog[claim.IsotopeID]
	if !ok {
		return ClaimResult{
			Status: ClaimStatusInsufficientData,
			Reason: fmt.Sprintf("isotope %s is not present in Track A catalog", claim.IsotopeID),
		}
	}
	if err := isotope.Validate(); err != nil {
		return ClaimResult{Status: ClaimStatusInvalidClaim, Reason: err.Error()}
	}
	if isotope.HalfLife <= 0 {
		return ClaimResult{
			Status:   ClaimStatusInsufficientData,
			Reason:   fmt.Sprintf("isotope %s has no evaluated half-life in Track A catalog", claim.IsotopeID),
			Evidence: []string{isotope.CitationLink},
		}
	}

	if isotope.HalfLife < claim.MinimumHalfLife {
		return ClaimResult{
			Status:   ClaimStatusStabilityIncongruent,
			Reason:   fmt.Sprintf("%s known half-life %s is shorter than claimed minimum %s", claim.IsotopeID, isotope.HalfLife, claim.MinimumHalfLife),
			Evidence: []string{isotope.CitationLink},
		}
	}

	return ClaimResult{
		Status:   ClaimStatusSupportedByTrackA,
		Reason:   fmt.Sprintf("%s known half-life %s meets claimed minimum %s", claim.IsotopeID, isotope.HalfLife, claim.MinimumHalfLife),
		Evidence: []string{isotope.CitationLink},
	}
}

func validateMechanismClaim(claim Claim) ClaimResult {
	if claim.Mechanism == "" {
		return ClaimResult{Status: ClaimStatusInvalidClaim, Reason: "mechanism claim requires mechanism text"}
	}
	return ClaimResult{
		Status: ClaimStatusOutsideSupportedModel,
		Reason: fmt.Sprintf("%q is outside Standard Model-compatible nuclear physics implemented by this project", claim.Mechanism),
	}
}
```

- [ ] **Step 4: Verify claim tests pass**

Run:

```bash
CGO_ENABLED=0 go test ./internal/physics
```

Expected: PASS.

- [ ] **Step 5: Commit claim validator**

Run:

```bash
git add internal/physics/claims.go internal/physics/claims_test.go
git commit -m "feat: add track a claim validator"
```

## Task 3: Add Citation And Lore Scaffolding

**Files:**
- Create: `citations/README.md`
- Create: `citations/TEMPLATE.md`
- Create: `citations/facts.md`
- Create: `citations/disputed.md`
- Create: `citations/refs.bib`
- Create: `citations/queue/to-fetch.md`
- Create: `citations/queue/to-read.md`
- Create: `citations/queue/to-extract.md`
- Create: `citations/pdfs/README.md`
- Create: `citations/pdfs/.gitignore`
- Create: `citations/papers/.gitkeep`
- Create: `citations/reports/.gitkeep`
- Create: `docs/lore/README.md`
- Create: `docs/lore/intake-notes/2026-04-27-user-submitted-element-115-personnel-discourse.md`
- Create: `docs/lore/records/README.md`

- [ ] **Step 1: Create citation directory files**

Use `apply_patch` to create `citations/README.md`:

```markdown
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
```

Use `apply_patch` to create `citations/TEMPLATE.md`:

```markdown
# {Title}

**Key:** `{firstauthorYEARkeyword}`
**DOI:** `{doi-or-none}`
**arXiv:** `{arxiv-id-or-none}`
**URL:** `{stable-url-or-none}`
**Year:** `{year}`
**Venue:** `{journal-report-or-institution}`
**Authors:** `{author-list-or-institution}`
**Tags:** `{track-a | track-b} {topic-tags}`
**Status:** `{to-read | skimmed | read | deep-read}`
**PDF:** `{pdfs/key.pdf | external URL | not stored}`
**Added:** `{YYYY-MM-DD}`

---

## TL;DR

One sentence describing what this source contributes.

## Relevance

Explain whether this source supports Track A physics, Track B context, or a disputed claim.

## Key Facts

Only include facts checked against the source.

## Limitations

Record author-stated limitations and project-relevant limits.

## Cited In

- [ ] `{file:line}` - `{claim supported}`

## Related Sources

- `{key}` - `{relationship}`

## BibTeX

```bibtex
% Add a verified BibTeX entry before using this source as a public citation.
```
```

Use `apply_patch` to create `citations/facts.md`:

```markdown
# Citable Facts Database

**Last updated:** 2026-04-27
**Total verified facts:** 0

This file is the cross-source index for facts used by Moscovium Statera Go. It starts empty on purpose. Add facts only after reading the source and recording the same fact in `citations/papers/{key}.md`.

## Rules

- Preserve exact values and units from the source.
- Include source location such as page, section, table, or figure.
- Include experimental or simulation conditions.
- Record evidence level: evaluated data, peer-reviewed, preprint, government document, institution statement, mainstream reporting, or other.
- Track B context facts must be marked `not-for-simulation`.

## Index

1. [Moscovium Isotope Data](#moscovium-isotope-data)
2. [Discovery And Naming](#discovery-and-naming)
3. [Decay Chains](#decay-chains)
4. [Theoretical Stability](#theoretical-stability)
5. [Contextual Events](#contextual-events)

## Moscovium Isotope Data

No verified facts recorded yet.

## Discovery And Naming

No verified facts recorded yet.

## Decay Chains

No verified facts recorded yet.

## Theoretical Stability

No verified facts recorded yet.

## Contextual Events

No verified facts recorded yet.
```

Use `apply_patch` to create `citations/disputed.md`:

```markdown
# Disputed And Weak Claims

This file tracks claims that are contested, weakly sourced, or easy to overstate.

Do not use this file to smuggle unverified claims into the project. Use it to make uncertainty explicit.

## Entry Format

```markdown
## {Claim}

**Status:** `{contested | weak evidence | source mismatch | unresolved}`

| Position | Source | Evidence Level | Notes |
|---|---|---|---|
| Supports | `{key}` | `{peer-reviewed | government document | mainstream reporting | other}` | `{specific location and caveat}` |
| Contradicts | `{key}` | `{peer-reviewed | government document | mainstream reporting | other}` | `{specific location and caveat}` |

**Project handling:** `{how the repository labels or avoids this claim}`
**Next action:** `{source to read, calculation to run, or wording to change}`
```

## Active Entries

No disputed claims recorded yet.
```

Use `apply_patch` to create `citations/refs.bib`:

```bibtex
% Compiled bibliography for Moscovium Statera Go.
% Add verified BibTeX entries from citations/papers/*.md.
```

Use `apply_patch` to create these queue files:

`citations/queue/to-fetch.md`

```markdown
# Sources To Fetch

- Oganessian et al. 2022, "First experiment at the Super Heavy Element Factory: High cross section of 288Mc in the 243Am+48Ca reaction and identification of the new isotope 264Lr", DOI: 10.1103/PhysRevC.106.L031301
- IUPAC 2016, "Names and symbols of the elements with atomic numbers 113, 115, 117 and 118", DOI: 10.1515/pac-2016-0501
- House Oversight Committee letter to FBI Director, April 20, 2026, URL: https://oversight.house.gov/wp-content/uploads/2026/04/FBI-Missing-Scientists-Letter_4.20.26.pdf
```

`citations/queue/to-read.md`

```markdown
# Sources To Read

No sources queued for reading yet.
```

`citations/queue/to-extract.md`

```markdown
# Sources To Extract

No sources queued for fact extraction yet.
```

Use `apply_patch` to create `citations/pdfs/README.md`:

```markdown
# Local PDFs

Store only open-access, accepted-manuscript, public-domain, or otherwise legally archivable PDFs here.

Do not commit paywalled publisher PDFs. For paywalled papers, store metadata and DOI links in `citations/papers/`.
```

Use `apply_patch` to create `citations/pdfs/.gitignore`:

```gitignore
*.pdf
!README.md
```

Use `apply_patch` to create empty keep files:

```text
citations/papers/.gitkeep
citations/reports/.gitkeep
```

- [ ] **Step 2: Create lore archive files**

Use `apply_patch` to create `docs/lore/README.md`:

```markdown
# Lore And Context Archive

This directory preserves Track B context: public discourse, media claims, government oversight documents, and user-submitted narratives that explain why people ask about Element 115.

Track B is not simulation data. Records here must never seed `data/research.seed.json`, physics constants, or model defaults.

## Labels

- `verified-event`: the event is supported by a reliable source.
- `media-claim`: a source reports or repeats a claim without proving it.
- `unsupported-linkage`: a claim connects facts without sufficient evidence.
- `physics-claim`: a claim can be expressed as a physics assertion for validation.
- `not-for-simulation`: the record is excluded from simulation defaults.

## Layout

| Path | Purpose |
|---|---|
| `intake-notes/` | Quarantined user-submitted material awaiting source review |
| `records/` | Normalized Track B records after source review |
```

Use `apply_patch` to create `docs/lore/intake-notes/2026-04-27-user-submitted-element-115-personnel-discourse.md`:

```markdown
# User-Submitted Element 115 Personnel Discourse Intake

**Date captured:** 2026-04-27
**Status:** quarantined user submission
**Labels:** `media-claim`, `unsupported-linkage`, `physics-claim`, `not-for-simulation`

## Boundary Notice

This file preserves a user-submitted narrative for source review. It is not evidence, not evaluated nuclear data, and not a source for simulation defaults. Named-person claims must be independently verified before they appear in normalized Track B records.

## Preserved Claim Inventory

- The submission argues that scientific authority and alternative historical narratives now include a counter-expert class of credentialed dissidents.
- The submission asks the repository to save and research claims involving Element 115, public discourse, missing or deceased scientific personnel, Apollo skepticism, September 11 skepticism, and COVID-19 response skepticism.
- The submission links Element 115 discourse to alleged disappearances or deaths of personnel connected in public reporting to aerospace, nuclear physics, materials science, astronomy, or government research.
- The submission specifically asks for a dual-track architecture: Track A for evaluated nuclear physics and Track B for lore, media claims, and contextual events.
- The submission proposes a validator that tests lore-derived physics claims against source-backed Moscovium isotope data.

## Handling Decision

The repository will record discourse about claims rather than asserting unsupported claims. Track A remains limited to evaluated nuclear data and peer-reviewed physics. Track B may preserve contextual records only when each record carries source metadata and `not-for-simulation`.
```

Use `apply_patch` to create `docs/lore/records/README.md`:

```markdown
# Normalized Context Records

Records in this directory are reviewed Track B entries. Each record must cite sources, distinguish verified events from media claims, and include `not-for-simulation`.

Do not add user-submitted claims directly here. Move material from `docs/lore/intake-notes/` only after source review.
```

- [ ] **Step 3: Verify Markdown scaffold is present**

Run:

```bash
find citations docs/lore -maxdepth 3 -type f | sort
```

Expected output includes all files listed in this task.

- [ ] **Step 4: Commit citation and lore scaffold**

Run:

```bash
git add citations docs/lore
git commit -m "docs: scaffold citation and lore intake"
```

## Task 4: Add Seed Source Records

**Files:**
- Create: `citations/papers/oganessian2022mcfactory.md`
- Create: `citations/papers/iupac2016names.md`
- Create: `citations/papers/houseoversight2026missing-scientists-letter.md`
- Create: `docs/lore/records/mccasland-2026-missing-person.md`
- Create: `docs/lore/records/grillmair-2026-homicide.md`

- [ ] **Step 1: Add Moscovium Track A paper record**

Use `apply_patch` to create `citations/papers/oganessian2022mcfactory.md`:

```markdown
# First experiment at the Super Heavy Element Factory: High cross section of 288Mc in the 243Am+48Ca reaction and identification of the new isotope 264Lr

**Key:** `oganessian2022mcfactory`
**DOI:** `10.1103/PhysRevC.106.L031301`
**arXiv:** `none`
**URL:** `https://journals.aps.org/prc/abstract/10.1103/PhysRevC.106.L031301`
**Year:** `2022`
**Venue:** `Physical Review C`
**Authors:** `Yu. Ts. Oganessian; V. K. Utyonkov; N. D. Kovrizhnykh; F. Sh. Abdullin; S. N. Dmitriev; D. Ibadullayev; M. G. Itkis; D. A. Kuznetsov; O. V. Petrushkin; A. V. Podshibiakin; A. N. Polyakov; A. G. Popeko; R. N. Sagaidak; L. Schlattauer; I. V. Shirokovski; V. D. Shubin; M. V. Shumeiko; D. I. Solovyev; Yu. S. Tsyganov; A. A. Voinov; V. G. Subbotin; A. Yu. Bodrov; A. V. Sabel'nikov; A. V. Khalkin; V. B. Zlokazov; K. P. Rykaczewski; T. T. King; J. B. Roberto; N. T. Brewer; R. K. Grzywacz; Z. G. Gan; Z. Y. Zhang; M. H. Huang; H. B. Yang`
**Tags:** `track-a moscovium 243Am 48Ca decay-chain`
**Status:** `to-read`
**PDF:** `not stored`
**Added:** `2026-04-27`

---

## TL;DR

Reports DGFRS-2 production and decay-chain observations for Moscovium isotopes in the 243Am+48Ca reaction.

## Relevance

This is Track A source material for Moscovium isotope production, decay chains, and cross-section context. Values must not be promoted to `citations/facts.md` or `data/research.seed.json` until checked against the paper.

## Key Facts

No extracted facts yet.

## Limitations

The record is metadata-only until read.

## Cited In

- [ ] `data/research.seed.json` - DOI trail for 288Mc seed record

## Related Sources

- `iupac2016names` - element naming and recognition context

## BibTeX

```bibtex
@article{oganessian2022mcfactory,
  title = {First experiment at the Super Heavy Element Factory: High cross section of {288Mc} in the {243Am}+{48Ca} reaction and identification of the new isotope {264Lr}},
  author = {Oganessian, Yu. Ts. and Utyonkov, V. K. and Kovrizhnykh, N. D. and Abdullin, F. Sh. and Dmitriev, S. N. and Ibadullayev, D. and Itkis, M. G. and Kuznetsov, D. A. and Petrushkin, O. V. and Podshibiakin, A. V. and Polyakov, A. N. and Popeko, A. G. and Sagaidak, R. N. and Schlattauer, L. and Shirokovski, I. V. and Shubin, V. D. and Shumeiko, M. V. and Solovyev, D. I. and Tsyganov, Yu. S. and Voinov, A. A. and Subbotin, V. G. and Bodrov, A. Yu. and Sabel'nikov, A. V. and Khalkin, A. V. and Zlokazov, V. B. and Rykaczewski, K. P. and King, T. T. and Roberto, J. B. and Brewer, N. T. and Grzywacz, R. K. and Gan, Z. G. and Zhang, Z. Y. and Huang, M. H. and Yang, H. B.},
  journal = {Physical Review C},
  volume = {106},
  pages = {L031301},
  year = {2022},
  doi = {10.1103/PhysRevC.106.L031301}
}
```
```

- [ ] **Step 2: Add IUPAC naming source record**

Use `apply_patch` to create `citations/papers/iupac2016names.md`:

```markdown
# Names and symbols of the elements with atomic numbers 113, 115, 117 and 118

**Key:** `iupac2016names`
**DOI:** `10.1515/pac-2016-0501`
**arXiv:** `none`
**URL:** `https://www.degruyterbrill.com/document/doi/10.1515/pac-2016-0501/pdf`
**Year:** `2016`
**Venue:** `Pure and Applied Chemistry`
**Authors:** `IUPAC Inorganic Chemistry Division`
**Tags:** `track-a iupac naming moscovium`
**Status:** `to-read`
**PDF:** `not stored`
**Added:** `2026-04-27`

---

## TL;DR

IUPAC recommendation assigning the name moscovium and symbol Mc to element 115.

## Relevance

This source supports nomenclature and standards-body provenance for element 115.

## Key Facts

No extracted facts yet.

## Limitations

This is a nomenclature source, not an isotope-data source.

## Cited In

- [ ] `README.md` - project description and element naming

## Related Sources

- `oganessian2022mcfactory` - isotope production and decay-chain source

## BibTeX

```bibtex
@article{iupac2016names,
  title = {Names and symbols of the elements with atomic numbers 113, 115, 117 and 118},
  author = {{IUPAC Inorganic Chemistry Division}},
  journal = {Pure and Applied Chemistry},
  year = {2016},
  doi = {10.1515/pac-2016-0501}
}
```
```

- [ ] **Step 3: Add House Oversight Track B source record**

Use `apply_patch` to create `citations/papers/houseoversight2026missing-scientists-letter.md`:

```markdown
# April 20, 2026 House Oversight letter to FBI Director on missing and deceased scientists

**Key:** `houseoversight2026missingScientistsLetter`
**DOI:** `none`
**arXiv:** `none`
**URL:** `https://oversight.house.gov/wp-content/uploads/2026/04/FBI-Missing-Scientists-Letter_4.20.26.pdf`
**Year:** `2026`
**Venue:** `U.S. House Committee on Oversight and Government Reform correspondence`
**Authors:** `James Comer; Eric Burlison`
**Tags:** `track-b government-document missing-persons media-claims not-for-simulation`
**Status:** `to-read`
**PDF:** `not stored`
**Added:** `2026-04-27`

---

## TL;DR

Congressional correspondence requesting a briefing on unconfirmed public reporting about missing or deceased individuals connected in reports to sensitive scientific information.

## Relevance

This is Track B context. It may support claims about the existence of a congressional inquiry, but it does not validate Element 115, UAP, or propulsion linkages.

## Key Facts

No extracted facts yet.

## Limitations

The letter explicitly frames the matter as unconfirmed public reporting. This record is `not-for-simulation`.

## Cited In

- [ ] `docs/lore/records/mccasland-2026-missing-person.md` - contextual source trail

## Related Sources

- `mccasland-2026-missing-person` - normalized Track B event record

## BibTeX

```bibtex
@misc{houseoversight2026missingScientistsLetter,
  title = {Letter to FBI Director Kash Patel on missing and deceased scientists},
  author = {Comer, James and Burlison, Eric},
  year = {2026},
  month = apr,
  url = {https://oversight.house.gov/wp-content/uploads/2026/04/FBI-Missing-Scientists-Letter_4.20.26.pdf}
}
```
```

- [ ] **Step 4: Add normalized McCasland context record**

Use `apply_patch` to create `docs/lore/records/mccasland-2026-missing-person.md`:

```markdown
# William Neil McCasland Missing-Person Context Record

**Date added:** 2026-04-27
**Event date:** 2026-02-27
**Labels:** `verified-event`, `media-claim`, `unsupported-linkage`, `not-for-simulation`
**Simulation use:** prohibited

## Verified Event

Reliable reporting attributed to the Bernalillo County Sheriff's Office states that retired U.S. Air Force Maj. Gen. William Neil McCasland was last seen at or near his Albuquerque home on February 27, 2026, and that a Silver Alert was issued.

## Unsupported Linkages

Public discourse has connected McCasland's disappearance to UAP or Element 115 narratives. This repository does not treat those linkages as verified. They remain `unsupported-linkage` unless supported by reliable sources.

## Sources

- `houseoversight2026missingScientistsLetter` - congressional request for briefing, explicitly based on unconfirmed public reporting.
- CNN/KVIA report, March 11, 2026: https://kvia.com/news/top-stories/2026/03/11/fbi-joins-search-for-retired-air-force-major-general-missing-for-nearly-2-weeks-2/

## Project Handling

This record may appear in a future Context tab. It must not feed `physics.Catalog`, `data/research.seed.json`, or stability calculations.
```

- [ ] **Step 5: Add normalized Grillmair context record**

Use `apply_patch` to create `docs/lore/records/grillmair-2026-homicide.md`:

```markdown
# Carl Grillmair Homicide Context Record

**Date added:** 2026-04-27
**Event date:** 2026-02-16
**Labels:** `verified-event`, `media-claim`, `unsupported-linkage`, `not-for-simulation`
**Simulation use:** prohibited

## Verified Event

Caltech reported that Carl Grillmair, an astronomer at Caltech's IPAC science and data center, died on February 16, 2026. Reliable news reporting states that his death was ruled a homicide and that a suspect was charged.

## Unsupported Linkages

Claims connecting Grillmair's death to Element 115, UAP, or hidden-technology narratives are not validated by this record. They remain `unsupported-linkage`.

## Sources

- Caltech memorial notice: https://www.astro.caltech.edu/news/caltech-mourns-the-passing-of-carl-grillmair-19592026
- Los Angeles Times report, February 19, 2026: https://www.latimes.com/california/story/2026-02-19/caltech-astrophysicist-fatally-shot-on-porch-in-antelope-valley

## Project Handling

This record may appear in a future Context tab. It must not feed `physics.Catalog`, `data/research.seed.json`, or stability calculations.
```

- [ ] **Step 6: Commit source records**

Run:

```bash
git add citations/papers docs/lore/records
git commit -m "docs: add initial moscovium source records"
```

## Task 5: Download Legally Archivable PDFs

**Files:**
- Create if downloadable: `citations/pdfs/iupac2016names.pdf`
- Create if downloadable: `citations/pdfs/houseoversight2026missing-scientists-letter.pdf`
- Create if downloadable: `citations/pdfs/oganessian2022mcfactory-accepted.pdf`

- [ ] **Step 1: Download open/public PDFs**

Run:

```bash
curl -L -o citations/pdfs/iupac2016names.pdf "https://www.degruyterbrill.com/document/doi/10.1515/pac-2016-0501/pdf"
curl -L -o citations/pdfs/houseoversight2026missing-scientists-letter.pdf "https://oversight.house.gov/wp-content/uploads/2026/04/FBI-Missing-Scientists-Letter_4.20.26.pdf"
curl -L -o citations/pdfs/oganessian2022mcfactory-accepted.pdf "https://link.aps.org/accepted/10.1103/PhysRevC.106.L031301"
```

- [ ] **Step 2: Verify downloads are PDFs**

Run:

```bash
file citations/pdfs/*.pdf
```

Expected: each downloaded file reports `PDF document`. If any file reports HTML or text, delete non-PDF downloads and keep only the metadata record:

```bash
for pdf in citations/pdfs/*.pdf; do
  if [ "$(file -b --mime-type "$pdf")" != "application/pdf" ]; then
    rm "$pdf"
  fi
done
```

- [ ] **Step 3: Do not commit PDF files**

Run:

```bash
git status --short citations/pdfs
```

Expected: no PDF files are staged because `citations/pdfs/.gitignore` ignores `*.pdf`. `citations/pdfs/README.md` and `.gitignore` are already committed from Task 3.

## Task 6: Full Verification

**Files:**
- No new files.

- [ ] **Step 1: Run Go tests**

Run:

```bash
CGO_ENABLED=0 go test ./...
```

Expected: PASS.

- [ ] **Step 2: Run CLI smoke path**

Run:

```bash
go run ./cmd/statera
```

Expected output:

```text
288Mc -> 284Nh -> 280Rg -> 276Mt -> 272Bh -> 268Db -> 264Lr
```

- [ ] **Step 3: Run UI smoke path**

Run:

```bash
CGO_ENABLED=0 go run ./cmd/statera-ui
```

Expected: app starts without compilation errors. If local graphics or windowing constraints block execution, capture the exact error in the final summary.

- [ ] **Step 4: Check whitespace**

Run:

```bash
git diff --check
```

Expected: no output.

- [ ] **Step 5: Final status check**

Run:

```bash
git status --short
```

Expected: clean except ignored PDFs in `citations/pdfs/` if downloaded.

## Self-Review Notes

- Spec coverage: this plan implements the first implementation slice only: traversal, stability validator, citation/lore scaffold, and seed records. RAG, solvers, and UI expansion are intentionally out of scope.
- Placeholder scan: source templates contain brace-delimited fields only inside `citations/TEMPLATE.md`, where they are the intended template syntax.
- Type consistency: `Claim`, `ClaimKind`, `ClaimStatus`, and `ValidateClaim` names match across tests and implementation steps.
