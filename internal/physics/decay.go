package physics

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type Isotope struct {
	Symbol       string
	Z            int
	A            int
	HalfLife     time.Duration
	QAlphaMeV    float64
	Daughter     string
	CitationLink string
}

type Catalog map[string]Isotope

func (i Isotope) ID() string {
	return fmt.Sprintf("%d%s", i.A, i.Symbol)
}

func (i Isotope) Validate() error {
	if i.Symbol == "" {
		return errors.New("isotope symbol is required")
	}
	if i.Z <= 0 {
		return fmt.Errorf("%s has invalid Z=%d", i.ID(), i.Z)
	}
	expectedZ, ok := AtomicNumberForSymbol(i.Symbol)
	if !ok {
		return fmt.Errorf("%s has unknown element symbol %s", i.ID(), i.Symbol)
	}
	if i.Z != expectedZ {
		return fmt.Errorf("%s has Z=%d, want %d for symbol %s", i.ID(), i.Z, expectedZ, i.Symbol)
	}
	if i.A < i.Z {
		return fmt.Errorf("%s has mass number below atomic number", i.ID())
	}
	if strings.TrimSpace(i.CitationLink) == "" {
		return fmt.Errorf("%s missing citation link", i.ID())
	}
	return nil
}

func ValidateCatalog(catalog Catalog) error {
	for id, isotope := range catalog {
		if _, err := ParseNuclideID(id); err != nil {
			return err
		}
		if isotope.ID() != id {
			return fmt.Errorf("catalog key %s does not match isotope ID %s", id, isotope.ID())
		}
		if err := isotope.Validate(); err != nil {
			detail := strings.TrimPrefix(err.Error(), isotope.ID()+" ")
			return fmt.Errorf("catalog key %s %s", id, detail)
		}
		if isotope.Daughter != "" {
			if err := ValidateAlphaDaughterID(id, isotope.Daughter); err != nil {
				return err
			}
		}
	}
	return nil
}

type SimulationOptions struct {
	Samples int   `json:"samples"`
	Seed    int64 `json:"seed"`
}

type DecaySimulationRecord struct {
	ID                      string  `json:"id"`
	SourcePath              string  `json:"source_path"`
	HalfLifeSeconds         float64 `json:"half_life_seconds"`
	DecayConstantPerSecond  float64 `json:"decay_constant_per_second"`
	MeanLifeSeconds         float64 `json:"mean_life_seconds"`
	MonteCarloP05Seconds    float64 `json:"monte_carlo_p05_seconds"`
	MonteCarloMedianSeconds float64 `json:"monte_carlo_median_seconds"`
	MonteCarloP95Seconds    float64 `json:"monte_carlo_p95_seconds"`
}

type DecaySimulationReport struct {
	StartID     string                  `json:"start_id"`
	Method      string                  `json:"method"`
	SampleCount int                     `json:"sample_count"`
	Seed        int64                   `json:"seed"`
	Records     []DecaySimulationRecord `json:"records"`
}

func DecayChain(start string, catalog Catalog) ([]Isotope, error) {
	if start == "" {
		return nil, errors.New("start isotope is required")
	}

	visited := make(map[string]bool)
	chain := make([]Isotope, 0, 8)
	if err := traverse(start, "", catalog, visited, nil, &chain); err != nil {
		return nil, err
	}
	return chain, nil
}

// DecayChainGraceful traverses the decay chain and stops when a daughter
// is not found in the catalog, returning the partial chain without an error.
// Use this for display/UI paths where the catalog boundary is expected.
func DecayChainGraceful(start string, catalog Catalog) []Isotope {
	if start == "" {
		return nil
	}
	visited := make(map[string]bool)
	chain := make([]Isotope, 0, 8)
	_ = traverseGraceful(start, catalog, visited, &chain)
	return chain
}

func traverseGraceful(id string, catalog Catalog, visited map[string]bool, chain *[]Isotope) error {
	if visited[id] {
		return nil
	}
	isotope, ok := catalog[id]
	if !ok {
		return fmt.Errorf("catalog boundary reached at %s", id)
	}
	if err := isotope.Validate(); err != nil {
		return nil
	}
	visited[id] = true
	*chain = append(*chain, isotope)
	if isotope.Daughter == "" {
		return nil
	}
	if err := ValidateAlphaDaughterID(id, isotope.Daughter); err != nil {
		return fmt.Errorf("catalog boundary reached at %s: %w", id, err)
	}
	return traverseGraceful(isotope.Daughter, catalog, visited, chain)
}

