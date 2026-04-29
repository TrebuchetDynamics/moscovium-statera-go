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

func TestDefaultBlockedSourcesListsCurrentSourceAccessBlockers(t *testing.T) {
	blocked := DefaultBlockedSources()
	if got, want := len(blocked), 2; got != want {
		t.Fatalf("DefaultBlockedSources count = %d, want %d", got, want)
	}

	byKey := map[string]BlockedSource{}
	for _, source := range blocked {
		byKey[source.Key] = source
		if strings.TrimSpace(source.Title) == "" {
			t.Fatalf("%s missing title", source.Key)
		}
		if strings.TrimSpace(source.Reason) == "" || !strings.Contains(source.Reason, "must not") {
			t.Fatalf("%s reason = %q, want source-access blocker text", source.Key, source.Reason)
		}
	}
	if byKey["royer2008alphaAnalytic"].DOI != "10.1103/PhysRevC.77.037602" {
		t.Fatalf("Royer blocked-source DOI = %q", byKey["royer2008alphaAnalytic"].DOI)
	}
	if byKey["wang2015alphaSystematics"].DOI != "10.1103/PhysRevC.92.064301" {
		t.Fatalf("Wang blocked-source DOI = %q", byKey["wang2015alphaSystematics"].DOI)
	}

	blocked[0].Key = "mutated"
	if DefaultBlockedSources()[0].Key == "mutated" {
		t.Fatal("DefaultBlockedSources did not return a defensive slice copy")
	}
}

func TestProvenanceGraphIncludesBlockedSourceNodes(t *testing.T) {
	blocked := DefaultBlockedSources()

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

func TestValidateProvenanceGraphIntegrityRejectsDuplicateEdgesAndMissingNodes(t *testing.T) {
	valid, err := ProvenanceGraphFromWorkbook(validWorkbookRecords(t), DefaultBlockedSources())
	if err != nil {
		t.Fatalf("ProvenanceGraphFromWorkbook returned error: %v", err)
	}
	if err := ValidateProvenanceGraphIntegrity(valid); err != nil {
		t.Fatalf("valid provenance graph rejected: %v", err)
	}

	duplicateEdge := valid
	duplicateEdge.Edges = append(append([]GraphEdge{}, valid.Edges...), valid.Edges[0])
	if err := ValidateProvenanceGraphIntegrity(duplicateEdge); err == nil || !strings.Contains(err.Error(), "duplicate graph edge") {
		t.Fatalf("duplicate edge validation error = %v, want duplicate graph edge", err)
	}

	missingFrom := valid
	missingFrom.Edges = append(append([]GraphEdge{}, valid.Edges...), GraphEdge{From: "isotope:999Xx", To: "doi:10.1103/PhysRevC.77.037602", Type: EdgeSupportedByDOI})
	if err := ValidateProvenanceGraphIntegrity(missingFrom); err == nil || !strings.Contains(err.Error(), "missing from node isotope:999Xx") {
		t.Fatalf("missing-from validation error = %v, want missing from node", err)
	}

	missingTo := valid
	missingTo.Edges = append(append([]GraphEdge{}, valid.Edges...), GraphEdge{From: "isotope:288Mc", To: "doi:10.0000/missing", Type: EdgeSupportedByDOI})
	if err := ValidateProvenanceGraphIntegrity(missingTo); err == nil || !strings.Contains(err.Error(), "missing to node doi:10.0000/missing") {
		t.Fatalf("missing-to validation error = %v, want missing to node", err)
	}
}

func TestProvenanceGraphNodeTableRowsExposeDeterministicAuditFields(t *testing.T) {
	graph, err := ProvenanceGraphFromWorkbook(validWorkbookRecords(t), []BlockedSource{
		{Key: "royer2008alphaAnalytic", Title: "Royer 2008", DOI: "10.1103/PhysRevC.77.037602", SourcePath: "citations/papers/royer2008alpha-analytic.md", Reason: "APS article/PDF content unavailable; coefficients must not be changed from metadata alone"},
	})
	if err != nil {
		t.Fatalf("ProvenanceGraphFromWorkbook returned error: %v", err)
	}
	graph.AddNode(GraphNode{ID: "doi:10.0000/orphan", Type: NodeDOI, Label: "10.0000/orphan", DOI: "10.0000/orphan"})

	rows := graph.NodeTableRows()
	if got, want := len(rows), graph.NodeCount(); got != want {
		t.Fatalf("row count = %d, want node count %d", got, want)
	}
	for i := 1; i < len(rows); i++ {
		if rows[i-1].NodeID > rows[i].NodeID {
			t.Fatalf("rows not sorted by NodeID at %d: %q > %q", i, rows[i-1].NodeID, rows[i].NodeID)
		}
	}

	byID := map[string]ProvenanceNodeTableRow{}
	for _, row := range rows {
		byID[row.NodeID] = row
	}

	isotope := byID["isotope:288Mc"]
	if isotope.NodeType != NodeIsotope || isotope.Status != "accepted" || isotope.SourcePath != "data/research.seed.json" {
		t.Fatalf("isotope:288Mc row missing identity/status/source fields: %+v", isotope)
	}
	if isotope.OutgoingEdges != 5 || isotope.IncomingEdges != 0 || isotope.Orphan {
		t.Fatalf("isotope:288Mc edge/orphan fields = incoming %d outgoing %d orphan %v, want 0 5 false", isotope.IncomingEdges, isotope.OutgoingEdges, isotope.Orphan)
	}

	doi := byID["doi:10.1103/PhysRevC.77.037602"]
	if doi.NodeType != NodeDOI || doi.DOIOrURL != "10.1103/PhysRevC.77.037602" || doi.IncomingEdges != 1 || doi.OutgoingEdges != 0 || doi.Orphan {
		t.Fatalf("doi row fields not preserved: %+v", doi)
	}

	blocked := byID["blocked_source:royer2008alphaAnalytic"]
	if blocked.NodeType != NodeBlockedSource || blocked.Status != "blocked" || blocked.SourcePath != "citations/papers/royer2008alpha-analytic.md" || blocked.DOIOrURL != "10.1103/PhysRevC.77.037602" {
		t.Fatalf("blocked-source row fields not preserved: %+v", blocked)
	}
	if blocked.OutgoingEdges != 2 || blocked.IncomingEdges != 0 || blocked.Orphan {
		t.Fatalf("blocked-source edge/orphan fields = incoming %d outgoing %d orphan %v, want 0 2 false", blocked.IncomingEdges, blocked.OutgoingEdges, blocked.Orphan)
	}

	orphan := byID["doi:10.0000/orphan"]
	if !orphan.Orphan || orphan.IncomingEdges != 0 || orphan.OutgoingEdges != 0 {
		t.Fatalf("orphan row fields = incoming %d outgoing %d orphan %v, want 0 0 true", orphan.IncomingEdges, orphan.OutgoingEdges, orphan.Orphan)
	}
}
