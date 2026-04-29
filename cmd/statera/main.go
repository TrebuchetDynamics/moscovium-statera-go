package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/physics"
	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/research"
)

func main() {
	options := provenanceNodeTableOptions{}
	flag.StringVar(&options.NodeType, "provenance-node-type", "", "optional provenance node table filter by node_type")
	flag.StringVar(&options.Status, "provenance-status", "", "optional provenance node table filter by status")
	flag.StringVar(&options.Format, "provenance-format", "tsv", "provenance node table output format: tsv or json")
	flag.Parse()

	if strings.TrimSpace(options.Format) == "json" {
		report, err := provenanceNodeTableReportWithOptions("data/research.seed.json", options)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(report)
		return
	}

	catalog := physics.Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, HalfLife: 170 * time.Millisecond, QAlphaMeV: 10.75, Daughter: "284Nh", CitationLink: "https://www.nndc.bnl.gov/ensnds/288/Mc/adopted.pdf"},
		"284Nh": {Symbol: "Nh", Z: 113, A: 284, Daughter: "280Rg", CitationLink: "data/research.seed.json"},
		"280Rg": {Symbol: "Rg", Z: 111, A: 280, Daughter: "276Mt", CitationLink: "data/research.seed.json"},
		"276Mt": {Symbol: "Mt", Z: 109, A: 276, Daughter: "272Bh", CitationLink: "data/research.seed.json"},
		"272Bh": {Symbol: "Bh", Z: 107, A: 272, Daughter: "268Db", CitationLink: "data/research.seed.json"},
		"268Db": {Symbol: "Db", Z: 105, A: 268, Daughter: "264Lr", CitationLink: "data/research.seed.json"},
		"264Lr": {Symbol: "Lr", Z: 103, A: 264, CitationLink: "data/research.seed.json"},
	}

	chain, err := physics.DecayChain("288Mc", catalog)
	if err != nil {
		log.Fatal(err)
	}

	for index, isotope := range chain {
		if index > 0 {
			fmt.Print(" -> ")
		}
		fmt.Print(isotope.ID())
	}
	fmt.Println()

	report, err := provenanceNodeTableReportWithOptions("data/research.seed.json", options)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(report)
}

type provenanceNodeTableOptions struct {
	NodeType string
	Status   string
	Format   string
}

func provenanceNodeTableReport(seedPath string) (string, error) {
	return provenanceNodeTableReportWithOptions(seedPath, provenanceNodeTableOptions{})
}

func provenanceNodeTableReportWithOptions(seedPath string, options provenanceNodeTableOptions) (string, error) {
	raw, err := os.ReadFile(seedPath)
	if err != nil {
		return "", err
	}
	var seed research.ResearchSeed
	if err := json.Unmarshal(raw, &seed); err != nil {
		return "", err
	}
	workbook, err := research.WorkbookFromSeed(seed, "data/research.seed.json")
	if err != nil {
		return "", err
	}
	graph, err := research.ProvenanceGraphFromWorkbook(workbook, research.DefaultBlockedSources())
	if err != nil {
		return "", err
	}

	rows := filterProvenanceNodeTableRows(graph.NodeTableRows(), options)
	if strings.TrimSpace(options.Format) == "json" {
		return provenanceNodeTableJSONReport(rows, options)
	}
	var builder strings.Builder
	fmt.Fprintf(&builder, "provenance_node_table rows=%d columns=8", len(rows))
	if filterSummary := provenanceNodeTableFilterSummary(options); filterSummary != "" {
		fmt.Fprintf(&builder, " filters=%s", filterSummary)
	}
	builder.WriteString("\n")
	builder.WriteString("node_id\tnode_type\tstatus\tsource_path\tdoi_or_url\tincoming_edges\toutgoing_edges\torphan\n")
	for _, row := range rows {
		fmt.Fprintf(
			&builder,
			"%s\t%s\t%s\t%s\t%s\t%d\t%d\t%t\n",
			row.NodeID,
			row.NodeType,
			row.Status,
			row.SourcePath,
			row.DOIOrURL,
			row.IncomingEdges,
			row.OutgoingEdges,
			row.Orphan,
		)
	}
	return builder.String(), nil
}

func filterProvenanceNodeTableRows(rows []research.ProvenanceNodeTableRow, options provenanceNodeTableOptions) []research.ProvenanceNodeTableRow {
	nodeType := strings.TrimSpace(options.NodeType)
	status := strings.TrimSpace(options.Status)
	if nodeType == "" && status == "" {
		return rows
	}
	filtered := make([]research.ProvenanceNodeTableRow, 0, len(rows))
	for _, row := range rows {
		if nodeType != "" && string(row.NodeType) != nodeType {
			continue
		}
		if status != "" && row.Status != status {
			continue
		}
		filtered = append(filtered, row)
	}
	return filtered
}

func provenanceNodeTableJSONReport(rows []research.ProvenanceNodeTableRow, options provenanceNodeTableOptions) (string, error) {
	type provenanceNodeTableJSONFilters struct {
		NodeType string `json:"node_type,omitempty"`
		Status   string `json:"status,omitempty"`
	}
	type provenanceNodeTableJSONRecord struct {
		NodeID        string `json:"node_id"`
		NodeType      string `json:"node_type"`
		Status        string `json:"status"`
		SourcePath    string `json:"source_path"`
		DOIOrURL      string `json:"doi_or_url"`
		IncomingEdges int    `json:"incoming_edges"`
		OutgoingEdges int    `json:"outgoing_edges"`
		Orphan        bool   `json:"orphan"`
	}
	type provenanceNodeTableJSONPayload struct {
		ReportType string                          `json:"report_type"`
		Rows       int                             `json:"rows"`
		Columns    int                             `json:"columns"`
		Filters    provenanceNodeTableJSONFilters  `json:"filters"`
		Records    []provenanceNodeTableJSONRecord `json:"records"`
	}
	payload := provenanceNodeTableJSONPayload{
		ReportType: "provenance_node_table",
		Rows:       len(rows),
		Columns:    8,
		Filters: provenanceNodeTableJSONFilters{
			NodeType: strings.TrimSpace(options.NodeType),
			Status:   strings.TrimSpace(options.Status),
		},
		Records: make([]provenanceNodeTableJSONRecord, 0, len(rows)),
	}
	for _, row := range rows {
		payload.Records = append(payload.Records, provenanceNodeTableJSONRecord{
			NodeID:        row.NodeID,
			NodeType:      string(row.NodeType),
			Status:        row.Status,
			SourcePath:    row.SourcePath,
			DOIOrURL:      row.DOIOrURL,
			IncomingEdges: row.IncomingEdges,
			OutgoingEdges: row.OutgoingEdges,
			Orphan:        row.Orphan,
		})
	}
	raw, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(raw) + "\n", nil
}

func provenanceNodeTableFilterSummary(options provenanceNodeTableOptions) string {
	parts := []string{}
	if nodeType := strings.TrimSpace(options.NodeType); nodeType != "" {
		parts = append(parts, "node_type="+nodeType)
	}
	if status := strings.TrimSpace(options.Status); status != "" {
		parts = append(parts, "status="+status)
	}
	return strings.Join(parts, ",")
}
