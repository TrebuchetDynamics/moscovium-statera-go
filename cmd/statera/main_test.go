package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestProvenanceNodeTableReportIsDeterministicAndSourceBacked(t *testing.T) {
	report, err := provenanceNodeTableReport("../../data/research.seed.json")
	if err != nil {
		t.Fatalf("provenanceNodeTableReport returned error: %v", err)
	}

	for _, want := range []string{
		"provenance_node_table rows=17 columns=8",
		"node_id\tnode_type\tstatus\tsource_path\tdoi_or_url\tincoming_edges\toutgoing_edges\torphan",
		"blocked_source:royer2008alphaAnalytic\tblocked_source\tblocked\tcitations/papers/royer2008alpha-analytic.md\t10.1103/PhysRevC.77.037602\t0\t2\tfalse",
		"doi:10.1103/PhysRevC.77.037602\tdoi\t\t\t10.1103/PhysRevC.77.037602\t1\t0\tfalse",
		"isotope:288Mc\tisotope\taccepted\tdata/research.seed.json\t\t0\t5\tfalse",
		"source_path:data/research.seed.json\tsource_path\t\tdata/research.seed.json\t\t2\t0\tfalse",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q\nfull report:\n%s", want, report)
		}
	}

	lines := strings.Split(strings.TrimSpace(report), "\n")
	if got, want := len(lines), 19; got != want {
		t.Fatalf("line count = %d, want %d", got, want)
	}
	for i := 3; i < len(lines); i++ {
		previousID := strings.Split(lines[i-1], "\t")[0]
		currentID := strings.Split(lines[i], "\t")[0]
		if previousID > currentID {
			t.Fatalf("rows are not sorted by node ID at line %d: %q > %q", i+1, previousID, currentID)
		}
	}
}

func TestProvenanceNodeTableReportFiltersByTypeAndStatus(t *testing.T) {
	report, err := provenanceNodeTableReportWithOptions("../../data/research.seed.json", provenanceNodeTableOptions{NodeType: "blocked_source", Status: "blocked"})
	if err != nil {
		t.Fatalf("provenanceNodeTableReportWithOptions returned error: %v", err)
	}

	for _, want := range []string{
		"provenance_node_table rows=2 columns=8 filters=node_type=blocked_source,status=blocked",
		"blocked_source:royer2008alphaAnalytic\tblocked_source\tblocked\tcitations/papers/royer2008alpha-analytic.md\t10.1103/PhysRevC.77.037602\t0\t2\tfalse",
		"blocked_source:wang2015alphaSystematics\tblocked_source\tblocked\tcitations/papers/wang2015alpha-systematics.md\t10.1103/PhysRevC.92.064301\t0\t2\tfalse",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("filtered report missing %q\nfull report:\n%s", want, report)
		}
	}
	for _, unwanted := range []string{"isotope:288Mc", "doi:10.1103/PhysRevC.77.037602"} {
		if strings.Contains(report, unwanted) {
			t.Fatalf("filtered report unexpectedly contains %q\nfull report:\n%s", unwanted, report)
		}
	}

	lines := strings.Split(strings.TrimSpace(report), "\n")
	if got, want := len(lines), 4; got != want {
		t.Fatalf("filtered line count = %d, want %d", got, want)
	}
}

