package research

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestResearchSeedContainsVerifiedMoscoviumRecords(t *testing.T) {
	raw, err := os.ReadFile("../../data/research.seed.json")
	if err != nil {
		t.Fatalf("read seed: %v", err)
	}

	seed := readResearchSeed(t, raw)

	seen := map[string]bool{}
	for _, record := range seed.Records {
		if record.Z != 115 {
			t.Fatalf("%s has Z=%d, want 115", record.ID, record.Z)
		}
		if record.A != 288 && record.A != 290 {
			t.Fatalf("unexpected seed isotope: %s", record.ID)
		}
		if record.Daughter == "" {
			t.Fatalf("%s missing daughter", record.ID)
		}
		if len(record.CitationURLs) == 0 {
			t.Fatalf("%s missing citation URLs", record.ID)
		}
		if len(record.DOIs) == 0 {
			t.Fatalf("%s missing DOI trail", record.ID)
		}
		seen[record.ID] = true
	}

	for _, id := range []string{"288Mc", "290Mc"} {
		if !seen[id] {
			t.Fatalf("seed missing %s", id)
		}
	}
}

func TestValidateResearchSeedRequiresCompleteProvenance(t *testing.T) {
	raw, err := os.ReadFile("../../data/research.seed.json")
	if err != nil {
		t.Fatalf("read seed: %v", err)
	}
	seed := readResearchSeed(t, raw)

	if err := ValidateResearchSeedProvenance(seed); err != nil {
		t.Fatalf("valid seed rejected: %v", err)
	}

	cases := []struct {
		name   string
		mutate func(*ResearchSeed)
		want   string
	}{
		{
			name:   "missing citation URL",
			mutate: func(seed *ResearchSeed) { seed.Records[0].CitationURLs = nil },
			want:   "288Mc missing citation_urls",
		},
		{
			name:   "missing DOI",
			mutate: func(seed *ResearchSeed) { seed.Records[0].DOIs = nil },
			want:   "288Mc missing DOI trail",
		},
		{
			name:   "missing evidence level",
			mutate: func(seed *ResearchSeed) { seed.Records[0].EvidenceLevel = "" },
			want:   "288Mc missing evidence_level",
		},
		{
			name:   "ID does not match A and symbol",
			mutate: func(seed *ResearchSeed) { seed.Records[0].ID = "289Mc" },
			want:   "289Mc ID mismatch, want 288Mc from A=288 and symbol=Mc",
		},
		{
			name: "duplicate record ID",
			mutate: func(seed *ResearchSeed) {
				seed.Records[1].ID = seed.Records[0].ID
				seed.Records[1].A = seed.Records[0].A
				seed.Records[1].N = seed.Records[0].N
				seed.Records[1].Daughter = seed.Records[0].Daughter
			},
			want: "duplicate seed record ID 288Mc",
		},
		{
			name:   "malformed DOI",
			mutate: func(seed *ResearchSeed) { seed.Records[0].DOIs[0] = "PhysRevC.106.L031301" },
			want:   "288Mc DOI",
		},
		{
			name:   "inconsistent neutron count",
			mutate: func(seed *ResearchSeed) { seed.Records[0].N = seed.Records[0].A - seed.Records[0].Z + 1 },
			want:   "288Mc has N=174, want A-Z=173",
		},
		{
			name:   "atomic number does not match nuclide symbol",
			mutate: func(seed *ResearchSeed) { seed.Records[0].Z = 113 },
			want:   "288Mc has Z=113, want 115 for symbol Mc",
		},
		{
			name:   "element name does not match nuclide symbol",
			mutate: func(seed *ResearchSeed) { seed.Records[0].Element = "Nihonium" },
			want:   "288Mc element name \"Nihonium\" does not match symbol Mc, want Moscovium",
		},
		{
			name:   "non-positive half-life",
			mutate: func(seed *ResearchSeed) { seed.Records[0].HalfLifeSeconds = 0 },
			want:   "288Mc half_life_seconds must be positive",
		},
		{
			name:   "non-positive Q alpha",
			mutate: func(seed *ResearchSeed) { seed.Records[0].QAlphaMeV = 0 },
			want:   "288Mc q_alpha_mev must be positive",
		},
		{
			name:   "negative Q alpha uncertainty",
			mutate: func(seed *ResearchSeed) { seed.Records[0].QAlphaUncertaintyMeV = -0.01 },
			want:   "288Mc q_alpha_uncertainty_mev must be non-negative",
		},
		{
			name:   "negative alpha energy uncertainty",
			mutate: func(seed *ResearchSeed) { seed.Records[0].AlphaEnergyUncertaintyMeV = -0.01 },
			want:   "288Mc alpha_energy_uncertainty_mev must be non-negative",
		},
		{
			name:   "non-positive alpha energy",
			mutate: func(seed *ResearchSeed) { seed.Records[0].AlphaEnergyMeV = -10.65 },
			want:   "288Mc alpha_energy_mev must be positive when present",
		},
		{
			name:   "alpha energy range must have two bounds",
			mutate: func(seed *ResearchSeed) { seed.Records[1].AlphaEnergyRangeMeV = []float64{9.78} },
			want:   "290Mc alpha_energy_range_mev must contain exactly lower and upper bounds",
		},
		{
			name:   "alpha energy range must be ascending",
			mutate: func(seed *ResearchSeed) { seed.Records[1].AlphaEnergyRangeMeV = []float64{10.31, 9.78} },
			want:   "290Mc alpha_energy_range_mev lower bound 10.31 MeV must be less than upper bound 9.78 MeV",
		},
		{
			name:   "alpha energy range values must be positive",
			mutate: func(seed *ResearchSeed) { seed.Records[1].AlphaEnergyRangeMeV = []float64{0, 10.31} },
			want:   "290Mc alpha_energy_range_mev values must be positive",
		},
		{
			name:   "unsupported decay mode",
			mutate: func(seed *ResearchSeed) { seed.Records[0].DecayMode = "" },
			want:   "288Mc decay_mode must be alpha",
		},
		{
			name:   "missing alpha daughter",
			mutate: func(seed *ResearchSeed) { seed.Records[0].Daughter = "" },
			want:   "288Mc missing alpha daughter",
		},
		{
			name:   "daughter mass not alpha decay",
			mutate: func(seed *ResearchSeed) { seed.Records[0].Daughter = "285Nh" },
			want:   "288Mc alpha daughter 285Nh has A=285, want 284",
		},
		{
			name:   "daughter symbol not alpha decay",
			mutate: func(seed *ResearchSeed) { seed.Records[0].Daughter = "284Fl" },
			want:   "288Mc alpha daughter 284Fl has Z=114, want 113",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			candidate := seed.Clone()
			tc.mutate(&candidate)
			err := ValidateResearchSeedProvenance(candidate)
			if err == nil {
				t.Fatal("ValidateResearchSeedProvenance accepted invalid provenance")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tc.want)
			}
		})
	}
}

