package research

import (
	"encoding/json"
	"os"
	"testing"
)

func TestResearchSeedContainsVerifiedMoscoviumRecords(t *testing.T) {
	raw, err := os.ReadFile("../../data/research.seed.json")
	if err != nil {
		t.Fatalf("read seed: %v", err)
	}

	var seed struct {
		Records []struct {
			ID           string   `json:"id"`
			Z            int      `json:"z"`
			A            int      `json:"a"`
			Daughter     string   `json:"daughter"`
			CitationURLs []string `json:"citation_urls"`
			DOIs         []string `json:"dois"`
		} `json:"records"`
	}
	if err := json.Unmarshal(raw, &seed); err != nil {
		t.Fatalf("unmarshal seed: %v", err)
	}

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
