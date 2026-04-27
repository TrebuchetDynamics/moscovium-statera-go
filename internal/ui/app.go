package ui

import (
	"time"

	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/physics"
)

type ModuleSpec struct {
	Name        string
	Description string
}

type Spec struct {
	Title   string
	Width   int
	Height  int
	Modules []ModuleSpec
}

type AppModel struct {
	Spec           Spec
	Views          []ViewSpec
	Summary        Summary
	Isotopes       []IsotopeRecord
	ClaimExamples  []ClaimExample
	SourceRecords  []SourceRecord
	ContextRecords []ContextRecord
}

type ViewSpec struct {
	Name        string
	Description string
}

type Summary struct {
	VerifiedIsotopes int
	CitationRecords  int
	ContextRecords   int
	DecayChain       []string
	BoundaryNotice   string
}

type IsotopeRecord struct {
	ID            string
	Z             int
	A             int
	HalfLife      time.Duration
	Daughter      string
	CitationCount int
	SourcePath    string
}

type ClaimExample struct {
	Label  string
	Claim  physics.Claim
	Result physics.ClaimResult
}

type SourceRecord struct {
	Key        string
	Track      string
	Status     string
	Identifier string
	PDF        string
	Relevance  string
	SourcePath string
}

type ContextRecord struct {
	Title         string
	EventDate     string
	Labels        []string
	SourceCount   int
	SimulationUse string
	SourcePath    string
}

func DefaultSpec() Spec {
	return Spec{
		Title:  "Moscovium Statera Go",
		Width:  1180,
		Height: 760,
		Modules: []ModuleSpec{
			{
				Name:        "Dashboard",
				Description: "Track A/Track B status and current decay-chain smoke path.",
			},
			{
				Name:        "Physics",
				Description: "Evaluated isotope records and deterministic claim-validator examples.",
			},
			{
				Name:        "Sources",
				Description: "Citation records with read status and DOI or URL provenance.",
			},
			{
				Name:        "Context",
				Description: "Quarantined Track B records marked not-for-simulation.",
			},
		},
	}
}

func DefaultModel() AppModel {
	spec := DefaultSpec()
	catalog := defaultCatalog()
	decayChain := decayChainIDs("288Mc", catalog)
	isotopes := isotopeRecords(catalog)
	sourceRecords := defaultSourceRecords()
	contextRecords := defaultContextRecords()
	claimExamples := defaultClaimExamples(catalog)

	return AppModel{
		Spec:  spec,
		Views: viewSpecs(spec.Modules),
		Summary: Summary{
			VerifiedIsotopes: len(isotopes),
			CitationRecords:  len(sourceRecords),
			ContextRecords:   len(contextRecords),
			DecayChain:       decayChain,
			BoundaryNotice:   "Track B context is excluded from simulation and cannot seed physics defaults.",
		},
		Isotopes:       isotopes,
		ClaimExamples:  claimExamples,
		SourceRecords:  sourceRecords,
		ContextRecords: contextRecords,
	}
}

func defaultCatalog() physics.Catalog {
	return physics.Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, HalfLife: 170 * time.Millisecond, QAlphaMeV: 10.75, Daughter: "284Nh", CitationLink: "data/research.seed.json"},
		"290Mc": {Symbol: "Mc", Z: 115, A: 290, HalfLife: 650 * time.Millisecond, QAlphaMeV: 10.45, Daughter: "286Nh", CitationLink: "data/research.seed.json"},
		"284Nh": {Symbol: "Nh", Z: 113, A: 284, Daughter: "280Rg", CitationLink: "data/research.seed.json"},
		"280Rg": {Symbol: "Rg", Z: 111, A: 280, Daughter: "276Mt", CitationLink: "data/research.seed.json"},
		"276Mt": {Symbol: "Mt", Z: 109, A: 276, Daughter: "272Bh", CitationLink: "data/research.seed.json"},
		"272Bh": {Symbol: "Bh", Z: 107, A: 272, Daughter: "268Db", CitationLink: "data/research.seed.json"},
		"268Db": {Symbol: "Db", Z: 105, A: 268, Daughter: "264Lr", CitationLink: "data/research.seed.json"},
		"264Lr": {Symbol: "Lr", Z: 103, A: 264, CitationLink: "data/research.seed.json"},
	}
}

func viewSpecs(modules []ModuleSpec) []ViewSpec {
	views := make([]ViewSpec, 0, len(modules))
	for _, module := range modules {
		views = append(views, ViewSpec{Name: module.Name, Description: module.Description})
	}
	return views
}