func DecaySimulationSummary(start string, catalog Catalog, options SimulationOptions) (DecaySimulationReport, error) {
	if options.Samples <= 0 {
		return DecaySimulationReport{}, errors.New("simulation sample count must be positive")
	}
	if options.Seed == 0 {
		return DecaySimulationReport{}, errors.New("simulation seed must be non-zero")
	}
	chain, err := DecayChain(start, catalog)
	if err != nil {
		return DecaySimulationReport{}, err
	}

	rng := rand.New(rand.NewSource(options.Seed))
	report := DecaySimulationReport{
		StartID:     start,
		Method:      "exponential_decay_fixed_seed",
		SampleCount: options.Samples,
		Seed:        options.Seed,
		Records:     make([]DecaySimulationRecord, 0, len(chain)),
	}
	for _, isotope := range chain {
		halfLifeSeconds := isotope.HalfLife.Seconds()
		if halfLifeSeconds <= 0 {
			return DecaySimulationReport{}, fmt.Errorf("%s half-life must be positive for decay simulation", isotope.ID())
		}
		lambda := math.Ln2 / halfLifeSeconds
		meanLife := 1 / lambda
		samples := exponentialDecaySamples(rng, meanLife, options.Samples)
		report.Records = append(report.Records, DecaySimulationRecord{
			ID:                      isotope.ID(),
			SourcePath:              isotope.CitationLink,
			HalfLifeSeconds:         halfLifeSeconds,
			DecayConstantPerSecond:  lambda,
			MeanLifeSeconds:         meanLife,
			MonteCarloP05Seconds:    quantile(samples, 0.05),
			MonteCarloMedianSeconds: quantile(samples, 0.50),
			MonteCarloP95Seconds:    quantile(samples, 0.95),
		})
	}
	return report, nil
}

func exponentialDecaySamples(rng *rand.Rand, meanLife float64, count int) []float64 {
	samples := make([]float64, count)
	for i := range samples {
		u := rng.Float64()
		for u <= 0 {
			u = rng.Float64()
		}
		samples[i] = -meanLife * math.Log(1-u)
	}
	sort.Float64s(samples)
	return samples
}

func quantile(sortedSamples []float64, q float64) float64 {
	if len(sortedSamples) == 0 {
		return 0
	}
	index := int(math.Round(q * float64(len(sortedSamples)-1)))
	if index < 0 {
		index = 0
	}
	if index >= len(sortedSamples) {
		index = len(sortedSamples) - 1
	}
	return sortedSamples[index]
}

func traverse(id string, parent string, catalog Catalog, visited map[string]bool, path []string, chain *[]Isotope) error {
	if visited[id] {
		return fmt.Errorf("cycle detected: %s", formatCycle(path, id))
	}

	isotope, ok := catalog[id]
	if !ok {
		if parent != "" {
			return fmt.Errorf("daughter %s referenced by %s not found in catalog", id, parent)
		}
		return fmt.Errorf("isotope %s not found in catalog", id)
	}
	if isotope.ID() != id {
		return fmt.Errorf("catalog key %s does not match isotope ID %s", id, isotope.ID())
	}
	if err := isotope.Validate(); err != nil {
		return err
	}

	visited[id] = true
	path = append(path, id)
	*chain = append(*chain, isotope)
	if isotope.Daughter == "" {
		return nil
	}
	if !visited[isotope.Daughter] {
		if err := ValidateAlphaDaughterID(id, isotope.Daughter); err != nil {
			return err
		}
	}
	return traverse(isotope.Daughter, id, catalog, visited, path, chain)
}

func formatCycle(path []string, id string) string {
	start := 0
	for index, pathID := range path {
		if pathID == id {
			start = index
			break
		}
	}

	cycle := append([]string{}, path[start:]...)
	cycle = append(cycle, id)
	return strings.Join(cycle, " -> ")
}
