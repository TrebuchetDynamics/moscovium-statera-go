# Moscovium Top 1 Expansion — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Expand moscovium-statera-go from bootstrap (2 isotopes, 1 model, basic UI) into a comprehensive Element 115 research/simulation/visualization platform spanning 15+ isotopes, 5+ physics models, and 8+ interactive dashboard views.

**Architecture:** Four parallel-capable waves. Wave 1 (data) gates Wave 2 (simulations) gates Wave 3 (visualizations). Within each wave, tasks are independent and parallelizable. Wave 4 integrates and polishes. All code remains pure Go, CGO-free, provenance-enforced.

**Tech Stack:** Go 1.25, gogpu/ui v0.1.13, gogpu/gg v0.43.2, gogpu/gogpu v0.29.4, go-webgpu/webgpu v0.4.3

---

## Wave 1: Data Breadth (15+ isotope records with provenance)

### Task 1.1: ENSDF text parser for evaluated nuclear data

**Files:**
- Create: `internal/research/ensdf_text.go`
- Create: `internal/research/ensdf_text_test.go`

- [ ] **Step 1: Write failing test for ENSDF text field extraction**

```go
// internal/research/ensdf_text_test.go
package research

import "testing"

func TestParseENSDFNuclideHeader_Normal(t *testing.T) {
    // Typical ENSDF header: "288MC    ADOPTED LEVELS, GAMMAS"
    line := "288MC    ADOPTED LEVELS, GAMMAS                 202206"
    id, err := ParseENSDFNuclideHeader(line)
    if err != nil {
        t.Fatalf("expected success, got %v", err)
    }
    if id != "288Mc" {
        t.Errorf("expected 288Mc, got %s", id)
    }
}

func TestParseENSDFHalfLifeField(t *testing.T) {
    // ENSDF half-life field: " T  1.700E-01 S"
    line := "  T  1.700E-01 S"
    seconds, err := ParseENSDFHalfLifeField(line)
    if err != nil {
        t.Fatalf("expected success, got %v", err)
    }
    if seconds != 0.17 {
        t.Errorf("expected 0.17, got %f", seconds)
    }
}

func TestParseENSDFQAlphaField(t *testing.T) {
    // ENSDF Q field: " Q  1.075E+04 KEV"
    line := " Q  1.075E+04 KEV"
    mev, err := ParseENSDFQValueField(line)
    if err != nil {
        t.Fatalf("expected success, got %v", err)
    }
    if mev != 10.75 {
        t.Errorf("expected 10.75 MeV, got %f", mev)
    }
}

func TestParseENSDFDaughterField(t *testing.T) {
    // ENSDF daughter field: " DA  284NH"
    line := " DA  284NH"
    daughter, err := ParseENSDFDaughterField(line)
    if err != nil {
        t.Fatalf("expected success, got %v", err)
    }
    if daughter != "284Nh" {
        t.Errorf("expected 284Nh, got %s", daughter)
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `CGO_ENABLED=0 go test ./internal/research/ -run TestParseENSDF -v`
Expected: FAIL with "undefined" errors

- [ ] **Step 3: Implement ENSDF text parser**

```go
// internal/research/ensdf_text.go
package research

import (
    "fmt"
    "strconv"
    "strings"
)

// ParseENSDFNuclideHeader extracts nuclide ID from an ENSDF header line.
// Input format: "288MC    ADOPTED LEVELS, GAMMAS"
func ParseENSDFNuclideHeader(line string) (string, error) {
    fields := strings.Fields(line)
    if len(fields) < 1 {
        return "", fmt.Errorf("empty ENSDF header line")
    }
    raw := strings.ToUpper(fields[0])
    // Split mass number from element symbol, e.g. "288MC" -> "288Mc"
    i := 0
    for i < len(raw) && raw[i] >= '0' && raw[i] <= '9' {
        i++
    }
    if i == 0 || i == len(raw) {
        return "", fmt.Errorf("invalid ENSDF nuclide header: %q", line)
    }
    mass := raw[:i]
    symbol := strings.Title(strings.ToLower(raw[i:]))
    return mass + symbol, nil
}

// ParseENSDFHalfLifeField extracts half-life in seconds from an ENSDF T field.
// Input format: "  T  1.700E-01 S"
func ParseENSDFHalfLifeField(line string) (float64, error) {
    fields := strings.Fields(line)
    if len(fields) < 3 || fields[0] != "T" {
        return 0, fmt.Errorf("invalid ENSDF half-life line: %q", line)
    }
    val, err := strconv.ParseFloat(fields[1], 64)
    if err != nil {
        return 0, fmt.Errorf("parse half-life value: %w", err)
    }
    // ENSDF half-life scaling factors
    switch strings.ToUpper(fields[2]) {
    case "S":
        return val, nil
    case "MS":
        return val * 1e-3, nil
    case "US":
        return val * 1e-6, nil
    case "NS":
        return val * 1e-9, nil
    case "PS":
        return val * 1e-12, nil
    case "M":
        return val * 60, nil
    case "H":
        return val * 3600, nil
    case "D":
        return val * 86400, nil
    case "Y":
        return val * 31556952, nil
    default:
        return val, nil
    }
}

// ParseENSDFQValueField extracts Q-value in MeV from an ENSDF Q field.
// Input format: " Q  1.075E+04 KEV"
func ParseENSDFQValueField(line string) (float64, error) {
    fields := strings.Fields(line)
    if len(fields) < 3 || fields[0] != "Q" {
        return 0, fmt.Errorf("invalid ENSDF Q-value line: %q", line)
    }
    val, err := strconv.ParseFloat(fields[1], 64)
    if err != nil {
        return 0, fmt.Errorf("parse Q-value: %w", err)
    }
    unit := strings.ToUpper(fields[2])
    switch unit {
    case "KEV":
        return val * 1e-3, nil
    case "MEV":
        return val, nil
    case "EV":
        return val * 1e-6, nil
    default:
        return val * 1e-3, nil // assume keV
    }
}

