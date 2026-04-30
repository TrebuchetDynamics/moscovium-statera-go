package research

import (
	"strings"
	"testing"
)

func validENSDFCandidateSet() ENSDFCandidateSet {
	return ENSDFCandidateSet{
		Schema: ENSDFCandidateSchemaV1,
		Source: ENSDFCandidateSource{
			SourceID:        "nndc-288mc-adopted-20260429",
			SourceType:      "ensdf_dataset",
			SourceURL:       "https://www.nndc.bnl.gov/ensdf/",
			RetrievedAt:     "2026-04-29T18:00:00Z",
			RetrievalStatus: "HTTP 200 metadata-only test fixture",
			ContentSHA256:   "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			ContentBytes:    1024,
			SourcePath:      "artifacts/source-cache/ensdf/288mc.txt",
		},
		Candidates: []ENSDFCandidateRecord{
			{
				CandidateID:     "candidate-288mc-alpha",
				ParentNuclideID: "288Mc",
				Element:         "Moscovium",
				Symbol:          "Mc",
				Z:               115,
				A:               288,
				N:               173,
				DatasetKind:     "adopted",
				DecayMode:       "alpha",
				Daughter:        "284Nh",
				HalfLife: ENSDFCandidateMeasurement{
					ValueSeconds:       float64Ptr(0.17),
					UncertaintySeconds: float64Ptr(0.03),
					ENSDFQualifier:     "=",
					RawText:            "0.17 s test fixture",
				},
				QAlpha: ENSDFCandidateMeasurement{
					ValueMeV:       float64Ptr(10.65),
					UncertaintyMeV: float64Ptr(0.08),
					ENSDFQualifier: "=",
					RawText:        "10.65 MeV test fixture",
				},
				AlphaEnergy: ENSDFCandidateMeasurement{
					ValueMeV:       float64Ptr(10.50),
					UncertaintyMeV: float64Ptr(0.07),
					RangeMeV:       []float64{10.40, 10.60},
					ENSDFQualifier: "=",
					RawText:        "10.50 MeV test fixture",
				},
				CitationURLs:  []string{"https://www.nndc.bnl.gov/ensdf/"},
				DOIs:          []string{"10.1103/PhysRevC.106.L031301"},
				References:    []string{"reference text preserved"},
				ParseWarnings: nil,
				ReviewStatus:  ENSDFCandidateStatusCandidateOnly,
				EvidenceLevel: "candidate-only-test-fixture",
			},
		},
	}
}

func float64Ptr(v float64) *float64 { return &v }

func TestValidateENSDFCandidateSetAcceptsStructurallyValidCandidateOnlyRecords(t *testing.T) {
	set := validENSDFCandidateSet()

	if err := ValidateENSDFCandidateSet(set); err != nil {
		t.Fatalf("valid candidate set rejected: %v", err)
	}
}

func TestValidateENSDFCandidateSetRejectsProvenanceIdentityAndStatusLoss(t *testing.T) {
	base := validENSDFCandidateSet()

	cases := []struct {
		name   string
		mutate func(*ENSDFCandidateSet)
		want   string
	}{
		{
			name:   "unsupported schema",
			mutate: func(set *ENSDFCandidateSet) { set.Schema = "moscovium-statera-go/ensdf-candidate/v2" },
			want:   "ENSDF candidate schema \"moscovium-statera-go/ensdf-candidate/v2\" unsupported, want moscovium-statera-go/ensdf-candidate/v1",
		},
		{
			name:   "missing source URL",
			mutate: func(set *ENSDFCandidateSet) { set.Source.SourceURL = "" },
			want:   "ENSDF candidate source_url must be non-blank HTTPS",
		},
		{
			name:   "source path requires content hash",
			mutate: func(set *ENSDFCandidateSet) { set.Source.ContentSHA256 = "" },
			want:   "ENSDF candidate source content_sha256 must be a lowercase SHA-256 hex string when source_path is present",
		},
		{
			name: "duplicate candidate ID",
			mutate: func(set *ENSDFCandidateSet) {
				set.Candidates = append(set.Candidates, set.Candidates[0])
			},
			want: "duplicate ENSDF candidate_id candidate-288mc-alpha",
		},
		{
			name:   "parent nuclide identity mismatch",
			mutate: func(set *ENSDFCandidateSet) { set.Candidates[0].Z = 113 },
			want:   "candidate-288mc-alpha parent 288Mc has Z=113, want 115 for symbol Mc",
		},
		{
			name:   "wrong element label",
			mutate: func(set *ENSDFCandidateSet) { set.Candidates[0].Element = "Nihonium" },
			want:   "candidate-288mc-alpha element name \"Nihonium\" does not match symbol Mc, want Moscovium",
		},
		{
			name:   "alpha daughter mismatch",
			mutate: func(set *ENSDFCandidateSet) { set.Candidates[0].Daughter = "285Nh" },
			want:   "288Mc alpha daughter 285Nh has A=285, want 284",
		},
		{
			name: "duplicate citation URL",
			mutate: func(set *ENSDFCandidateSet) {
				set.Candidates[0].CitationURLs = append(set.Candidates[0].CitationURLs, " "+set.Candidates[0].CitationURLs[0]+" ")
			},
			want: "candidate-288mc-alpha duplicate citation_urls entry",
		},
		{
			name:   "malformed DOI",
			mutate: func(set *ENSDFCandidateSet) { set.Candidates[0].DOIs[0] = "PhysRevC.106.L031301" },
			want:   "candidate-288mc-alpha DOI",
		},
		{
			name:   "unsupported qualifier",
			mutate: func(set *ENSDFCandidateSet) { set.Candidates[0].HalfLife.ENSDFQualifier = "ABOUT" },
			want:   "candidate-288mc-alpha half_life ensdf_qualifier \"ABOUT\" unsupported",
		},
		{
			name:   "limit qualifier requires human review",
			mutate: func(set *ENSDFCandidateSet) { set.Candidates[0].HalfLife.ENSDFQualifier = "LT" },
			want:   "candidate-288mc-alpha half_life qualifier LT requires review_status requires_human_review",
		},
		{
			name:   "parse warnings require human review",
			mutate: func(set *ENSDFCandidateSet) { set.Candidates[0].ParseWarnings = []string{"ambiguous token"} },
			want:   "candidate-288mc-alpha parse_warnings require review_status requires_human_review",
		},
		{
			name:   "accepted status blocked at candidate layer",
			mutate: func(set *ENSDFCandidateSet) { set.Candidates[0].ReviewStatus = ENSDFCandidateStatusAccepted },
			want:   "candidate-288mc-alpha review_status accepted is not allowed in candidate-only validation",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			candidate := base.Clone()
			tc.mutate(&candidate)
			err := ValidateENSDFCandidateSet(candidate)
			if err == nil {
				t.Fatal("ValidateENSDFCandidateSet accepted invalid candidate")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tc.want)
			}
		})
	}
}
