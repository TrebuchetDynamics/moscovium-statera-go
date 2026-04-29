package research

import (
	"fmt"
	"sort"
	"strings"
)

type GraphNodeType string

type GraphEdgeType string

const (
	NodeIsotope       GraphNodeType = "isotope"
	NodeSourcePath    GraphNodeType = "source_path"
	NodeCitationURL   GraphNodeType = "citation_url"
	NodeDOI           GraphNodeType = "doi"
	NodeBlockedSource GraphNodeType = "blocked_source"

	EdgeDerivedFromSourcePath  GraphEdgeType = "derived_from_source_path"
	EdgeSupportedByCitationURL GraphEdgeType = "supported_by_citation_url"
	EdgeSupportedByDOI         GraphEdgeType = "supported_by_doi"
	EdgeBlockedBySourceAccess  GraphEdgeType = "blocked_by_source_access"
	EdgeDocumentedAtSourcePath GraphEdgeType = "documented_at_source_path"
)

type GraphNode struct {
	ID         string
	Type       GraphNodeType
	Label      string
	Status     string
	SourcePath string
	URL        string
	DOI        string
	Note       string
}

type GraphEdge struct {
	From string
	To   string
	Type GraphEdgeType
}

type ProvenanceGraph struct {
	Nodes map[string]GraphNode
	Edges []GraphEdge
}

type BlockedSource struct {
	Key        string
	Title      string
	DOI        string
	SourcePath string
	Reason     string
}

type ProvenanceNodeTableRow struct {
	NodeID        string
	NodeType      GraphNodeType
	Status        string
	SourcePath    string
	DOIOrURL      string
	IncomingEdges int
	OutgoingEdges int
	Orphan        bool
}

func ProvenanceGraphFromWorkbook(records []WorkbookRecord, blockedSources []BlockedSource) (ProvenanceGraph, error) {
	if err := ValidateWorkbookRecords(records); err != nil {
		return ProvenanceGraph{}, err
	}
	graph := ProvenanceGraph{Nodes: map[string]GraphNode{}, Edges: []GraphEdge{}}
	for _, record := range records {
		isotopeID := "isotope:" + strings.TrimSpace(record.ID)
		graph.AddNode(GraphNode{ID: isotopeID, Type: NodeIsotope, Label: strings.TrimSpace(record.ID), Status: "accepted", SourcePath: strings.TrimSpace(record.SourcePath), Note: strings.TrimSpace(record.EvidenceLevel)})

		sourcePath := strings.TrimSpace(record.SourcePath)
		sourcePathID := "source_path:" + sourcePath
		graph.AddNode(GraphNode{ID: sourcePathID, Type: NodeSourcePath, Label: sourcePath, SourcePath: sourcePath})
		graph.AddEdge(GraphEdge{From: isotopeID, To: sourcePathID, Type: EdgeDerivedFromSourcePath})

		for _, citationURL := range record.CitationURLs {
			trimmedURL := strings.TrimSpace(citationURL)
			citationID := "citation_url:" + trimmedURL
			graph.AddNode(GraphNode{ID: citationID, Type: NodeCitationURL, Label: trimmedURL, URL: trimmedURL})
			graph.AddEdge(GraphEdge{From: isotopeID, To: citationID, Type: EdgeSupportedByCitationURL})
		}
		for _, doi := range record.DOIs {
			trimmedDOI := strings.TrimSpace(doi)
			doiID := "doi:" + trimmedDOI
			graph.AddNode(GraphNode{ID: doiID, Type: NodeDOI, Label: trimmedDOI, DOI: trimmedDOI})
			graph.AddEdge(GraphEdge{From: isotopeID, To: doiID, Type: EdgeSupportedByDOI})
		}
	}
	for _, source := range blockedSources {
		key := strings.TrimSpace(source.Key)
		if key == "" {
			return ProvenanceGraph{}, fmt.Errorf("blocked source key must not be blank")
		}
		blockedID := "blocked_source:" + key
		graph.AddNode(GraphNode{ID: blockedID, Type: NodeBlockedSource, Label: strings.TrimSpace(source.Title), Status: "blocked", SourcePath: strings.TrimSpace(source.SourcePath), DOI: strings.TrimSpace(source.DOI), Note: strings.TrimSpace(source.Reason)})
		if strings.TrimSpace(source.DOI) != "" {
			doiID := "doi:" + strings.TrimSpace(source.DOI)
			graph.AddNode(GraphNode{ID: doiID, Type: NodeDOI, Label: strings.TrimSpace(source.DOI), DOI: strings.TrimSpace(source.DOI)})
			graph.AddEdge(GraphEdge{From: blockedID, To: doiID, Type: EdgeBlockedBySourceAccess})
		}
		if strings.TrimSpace(source.SourcePath) != "" {
			sourcePathID := "source_path:" + strings.TrimSpace(source.SourcePath)
			graph.AddNode(GraphNode{ID: sourcePathID, Type: NodeSourcePath, Label: strings.TrimSpace(source.SourcePath), SourcePath: strings.TrimSpace(source.SourcePath)})
			graph.AddEdge(GraphEdge{From: blockedID, To: sourcePathID, Type: EdgeDocumentedAtSourcePath})
		}
	}
	return graph, nil
}

