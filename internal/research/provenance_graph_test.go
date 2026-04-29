package research

import (
	"strings"
	"testing"
)

func TestProvenanceGraphFromWorkbookBuildsEvidenceEdges(t *testing.T) {
	records := validWorkbookRecords(t)
	graph, err := ProvenanceGraphFromWorkbook(records, nil)
	if err != nil {
		t.Fatalf("ProvenanceGraphFromWorkbook returned error: %v", err)
	}

	if got, want := graph.NodeCount(), 11; got != want {
		t.Fatalf("node count = %d, want %d", got, want)
	}
	if got, want := graph.EdgeCount(), 10; got != want {
		t.Fatalf("edge count = %d, want %d", got, want)
	}

	for _, record := range records {
		isotopeNodeID := "isotope:" + record.ID
		if !graph.HasNode(isotopeNodeID) {
			t.Fatalf("missing isotope node %s", isotopeNodeID)
		}
		if !graph.HasEdge(isotopeNodeID, "source_path:data/research.seed.json", EdgeDerivedFromSourcePath) {
			t.Fatalf("missing source-path edge for %s", record.ID)
		}
		for _, citationURL := range record.CitationURLs {
			if !graph.HasEdge(isotopeNodeID, "citation_url:"+citationURL, EdgeSupportedByCitationURL) {
				t.Fatalf("missing citation URL edge for %s -> %s", record.ID, citationURL)
			}
		}
		for _, doi := range record.DOIs {
			if !graph.HasEdge(isotopeNodeID, "doi:"+doi, EdgeSupportedByDOI) {
				t.Fatalf("missing DOI edge for %s -> %s", record.ID, doi)
			}
		}
	}
}

func TestProvenanceGraphIncludesBlockedSourceNodes(t *testing.T) {
	blocked := []BlockedSource{
		{
			Key:        "royer2008alphaAnalytic",
			Title:      "Recent alpha decay half-lives and analytic expression predictions including superheavy nuclei",
			DOI:        "10.1103/PhysRevC.77.037602",
			SourcePath: "citations/papers/royer2008alpha-analytic.md",
			Reason:     "APS article/PDF content unavailable in prior source triage; coefficients must not be changed from metadata alone",
		},
		{
			Key:        "wang2015alphaSystematics",
			Title:      "Systematic study of alpha-decay energies and half-lives of superheavy nuclei",
			DOI:        "10.1103/PhysRevC.92.064301",
			SourcePath: "citations/papers/wang2015alpha-systematics.md",
			Reason:     "APS article/PDF content unavailable in prior source triage; second-model formulas must not be added from metadata alone",
		},
	}

	graph, err := ProvenanceGraphFromWorkbook(validWorkbookRecords(t), blocked)
	if err != nil {
		t.Fatalf("ProvenanceGraphFromWorkbook returned error: %v", err)
	}
	if got, want := graph.NodeCount(), 17; got != want {
		t.Fatalf("node count = %d, want %d", got, want)
	}
	if got, want := graph.EdgeCount(), 14; got != want {
		t.Fatalf("edge count = %d, want %d", got, want)
	}

	for _, source := range blocked {
		blockedNodeID := "blocked_source:" + source.Key
		node, ok := graph.Node(blockedNodeID)
		if !ok {
			t.Fatalf("missing blocked source node %s", blockedNodeID)
		}
		if node.Type != NodeBlockedSource {
			t.Fatalf("%s type = %q, want %q", blockedNodeID, node.Type, NodeBlockedSource)
		}
		if node.Status != "blocked" {
			t.Fatalf("%s status = %q, want blocked", blockedNodeID, node.Status)
		}
		if !strings.Contains(node.Note, "must not") {
			t.Fatalf("%s note = %q, want blocking reason", blockedNodeID, node.Note)
		}
		if !graph.HasEdge(blockedNodeID, "doi:"+source.DOI, EdgeBlockedBySourceAccess) {
			t.Fatalf("missing blocked DOI edge for %s -> %s", source.Key, source.DOI)
		}
		if !graph.HasEdge(blockedNodeID, "source_path:"+source.SourcePath, EdgeDocumentedAtSourcePath) {
			t.Fatalf("missing blocked source-path edge for %s -> %s", source.Key, source.SourcePath)
		}
	}
}

func TestProvenanceGraphTextSummaryReportsCountsAndBlockedSources(t *testing.T) {
	graph, err := ProvenanceGraphFromWorkbook(validWorkbookRecords(t), []BlockedSource{
		{Key: "royer2008alphaAnalytic", Title: "Royer 2008", DOI: "10.1103/PhysRevC.77.037602", SourcePath: "citations/papers/royer2008alpha-analytic.md", Reason: "APS article/PDF content unavailable; coefficients must not be changed from metadata alone"},
		{Key: "wang2015alphaSystematics", Title: "Wang 2015", DOI: "10.1103/PhysRevC.92.064301", SourcePath: "citations/papers/wang2015alpha-systematics.md", Reason: "APS article/PDF content unavailable; second-model formulas must not be added from metadata alone"},
	})
	if err != nil {
		t.Fatalf("ProvenanceGraphFromWorkbook returned error: %v", err)
	}

	summary := graph.TextSummary()
	for _, want := range []string{
		"nodes=17",
		"edges=14",
		"isotopes=2",
		"doi_nodes=6",
		"citation_url_nodes=4",
		"blocked_source_nodes=2",
		"orphans=0",
		"blocked_sources=royer2008alphaAnalytic,wang2015alphaSystematics",
	} {
		if !strings.Contains(summary, want) {
			t.Fatalf("TextSummary() = %q, want substring %q", summary, want)
		}
	}
}

func TestProvenanceGraphReportsOrphans(t *testing.T) {
	graph, err := ProvenanceGraphFromWorkbook(validWorkbookRecords(t), nil)
	if err != nil {
		t.Fatalf("ProvenanceGraphFromWorkbook returned error: %v", err)
	}
	graph.AddNode(GraphNode{ID: "doi:10.0000/orphan", Type: NodeDOI, Label: "10.0000/orphan"})

	orphans := graph.OrphanNodes()
	if got, want := len(orphans), 1; got != want {
		t.Fatalf("orphan count = %d, want %d: %#v", got, want, orphans)
	}
	if orphans[0].ID != "doi:10.0000/orphan" {
		t.Fatalf("orphan ID = %q, want doi:10.0000/orphan", orphans[0].ID)
	}
}
