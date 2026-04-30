package research

import "strings"

type SourceTriageClass string

const (
	SourceTriageClassEvaluatedData           SourceTriageClass = "evaluated-data"
	SourceTriageClassExperimentalPaper       SourceTriageClass = "experimental-paper"
	SourceTriageClassTheoreticalModel        SourceTriageClass = "theoretical-model"
	SourceTriageClassReview                  SourceTriageClass = "review"
	SourceTriageClassContextNotForSimulation SourceTriageClass = "context-not-for-simulation"
	SourceTriageClassBlockedMetadataOnly     SourceTriageClass = "blocked-metadata-only"
	SourceTriageClassUnknownNeedsHumanReview SourceTriageClass = "unknown-needs-human-review"
)

type SourceTriageInput struct {
	RecordID     string
	Title        string
	DOI          string
	SourceURL    string
	AccessStatus string
	Text         string
}

type SourceTriageDecision struct {
	RecordID            string
	Class               SourceTriageClass
	CandidateOnly       bool
	RequiresHumanReview bool
}

func ClassifySourceTriage(input SourceTriageInput) SourceTriageDecision {
	class := classifySourceTriageClass(input)
	return SourceTriageDecision{
		RecordID:            strings.TrimSpace(input.RecordID),
		Class:               class,
		CandidateOnly:       true,
		RequiresHumanReview: classRequiresHumanReview(class),
	}
}

func classifySourceTriageClass(input SourceTriageInput) SourceTriageClass {
	accessStatus := strings.ToLower(strings.TrimSpace(input.AccessStatus))
	if accessStatus == "blocked-access" || accessStatus == "metadata-only" || accessStatus == "blocked" {
		return SourceTriageClassBlockedMetadataOnly
	}

	text := strings.ToLower(strings.Join([]string{input.RecordID, input.Title, input.SourceURL, input.Text}, " "))
	if containsAny(text, "oversight.house.gov", "house oversight", "hearing transcript", "committee oversight") {
		return SourceTriageClassContextNotForSimulation
	}
	if containsAny(text, "nudat", "ensdf", "nndc", "evaluated data", "adopted levels") {
		return SourceTriageClassEvaluatedData
	}
	if containsAny(text, "theoretical model", "analytical formula", "analytic", "coefficients", "royer") {
		return SourceTriageClassTheoreticalModel
	}
	if containsAny(text, "experimental", "experiment", "decay spectroscopy", "synthesis") {
		return SourceTriageClassExperimentalPaper
	}
	if containsAny(text, "review", "survey") {
		return SourceTriageClassReview
	}
	return SourceTriageClassUnknownNeedsHumanReview
}

func classRequiresHumanReview(class SourceTriageClass) bool {
	switch class {
	case SourceTriageClassBlockedMetadataOnly, SourceTriageClassContextNotForSimulation, SourceTriageClassUnknownNeedsHumanReview:
		return true
	default:
		return false
	}
}

func containsAny(text string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(text, needle) {
			return true
		}
	}
	return false
}