func TestResearchSeedCatalogConvertsVerifiedRecordsWithProvenance(t *testing.T) {
	raw, err := os.ReadFile("../../data/research.seed.json")
	if err != nil {
		t.Fatalf("read seed: %v", err)
	}
	seed := readResearchSeed(t, raw)

	catalog, err := seed.Catalog("data/research.seed.json")
	if err != nil {
		t.Fatalf("Catalog returned error: %v", err)
	}

	if len(catalog) != 2 {
		t.Fatalf("catalog record count = %d, want 2", len(catalog))
	}

	mc288, ok := catalog["288Mc"]
	if !ok {
		t.Fatal("catalog missing 288Mc")
	}
	if mc288.Symbol != "Mc" || mc288.Z != 115 || mc288.A != 288 {
		t.Fatalf("288Mc identity = symbol %s Z %d A %d, want Mc 115 288", mc288.Symbol, mc288.Z, mc288.A)
	}
	if mc288.HalfLife.String() != "170ms" {
		t.Fatalf("288Mc half-life = %s, want 170ms", mc288.HalfLife)
	}
	if mc288.QAlphaMeV != 10.75 {
		t.Fatalf("288Mc Q alpha = %.2f MeV, want 10.75 MeV", mc288.QAlphaMeV)
	}
	if mc288.Daughter != "284Nh" {
		t.Fatalf("288Mc daughter = %q, want 284Nh", mc288.Daughter)
	}
	if mc288.CitationLink != "https://www.nndc.bnl.gov/ensnds/288/Mc/adopted.pdf" {
		t.Fatalf("288Mc citation link = %q, want first source URL", mc288.CitationLink)
	}

	mc290 := catalog["290Mc"]
	if mc290.HalfLife.String() != "650ms" {
		t.Fatalf("290Mc half-life = %s, want 650ms", mc290.HalfLife)
	}
}

func TestResearchSeedCatalogRejectsInvalidSeedBeforeConversion(t *testing.T) {
	raw, err := os.ReadFile("../../data/research.seed.json")
	if err != nil {
		t.Fatalf("read seed: %v", err)
	}
	seed := readResearchSeed(t, raw)
	seed.Records[0].Z = 113

	_, err = seed.Catalog("data/research.seed.json")
	if err == nil {
		t.Fatal("Catalog accepted invalid seed")
	}
	if !strings.Contains(err.Error(), "288Mc has Z=113, want 115 for symbol Mc") {
		t.Fatalf("error = %q, want atomic-number/symbol mismatch", err)
	}
}

func readResearchSeed(t *testing.T, raw []byte) ResearchSeed {
	t.Helper()
	var seed ResearchSeed
	if err := json.Unmarshal(raw, &seed); err != nil {
		t.Fatalf("unmarshal seed: %v", err)
	}
	return seed
}
