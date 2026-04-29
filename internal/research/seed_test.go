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
			name:   "unsupported decay mode",
			mutate: func(seed *ResearchSeed) { seed.Records[0].DecayMode = "" },
			want:   "288Mc decay_mode must be alpha",
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

func readResearchSeed(t *testing.T, raw []byte) ResearchSeed {
	t.Helper()
	var seed ResearchSeed
	if err := json.Unmarshal(raw, &seed); err != nil {
		t.Fatalf("unmarshal seed: %v", err)
	}
	return seed
}