func (graph ProvenanceGraph) AddNode(node GraphNode) {
	if graph.Nodes == nil {
		graph.Nodes = map[string]GraphNode{}
	}
	graph.Nodes[node.ID] = node
}

func (graph *ProvenanceGraph) AddEdge(edge GraphEdge) {
	graph.Edges = append(graph.Edges, edge)
}

func (graph ProvenanceGraph) Node(id string) (GraphNode, bool) {
	node, ok := graph.Nodes[id]
	return node, ok
}

func (graph ProvenanceGraph) HasNode(id string) bool {
	_, ok := graph.Nodes[id]
	return ok
}

func (graph ProvenanceGraph) HasEdge(from, to string, edgeType GraphEdgeType) bool {
	for _, edge := range graph.Edges {
		if edge.From == from && edge.To == to && edge.Type == edgeType {
			return true
		}
	}
	return false
}

func (graph ProvenanceGraph) NodeCount() int {
	return len(graph.Nodes)
}

func (graph ProvenanceGraph) EdgeCount() int {
	return len(graph.Edges)
}

func (graph ProvenanceGraph) TextSummary() string {
	counts := map[GraphNodeType]int{}
	blockedKeys := []string{}
	for _, node := range graph.Nodes {
		counts[node.Type]++
		if node.Type == NodeBlockedSource {
			blockedKeys = append(blockedKeys, strings.TrimPrefix(node.ID, "blocked_source:"))
		}
	}
	sort.Strings(blockedKeys)
	return fmt.Sprintf(
		"provenance_graph nodes=%d edges=%d isotopes=%d doi_nodes=%d citation_url_nodes=%d blocked_source_nodes=%d orphans=%d blocked_sources=%s",
		graph.NodeCount(),
		graph.EdgeCount(),
		counts[NodeIsotope],
		counts[NodeDOI],
		counts[NodeCitationURL],
		counts[NodeBlockedSource],
		len(graph.OrphanNodes()),
		strings.Join(blockedKeys, ","),
	)
}

func (graph ProvenanceGraph) OrphanNodes() []GraphNode {
	connected := map[string]bool{}
	for _, edge := range graph.Edges {
		connected[edge.From] = true
		connected[edge.To] = true
	}
	orphans := []GraphNode{}
	for id, node := range graph.Nodes {
		if !connected[id] {
			orphans = append(orphans, node)
		}
	}
	sort.Slice(orphans, func(i, j int) bool { return orphans[i].ID < orphans[j].ID })
	return orphans
}

func (graph ProvenanceGraph) NodeTableRows() []ProvenanceNodeTableRow {
	incoming := map[string]int{}
	outgoing := map[string]int{}
	connected := map[string]bool{}
	for _, edge := range graph.Edges {
		outgoing[edge.From]++
		incoming[edge.To]++
		connected[edge.From] = true
		connected[edge.To] = true
	}

	ids := make([]string, 0, len(graph.Nodes))
	for id := range graph.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	rows := make([]ProvenanceNodeTableRow, 0, len(ids))
	for _, id := range ids {
		node := graph.Nodes[id]
		doiOrURL := strings.TrimSpace(node.DOI)
		if doiOrURL == "" {
			doiOrURL = strings.TrimSpace(node.URL)
		}
		rows = append(rows, ProvenanceNodeTableRow{
			NodeID:        id,
			NodeType:      node.Type,
			Status:        strings.TrimSpace(node.Status),
			SourcePath:    strings.TrimSpace(node.SourcePath),
			DOIOrURL:      doiOrURL,
			IncomingEdges: incoming[id],
			OutgoingEdges: outgoing[id],
			Orphan:        !connected[id],
		})
	}
	return rows
}
