package research

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/physics"
)

const ENSDFCandidateSchemaV1 = "moscovium-statera-go/ensdf-candidate/v1"

const (
	ENSDFCandidateStatusCandidateOnly       = "candidate_only"
	ENSDFCandidateStatusRequiresHumanReview = "requires_human_review"
	ENSDFCandidateStatusBlockedSource       = "blocked_source"
	ENSDFCandidateStatusBlockedIdentity     = "blocked_identity"
	ENSDFCandidateStatusBlockedProvenance   = "blocked_provenance"
	ENSDFCandidateStatusAccepted            = "accepted"
)

type ENSDFCandidateSet struct {
	Schema     string                 `json:"schema"`
	Source     ENSDFCandidateSource   `json:"source"`
	Candidates []ENSDFCandidateRecord `json:"candidates"`
}

type ENSDFCandidateSource struct {
	SourceID        string `json:"source_id"`
	SourceType      string `json:"source_type"`
	SourceURL       string `json:"source_url"`
	RetrievedAt     string `json:"retrieved_at"`
	RetrievalStatus string `json:"retrieval_status"`
	ContentSHA256   string `json:"content_sha256,omitempty"`
	ContentBytes    int64  `json:"content_bytes,omitempty"`
	SourcePath      string `json:"source_path,omitempty"`
}

type ENSDFCandidateRecord struct {
	CandidateID     string                    `json:"candidate_id"`
	ParentNuclideID string                    `json:"parent_nuclide_id"`
	Element         string                    `json:"element"`
	Symbol          string                    `json:"symbol"`
	Z               int                       `json:"z"`
	A               int                       `json:"a"`
	N               int                       `json:"n"`
	DatasetKind     string                    `json:"dataset_kind"`
	DecayMode       string                    `json:"decay_mode"`
	Daughter        string                    `json:"daughter,omitempty"`
	HalfLife        ENSDFCandidateMeasurement `json:"half_life"`
	QAlpha          ENSDFCandidateMeasurement `json:"q_alpha"`
	AlphaEnergy     ENSDFCandidateMeasurement `json:"alpha_energy"`
	CitationURLs    []string                  `json:"citation_urls"`
	DOIs            []string                  `json:"dois"`
	References      []string                  `json:"references,omitempty"`
	ParseWarnings   []string                  `json:"parse_warnings,omitempty"`
	ReviewStatus    string                    `json:"review_status"`
	EvidenceLevel   string                    `json:"evidence_level"`
}

type ENSDFCandidateMeasurement struct {
	ValueSeconds       *float64  `json:"value_seconds,omitempty"`
	UncertaintySeconds *float64  `json:"uncertainty_seconds,omitempty"`
	LowerSeconds       *float64  `json:"lower_seconds,omitempty"`
	UpperSeconds       *float64  `json:"upper_seconds,omitempty"`
	ValueMeV           *float64  `json:"value_mev,omitempty"`
	UncertaintyMeV     *float64  `json:"uncertainty_mev,omitempty"`
	RangeMeV           []float64 `json:"range_mev,omitempty"`
	ENSDFQualifier     string    `json:"ensdf_qualifier"`
	RawText            string    `json:"raw_text"`
}

func (set ENSDFCandidateSet) Clone() ENSDFCandidateSet {
	clone := set
	clone.Candidates = append([]ENSDFCandidateRecord(nil), set.Candidates...)
	for i := range clone.Candidates {
		clone.Candidates[i].HalfLife = clone.Candidates[i].HalfLife.clone()
		clone.Candidates[i].QAlpha = clone.Candidates[i].QAlpha.clone()
		clone.Candidates[i].AlphaEnergy = clone.Candidates[i].AlphaEnergy.clone()
		clone.Candidates[i].CitationURLs = append([]string(nil), set.Candidates[i].CitationURLs...)
		clone.Candidates[i].DOIs = append([]string(nil), set.Candidates[i].DOIs...)
		clone.Candidates[i].References = append([]string(nil), set.Candidates[i].References...)
		clone.Candidates[i].ParseWarnings = append([]string(nil), set.Candidates[i].ParseWarnings...)
	}
	return clone
}

