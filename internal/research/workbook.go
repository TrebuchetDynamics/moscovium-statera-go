package research

import "strings"

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
	return records, nil
}
