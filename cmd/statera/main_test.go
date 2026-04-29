package main

import (
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