func (measurement ENSDFCandidateMeasurement) clone() ENSDFCandidateMeasurement {
	clone := measurement
	clone.ValueSeconds = cloneFloat64Ptr(measurement.ValueSeconds)
	clone.UncertaintySeconds = cloneFloat64Ptr(measurement.UncertaintySeconds)
	clone.LowerSeconds = cloneFloat64Ptr(measurement.LowerSeconds)
	clone.UpperSeconds = cloneFloat64Ptr(measurement.UpperSeconds)
	clone.ValueMeV = cloneFloat64Ptr(measurement.ValueMeV)
	clone.UncertaintyMeV = cloneFloat64Ptr(measurement.UncertaintyMeV)
	clone.RangeMeV = append([]float64(nil), measurement.RangeMeV...)
	return clone
}

func cloneFloat64Ptr(v *float64) *float64 {
	if v == nil {
		return nil
	}
	clone := *v
	return &clone
}

func ValidateENSDFCandidateSet(set ENSDFCandidateSet) error {
	if strings.TrimSpace(set.Schema) != ENSDFCandidateSchemaV1 {
		return fmt.Errorf("ENSDF candidate schema %q unsupported, want %s", set.Schema, ENSDFCandidateSchemaV1)
	}
	if err := validateENSDFCandidateSource(set.Source); err != nil {
		return err
	}
	if len(set.Candidates) == 0 {
		return fmt.Errorf("ENSDF candidates must not be empty")
	}

	seenIDs := map[string]bool{}
	for _, candidate := range set.Candidates {
		id := strings.TrimSpace(candidate.CandidateID)
		if id == "" {
			return fmt.Errorf("ENSDF candidate_id must not be blank")
		}
		if seenIDs[id] {
			return fmt.Errorf("duplicate ENSDF candidate_id %s", id)
		}
		seenIDs[id] = true
		if err := validateENSDFCandidateRecord(candidate); err != nil {
			return err
		}
	}
	return nil
}

var sha256HexPattern = regexp.MustCompile(`^[a-f0-9]{64}$`)

func validateENSDFCandidateSource(source ENSDFCandidateSource) error {
	if strings.TrimSpace(source.SourceID) == "" {
		return fmt.Errorf("ENSDF candidate source_id must not be blank")
	}
	if strings.TrimSpace(source.SourceType) == "" {
		return fmt.Errorf("ENSDF candidate source_type must not be blank")
	}
	if !strings.HasPrefix(strings.TrimSpace(source.SourceURL), "https://") {
		return fmt.Errorf("ENSDF candidate source_url must be non-blank HTTPS")
	}
	if strings.TrimSpace(source.RetrievedAt) == "" {
		return fmt.Errorf("ENSDF candidate retrieved_at must not be blank")
	}
	if strings.TrimSpace(source.RetrievalStatus) == "" {
		return fmt.Errorf("ENSDF candidate retrieval_status must not be blank")
	}
	if strings.TrimSpace(source.SourcePath) != "" {
		if !sha256HexPattern.MatchString(strings.TrimSpace(source.ContentSHA256)) {
			return fmt.Errorf("ENSDF candidate source content_sha256 must be a lowercase SHA-256 hex string when source_path is present")
		}
		if source.ContentBytes <= 0 {
			return fmt.Errorf("ENSDF candidate source content_bytes must be positive when source_path is present")
		}
	}
	return nil
}

