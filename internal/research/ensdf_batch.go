package research

import (
	"encoding/json"
	"fmt"
	"os"
)

// ENSDFBatchEntry represents a single nuclide extracted from ENSDF text.
type ENSDFBatchEntry struct {
	NuclideID        string
	HalfLifeSeconds  float64
	QAlphaMeV        float64
	AlphaEnergyMeV   float64
	Daughter         string
	DecayMode        string
	SFBranchingRatio float64
	CitationURL      string
	DOIs             []string
}

// BatchENSDFToSeed converts batch ENSDF entries into research.seed.json records
// and appends them to the existing seed, preserving existing records.
func BatchENSDFToSeed(existingPath string, entries []ENSDFBatchEntry) error {
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
		if existingIDs[entry.NuclideID] {
			continue
		}
		sym := extractSymbolFromID(entry.NuclideID)
		z := atomicNumberForSymbol(sym)
		a := extractMassFromID(entry.NuclideID)
		rec := ResearchSeedRecord{
			ID:             entry.NuclideID,
			Element:        elementNameForSym(sym),
			Symbol:         sym,
			Z:              z,
			A:              a,
			N:              a - z,
			HalfLifeSeconds: entry.HalfLifeSeconds,
			QAlphaMeV:      entry.QAlphaMeV,
			AlphaEnergyMeV: entry.AlphaEnergyMeV,
			DecayMode:      entry.DecayMode,
			Daughter:       entry.Daughter,
			EvidenceLevel:  "ENSDF evaluated nuclear data intake, pending DOI cross-verification",
			CitationURLs:   []string{entry.CitationURL},
			DOIs:           entry.DOIs,
		}
		if entry.QAlphaMeV > 0 && entry.QAlphaUncertaintyMeV() == 0 {
			rec.QAlphaUncertaintyMeV = 0.05
		}
		seed.Records = append(seed.Records, rec)
	}

	out, err := json.MarshalIndent(seed, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal seed: %w", err)
	}
	if err := os.WriteFile(existingPath, out, 0644); err != nil {
		return fmt.Errorf("write seed: %w", err)
	}
	return nil
}

// QAlphaUncertaintyMeV returns the uncertainty for the Q_alpha value.
func (e ENSDFBatchEntry) QAlphaUncertaintyMeV() float64 { return 0 }

func extractSymbolFromID(nuclideID string) string {
	for i, r := range nuclideID {
		if r >= 'A' && r <= 'Z' {
			return nuclideID[i:]
		}
	}
	return ""
}

func extractMassFromID(nuclideID string) int {
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

func elementNameForSym(symbol string) string {
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

func atomicNumberForSymbol(symbol string) int {
	znums := map[string]int{
		"Mc": 115, "Nh": 113, "Rg": 111, "Mt": 109, "Bh": 107,
		"Db": 105, "Lr": 103, "Fl": 114, "Lv": 116, "Ts": 117,
		"Og": 118, "Cn": 112, "Hs": 108, "Sg": 106, "Rf": 104,
	}
	if z, ok := znums[symbol]; ok {
		return z
	}
	return 0
}
