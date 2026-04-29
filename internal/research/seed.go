package research

import (
	"fmt"
	"strings"
	"time"

	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/physics"
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
	seenIDs := map[string]bool{}
	for _, record := range seed.Records {
		nuclide, err := physics.ParseNuclideID(record.ID)
		if err != nil {
			return err
		}
		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate seed record ID %s", record.ID)
		}
		seenIDs[record.ID] = true
		expectedID := fmt.Sprintf("%d%s", record.A, record.Symbol)
		if record.ID != expectedID {
			return fmt.Errorf("%s ID mismatch, want %s from A=%d and symbol=%s", record.ID, expectedID, record.A, record.Symbol)
		}
		if nuclide.Symbol != record.Symbol {
			return fmt.Errorf("%s symbol mismatch, want %s from ID", record.ID, nuclide.Symbol)
		}
		if nuclide.A != record.A {
			return fmt.Errorf("%s A mismatch, want %d from ID", record.ID, nuclide.A)
		}
		if record.Z != nuclide.Z {
			return fmt.Errorf("%s has Z=%d, want %d for symbol %s", record.ID, record.Z, nuclide.Z, record.Symbol)
		}
		expectedElement, ok := physics.ElementNameForSymbol(record.Symbol)
		if !ok {
			return fmt.Errorf("%s has unknown element symbol %s", record.ID, record.Symbol)
		}
		if strings.TrimSpace(record.Element) != expectedElement {
			return fmt.Errorf("%s element name %q does not match symbol %s, want %s", record.ID, record.Element, record.Symbol, expectedElement)
		}
		if record.N != record.A-record.Z {
			return fmt.Errorf("%s has N=%d, want A-Z=%d", record.ID, record.N, record.A-record.Z)
		}
		if record.HalfLifeSeconds <= 0 {
			return fmt.Errorf("%s half_life_seconds must be positive", record.ID)
		}
		if err := validateHalfLifeInterval(record); err != nil {
			return err
		}
		if record.QAlphaMeV <= 0 {
			return fmt.Errorf("%s q_alpha_mev must be positive", record.ID)
		}
		if record.QAlphaUncertaintyMeV < 0 {
			return fmt.Errorf("%s q_alpha_uncertainty_mev must be non-negative", record.ID)
		}
		if err := validateAlphaEnergy(record); err != nil {
			return err
		}
		if strings.TrimSpace(record.DecayMode) != "alpha" {
			return fmt.Errorf("%s decay_mode must be alpha", record.ID)
		}
		if err := validateAlphaDaughter(record); err != nil {
			return err
		}
		if strings.TrimSpace(record.EvidenceLevel) == "" {
			return fmt.Errorf("%s missing evidence_level", record.ID)
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

func (seed ResearchSeed) Catalog(sourcePath string) (physics.Catalog, error) {
	if err := ValidateResearchSeedProvenance(seed); err != nil {
		return nil, err
	}

	catalog := make(physics.Catalog, len(seed.Records))
	for _, record := range seed.Records {
		citationLink := strings.TrimSpace(sourcePath)
		if len(record.CitationURLs) > 0 {
			citationLink = strings.TrimSpace(record.CitationURLs[0])
		}
		catalog[record.ID] = physics.Isotope{
			Symbol:       record.Symbol,
			Z:            record.Z,
			A:            record.A,
			HalfLife:     time.Duration(record.HalfLifeSeconds * float64(time.Second)),
			QAlphaMeV:    record.QAlphaMeV,
			Daughter:     strings.TrimSpace(record.Daughter),
			CitationLink: citationLink,
		}
	}
	return catalog, nil
}

func validateHalfLifeInterval(record ResearchSeedRecord) error {
	lower := record.HalfLifeLowerSeconds
	upper := record.HalfLifeUpperSeconds
	if lower == 0 && upper == 0 {
		return nil
	}
	if lower == 0 || upper == 0 {
		return fmt.Errorf("%s half-life interval must include both lower and upper bounds when either is present", record.ID)
	}
	if lower <= 0 || upper <= 0 {
		return fmt.Errorf("%s half_life_lower_seconds and half_life_upper_seconds must be positive when present", record.ID)
	}
	if lower >= record.HalfLifeSeconds || upper <= record.HalfLifeSeconds {
		return fmt.Errorf("%s half-life interval [%.2f, %.2f] s must bracket nominal half_life_seconds %.2f s", record.ID, lower, upper, record.HalfLifeSeconds)
	}
	return nil
}

func validateAlphaEnergy(record ResearchSeedRecord) error {
	if record.AlphaEnergyUncertaintyMeV < 0 {
		return fmt.Errorf("%s alpha_energy_uncertainty_mev must be non-negative", record.ID)
	}
	if record.AlphaEnergyMeV < 0 {
		return fmt.Errorf("%s alpha_energy_mev must be positive when present", record.ID)
	}
	if len(record.AlphaEnergyRangeMeV) == 0 {
		return nil
	}
	if len(record.AlphaEnergyRangeMeV) != 2 {
		return fmt.Errorf("%s alpha_energy_range_mev must contain exactly lower and upper bounds", record.ID)
	}
	lower := record.AlphaEnergyRangeMeV[0]
	upper := record.AlphaEnergyRangeMeV[1]
	if lower <= 0 || upper <= 0 {
		return fmt.Errorf("%s alpha_energy_range_mev values must be positive", record.ID)
	}
	if lower >= upper {
		return fmt.Errorf("%s alpha_energy_range_mev lower bound %.2f MeV must be less than upper bound %.2f MeV", record.ID, lower, upper)
	}
	return nil
}

func validateAlphaDaughter(record ResearchSeedRecord) error {
	daughter := strings.TrimSpace(record.Daughter)
	if daughter == "" {
		return fmt.Errorf("%s missing alpha daughter", record.ID)
	}

	if err := physics.ValidateAlphaDaughterID(record.ID, daughter); err != nil {
		return err
	}
	return nil
}
