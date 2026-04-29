package research

import (
	"fmt"
	"strings"

	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/physics"
)

// WorkbookRecord is a UI-neutral, source-backed isotope workbook row derived
// from validated research seed records. It must not add or infer scientific
// values beyond the source seed record.
type WorkbookRecord struct {
	ID              string
	Element         string
	Symbol          string
	Z               int
	A               int
	N               int
	HalfLifeSeconds float64
	QAlphaMeV       float64
	Daughter        string
	EvidenceLevel   string
	SourcePath      string
	CitationURLs    []string
	DOIs            []string
}

func WorkbookFromSeed(seed ResearchSeed, sourcePath string) ([]WorkbookRecord, error) {
	if err := ValidateResearchSeedProvenance(seed); err != nil {
		return nil, err
	}

	records := make([]WorkbookRecord, 0, len(seed.Records))
	trimmedSourcePath := strings.TrimSpace(sourcePath)
	for _, seedRecord := range seed.Records {
		records = append(records, WorkbookRecord{
			ID:              seedRecord.ID,
			Element:         seedRecord.Element,
			Symbol:          seedRecord.Symbol,
			Z:               seedRecord.Z,
			A:               seedRecord.A,
			N:               seedRecord.N,
			HalfLifeSeconds: seedRecord.HalfLifeSeconds,
			QAlphaMeV:       seedRecord.QAlphaMeV,
			Daughter:        strings.TrimSpace(seedRecord.Daughter),
			EvidenceLevel:   strings.TrimSpace(seedRecord.EvidenceLevel),
			SourcePath:      trimmedSourcePath,
			CitationURLs:    append([]string(nil), seedRecord.CitationURLs...),
			DOIs:            append([]string(nil), seedRecord.DOIs...),
		})
	}
	if err := ValidateWorkbookRecords(records); err != nil {
		return nil, err
	}
	return records, nil
}

// ValidateWorkbookRecords rejects workbook rows that lost source provenance or
// nuclide identity consistency during conversion or UI-neutral handling.
func ValidateWorkbookRecords(records []WorkbookRecord) error {
	if len(records) == 0 {
		return fmt.Errorf("workbook records must not be empty")
	}
	for _, record := range records {
		id := strings.TrimSpace(record.ID)
		if id == "" {
			return fmt.Errorf("workbook record ID must not be blank")
		}
		if _, err := physics.ParseNuclideID(id); err != nil {
			return err
		}
		if record.N != record.A-record.Z {
			return fmt.Errorf("%s has N=%d, want A-Z=%d", id, record.N, record.A-record.Z)
		}
		if strings.TrimSpace(record.EvidenceLevel) == "" {
			return fmt.Errorf("%s evidence level must not be blank", id)
		}
		if strings.TrimSpace(record.SourcePath) == "" {
			return fmt.Errorf("%s source path must not be blank", id)
		}
		if len(record.CitationURLs) == 0 {
			return fmt.Errorf("%s missing citation URLs", id)
		}
		for i, citationURL := range record.CitationURLs {
			if strings.TrimSpace(citationURL) == "" {
				return fmt.Errorf("%s citation URLs[%d] must not be blank", id, i)
			}
		}
		if len(record.DOIs) == 0 {
			return fmt.Errorf("%s missing DOI trail", id)
		}
		for i, doi := range record.DOIs {
			if strings.TrimSpace(doi) == "" {
				return fmt.Errorf("%s DOIs[%d] must not be blank", id, i)
			}
		}
		daughter := strings.TrimSpace(record.Daughter)
		if daughter != "" {
			if err := physics.ValidateAlphaDaughterID(id, daughter); err != nil {
				return err
			}
		}
	}
	return nil
}
