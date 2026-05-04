package research

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestWorkbookFromSeedPreservesProvenance(t *testing.T) {
	workbook := validWorkbookRecords(t)
	if len(workbook) < 2 {
		t.Fatalf("workbook record count = %d, want at least 2", len(workbook))
	}

	byID := map[string]WorkbookRecord{}
	for _, record := range workbook {
		byID[record.ID] = record
		if record.ID == "" {
			t.Fatal("workbook record missing ID")
		}
		if record.Z == 0 {
			t.Fatalf("%s missing Z", record.ID)
		}
		if record.A == 0 {
			t.Fatalf("%s missing A", record.ID)
		}
		if record.N == 0 {
			t.Fatalf("%s missing N", record.ID)
		}
		if record.HalfLifeSeconds <= 0 {
			t.Fatalf("%s half-life seconds = %.6g, want positive", record.ID, record.HalfLifeSeconds)
		}
		if record.QAlphaMeV <= 0 {
			t.Fatalf("%s Q-alpha MeV = %.6g, want positive", record.ID, record.QAlphaMeV)
		}
		if record.Daughter == "" {
			t.Fatalf("%s missing daughter", record.ID)
		}
		if record.EvidenceLevel == "" {
			t.Fatalf("%s missing evidence level", record.ID)
		}
		if record.SourcePath != "data/research.seed.json" {
			t.Fatalf("%s source path = %q, want data/research.seed.json", record.ID, record.SourcePath)
		}
		if len(record.CitationURLs) == 0 {
			t.Fatalf("%s missing citation URLs", record.ID)
		}
		if len(record.DOIs) == 0 {
			t.Fatalf("%s missing DOI trail", record.ID)
		}
	}

	mc288, ok := byID["288Mc"]
	if !ok {
		t.Fatal("workbook missing 288Mc")
	}
	if mc288.Z != 115 || mc288.A != 288 || mc288.N != 173 {
		t.Fatalf("288Mc identity = Z=%d A=%d N=%d, want 115 288 173", mc288.Z, mc288.A, mc288.N)
	}
	if mc288.HalfLifeSeconds != 0.17 {
		t.Fatalf("288Mc half-life seconds = %.6g, want 0.17", mc288.HalfLifeSeconds)
	}
	if mc288.QAlphaMeV != 10.75 {
		t.Fatalf("288Mc Q-alpha MeV = %.6g, want 10.75", mc288.QAlphaMeV)
	}
	if mc288.Daughter != "284Nh" {
		t.Fatalf("288Mc daughter = %q, want 284Nh", mc288.Daughter)
	}

	mc290, ok := byID["290Mc"]
	if !ok {
		t.Fatal("workbook missing 290Mc")
	}
	if mc290.Z != 115 || mc290.A != 290 || mc290.N != 175 {
		t.Fatalf("290Mc identity = Z=%d A=%d N=%d, want 115 290 175", mc290.Z, mc290.A, mc290.N)
	}
}

func TestValidateWorkbookRejectsProvenanceLoss(t *testing.T) {
	valid := validWorkbookRecords(t)
	if err := ValidateWorkbookRecords(valid); err != nil {
		t.Fatalf("valid workbook rejected: %v", err)
	}

	cases := []struct {
		name   string
		mutate func([]WorkbookRecord)
		want   string
	}{
		{
			name:   "missing citation URLs rejects",
			mutate: func(records []WorkbookRecord) { records[0].CitationURLs = nil },
			want:   "288Mc missing citation URLs",
		},
		{
			name:   "missing DOI rejects",
			mutate: func(records []WorkbookRecord) { records[0].DOIs = nil },
			want:   "288Mc missing DOI trail",
		},
		{
			name:   "inconsistent neutron count rejects",
			mutate: func(records []WorkbookRecord) { records[0].N = records[0].A - records[0].Z + 1 },
			want:   "288Mc has N=174, want A-Z=173",
		},
		{
			name:   "blank evidence level rejects",
			mutate: func(records []WorkbookRecord) { records[0].EvidenceLevel = " \t" },
			want:   "288Mc evidence level must not be blank",
		},
		{
			name:   "blank source path rejects",
			mutate: func(records []WorkbookRecord) { records[0].SourcePath = " \n" },
			want:   "288Mc source path must not be blank",
		},
		{
			name:   "malformed daughter ID rejects when present",
			mutate: func(records []WorkbookRecord) { records[0].Daughter = "284" },
			want:   "must include an element symbol",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			candidate := cloneWorkbookRecords(valid)
			tc.mutate(candidate)
			err := ValidateWorkbookRecords(candidate)
			if err == nil {
				t.Fatal("ValidateWorkbookRecords accepted invalid workbook")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tc.want)
			}
		})
	}
}

func validWorkbookRecords(t *testing.T) []WorkbookRecord {
	t.Helper()
	raw, err := os.ReadFile("../../data/research.seed.json")
	if err != nil {
		t.Fatalf("read seed: %v", err)
	}

	var seed ResearchSeed
	if err := json.Unmarshal(raw, &seed); err != nil {
		t.Fatalf("unmarshal seed: %v", err)
	}

	workbook, err := WorkbookFromSeed(seed, "data/research.seed.json")
	if err != nil {
		t.Fatalf("WorkbookFromSeed returned error: %v", err)
	}
	return workbook
}

func cloneWorkbookRecords(records []WorkbookRecord) []WorkbookRecord {
	clone := append([]WorkbookRecord(nil), records...)
	for i := range records {
		clone[i].CitationURLs = append([]string(nil), records[i].CitationURLs...)
		clone[i].DOIs = append([]string(nil), records[i].DOIs...)
	}
	return clone
}
