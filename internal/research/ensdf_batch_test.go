package research

import (
	"encoding/json"
	"os"
	"testing"
)

func TestBatchENSDFToSeed_AppendsNewRecords(t *testing.T) {
	dir := t.TempDir()
	seedPath := dir + "/research.seed.json"
	original := `{"schema":"moscovium-statera-go/research-seed/v1","notes":["test"],"records":[]}`
	if err := os.WriteFile(seedPath, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	entries := []ENSDFBatchEntry{
		{
			NuclideID: "284Nh", HalfLifeSeconds: 0.98, QAlphaMeV: 10.03,
			AlphaEnergyMeV: 9.93, Daughter: "280Rg", DecayMode: "alpha",
			CitationURL: "https://www.nndc.bnl.gov/ensnds/284/Nh/adopted.pdf",
			DOIs:        []string{"10.1103/PhysRevC.99.054306"},
		},
	}
	if err := BatchENSDFToSeed(seedPath, entries); err != nil {
		t.Fatalf("batch intake: %v", err)
	}

	raw, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatal(err)
	}
	var seed ResearchSeed
	if err := json.Unmarshal(raw, &seed); err != nil {
		t.Fatal(err)
	}
	if len(seed.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(seed.Records))
	}
	if seed.Records[0].ID != "284Nh" {
		t.Errorf("expected 284Nh, got %s", seed.Records[0].ID)
	}
	if seed.Records[0].Symbol != "Nh" {
		t.Errorf("expected Nh, got %s", seed.Records[0].Symbol)
	}
	if seed.Records[0].Z != 113 {
		t.Errorf("expected Z=113, got Z=%d", seed.Records[0].Z)
	}
	if len(seed.Records[0].DOIs) != 1 {
		t.Errorf("expected 1 DOI, got %d", len(seed.Records[0].DOIs))
	}
}

func TestBatchENSDFToSeed_Idempotent(t *testing.T) {
	dir := t.TempDir()
	seedPath := dir + "/research.seed.json"
	original := `{"schema":"moscovium-statera-go/research-seed/v1","notes":["test"],"records":[{"id":"288Mc","element":"Moscovium","symbol":"Mc","z":115,"a":288,"n":173,"half_life_seconds":0.17,"q_alpha_mev":10.75,"q_alpha_uncertainty_mev":0.05,"alpha_energy_mev":10.65,"decay_mode":"alpha","daughter":"284Nh","evidence_level":"test","citation_urls":["https://www.nndc.bnl.gov/ensnds"],"dois":["10.1103/PhysRevC.106.L031301"]}]}`
	if err := os.WriteFile(seedPath, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	// Try to add 288Mc again — should be idempotent
	entries := []ENSDFBatchEntry{
		{
			NuclideID: "288Mc", HalfLifeSeconds: 0.99, QAlphaMeV: 10.75,
			AlphaEnergyMeV: 10.65, Daughter: "284Nh", DecayMode: "alpha",
			CitationURL: "https://test",
			DOIs:        []string{"10.0000/test"},
		},
	}
	if err := BatchENSDFToSeed(seedPath, entries); err != nil {
		t.Fatalf("batch intake: %v", err)
	}

	raw, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatal(err)
	}
	var seed ResearchSeed
	if err := json.Unmarshal(raw, &seed); err != nil {
		t.Fatal(err)
	}
	// Should still have exactly 1 record (no duplicate)
	if len(seed.Records) != 1 {
		t.Fatalf("expected 1 record (idempotent), got %d", len(seed.Records))
	}
	// Original values preserved
	if seed.Records[0].HalfLifeSeconds != 0.17 {
		t.Errorf("expected original half-life 0.17, got %f", seed.Records[0].HalfLifeSeconds)
	}
}