func validateENSDFCandidateRecord(candidate ENSDFCandidateRecord) error {
	nuclide, err := physics.ParseNuclideID(candidate.ParentNuclideID)
	if err != nil {
		return err
	}
	id := strings.TrimSpace(candidate.CandidateID)
	if candidate.Symbol != nuclide.Symbol {
		return fmt.Errorf("%s parent %s symbol mismatch, want %s from ID", id, candidate.ParentNuclideID, nuclide.Symbol)
	}
	if candidate.A != nuclide.A {
		return fmt.Errorf("%s parent %s has A=%d, want %d from ID", id, candidate.ParentNuclideID, candidate.A, nuclide.A)
	}
	if candidate.Z != nuclide.Z {
		return fmt.Errorf("%s parent %s has Z=%d, want %d for symbol %s", id, candidate.ParentNuclideID, candidate.Z, nuclide.Z, candidate.Symbol)
	}
	if candidate.N != candidate.A-candidate.Z {
		return fmt.Errorf("%s parent %s has N=%d, want A-Z=%d", id, candidate.ParentNuclideID, candidate.N, candidate.A-candidate.Z)
	}
	expectedElement, ok := physics.ElementNameForSymbol(candidate.Symbol)
	if !ok {
		return fmt.Errorf("%s parent %s has unknown element symbol %s", id, candidate.ParentNuclideID, candidate.Symbol)
	}
	if strings.TrimSpace(candidate.Element) != expectedElement {
		return fmt.Errorf("%s element name %q does not match symbol %s, want %s", id, candidate.Element, candidate.Symbol, expectedElement)
	}
	if strings.TrimSpace(candidate.DecayMode) == "alpha" && strings.TrimSpace(candidate.Daughter) != "" {
		if err := physics.ValidateAlphaDaughterID(candidate.ParentNuclideID, candidate.Daughter); err != nil {
			return err
		}
	}
	if err := validateENSDFMeasurement(id, "half_life", candidate.HalfLife, candidate.ReviewStatus); err != nil {
		return err
	}
	if err := validateENSDFMeasurement(id, "q_alpha", candidate.QAlpha, candidate.ReviewStatus); err != nil {
		return err
	}
	if err := validateENSDFMeasurement(id, "alpha_energy", candidate.AlphaEnergy, candidate.ReviewStatus); err != nil {
		return err
	}
	if err := validateStringSet(id, "citation_urls", candidate.CitationURLs, true); err != nil {
		return err
	}
	if err := validateStringSet(id, "dois", candidate.DOIs, false); err != nil {
		return err
	}
	if strings.TrimSpace(candidate.EvidenceLevel) == "" {
		return fmt.Errorf("%s evidence_level must not be blank", id)
	}
	status := strings.TrimSpace(candidate.ReviewStatus)
	if status == ENSDFCandidateStatusAccepted {
		return fmt.Errorf("%s review_status accepted is not allowed in candidate-only validation", id)
	}
	if !allowedENSDFCandidateStatus(status) {
		return fmt.Errorf("%s review_status %q unsupported", id, candidate.ReviewStatus)
	}
	for i, warning := range candidate.ParseWarnings {
		if strings.TrimSpace(warning) == "" {
			return fmt.Errorf("%s parse_warnings[%d] must not be blank", id, i)
		}
	}
	if len(candidate.ParseWarnings) > 0 && status != ENSDFCandidateStatusRequiresHumanReview {
		return fmt.Errorf("%s parse_warnings require review_status requires_human_review", id)
	}
	return nil
}

