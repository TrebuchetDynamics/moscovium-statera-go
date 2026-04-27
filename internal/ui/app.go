package ui

import (
	"math"
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
	Spec             Spec
	Views            []ViewSpec
	Summary          Summary
	Isotopes         []IsotopeRecord
	ClaimExamples    []ClaimExample
	SourceRecords    []SourceRecord
	EducationLessons []LessonRecord
	ResearchItems    []ResearchItem
	DesignScenarios  []DesignScenario
	ContextRecords   []ContextRecord
	AlphaSystematics []AlphaSystematicsRecord
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

type LessonRecord struct {
	Title      string
	Objective  string
	Concept    string
	SourcePath string
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

type ResearchItem struct {
	Key        string
	Track      string
	Status     string
	Identifier string
	PDF        string
	Relevance  string
	SourcePath string
}

type DesignScenario struct {
	Name          string
	Goal          string
	Inputs        []string
	Result        physics.ClaimResult
	SimulationUse string
	Constraint    string
	SourcePath    string
}

type ContextRecord struct {
	Title         string
	EventDate     string
	Labels        []string
	SourceCount   int
	SimulationUse string
	SourcePath    string
}

type AlphaSystematicsRecord struct {
	IsotopeID         string
	Z                 int
	A                 int
	ParityClass       string
	QAlphaMeV         float64
	EvaluatedHalfLife time.Duration
	PredictedHalfLife time.Duration
	LogResidual       float64
	ModelName         string
	ModelReference    string
	EvidenceClass     string
	SourcePath        string
}

func DefaultSpec() Spec {
	return Spec{
		Title:  "Moscovium Statera Go",
		Width:  1180,
		Height: 760,
		Modules: []ModuleSpec{
			{
				Name:        "Education",
				Description: "Guided lessons for evaluated isotope data, decay traversal, and validator outcomes.",
			},
			{
				Name:        "Research",
				Description: "Citation and source records with provenance, status, and Track A/Track B labels.",
			},
			{
				Name:        "Design",
				Description: "Constrained scenario sandbox backed by the existing physics validator.",
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
	researchItems := defaultResearchItems()
	contextRecords := defaultContextRecords()
	claimExamples := defaultClaimExamples(catalog)
	designScenarios := defaultDesignScenarios(catalog)
	alphaSystematics := defaultAlphaSystematicsRecords(catalog)

	return AppModel{
		Spec:  spec,
		Views: viewSpecs(spec.Modules),
		Summary: Summary{
			VerifiedIsotopes: len(isotopes),
			CitationRecords:  len(researchItems),
			ContextRecords:   len(contextRecords),
			DecayChain:       decayChain,
			BoundaryNotice:   "Track B context is excluded from simulation and cannot seed physics defaults.",
		},
		Isotopes:         isotopes,
		ClaimExamples:    claimExamples,
		SourceRecords:    sourceRecords,
		EducationLessons: defaultEducationLessons(),
		ResearchItems:    researchItems,
		DesignScenarios:  designScenarios,
		ContextRecords:   contextRecords,
		AlphaSystematics: alphaSystematics,
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

func defaultEducationLessons() []LessonRecord {
	return []LessonRecord{
		{
			Title:      "Evaluated Isotope Records",
			Objective:  "Read a Track A isotope record as a small set of auditable fields: Z, A, half-life, Q alpha, daughter, and provenance.",
			Concept:    "Simulation input starts from evaluated nuclear data, not from public narratives or unsourced claims.",
			SourcePath: "data/research.seed.json",
		},
		{
			Title:      "Decay-Chain Traversal",
			Objective:  "Follow the current 288Mc smoke path through daughter records until the catalog chain terminates.",
			Concept:    "Traversal is deterministic: each daughter ID must exist in the catalog or the chain stops with an error.",
			SourcePath: "internal/physics/decay.go",
		},
		{
			Title:      "Half-Life Checks",
			Objective:  "Compare a minimum half-life claim with the catalog half-life for the named isotope.",
			Concept:    "A short-lived isotope can support a millisecond-scale minimum while failing hour-scale stability claims.",
			SourcePath: "internal/physics/claims.go",
		},
		{
			Title:      "Supported Model Boundary",
			Objective:  "Identify mechanism claims that have no standard-model path in the current validator.",
			Concept:    "The validator can reject unsupported mechanisms without treating context records as physics evidence.",
			SourcePath: "docs/research-charter.md",
		},
		{
			Title:      "Peer-Reviewed-Model Boundary",
			Objective:  "Read alpha half-life predictions as model output, not as evaluated data.",
			Concept:    "Royer-style analytic formulas live in the peer-reviewed-model evidence class. Their output never substitutes for evaluated half-lives and must always carry the model name and DOI.",
			SourcePath: "internal/physics/alpha.go",
		},
	}
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

func defaultResearchItems() []ResearchItem {
	items := defaultSourceRecords()
	records := make([]ResearchItem, 0, len(items))
	for _, item := range items {
		records = append(records, ResearchItem{
			Key:        item.Key,
			Track:      item.Track,
			Status:     item.Status,
			Identifier: item.Identifier,
			PDF:        item.PDF,
			Relevance:  item.Relevance,
			SourcePath: item.SourcePath,
		})
	}
	return records
}

func defaultDesignScenarios(catalog physics.Catalog) []DesignScenario {
	scenarios := []struct {
		name       string
		goal       string
		inputs     []string
		claim      physics.Claim
		constraint string
		sourcePath string
	}{
		{
			name:   "Millisecond Stability Check",
			goal:   "Confirm that the current catalog can support a 288Mc half-life threshold at or below the evaluated millisecond scale.",
			inputs: []string{"288Mc", "minimum half-life 100ms", "Track A catalog"},
			claim: physics.Claim{
				Kind:            physics.ClaimKindMinimumHalfLife,
				IsotopeID:       "288Mc",
				MinimumHalfLife: 100 * time.Millisecond,
			},
			constraint: "Allowed only because the threshold is checked against the Track A half-life record.",
			sourcePath: "data/research.seed.json",
		},
		{
			name:   "Room-Scale Stability Rejection",
			goal:   "Show that hour-scale stability claims are incongruent with the current 288Mc record.",
			inputs: []string{"288Mc", "minimum half-life 1h", "Track A catalog"},
			claim: physics.Claim{
				Kind:            physics.ClaimKindMinimumHalfLife,
				IsotopeID:       "288Mc",
				MinimumHalfLife: time.Hour,
			},
			constraint: "Blocked because the requested half-life exceeds the evaluated catalog value used by the app.",
			sourcePath: "data/research.seed.json",
		},
		{
			name:   "Unsupported Mechanism Rejection",
			goal:   "Show that a non-standard mechanism claim cannot become a simulation parameter.",
			inputs: []string{"mechanism", "antigravity propulsion", "supported-model boundary"},
			claim: physics.Claim{
				Kind:      physics.ClaimKindMechanism,
				Mechanism: "antigravity propulsion",
			},
			constraint: "Blocked because the current validator has no supported standard-model mechanism for this claim type.",
			sourcePath: "docs/research-charter.md",
		},
		{
			name:   "Mixed Claim Rejection",
			goal:   "Show that a half-life check cannot be combined with mechanism text in one simulation input.",
			inputs: []string{"288Mc", "minimum half-life 100ms", "mechanism text"},
			claim: physics.Claim{
				Kind:            physics.ClaimKindMinimumHalfLife,
				IsotopeID:       "288Mc",
				MinimumHalfLife: 100 * time.Millisecond,
				Mechanism:       "antigravity propulsion",
			},
			constraint: "Blocked because the validator requires a single coherent claim shape.",
			sourcePath: "internal/physics/claims.go",
		},
	}

	records := make([]DesignScenario, 0, len(scenarios))
	for _, scenario := range scenarios {
		result := physics.ValidateClaim(scenario.claim, catalog)
		simulationUse := "blocked"
		if result.Status == physics.ClaimStatusSupportedByTrackA {
			simulationUse = "allowed"
		}
		records = append(records, DesignScenario{
			Name:          scenario.name,
			Goal:          scenario.goal,
			Inputs:        scenario.inputs,
			Result:        result,
			SimulationUse: simulationUse,
			Constraint:    scenario.constraint,
			SourcePath:    scenario.sourcePath,
		})
	}
	return records
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

func defaultAlphaSystematicsRecords(catalog physics.Catalog) []AlphaSystematicsRecord {
	model := physics.RoyerModel()
	ids := []string{"288Mc", "290Mc"}
	records := make([]AlphaSystematicsRecord, 0, len(ids))
	for _, id := range ids {
		isotope, ok := catalog[id]
		if !ok || isotope.QAlphaMeV <= 0 {
			continue
		}
		prediction, err := model.Predict(isotope.Z, isotope.A, isotope.QAlphaMeV)
		if err != nil {
			continue
		}
		residual := math.NaN()
		if isotope.HalfLife > 0 {
			evaluatedSeconds := isotope.HalfLife.Seconds()
			if evaluatedSeconds > 0 {
				residual = prediction.LogHalfLifeSeconds - math.Log10(evaluatedSeconds)
			}
		}
		records = append(records, AlphaSystematicsRecord{
			IsotopeID:         id,
			Z:                 isotope.Z,
			A:                 isotope.A,
			ParityClass:       string(prediction.ParityClass),
			QAlphaMeV:         isotope.QAlphaMeV,
			EvaluatedHalfLife: isotope.HalfLife,
			PredictedHalfLife: prediction.HalfLife,
			LogResidual:       residual,
			ModelName:         model.Name,
			ModelReference:    model.Reference,
			EvidenceClass:     string(prediction.EvidenceClass),
			SourcePath:        "internal/physics/alpha.go",
		})
	}
	return records
}
