package research

import "testing"

func TestClassifySourceTriageBaselineUsesDeterministicEvidenceClasses(t *testing.T) {
	tests := []struct {
		name   string
		input  SourceTriageInput
		want   SourceTriageClass
		review bool
	}{
		{
			name: "NuDat ENSDF page is evaluated data",
			input: SourceTriageInput{
				RecordID:     "nudat3ensdf-288mc",
				Title:        "NuDat 3 ENSDF adopted levels for 288Mc",
				SourceURL:    "https://www.nndc.bnl.gov/nudat3/",
				AccessStatus: "source-content-readable",
				Text:         "NNDC NuDat ENSDF adopted levels decay radiation evaluated data",
			},
			want: SourceTriageClassEvaluatedData,
		},
		{
			name: "Royer source content readable is theoretical model",
			input: SourceTriageInput{
				RecordID:     "royer2008alpha-analytic",
				Title:        "Analytical expressions for alpha-decay half-lives and potential barriers",
				DOI:          "10.1103/PhysRevC.77.037602",
				SourceURL:    "https://doi.org/10.1103/PhysRevC.77.037602",
				AccessStatus: "source-content-readable",
				Text:         "Royer alpha-decay half-lives analytical formula theoretical model coefficients",
			},
			want: SourceTriageClassTheoreticalModel,
		},
		{
			name: "Royer blocked source is metadata only",
			input: SourceTriageInput{
				RecordID:     "royer2008alpha-analytic",
				Title:        "Analytical expressions for alpha-decay half-lives and potential barriers",
				DOI:          "10.1103/PhysRevC.77.037602",
				SourceURL:    "https://doi.org/10.1103/PhysRevC.77.037602",
				AccessStatus: "blocked-access",
				Text:         "APS abstract visible but source PDF returned HTTP 403",
			},
			want:   SourceTriageClassBlockedMetadataOnly,
			review: true,
		},
		{
			name: "House Oversight record is context not simulation",
			input: SourceTriageInput{
				RecordID:     "house-oversight-context",
				Title:        "House Oversight hearing transcript context",
				SourceURL:    "https://oversight.house.gov/",
				AccessStatus: "source-content-readable",
				Text:         "hearing transcript committee oversight public context not evaluated nuclear data",
			},
			want:   SourceTriageClassContextNotForSimulation,
			review: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifySourceTriage(tt.input)
			if got.Class != tt.want {
				t.Fatalf("ClassifySourceTriage class = %q, want %q", got.Class, tt.want)
			}
			if !got.CandidateOnly {
				t.Fatalf("ClassifySourceTriage CandidateOnly = false, want true")
			}
			if got.RequiresHumanReview != tt.review {
				t.Fatalf("ClassifySourceTriage RequiresHumanReview = %v, want %v", got.RequiresHumanReview, tt.review)
			}
			if got.RecordID != tt.input.RecordID {
				t.Fatalf("ClassifySourceTriage RecordID = %q, want %q", got.RecordID, tt.input.RecordID)
			}
		})
	}
}

func TestClassifySourceTriageUnknownRequiresHumanReview(t *testing.T) {
	got := ClassifySourceTriage(SourceTriageInput{
		RecordID:     "ambiguous-snippet",
		Title:        "Ambiguous source snippet",
		SourceURL:    "https://example.org/source",
		AccessStatus: "source-content-readable",
		Text:         "short ambiguous text",
	})
	if got.Class != SourceTriageClassUnknownNeedsHumanReview {
		t.Fatalf("ClassifySourceTriage class = %q, want %q", got.Class, SourceTriageClassUnknownNeedsHumanReview)
	}
	if !got.CandidateOnly || !got.RequiresHumanReview {
		t.Fatalf("ClassifySourceTriage candidate flags = candidate_only:%v requires_review:%v, want both true", got.CandidateOnly, got.RequiresHumanReview)
	}
}