func validateENSDFMeasurement(candidateID, field string, measurement ENSDFCandidateMeasurement, reviewStatus string) error {
	qualifier := strings.TrimSpace(measurement.ENSDFQualifier)
	if qualifier == "" {
		qualifier = "="
	}
	if !allowedENSDFQualifier(qualifier) {
		return fmt.Errorf("%s %s ensdf_qualifier %q unsupported", candidateID, field, measurement.ENSDFQualifier)
	}
	if requiresHumanReviewQualifier(qualifier) && strings.TrimSpace(reviewStatus) != ENSDFCandidateStatusRequiresHumanReview {
		return fmt.Errorf("%s %s qualifier %s requires review_status requires_human_review", candidateID, field, qualifier)
	}
	if measurement.ValueSeconds != nil && *measurement.ValueSeconds <= 0 {
		return fmt.Errorf("%s %s value_seconds must be positive when present", candidateID, field)
	}
	if measurement.UncertaintySeconds != nil && *measurement.UncertaintySeconds < 0 {
		return fmt.Errorf("%s %s uncertainty_seconds must be non-negative when present", candidateID, field)
	}
	if (measurement.LowerSeconds == nil) != (measurement.UpperSeconds == nil) {
		return fmt.Errorf("%s %s interval must include both lower_seconds and upper_seconds when either is present", candidateID, field)
	}
	if measurement.LowerSeconds != nil {
		if *measurement.LowerSeconds <= 0 || *measurement.UpperSeconds <= 0 {
			return fmt.Errorf("%s %s lower_seconds and upper_seconds must be positive when present", candidateID, field)
		}
		if *measurement.LowerSeconds >= *measurement.UpperSeconds {
			return fmt.Errorf("%s %s lower_seconds must be less than upper_seconds", candidateID, field)
		}
	}
	if measurement.ValueMeV != nil && *measurement.ValueMeV <= 0 {
		return fmt.Errorf("%s %s value_mev must be positive when present", candidateID, field)
	}
	if measurement.UncertaintyMeV != nil && *measurement.UncertaintyMeV < 0 {
		return fmt.Errorf("%s %s uncertainty_mev must be non-negative when present", candidateID, field)
	}
	if len(measurement.RangeMeV) > 0 {
		if len(measurement.RangeMeV) != 2 {
			return fmt.Errorf("%s %s range_mev must contain exactly lower and upper bounds", candidateID, field)
		}
		if measurement.RangeMeV[0] <= 0 || measurement.RangeMeV[1] <= 0 {
			return fmt.Errorf("%s %s range_mev values must be positive", candidateID, field)
		}
		if measurement.RangeMeV[0] >= measurement.RangeMeV[1] {
			return fmt.Errorf("%s %s range_mev lower bound must be less than upper bound", candidateID, field)
		}
	}
	return nil
}

func validateStringSet(candidateID, field string, values []string, requireHTTPS bool) error {
	if len(values) == 0 {
		return fmt.Errorf("%s missing %s", candidateID, field)
	}
	seen := map[string]bool{}
	for i, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return fmt.Errorf("%s %s[%d] must not be blank", candidateID, field, i)
		}
		if seen[trimmed] {
			return fmt.Errorf("%s duplicate %s entry %q", candidateID, field, trimmed)
		}
		seen[trimmed] = true
		if requireHTTPS && !strings.HasPrefix(trimmed, "https://") {
			return fmt.Errorf("%s %s[%d] must use https://", candidateID, field, i)
		}
		if !requireHTTPS && (!strings.HasPrefix(trimmed, "10.") || strings.ContainsAny(trimmed, " \t\n\r")) {
			return fmt.Errorf("%s DOI %q must be a DOI string beginning with 10. and containing no whitespace", candidateID, value)
		}
	}
	return nil
}

func allowedENSDFQualifier(qualifier string) bool {
	switch qualifier {
	case "=", "AP", "LT", "GT", "LE", "GE":
		return true
	default:
		return false
	}
}

func requiresHumanReviewQualifier(qualifier string) bool {
	switch qualifier {
	case "AP", "LT", "GT", "LE", "GE":
		return true
	default:
		return false
	}
}

func allowedENSDFCandidateStatus(status string) bool {
	switch status {
	case ENSDFCandidateStatusCandidateOnly, ENSDFCandidateStatusRequiresHumanReview, ENSDFCandidateStatusBlockedSource, ENSDFCandidateStatusBlockedIdentity, ENSDFCandidateStatusBlockedProvenance:
		return true
	default:
		return false
	}
}