func decayChainIDs(start string, catalog physics.Catalog) []string {
	chain, err := physics.DecayChain(start, catalog)
	if err != nil {
		return []string{err.Error()}
	}
	ids := make([]string, 0, len(chain))
	for _, isotope := range chain {
		ids = append(ids, isotope.ID())
	}
	return ids
}

func isotopeRecords(catalog physics.Catalog) []IsotopeRecord {
	ids := []string{"288Mc", "290Mc"}
	records := make([]IsotopeRecord, 0, len(ids))
	for _, id := range ids {
		isotope := catalog[id]
		records = append(records, IsotopeRecord{
			ID:            id,
			Z:             isotope.Z,
			A:             isotope.A,
			HalfLife:      isotope.HalfLife,
			Daughter:      isotope.Daughter,
			CitationCount: 2,
			SourcePath:    "data/research.seed.json",
		})
	}
	return records
}

func defaultClaimExamples(catalog physics.Catalog) []ClaimExample {
	examples := []ClaimExample{
		{
			Label: "288Mc minimum half-life >= 100ms",
			Claim: physics.Claim{
				Kind:            physics.ClaimKindMinimumHalfLife,
				IsotopeID:       "288Mc",
				MinimumHalfLife: 100 * time.Millisecond,
			},
		},
		{
			Label: "288Mc minimum half-life >= 1h",
			Claim: physics.Claim{
				Kind:            physics.ClaimKindMinimumHalfLife,
				IsotopeID:       "288Mc",
				MinimumHalfLife: time.Hour,
			},
		},
		{
			Label: "antigravity propulsion",
			Claim: physics.Claim{
				Kind:      physics.ClaimKindMechanism,
				Mechanism: "antigravity propulsion",
			},
		},
		{
			Label: "mixed half-life plus mechanism",
			Claim: physics.Claim{
				Kind:            physics.ClaimKindMinimumHalfLife,
				IsotopeID:       "288Mc",
				MinimumHalfLife: 100 * time.Millisecond,
				Mechanism:       "antigravity propulsion",
			},
		},
	}

	for i := range examples {
		examples[i].Result = physics.ValidateClaim(examples[i].Claim, catalog)
	}
	return examples
}

func defaultSourceRecords() []SourceRecord {
	return []SourceRecord{
		{
			Key:        "oganessian2022mcfactory",
			Track:      "Track A",
			Status:     "to-read",
			Identifier: "10.1103/PhysRevC.106.L031301",
			PDF:        "not stored",
			Relevance:  "Metadata record for a Physical Review C article on the 243Am+48Ca reaction and Moscovium-related decay-chain source material.",
			SourcePath: "citations/papers/oganessian2022mcfactory.md",
		},
		{
			Key:        "iupac2016names",
			Track:      "Track A",
			Status:     "to-read",
			Identifier: "10.1515/pac-2016-0501",
			PDF:        "not stored",
			Relevance:  "Metadata record for the IUPAC Recommendations 2016 article on names and symbols for elements 113, 115, 117, and 118.",
			SourcePath: "citations/papers/iupac2016names.md",
		},
		{
			Key:        "houseoversight2026missingScientistsLetter",
			Track:      "Track B",
			Status:     "to-read",
			Identifier: "https://oversight.house.gov/wp-content/uploads/2026/04/FBI-Missing-Scientists-Letter_4.20.26.pdf",
			PDF:        "not stored",
			Relevance:  "Context record for congressional correspondence; not validation of Element 115, UAP, or propulsion linkages.",
			SourcePath: "citations/papers/houseoversight2026missing-scientists-letter.md",
		},
	}
}

func defaultContextRecords() []ContextRecord {
	return []ContextRecord{
		{
			Title:         "William Neil McCasland Missing-Person Context Record",
			EventDate:     "2026-02-27",
			Labels:        []string{"verified-event", "media-claim", "unsupported-linkage", "not-for-simulation"},
			SourceCount:   2,
			SimulationUse: "prohibited",
			SourcePath:    "docs/lore/records/mccasland-2026-missing-person.md",
		},
		{
			Title:         "Carl Grillmair Homicide Context Record",
			EventDate:     "2026-02-16",
			Labels:        []string{"verified-event", "media-claim", "unsupported-linkage", "not-for-simulation"},
			SourceCount:   2,
			SimulationUse: "prohibited",
			SourcePath:    "docs/lore/records/grillmair-2026-homicide.md",
		},
	}
}