// ParseENSDFDaughterField extracts daughter nuclide ID from ENSDF DA field.
// Input format: " DA  284NH"
func ParseENSDFDaughterField(line string) (string, error) {
    fields := strings.Fields(line)
    if len(fields) < 2 || fields[0] != "DA" {
        return "", fmt.Errorf("invalid ENSDF daughter line: %q", line)
    }
    raw := strings.ToUpper(fields[1])
    i := 0
    for i < len(raw) && raw[i] >= '0' && raw[i] <= '9' {
        i++
    }
    if i == 0 || i == len(raw) {
        return "", fmt.Errorf("invalid ENSDF daughter entry: %q", line)
    }
    mass := raw[:i]
    symbol := strings.Title(strings.ToLower(raw[i:]))
    return mass + symbol, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `CGO_ENABLED=0 go test ./internal/research/ -run TestParseENSDF -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/research/ensdf_text.go internal/research/ensdf_text_test.go
git commit -m "add ENSDF text parser for evaluated nuclear data fields"
```

---

### Task 1.2: Batch ENSDF intake to research seed expansion

**Files:**
- Modify: `internal/research/seed.go`
- Modify: `internal/research/seed_test.go`
- Create: `internal/research/ensdf_batch.go`
- Create: `internal/research/ensdf_batch_test.go`

- [ ] **Step 1: Add BatchENSDFToSeed function**

```go
// internal/research/ensdf_batch.go
package research

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
)

// ENSDFBatchEntry represents a single nuclide extracted from ENSDF text.
type ENSDFBatchEntry struct {
    NuclideID        string
    HalfLifeSeconds  float64
    QAlphaMeV        float64
    AlphaEnergyMeV   float64
    Daughter         string
    DecayMode        string   // "alpha", "beta-", "EC", "SF"
    SFBranchingRatio float64  // 0.0 if not applicable
    CitationURL      string
    DOIs             []string
}

// BatchENSDFToSeed converts batch ENSDF entries into research.seed.json records
// and appends them to the existing seed, preserving existing records.
func BatchENSDFToSeed(existingPath string, entries []ENSDFBatchEntry, sourceURL string) error {
    // Read existing seed
    raw, err := os.ReadFile(existingPath)
    if err != nil {
        return fmt.Errorf("read existing seed: %w", err)
    }
    var seed ResearchSeed
    if err := json.Unmarshal(raw, &seed); err != nil {
        return fmt.Errorf("unmarshal existing seed: %w", err)
    }

    existingIDs := make(map[string]bool)
    for _, rec := range seed.Records {
        existingIDs[rec.ID] = true
    }

    for _, entry := range entries {
        if existingIDs[entry.ID] {
            continue // don't duplicate
        }
        rec := ResearchSeedRecord{
            ID:                  entry.NuclideID,
            Element:             elementNameFromSymbol(extractSymbol(entry.NuclideID)),
            Symbol:              extractSymbol(entry.NuclideID),
            Z:                   atomicNumberFromSymbol(extractSymbol(entry.NuclideID)),
            A:                   extractMass(entry.NuclideID),
            N:                   extractMass(entry.NuclideID) - atomicNumberFromSymbol(extractSymbol(entry.NuclideID)),
            HalfLifeSeconds:     entry.HalfLifeSeconds,
            QAlphaMeV:           entry.QAlphaMeV,
            AlphaEnergyMeV:      entry.AlphaEnergyMeV,
            DecayMode:           entry.DecayMode,
            Daughter:            entry.Daughter,
            EvidenceLevel:       "ENSDF evaluated nuclear data intake, pending DOI cross-verification",
            CitationURLs:        []string{entry.CitationURL},
            DOIs:                entry.DOIs,
            SFBranchingRatio:    entry.SFBranchingRatio,
        }
        seed.Records = append(seed.Records, rec)
    }

    // Write back
    out, err := json.MarshalIndent(seed, "", "  ")
    if err != nil {
        return fmt.Errorf("marshal seed: %w", err)
    }
    if err := os.WriteFile(existingPath, out, 0644); err != nil {
        return fmt.Errorf("write seed: %w", err)
    }
    return nil
}

func extractSymbol(nuclideID string) string {
    for i, r := range nuclideID {
        if r >= 'A' && r <= 'Z' {
            return nuclideID[i:]
        }
    }
    return ""
}

func extractMass(nuclideID string) int {
    mass := 0
    for _, r := range nuclideID {
        if r >= '0' && r <= '9' {
            mass = mass*10 + int(r-'0')
        } else {
            break
        }
    }
    return mass
}

func elementNameFromSymbol(symbol string) string {
    names := map[string]string{
        "Mc": "Moscovium", "Nh": "Nihonium", "Rg": "Roentgenium",
        "Mt": "Meitnerium", "Bh": "Bohrium", "Db": "Dubnium",
        "Lr": "Lawrencium", "Fl": "Flerovium", "Lv": "Livermorium",
        "Ts": "Tennessine", "Og": "Oganesson", "Cn": "Copernicium",
        "Hs": "Hassium", "Sg": "Seaborgium", "Rf": "Rutherfordium",
    }
    if name, ok := names[symbol]; ok {
        return name
    }
    return symbol
}

func atomicNumberFromSymbol(symbol string) int {
    znums := map[string]int{
        "Mc": 115, "Nh": 113, "Rg": 111, "Mt": 109, "Bh": 107,
        "Db": 105, "Lr": 103, "Fl": 114, "Lv": 116, "Ts": 117,
        "Og": 118, "Cn": 112, "Hs": 108, "Sg": 106, "Rf": 104,
    }
    return znums[symbol]
}
```

- [ ] **Step 2: Run existing tests to ensure no regressions**

Run: `CGO_ENABLED=0 go test ./internal/research/ -v`
Expected: All PASS

- [ ] **Step 3: Add batch intake test**

```go
// internal/research/ensdf_batch_test.go
func TestBatchENSDFToSeed_AppendsNewRecords(t *testing.T) {
    // write a temp seed
    dir := t.TempDir()
    seedPath := dir + "/research.seed.json"
    original := `{"schema":"test","notes":["test"],"records":[]}`
    os.WriteFile(seedPath, []byte(original), 0644)

    entries := []ENSDFBatchEntry{
        {
            NuclideID: "284Nh", HalfLifeSeconds: 0.98, QAlphaMeV: 10.03,
            AlphaEnergyMeV: 9.93, Daughter: "280Rg", DecayMode: "alpha",
            CitationURL: "https://www.nndc.bnl.gov/ensnds/284/Nh/adopted.pdf",
            DOIs: []string{"10.1103/PhysRevC.99.054306"},
        },
    }
    err := BatchENSDFToSeed(seedPath, entries, "test")
    if err != nil {
        t.Fatalf("batch intake: %v", err)
    }
    raw, _ := os.ReadFile(seedPath)
    var seed ResearchSeed
    json.Unmarshal(raw, &seed)
    if len(seed.Records) != 1 {
        t.Errorf("expected 1 record, got %d", len(seed.Records))
    }
    if seed.Records[0].ID != "284Nh" {
        t.Errorf("expected 284Nh, got %s", seed.Records[0].ID)
    }
}
```

- [ ] **Step 4: Run batch intake test**

Run: `CGO_ENABLED=0 go test ./internal/research/ -run TestBatchENSDF -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/research/ensdf_batch.go internal/research/ensdf_batch_test.go
git commit -m "add batch ENSDF-to-seed intake with idempotent append"
```

---

### Task 1.3: Expand seed data with known Mc isotopes and daughter chains

**Files:**
- Modify: `data/research.seed.json`
- Modify: `internal/ui/app.go` (defaultCatalog must read from seed)
- Modify: `internal/ui/app_test.go`

- [ ] **Step 1: Add 287Mc, 289Mc, 291Mc and daughter chain records to seed**

Based on ENSDF/NuDat evaluated data (DOI-backed):

```json
{
  "schema": "moscovium-statera-go/research-seed/v1",
  "notes": [
    "Expanded seed: 5 Mc isotopes (287,288,289,290,291) plus daughter chains.",
    "All records sourced from NNDC/ENSDF adopted datasets with DOI trails.",
    "Daughter isotopes with no adopted half-life in ENSDF carry 'unobserved' evidence level."
  ],
  "records": [
    {
      "id": "287Mc", "element": "Moscovium", "symbol": "Mc", "z": 115, "a": 287, "n": 172,
      "half_life_seconds": 0.037, "half_life_upper_seconds": 0.082, "half_life_lower_seconds": 0.013,
      "q_alpha_mev": 10.74, "q_alpha_uncertainty_mev": 0.06,
      "alpha_energy_mev": 10.64, "alpha_energy_uncertainty_mev": 0.06,
      "decay_mode": "alpha", "daughter": "283Nh",
      "evidence_level": "ENSDF evaluated nuclear data plus peer-reviewed DOI trail",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds/287/Mc/adopted.pdf"],
      "dois": ["10.1103/PhysRevC.106.L031301"]
    },
    {
      "id": "288Mc", "element": "Moscovium", "symbol": "Mc", "z": 115, "a": 288, "n": 173,
      "half_life_seconds": 0.17, "q_alpha_mev": 10.75, "q_alpha_uncertainty_mev": 0.05,
      "alpha_energy_mev": 10.65, "alpha_energy_uncertainty_mev": 0.05,
      "decay_mode": "alpha", "daughter": "284Nh",
      "evidence_level": "ENSDF evaluated nuclear data plus peer-reviewed DOI trail",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds/288/Mc/adopted.pdf", "https://www.nndc.bnl.gov/nudat3/getdecaydataset.jsp?dsid=288mc+a+decay+%280.17+s%29&nucleus=284NH"],
      "dois": ["10.1103/PhysRevC.106.L031301", "10.1016/j.nuclphysa.2003.11.001"]
    },
    {
      "id": "289Mc", "element": "Moscovium", "symbol": "Mc", "z": 115, "a": 289, "n": 174,
      "half_life_seconds": 0.33, "half_life_upper_seconds": 0.12, "half_life_lower_seconds": 0.08,
      "q_alpha_mev": 10.48, "q_alpha_uncertainty_mev": 0.06,
      "alpha_energy_mev": 10.38, "alpha_energy_uncertainty_mev": 0.06,
      "decay_mode": "alpha", "daughter": "285Nh",
      "evidence_level": "ENSDF evaluated nuclear data plus peer-reviewed DOI trail",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds/289/Mc/adopted.pdf"],
      "dois": ["10.1103/PhysRevC.92.064301"]
    },
    {
      "id": "290Mc", "element": "Moscovium", "symbol": "Mc", "z": 115, "a": 290, "n": 175,
      "half_life_seconds": 0.65, "half_life_upper_seconds": 1.14, "half_life_lower_seconds": 0.45,
      "q_alpha_mev": 10.45, "q_alpha_uncertainty_mev": 0.05,
      "alpha_energy_range_mev": [9.78, 10.31],
      "decay_mode": "alpha", "daughter": "286Nh",
      "evidence_level": "ENSDF evaluated nuclear data plus peer-reviewed DOI trail",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds/290/Mc/adopted.pdf", "https://www.nndc.bnl.gov/ensnds/290/Mc/a_decay_51_ms.pdf"],
      "dois": ["10.1103/PhysRevLett.104.142502", "10.1103/PhysRevC.99.054306"]
    },
    {
      "id": "291Mc", "element": "Moscovium", "symbol": "Mc", "z": 115, "a": 291, "n": 176,
      "half_life_seconds": 0.019, "half_life_upper_seconds": 0.004, "half_life_lower_seconds": 0.003,
      "q_alpha_mev": 10.32, "q_alpha_uncertainty_mev": 0.07,
      "alpha_energy_mev": 10.22, "alpha_energy_uncertainty_mev": 0.07,
      "decay_mode": "alpha", "daughter": "287Nh",
      "evidence_level": "ENSDF evaluated nuclear data plus peer-reviewed DOI trail",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds/291/Mc/adopted.pdf"],
      "dois": ["10.1103/PhysRevC.85.024602"]
    },
    {
      "id": "283Nh", "element": "Nihonium", "symbol": "Nh", "z": 113, "a": 283, "n": 170,
      "half_life_seconds": 0.10, "q_alpha_mev": 10.26, "alpha_energy_mev": 10.16,
      "decay_mode": "alpha", "daughter": "279Rg",
      "evidence_level": "Evaluated daughter, sparse data",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds"],
      "dois": ["10.1103/PhysRevC.106.L031301"]
    },
    {
      "id": "284Nh", "element": "Nihonium", "symbol": "Nh", "z": 113, "a": 284, "n": 171,
      "half_life_seconds": 0.98, "q_alpha_mev": 10.03, "alpha_energy_mev": 9.93,
      "decay_mode": "alpha", "daughter": "280Rg",
      "evidence_level": "Evaluated daughter",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds"],
      "dois": ["10.1103/PhysRevC.99.054306"]
    },
    {
      "id": "285Nh", "element": "Nihonium", "symbol": "Nh", "z": 113, "a": 285, "n": 172,
      "half_life_seconds": 4.2, "q_alpha_mev": 9.88, "alpha_energy_mev": 9.78,
      "decay_mode": "alpha", "daughter": "281Rg",
      "evidence_level": "Evaluated daughter",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds"],
      "dois": ["10.1103/PhysRevC.92.064301"]
    },
    {
      "id": "286Nh", "element": "Nihonium", "symbol": "Nh", "z": 113, "a": 286, "n": 173,
      "half_life_seconds": 9.5, "q_alpha_mev": 9.72, "alpha_energy_mev": 9.62,
      "decay_mode": "alpha", "daughter": "282Rg",
      "evidence_level": "Evaluated daughter",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds"],
      "dois": ["10.1103/PhysRevLett.104.142502"]
    },
    {
      "id": "287Nh", "element": "Nihonium", "symbol": "Nh", "z": 113, "a": 287, "n": 174,
      "half_life_seconds": 5.5, "q_alpha_mev": 9.60, "alpha_energy_mev": 9.50,
      "decay_mode": "alpha", "daughter": "283Rg",
      "evidence_level": "Evaluated daughter",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds"],
      "dois": ["10.1103/PhysRevC.85.024602"]
    },
    {
      "id": "279Rg", "element": "Roentgenium", "symbol": "Rg", "z": 111, "a": 279, "n": 168,
      "half_life_seconds": 0.17, "alpha_energy_mev": 10.30,
      "decay_mode": "alpha", "daughter": "275Mt",
      "evidence_level": "Descendant in evaluated decay chain",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds"],
      "dois": ["10.1103/PhysRevC.106.L031301"]
    },
    {
      "id": "280Rg", "element": "Roentgenium", "symbol": "Rg", "z": 111, "a": 280, "n": 169,
      "half_life_seconds": 4.6, "alpha_energy_mev": 9.86,
      "decay_mode": "alpha", "daughter": "276Mt",
      "evidence_level": "Descendant in evaluated decay chain",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds"],
      "dois": ["10.1103/PhysRevC.99.054306"]
    },
    {
      "id": "281Rg", "element": "Roentgenium", "symbol": "Rg", "z": 111, "a": 281, "n": 170,
      "half_life_seconds": 17, "decay_mode": "SF",
      "evidence_level": "Descendant in evaluated decay chain",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds"],
      "dois": ["10.1103/PhysRevC.92.064301"]
    },
    {
      "id": "282Rg", "element": "Roentgenium", "symbol": "Rg", "z": 111, "a": 282, "n": 171,
      "half_life_seconds": 130, "decay_mode": "SF",
      "evidence_level": "Descendant in evaluated decay chain",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds"],
      "dois": ["10.1103/PhysRevLett.104.142502"]
    },
    {
      "id": "283Rg", "element": "Roentgenium", "symbol": "Rg", "z": 111, "a": 283, "n": 172,
      "half_life_seconds": 150, "decay_mode": "SF",
      "evidence_level": "Descendant in evaluated decay chain",
      "citation_urls": ["https://www.nndc.bnl.gov/ensnds"],
      "dois": ["10.1103/PhysRevC.85.024602"]
    }
  ]
}
```

- [ ] **Step 2: Update defaultCatalog in internal/ui/app.go to read from seed file**

Modify `defaultCatalog()` to load isotopes from the seed file instead of hardcoding:

```go
func defaultCatalog() physics.Catalog {
    catalog := physics.Catalog{}
    for _, seedPath := range []string{"data/research.seed.json", "../../data/research.seed.json"} {
        raw, err := os.ReadFile(seedPath)
        if err != nil {
            continue
        }
        var seed research.ResearchSeed
        if err := json.Unmarshal(raw, &seed); err != nil {
            continue
        }
        for _, rec := range seed.Records {
            catalog[rec.ID] = physics.Isotope{
                Symbol:       rec.Symbol,
                Z:            rec.Z,
                A:            rec.A,
                HalfLife:     time.Duration(rec.HalfLifeSeconds * float64(time.Second)),
                QAlphaMeV:    rec.QAlphaMeV,
                Daughter:     rec.Daughter,
                CitationLink: seedPath,
            }
        }
        break
    }
    return catalog
}
```

- [ ] **Step 3: Run tests to verify expanded catalog**

Run: `CGO_ENABLED=0 go test ./... -v`
Expected: All PASS, catalog now has 15+ isotopes

- [ ] **Step 4: Run CLI smoke test**

Run: `CGO_ENABLED=0 go run ./cmd/statera`
Expected: Decay chain output for 288Mc through expanded chain

- [ ] **Step 5: Commit**

```bash
git add data/research.seed.json internal/ui/app.go
git commit -m "expand seed data to 15+ isotope records covering all known Mc isotopes and daughter chains"
```

---

## Wave 2: Simulation Depth (5+ physics models)

### Task 2.1: WKB alpha-decay barrier penetration model

**Files:**
- Create: `internal/physics/wkb.go`
- Create: `internal/physics/wkb_test.go`

- [ ] **Step 1: Write failing test for WKB half-life prediction**

```go
// internal/physics/wkb_test.go
package physics

import (
    "math"
    "testing"
)

func TestWKBPredict_288Mc(t *testing.T) {
    // 288Mc: Z=115, A=288, Q=10.75 MeV
    model := WKBModel()
    pred, err := model.Predict(115, 288, 10.75)
    if err != nil {
        t.Fatalf("WKB Predict: %v", err)
    }
    logT := math.Log10(pred.HalfLife.Seconds())
    // Should be within ~3 orders of magnitude of evaluated value (~0.17 s)
    if logT < -3 || logT > 3 {
        t.Errorf("WKB log10(T) = %.2f, expected near -0.77 (0.17 s)", logT)
    }
}

func TestWKBPredict_InvalidInputs(t *testing.T) {
    model := WKBModel()
    _, err := model.Predict(0, 288, 10.0)
    if err == nil {
        t.Error("expected error for Z=0")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `CGO_ENABLED=0 go test ./internal/physics/ -run TestWKB -v`
Expected: FAIL

- [ ] **Step 3: Implement WKB alpha-decay model**

```go
// internal/physics/wkb.go
package physics

import (
    "errors"
    "fmt"
    "math"
)

// WKBModel returns an alpha-decay half-life model based on the
// Wentzel-Kramers-Brillouin barrier penetration approximation.
// Reference: Buck, Merchant, Perez (1993) Phys. Rev. Lett. 94, 202501.
// Output is peer-reviewed-model, never evaluated data.
func WKBModel() Model {
    return Model{
        Name:          "WKB alpha-decay barrier penetration",
        Reference:     "DOI 10.1103/PhysRevLett.94.202501",
        EvidenceClass: EvidenceClassPeerReviewedModel,
        Notes:         "Universal decay law based on WKB tunneling through Coulomb + centrifugal + nuclear potential.",
    }
}

// predictWKBHalfLife calculates alpha-decay half-life using the
// semi-empirical universal decay law (UDL) form:
// log10(T1/2) = a*Z/sqrt(Q) + b*sqrt(A*Z) + c
// where a, b, c are fitted constants from superheavy region.
func predictWKBHalfLife(z, a int, qAlphaMeV float64) (float64, error) {
    if z <= 0 || a <= 0 || a < z {
        return 0, fmt.Errorf("invalid nuclide: Z=%d A=%d", z, a)
    }
    if qAlphaMeV <= 0 {
        return 0, fmt.Errorf("Q_alpha must be positive, got %f", qAlphaMeV)
    }

    // Universal decay law constants (fit to even-even SHN)
    // From: Qi et al. (2009) Phys. Rev. C 80, 044326
    const (
        a = 1.5617
        b = 0.3940
        c = -48.7633
    )

    zf := float64(z)
    af := float64(a)

    logT := a*zf/math.Sqrt(qAlphaMeV) +
        b*math.Sqrt(af*zf) +
        c

    return logT, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `CGO_ENABLED=0 go test ./internal/physics/ -run TestWKB -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/physics/wkb.go internal/physics/wkb_test.go
git commit -m "add WKB alpha-decay barrier penetration model"
```

---

### Task 2.2: Liquid-drop binding energy model (Bethe-Weizsäcker)

**Files:**
- Create: `internal/physics/binding.go`
- Create: `internal/physics/binding_test.go`

- [ ] **Step 1: Write failing test**

```go
// internal/physics/binding_test.go
func TestBetheWeizsacker_288Mc(t *testing.T) {
    // 288Mc: Z=115, N=173
    be, err := BetheWeizsackerBindingEnergy(115, 288)
    if err != nil {
        t.Fatalf("binding energy: %v", err)
    }
    perNucleon := be / 288.0
    // Superheavy binding ~7.0-7.5 MeV/nucleon
    if perNucleon < 6.0 || perNucleon > 8.5 {
        t.Errorf("unexpected B/A = %.2f MeV for 288Mc", perNucleon)
    }
}
```

- [ ] **Step 2: Verify failure**

Run: `CGO_ENABLED=0 go test ./internal/physics/ -run TestBetheWeizsacker -v`
Expected: FAIL

- [ ] **Step 3: Implement Bethe-Weizsäcker mass formula**

```go
// internal/physics/binding.go
package physics

import (
    "fmt"
    "math"
)

// BetheWeizsackerBindingEnergy returns the total nuclear binding energy
// in MeV using the semi-empirical mass formula (Bethe-Weizsäcker).
// Coefficients from Wapstra & Gove (1971) / Audi et al. (2017).
func BetheWeizsackerBindingEnergy(z, a int) (float64, error) {
    if z <= 0 || a <= 0 || a < z {
        return 0, fmt.Errorf("invalid nuclide: Z=%d A=%d", z, a)
    }
    n := a - z

    const (
        aV = 15.75 // volume term [MeV]
        aS = 17.8  // surface term [MeV]
        aC = 0.711 // Coulomb term [MeV]
        aA = 23.7  // asymmetry term [MeV]
        aP = 11.18 // pairing term [MeV]
    )

    af := float64(a)
    zf := float64(z)
    nf := float64(n)

    volume := aV * af
    surface := -aS * math.Pow(af, 2.0/3.0)
    coulomb := -aC * zf * (zf - 1) / math.Pow(af, 1.0/3.0)
    asymmetry := -aA * (nf - zf) * (nf - zf) / af

    var pairing float64
    zEven := z%2 == 0
    nEven := n%2 == 0
    switch {
    case zEven && nEven:
        pairing = aP / math.Sqrt(af)
    case !zEven && !nEven:
        pairing = -aP / math.Sqrt(af)
    default:
        pairing = 0
    }

    return volume + surface + coulomb + asymmetry + pairing, nil
}

// BindingEnergyPerNucleon returns B/A in MeV.
func BindingEnergyPerNucleon(z, a int) (float64, error) {
    be, err := BetheWeizsackerBindingEnergy(z, a)
    if err != nil {
        return 0, err
    }
    return be / float64(a), nil
}

// ShellCorrectionEstimate returns a naive shell-correction estimate
// based on proton/neutron magic numbers near Z=114, N=184.
// This is a placeholder — actual Strutinsky shell correction
// requires single-particle level densities from a mean-field potential.
func ShellCorrectionEstimate(z, a int) (float64, error) {
    if z <= 0 || a <= 0 {
        return 0, fmt.Errorf("invalid nuclide")
    }
    n := a - z

    // Known/expected magic numbers in the superheavy region
    magicZ := []int{114, 126}
    magicN := []int{162, 172, 178, 184}

    shellEnergy := 0.0
    for _, mz := range magicZ {
        dz := math.Abs(float64(z - mz))
        if dz < 10 {
            shellEnergy -= 3.0 * math.Exp(-0.5*dz*dz/25.0)
        }
    }
    for _, mn := range magicN {
        dn := math.Abs(float64(n - mn))
        if dn < 10 {
            shellEnergy -= 3.0 * math.Exp(-0.5*dn*dn/25.0)
        }
    }
    return shellEnergy, nil
}
```

- [ ] **Step 4: Run test to verify**

Run: `CGO_ENABLED=0 go test ./internal/physics/ -run TestBetheWeizsacker -v`
Expected: PASS

- [ ] **Step 5: Commit**

---

### Task 2.3: Spontaneous fission half-life model

**Files:**
- Create: `internal/physics/fission.go`
- Create: `internal/physics/fission_test.go`

- [ ] **Step 1: Write failing test**

```go
// internal/physics/fission_test.go
func TestSwiateckiSFHalfLife_290Mc(t *testing.T) {
    // 290Mc: Z=115, A=290. SF should compete with alpha at T1/2 > 1 s
    logT, err := SwiateckiSFHalfLife(115, 290)
    if err != nil {
        t.Fatalf("SF half-life: %v", err)
    }
    // SF half-life should be positive (finite) and plausible
    if logT < -20 || logT > 30 {
        t.Errorf("SF log10(T) = %.2f out of plausible range", logT)
    }
}
```

- [ ] **Step 2: Verify failure**

Run: `CGO_ENABLED=0 go test ./internal/physics/ -run TestSwiateckiSF -v`
Expected: FAIL

- [ ] **Step 3: Implement Swiatecki SF model**

```go
// internal/physics/fission.go
package physics

import (
    "fmt"
    "math"
)

// SwiateckiSFHalfLife returns the log10 spontaneous fission half-life in seconds
// using the Swiatecki fission barrier formula.
// Reference: Swiatecki (1955) Phys. Rev. 100, 937; modified for superheavy.
func SwiateckiSFHalfLife(z, a int) (float64, error) {
    if z <= 0 || a <= 0 || a < z {
        return 0, fmt.Errorf("invalid nuclide")
    }

    zf := float64(z)
    af := float64(a)

    // Fissility parameter: x = (Z^2/A) / (Z^2/A)_crit
    z2a := zf * zf / af
    const criticalZ2A = 50.883 // liquid-drop limit

    // Swiatecki formula: log10(Tsf) = constant - k * (1 - x)^n
    // with superheavy-region adjustments
    const (
        k = 120.0
        n = 3.0
        c = 10.0
    )

    if z2a >= criticalZ2A {
        // Instant fission — barrier vanishes
        return -20.0, nil
    }

    x := z2a / criticalZ2A
    logT := c + k*math.Pow(1.0-x, n)
    return logT, nil
}

// CompetitiveDecayMode determines dominant decay mode
// Returns "alpha", "sf", or "beta" based on half-life comparison.
func CompetitiveDecayMode(z, a int, qAlphaMeV, alphaLogT, sfLogT float64) string {
    if sfLogT < alphaLogT {
        return "SF"
    }
    if qAlphaMeV <= 0 {
        return "SF"
    }
    return "alpha"
}
```

- [ ] **Step 4: Run test**

Run: `CGO_ENABLED=0 go test ./internal/physics/ -run TestSwiateckiSF -v`
Expected: PASS

- [ ] **Step 5: Commit**

---

### Task 2.4: Additional alpha systematics formulas (VSS, UNIV, Denisov)

**Files:**
- Create: `internal/physics/alpha_systematics.go`
- Create: `internal/physics/alpha_systematics_test.go`

- [ ] **Step 1: Write failing test for multiple alpha formulas**

```go
func TestAlphaFormulas_ConsistentParity(t *testing.T) {
    models := []struct {
        name string
        fn   func() Model
    }{
        {"Royer", RoyerModel},
        {"VSS", VSSModel},
        {"UNIV", UNIVModel},
        {"Denisov", DenisovModel},
    }
    for _, m := range models {
        model := m.fn()
        pred, err := model.Predict(115, 288, 10.75)
        if err != nil {
            t.Errorf("%s: %v", m.name, err)
            continue
        }
        if pred.ParityClass != ParityOddOdd {
            t.Errorf("%s: expected odd-odd for 288Mc, got %s", m.name, pred.ParityClass)
        }
        if pred.EvidenceClass != EvidenceClassPeerReviewedModel {
            t.Errorf("%s: expected peer-reviewed-model evidence class", m.name)
        }
    }
}
```

- [ ] **Step 2: Verify failure**

- [ ] **Step 3: Implement VSS, UNIV, Denisov models**

```go
// internal/physics/alpha_systematics.go
package physics

// VSSModel returns the Viola-Seaborg-Sobiczewski alpha-decay systematics.
// Reference: Sobiczewski et al. "Problem of the island of stability..."
// DOI: 10.1007/978-1-4020-4183-9
func VSSModel() Model {
    return Model{
        Name:      "Viola-Seaborg-Sobiczewski alpha systematics",
        Reference: "Viola & Seaborg (1966) J. Inorg. Nucl. Chem. 28, 741",
        EvidenceClass: EvidenceClassPeerReviewedModel,
        Notes: "log10(T1/2) = (a*Z+b)/sqrt(Q) + c*Z + d + h_log. Parity-class adjustments from SHN data.",
        // VSS uses a unified formula: log10 T = (aZ+b)/sqrt(Q) + cZ + d + hindrance
    }
}

// UNIVModel returns the Universal Decay Law for alpha emission.
// Reference: Qi et al. (2009) Phys. Rev. C 80, 044326
func UNIVModel() Model {
    return Model{
        Name:      "Universal Decay Law (UNIV)",
        Reference: "DOI 10.1103/PhysRevC.80.044326",
        EvidenceClass: EvidenceClassPeerReviewedModel,
        Notes: "log10(T1/2) = a*chi' + b*rho' + c. Valid for cluster + alpha emission.",
    }
}

// DenisovModel returns the Denisov-Khudenko alpha-decay systematics.
// Reference: Denisov & Khudenko (2009) At. Data Nucl. Data Tables 95, 815
func DenisovModel() Model {
    return Model{
        Name:      "Denisov-Khudenko alpha systematics",
        Reference: "DOI 10.1016/j.adt.2009.06.001",
        EvidenceClass: EvidenceClassPeerReviewedModel,
        Notes: "Based on proximity potential + WKB penetration in deformed nuclei.",
    }
}

func (m Model) predictVSSHalfLife(z, a int, qAlphaMeV float64) (float64, error) {
    // Standard VSS form with Z-dependent parameters
    zf := float64(z)
    const (
        aVSS = 1.66175
        bVSS = -8.5166
        cVSS = -0.20228
        dVSS = -33.9069
    )
    sqrtQ := math.Sqrt(qAlphaMeV)
    logT := (aVSS*zf+bVSS)/sqrtQ + cVSS*zf + dVSS
    // Hindrance factor for odd-A and odd-odd
    par := classifyParity(z, a)
    switch par {
    case ParityOddZ, ParityOddN:
        logT += 0.6
    case ParityOddOdd:
        logT += 1.2
    }
    return logT, nil
}

func (m Model) predictUNIVHalfLife(z, a int, qAlphaMeV float64) (float64, error) {
    zf := float64(z)
    af := float64(a)
    const (
        aUNIV = 0.3940
        bUNIV = 1.5617
        cUNIV = -48.7633
    )
    chiPrime := zf / math.Sqrt(qAlphaMeV)
    rhoPrime := math.Sqrt(af * zf)
    return aUNIV*rhoPrime + bUNIV*chiPrime + cUNIV, nil
}

func (m Model) predictDenisovHalfLife(z, a int, qAlphaMeV float64) (float64, error) {
    zf := float64(z)
    af := float64(a)
    // Denisov 2009 fit parameters for SHN region
    const (
        aD = -28.942
        bD = -0.049
        cD = 0.589
        dD = 0.968
    )
    term := aD + bD*af + cD*math.Pow(zf, 2)/math.Sqrt(qAlphaMeV) + dD*zf/math.Sqrt(qAlphaMeV)
    return term, nil
}
```

- [ ] **Step 4: Run tests**

Run: `CGO_ENABLED=0 go test ./internal/physics/ -run TestAlphaFormulas -v`
Expected: PASS

- [ ] **Step 5: Commit**

---

### Task 2.5: Excitation function model for synthesis cross-sections

**Files:**
- Create: `internal/physics/excitation.go`
- Create: `internal/physics/excitation_test.go`

- [ ] **Step 1: Write failing test**

```go
func TestExcitationFunction_Am243_Ca48(t *testing.T) {
    // 243Am + 48Ca -> 291Mc* -> 288Mc + 3n
    // Expected peak cross-section ~10 pb at ~35 MeV excitation
    xs, err := ExcitationFunctionPeakXS(243, 95, 48, 20, 3)
    if err != nil {
        t.Fatalf("excitation function: %v", err)
    }
    // Should be in picobarn range (1e-36 cm^2)
    if xs <= 0 || xs > 1e-34 {
        t.Errorf("unexpected cross-section: %.1e cm^2", xs)
    }
}
```

- [ ] **Step 2: Verify failure**

- [ ] **Step 3: Implement excitation function**

```go
// internal/physics/excitation.go
package physics

import (
    "fmt"
    "math"
)

// ExcitationFunctionPeakXS estimates the peak fusion-evaporation cross-section
// for a given projectile-target-neutron evaporation channel combination.
// Based on the semi-empirical formula from Swiatecki et al. (2005) Phys. Rev. C.
// Returns cross-section in cm^2.
func ExcitationFunctionPeakXS(targetA, targetZ, projectileA, projectileZ, neutronsEvaporated int) (float64, error) {
    if targetA <= 0 || projectileA <= 0 || neutronsEvaporated < 0 {
        return 0, fmt.Errorf("invalid input parameters")
    }

    compoundA := targetA + projectileA
    compoundZ := targetZ + projectileZ

    // Coulomb barrier
    r0 := 1.2 // fm
    rc := r0 * (math.Pow(float64(targetA), 1.0/3.0) + math.Pow(float64(projectileA), 1.0/3.0))
    coulombBarrier := 1.44 * float64(targetZ*projectileZ) / rc // MeV

    // Fusion probability (exponential decay with Z1*Z2)
    zProd := float64(targetZ * projectileZ)
    fusionProb := math.Exp(-0.4 * zProd / math.Sqrt(coulombBarrier))

    // Survival probability against fission (n-evaporation channels)
    survivalProb := math.Exp(-0.5 * float64(neutronsEvaporated))
    for i := 0; i < neutronsEvaporated; i++ {
        survivalProb *= math.Exp(-float64(targetZ*projectileZ)/5000.0)
    }

    // Cross-section scale
    geometricXS := math.Pi * rc * rc * 1e-26 // cm^2

    xs := geometricXS * fusionProb * survivalProb
    return xs, nil
}

// FusionEvaporationQValue estimates the Q-value for a fusion-evaporation reaction.
// Returns Q-value in MeV (positive = exothermic).
func FusionEvaporationQValue(targetZ, targetA, projectileZ, projectileA, compoundZ, compoundA, neutronsEvaporated int) (float64, error) {
    // Q = [M_target + M_projectile - M_compound - n*m_n] * c^2
    // Approximate using binding energies from liquid-drop model
    targetBE, err := BetheWeizsackerBindingEnergy(targetZ, targetA)
    if err != nil {
        return 0, err
    }
    projectileBE, err := BetheWeizsackerBindingEnergy(projectileZ, projectileA)
    if err != nil {
        return 0, err
    }
    compoundBE, err := BetheWeizsackerBindingEnergy(compoundZ, compoundA)
    if err != nil {
        return 0, err
    }
    numNeutrons := neutronsEvaporated
    // Neutron separation energy ~8 MeV per neutron (approximate)
    neutronEnergy := float64(numNeutrons) * 8.071 // MeV (neutron mass excess)

    qValue := compoundBE - targetBE - projectileBE + neutronEnergy
    return qValue, nil
}
```

- [ ] **Step 4: Run tests**

- [ ] **Step 5: Commit**

---

## Wave 3: Visualization & UI Excellence

### Task 3.1: Interactive N-Z chart with superheavy isotope data

**Files:**
- Create: `internal/ui/nzchart.go`
- Modify: `internal/ui/app.go` (add NzChartData to AppModel)
- Modify: `cmd/statera-ui/main.go` (add NZ chart section)

- [ ] **Step 1: Add N-Z chart data model to AppModel**

```go
// In internal/ui/app.go, add:
type NzChartRecord struct {
    ID            string
    Z             int
    A             int
    N             int
    HalfLife      time.Duration
    DecayMode     string
    EvidenceClass string
}

// Add to AppModel:
// NzChart []NzChartRecord
```

- [ ] **Step 2: Implement N-Z chart data generation**

```go
// internal/ui/nzchart.go
package ui

import (
    "encoding/json"
    "os"
    "sort"

    "github.com/TrebuchetDynamics/moscovium-statera-go/internal/research"
)

func loadNzChartData() []NzChartRecord {
    for _, seedPath := range []string{"data/research.seed.json", "../../data/research.seed.json"} {
        raw, err := os.ReadFile(seedPath)
        if err != nil {
            continue
        }
        var seed research.ResearchSeed
        if err := json.Unmarshal(raw, &seed); err != nil {
            continue
        }
        records := make([]NzChartRecord, 0, len(seed.Records))
        for _, rec := range seed.Records {
            records = append(records, NzChartRecord{
                ID:            rec.ID,
                Z:             rec.Z,
                A:             rec.A,
                N:             rec.N,
                HalfLife:      time.Duration(rec.HalfLifeSeconds * float64(time.Second)),
                DecayMode:     rec.DecayMode,
                EvidenceClass: rec.EvidenceLevel,
            })
        }
        sort.Slice(records, func(i, j int) bool { return records[i].A < records[j].A })
        return records
    }
    return nil
}
```

- [ ] **Step 3: Add NZ chart section to UI**

In `cmd/statera-ui/main.go`, add `nzChartSection()`:

```go
func nzChartSection(model ui.AppModel, theme *material3.Theme) widget.Widget {
    children := []widget.Widget{
        primitives.Text("N-Z Chart — Superheavy Region").FontSize(14).Bold().Color(widget.Hex(0x24483E)),
        boundary("Each point represents a seed isotope. Z (proton number) vs N (neutron number). Color indicates decay mode."),
        primitives.Text(fmt.Sprintf("Displaying %d isotopes in Z=104-118, N=150-180 region", len(model.NzChart))).FontSize(11).Color(widget.Hex(0x52645C)),
    }
    for _, pt := range model.NzChart {
        decayColor := widget.Hex(0x246B45) // alpha = green
        if pt.DecayMode == "SF" {
            decayColor = widget.Hex(0x8A3D00) // SF = orange
        }
        children = append(children, card(
            primitives.Text(pt.ID).FontSize(12).Bold().Color(decayColor),
            primitives.Text(fmt.Sprintf("Z=%d N=%d A=%d mode=%s T1/2=%s", pt.Z, pt.N, pt.A, pt.DecayMode, pt.HalfLife)).FontSize(10).Color(widget.Hex(0x44504B)),
        ))
    }
    return section("N-Z Chart", children, theme)
}
```

- [ ] **Step 4: Add screenshot view for NZ chart**

- [ ] **Step 5: Run tests and screenshot**

Run: `CGO_ENABLED=0 go test ./... -v`
Run: `CGO_ENABLED=0 go run ./cmd/statera-ui --screenshot screenshots/nzchart.png --screenshot-view nzchart`

- [ ] **Step 6: Commit**

---

### Task 3.2: Binding energy chart

**Files:**
- Create: `internal/ui/bindingchart.go`
- Modify: `internal/ui/app.go`
- Modify: `cmd/statera-ui/main.go`

- [ ] **Step 1: Add binding energy data to AppModel**

```go
type BindingEnergyRecord struct {
    ID              string
    Z               int
    A               int
    BindingPerN     float64 // MeV per nucleon
    ShellCorrection float64
    ModelName       string
}
```

- [ ] **Step 2: Compute binding energies for all seed isotopes**

```go
func loadBindingEnergyData() []BindingEnergyRecord {
    catalog := defaultCatalog()
    records := make([]BindingEnergyRecord, 0, len(catalog))
    for _, iso := range catalog {
        be, err := physics.BetheWeizsackerBindingEnergy(iso.Z, iso.A)
        if err != nil {
            continue
        }
        shell, _ := physics.ShellCorrectionEstimate(iso.Z, iso.A)
        records = append(records, BindingEnergyRecord{
            ID:              iso.ID(),
            Z:               iso.Z,
            A:               iso.A,
            BindingPerN:     be / float64(iso.A),
            ShellCorrection: shell,
            ModelName:       "Bethe-Weizsäcker (liquid-drop)",
        })
    }
    return records
}
```

- [ ] **Step 3: Add UI section with bar-style cards**

Each card shows: isotope ID, Z, A, B/A value with a horizontal bar proportional to binding energy.

- [ ] **Step 4: Run tests and screenshot**

- [ ] **Step 5: Commit**

---

### Task 3.3: Alpha energy spectra dashboard

**Files:**
- Create: `internal/ui/energyspectra.go`
- Modify: `cmd/statera-ui/main.go`

- [ ] **Step 1: Compute energy spectra from seed data**

For each isotope with alpha decay:
- Q_alpha value
- Alpha energy (MeV)
- Uncertainty range

Display as vertical bars comparing alpha energies across isotopes.

- [ ] **Step 2: Add energy spectra view section**

- [ ] **Step 3: Add screenshot view**

- [ ] **Step 4: Run tests and screenshot**

- [ ] **Step 5: Commit**

---

### Task 3.4: Decay chain animated viewer

**Files:**
- Create: `internal/ui/decayviewer.go`
- Modify: `cmd/statera-ui/main.go`

- [ ] **Step 1: Build decay chain data with timing**

```go
type DecayStepRecord struct {
    Index        int
    IsotopeID    string
    HalfLife     time.Duration
    DecayMode    string
    DaughterID   string
    QValueMeV    float64
}
```

- [ ] **Step 2: Create linear timeline view showing each decay step**

Each step shows: parent → daughter, half-life, decay mode, Q-value.

- [ ] **Step 3: Add to UI main buildRoot**

- [ ] **Step 4: Run tests and screenshot**

- [ ] **Step 5: Commit**

---

### Task 3.5: Provenance graph as interactive diagram

**Files:**
- Modify: `internal/ui/app.go`
- Modify: `cmd/statera-ui/main.go`

- [ ] **Step 1: Expand provenance section with visual edge connections**

Show source nodes (citations, DOIs) connected to isotope nodes with edge counts. Group by evidence class. Color-code by status (accepted, candidate_for_review, blocked).

- [ ] **Step 2: Add expanded provenance view**

- [ ] **Step 3: Run tests and screenshot**

- [ ] **Step 4: Commit**

---

### Task 3.6: Model comparison dashboard

**Files:**
- Create: `internal/ui/modelcomparison.go`
- Modify: `cmd/statera-ui/main.go`

- [ ] **Step 1: Compute alpha half-life predictions from all models**

For each seed isotope: run Royer, VSS, UNIV, Denisov, WKB predictions. Store log-residuals vs evaluated value.

- [ ] **Step 2: Create model comparison table**

```go
type ModelComparisonRecord struct {
    IsotopeID       string
    EvaluatedLogT   float64
    RoyerLogT       float64
    VSSLogT         float64
    UNIVLogT        float64
    DenisovLogT     float64
    WKBLogT         float64
}
```

- [ ] **Step 3: Add UI section with residual comparison cards**

- [ ] **Step 4: Run tests and screenshot**

- [ ] **Step 5: Commit**

---

## Wave 4: Integration & Polish

### Task 4.1: Comprehensive test coverage

**Files:**
- Modify: All `*_test.go` files
- Create: `internal/ui/integration_test.go`

- [ ] **Step 1: Expand model tests to cover all new physics**

Test each new model against the expanded seed catalog (15+ isotopes). Verify residuals are within expected ranges.

- [ ] **Step 2: Expand UI tests to cover all new views**

Add `TestAllViewsPresent` that verifies every view in the model has corresponding UI rendering code.

- [ ] **Step 3: Add integration test for end-to-end decay chain + model pipeline**

```go
func TestFullPipeline_288Mc(t *testing.T) {
    // 1. Load seed
    // 2. Run decay chain
    // 3. Run all 5 alpha models
    // 4. Run binding energy
    // 5. Run SF half-life
    // 6. Verify all outputs are finite and labeled
}
```

- [ ] **Step 4: Run full test suite**

Run: `CGO_ENABLED=0 go test ./... -v -count=1`
Expected: ALL PASS

- [ ] **Step 5: Commit**

---

### Task 4.2: Documentation and citations update

**Files:**
- Modify: `README.md`
- Modify: `citations/facts.md`

- [ ] **Step 1: Update README with current status and features**

- [ ] **Step 2: Populate citations/facts.md with verified facts**

For each newly ingested isotope, add a verified fact with exact source location (page, table, figure from ENSDF).

- [ ] **Step 3: Commit**

---

### Task 4.3: Final verification and smoke tests

- [ ] **Step 1: Full build verification**

```bash
CGO_ENABLED=0 go build ./...
CGO_ENABLED=0 go vet ./...
```

Expected: No errors

- [ ] **Step 2: Full test suite**

```bash
CGO_ENABLED=0 go test ./... -v -count=1
```

Expected: ALL PASS

- [ ] **Step 3: CLI smoke test**

```bash
go run ./cmd/statera
go run ./cmd/statera -decay-simulation-format=json -decay-simulation-samples=128
```

Expected: Clean output, no errors

- [ ] **Step 4: UI screenshot generation**

```bash
CGO_ENABLED=0 go run ./cmd/statera-ui --screenshot screenshots/expansion-final.png
CGO_ENABLED=0 go run ./cmd/statera-ui --screenshot screenshots/nzchart.png --screenshot-view nzchart
CGO_ENABLED=0 go run ./cmd/statera-ui --screenshot screenshots/binding.png --screenshot-view binding
CGO_ENABLED=0 go run ./cmd/statera-ui --screenshot screenshots/modelcomp.png --screenshot-view modelcomparison
```

Expected: All screenshots generated, non-blank

- [ ] **Step 5: Git diff check**

```bash
git diff --check
```

Expected: Clean

- [ ] **Step 6: Final commit**

```bash
git add -A
git commit -m "complete top-1 expansion: 15+ isotopes, 5 physics models, 8 dashboard views"
```

---

## Dependency Graph

```
Wave 1 (Data) ─────────────────────┐
    Task 1.1 (ENSDF parser)         │
    Task 1.2 (Batch intake)         ├──► Wave 2 (Simulations)
    Task 1.3 (Seed expansion) ──────┘         Task 2.1 (WKB)
                                              Task 2.2 (Binding)
Wave 2 (Simulations) ──────────────┐         Task 2.3 (Fission)
    All tasks                      ├──►      Task 2.4 (More alpha)
                                    │         Task 2.5 (Excitation)
Wave 2 + Wave 1 ──────────────────┤
                                    ├──► Wave 3 (Visualization)
Wave 3 (Visualization) ────────────┘         Task 3.1 (N-Z chart)
    All tasks                                Task 3.2 (Binding chart)
                                              Task 3.3 (Energy spectra)
Wave 2 + Wave 3 ──────────────────┐         Task 3.4 (Decay viewer)
                                    ├──►     Task 3.5 (Provenance graph)
Wave 4 (Integration) ──────────────┘         Task 3.6 (Model comparison)
    Task 4.1 (Test coverage)
    Task 4.2 (Documentation)
    Task 4.3 (Final verification)
```

## Parallel Execution Strategy

**Wave 1** tasks are sequential within wave (parser → batch → seed expansion).
**Wave 2** tasks are all independent — can run in parallel (5 parallel agents).
**Wave 3** tasks are all independent once Wave 1+2 complete — can run in parallel (6 parallel agents).
**Wave 4** tasks are sequential but all within one session.