func TestDecaySimulationSeedJSONReportUsesOnlyAcceptedSeedRecords(t *testing.T) {
	report, err := decaySimulationSeedJSONReport("../../data/research.seed.json", 64, 20260429)
	if err != nil {
		t.Fatalf("decaySimulationSeedJSONReport returned error: %v", err)
	}

	var payload struct {
		ReportType  string `json:"report_type"`
		SourcePath  string `json:"source_path"`
		SampleCount int    `json:"sample_count"`
		Seed        int64  `json:"seed"`
		Records     []struct {
			ID                      string  `json:"id"`
			SourcePath              string  `json:"source_path"`
			HalfLifeSeconds         float64 `json:"half_life_seconds"`
			DecayConstantPerSecond  float64 `json:"decay_constant_per_second"`
			MeanLifeSeconds         float64 `json:"mean_life_seconds"`
			MonteCarloP05Seconds    float64 `json:"monte_carlo_p05_seconds"`
			MonteCarloMedianSeconds float64 `json:"monte_carlo_median_seconds"`
			MonteCarloP95Seconds    float64 `json:"monte_carlo_p95_seconds"`
		} `json:"records"`
	}
	if err := json.Unmarshal([]byte(report), &payload); err != nil {
		t.Fatalf("JSON report did not unmarshal: %v\nreport:\n%s", err, report)
	}
	if payload.ReportType != "decay_simulation_seed_summary" || payload.SourcePath != "data/research.seed.json" || payload.SampleCount != 64 || payload.Seed != 20260429 {
		t.Fatalf("metadata = %+v, want seed-backed decay simulation summary", payload)
	}
	if got, want := len(payload.Records), 2; got != want {
		t.Fatalf("record count = %d, want %d", got, want)
	}
	if payload.Records[0].ID != "288Mc" || payload.Records[1].ID != "290Mc" {
		t.Fatalf("record IDs = %q, %q; want sorted accepted seed IDs 288Mc, 290Mc", payload.Records[0].ID, payload.Records[1].ID)
	}
	for _, record := range payload.Records {
		if record.SourcePath == "" || record.HalfLifeSeconds <= 0 || record.DecayConstantPerSecond <= 0 || record.MeanLifeSeconds <= 0 {
			t.Fatalf("record has missing source or non-positive deterministic metrics: %+v", record)
		}
		if record.MonteCarloP05Seconds <= 0 || record.MonteCarloP05Seconds > record.MonteCarloMedianSeconds || record.MonteCarloMedianSeconds > record.MonteCarloP95Seconds {
			t.Fatalf("record Monte Carlo quantiles are invalid: %+v", record)
		}
	}
	if strings.Contains(report, "284Nh") || strings.Contains(report, "286Nh") {
		t.Fatalf("simulation report must not infer daughter isotope records without source-backed daughter half-lives:\n%s", report)
	}
}

func TestProvenanceNodeTableJSONReportIsDeterministicAndFilterAware(t *testing.T) {
	report, err := provenanceNodeTableReportWithOptions("../../data/research.seed.json", provenanceNodeTableOptions{NodeType: "blocked_source", Status: "blocked", Format: "json"})
	if err != nil {
		t.Fatalf("provenanceNodeTableReportWithOptions returned error: %v", err)
	}
	if strings.Contains(report, "node_id\tnode_type") {
		t.Fatalf("JSON report unexpectedly contains tab-separated header: %s", report)
	}

	var payload struct {
		ReportType string `json:"report_type"`
		Rows       int    `json:"rows"`
		Columns    int    `json:"columns"`
		Filters    struct {
			NodeType string `json:"node_type,omitempty"`
			Status   string `json:"status,omitempty"`
		} `json:"filters"`
		Records []struct {
			NodeID        string `json:"node_id"`
			NodeType      string `json:"node_type"`
			Status        string `json:"status"`
			SourcePath    string `json:"source_path"`
			DOIOrURL      string `json:"doi_or_url"`
			IncomingEdges int    `json:"incoming_edges"`
			OutgoingEdges int    `json:"outgoing_edges"`
			Orphan        bool   `json:"orphan"`
		} `json:"records"`
	}
	if err := json.Unmarshal([]byte(report), &payload); err != nil {
		t.Fatalf("JSON report did not unmarshal: %v\nreport:\n%s", err, report)
	}
	if payload.ReportType != "provenance_node_table" || payload.Rows != 2 || payload.Columns != 8 {
		t.Fatalf("summary = type %q rows %d columns %d, want provenance_node_table rows 2 columns 8", payload.ReportType, payload.Rows, payload.Columns)
	}
	if payload.Filters.NodeType != "blocked_source" || payload.Filters.Status != "blocked" {
		t.Fatalf("filters = node_type %q status %q, want blocked_source blocked", payload.Filters.NodeType, payload.Filters.Status)
	}
	if got, want := len(payload.Records), 2; got != want {
		t.Fatalf("record count = %d, want %d", got, want)
	}
	wantIDs := []string{"blocked_source:royer2008alphaAnalytic", "blocked_source:wang2015alphaSystematics"}
	for i, wantID := range wantIDs {
		record := payload.Records[i]
		if record.NodeID != wantID {
			t.Fatalf("record[%d].node_id = %q, want %q", i, record.NodeID, wantID)
		}
		if record.NodeType != "blocked_source" || record.Status != "blocked" || record.IncomingEdges != 0 || record.OutgoingEdges != 2 || record.Orphan {
			t.Fatalf("record[%d] fields = %+v, want blocked source with 0 incoming, 2 outgoing, orphan=false", i, record)
		}
	}
}
