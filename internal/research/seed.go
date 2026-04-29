package research

import (
	"fmt"
	"strings"
)

type ResearchSeed struct {
	Schema  string               `json:"schema"`
	Notes   []string             `json:"notes"`
	Records []ResearchSeedRecord `json:"records"`
}

type ResearchSeedRecord struct {
	ID                        string    `json:"id"`
	Element                   string    `json:"element"`
	Symbol                    string    `json:"symbol"`
	Z                         int       `json:"z"`
	A                         int       `json:"a"`
	N                         int       `json:"n"`
	HalfLifeSeconds           float64   `json:"half_life_seconds"`
	HalfLifeUpperSeconds      float64   `json:"half_life_upper_seconds,omitempty"`
	HalfLifeLowerSeconds      float64   `json:"half_life_lower_seconds,omitempty"`
	QAlphaMeV                 float64   `json:"q_alpha_mev"`
	QAlphaUncertaintyMeV      float64   `json:"q_alpha_uncertainty_mev"`
	AlphaEnergyMeV            float64   `json:"alpha_energy_mev,omitempty"`
	AlphaEnergyUncertaintyMeV float64   `json:"alpha_energy_uncertainty_mev,omitempty"`
	AlphaEnergyRangeMeV       []float64 `json:"alpha_energy_range_mev,omitempty"`
	DecayMode                 string    `json:"decay_mode"`
	Daughter                  string    `json:"daughter"`
	EvidenceLevel             string    `json:"evidence_level"`
	CitationURLs              []string  `json:"citation_urls"`
	DOIs                      []string  `json:"dois"`
}

func (seed ResearchSeed) Clone() ResearchSeed {
	clone := seed
	clone.Notes = append([]string(nil), seed.Notes...)
	clone.Records = append([]ResearchSeedRecord(nil), seed.Records...)
	for i := range seed.Records {
		clone.Records[i].AlphaEnergyRangeMeV = append([]float64(nil), seed.Records[i].AlphaEnergyRangeMeV...)
		clone.Records[i].CitationURLs = append([]string(nil), seed.Records[i].CitationURLs...)
		clone.Records[i].DOIs = append([]string(nil), seed.Records[i].DOIs...)
	}
	return clone
}

func ValidateResearchSeedProvenance(seed ResearchSeed) error {
	for _, record := range seed.Records {
		if record.N != record.A-record.Z {
			return fmt.Errorf("%s has N=%d, want A-Z=%d", record.ID, record.N, record.A-record.Z)
		}
		if len(record.CitationURLs) == 0 {
			return fmt.Errorf("%s missing citation_urls", record.ID)
		}
		for _, citationURL := range record.CitationURLs {
			if !strings.HasPrefix(strings.TrimSpace(citationURL), "https://") {
				return fmt.Errorf("%s citation URL %q must use https://", record.ID, citationURL)
			}
		}
		if len(record.DOIs) == 0 {
			return fmt.Errorf("%s missing DOI trail", record.ID)
		}
		for _, doi := range record.DOIs {
			if !strings.HasPrefix(strings.TrimSpace(doi), "10.") || strings.ContainsAny(strings.TrimSpace(doi), " \t\n\r") {
				return fmt.Errorf("%s DOI %q must be a DOI string beginning with 10. and containing no whitespace", record.ID, doi)
			}
		}
	}
	return nil
}
