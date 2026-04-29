package physics

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestDecayChainTraversesAlphaSequence(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, HalfLife: 170 * time.Millisecond, QAlphaMeV: 10.75, Daughter: "284Nh", CitationLink: "https://www.nndc.bnl.gov/ensnds/288/Mc/adopted.pdf"},
		"284Nh": {Symbol: "Nh", Z: 113, A: 284, Daughter: "280Rg", CitationLink: "seed"},
		"280Rg": {Symbol: "Rg", Z: 111, A: 280, Daughter: "276Mt", CitationLink: "seed"},
		"276Mt": {Symbol: "Mt", Z: 109, A: 276, Daughter: "272Bh", CitationLink: "seed"},
		"272Bh": {Symbol: "Bh", Z: 107, A: 272, Daughter: "268Db", CitationLink: "seed"},
		"268Db": {Symbol: "Db", Z: 105, A: 268, Daughter: "264Lr", CitationLink: "seed"},
		"264Lr": {Symbol: "Lr", Z: 103, A: 264, CitationLink: "seed"},
	}

	chain, err := DecayChain("288Mc", catalog)
	if err != nil {
		t.Fatalf("DecayChain returned error: %v", err)
	}

	got := make([]string, 0, len(chain))
	for _, isotope := range chain {
		got = append(got, isotope.ID())
	}

	want := []string{"288Mc", "284Nh", "280Rg", "276Mt", "272Bh", "268Db", "264Lr"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("chain IDs = %v, want %v", got, want)
	}
}

func TestDecayChainRejectsCycles(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, Daughter: "284Nh", CitationLink: "seed"},
		"284Nh": {Symbol: "Nh", Z: 113, A: 284, Daughter: "288Mc", CitationLink: "seed"},
	}

	_, err := DecayChain("288Mc", catalog)
	if err == nil {
		t.Fatal("DecayChain returned nil error for cyclic decay graph")
	}
}

func TestIsotopeRequiresCitation(t *testing.T) {
	isotope := Isotope{Symbol: "Mc", Z: 115, A: 288}
	if err := isotope.Validate(); err == nil {
		t.Fatal("Validate returned nil error for isotope without citation")
	}
}

func TestIsotopeRejectsBlankCitation(t *testing.T) {
	isotope := Isotope{Symbol: "Mc", Z: 115, A: 288, CitationLink: " \t\n"}
	err := isotope.Validate()
	if err == nil {
		t.Fatal("Validate returned nil error for blank citation")
	}
	if !strings.Contains(err.Error(), "288Mc missing citation link") {
		t.Fatalf("error = %q, want missing citation link", err)
	}
}

func TestIsotopeRejectsAtomicNumberSymbolMismatch(t *testing.T) {
	isotope := Isotope{Symbol: "Mc", Z: 113, A: 288, CitationLink: "seed"}
	err := isotope.Validate()
	if err == nil {
		t.Fatal("Validate returned nil error for atomic-number/symbol mismatch")
	}
	if !strings.Contains(err.Error(), "288Mc has Z=113, want 115 for symbol Mc") {
		t.Fatalf("error = %q, want atomic-number/symbol mismatch", err)
	}
}

func TestDecayChainRejectsCatalogKeyMismatch(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 290, CitationLink: "seed"},
	}

	_, err := DecayChain("288Mc", catalog)
	if err == nil {
		t.Fatal("DecayChain returned nil error for catalog key mismatch")
	}
	if !strings.Contains(err.Error(), "catalog key 288Mc does not match isotope ID 290Mc") {
		t.Fatalf("error = %q, want catalog key mismatch", err)
	}
}

func TestValidateCatalogRejectsKeyIdentityAndDaughterErrorsBeforeTraversal(t *testing.T) {
	validCatalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, Daughter: "284Nh", CitationLink: "seed"},
		"284Nh": {Symbol: "Nh", Z: 113, A: 284, CitationLink: "seed"},
	}
	if err := ValidateCatalog(validCatalog); err != nil {
		t.Fatalf("ValidateCatalog rejected valid catalog: %v", err)
	}

	cases := []struct {
		name   string
		mutate func(Catalog)
		want   string
	}{
		{
			name: "malformed catalog key",
			mutate: func(catalog Catalog) {
				catalog["Mc288"] = catalog["288Mc"]
				delete(catalog, "288Mc")
			},
			want: "nuclide ID \"Mc288\" must start with a mass number",
		},
		{
			name: "key isotope identity mismatch",
			mutate: func(catalog Catalog) {
				catalog["288Mc"] = Isotope{Symbol: "Mc", Z: 115, A: 290, Daughter: "284Nh", CitationLink: "seed"}
			},
			want: "catalog key 288Mc does not match isotope ID 290Mc",
		},
		{
			name: "key atomic number mismatch",
			mutate: func(catalog Catalog) {
				catalog["288Mc"] = Isotope{Symbol: "Mc", Z: 113, A: 288, Daughter: "284Nh", CitationLink: "seed"}
			},
			want: "catalog key 288Mc has Z=113, want 115 for symbol Mc",
		},
		{
			name: "non-alpha daughter",
			mutate: func(catalog Catalog) {
				catalog["288Mc"] = Isotope{Symbol: "Mc", Z: 115, A: 288, Daughter: "285Nh", CitationLink: "seed"}
				catalog["285Nh"] = Isotope{Symbol: "Nh", Z: 113, A: 285, CitationLink: "seed"}
				delete(catalog, "284Nh")
			},
			want: "288Mc alpha daughter 285Nh has A=285, want 284",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			catalog := cloneCatalog(validCatalog)
			tc.mutate(catalog)
			err := ValidateCatalog(catalog)
			if err == nil {
				t.Fatal("ValidateCatalog accepted invalid catalog")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tc.want)
			}
		})
	}
}

func cloneCatalog(catalog Catalog) Catalog {
	clone := make(Catalog, len(catalog))
	for id, isotope := range catalog {
		clone[id] = isotope
	}
	return clone
}

func TestDecayChainReportsMissingDaughterParent(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, Daughter: "284Nh", CitationLink: "seed"},
	}

	_, err := DecayChain("288Mc", catalog)
	if err == nil {
		t.Fatal("DecayChain returned nil error for missing daughter")
	}
	if !strings.Contains(err.Error(), "daughter 284Nh referenced by 288Mc not found in catalog") {
		t.Fatalf("error = %q, want missing daughter with parent context", err)
	}
}

func TestDecayChainRejectsNonAlphaDaughterTransition(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, Daughter: "285Nh", CitationLink: "seed"},
		"285Nh": {Symbol: "Nh", Z: 113, A: 285, CitationLink: "seed"},
	}

	_, err := DecayChain("288Mc", catalog)
	if err == nil {
		t.Fatal("DecayChain returned nil error for non-alpha daughter transition")
	}
	if !strings.Contains(err.Error(), "288Mc alpha daughter 285Nh has A=285, want 284") {
		t.Fatalf("error = %q, want alpha daughter mass mismatch", err)
	}
}

func TestDecayChainReportsCyclePath(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, Daughter: "284Nh", CitationLink: "seed"},
		"284Nh": {Symbol: "Nh", Z: 113, A: 284, Daughter: "288Mc", CitationLink: "seed"},
	}

	_, err := DecayChain("288Mc", catalog)
	if err == nil {
		t.Fatal("DecayChain returned nil error for cyclic decay graph")
	}
	if !strings.Contains(err.Error(), "cycle detected: 288Mc -> 284Nh -> 288Mc") {
		t.Fatalf("error = %q, want cycle path", err)
	}
}

func TestDecaySimulationSummaryComputesSourceBackedMetricsAndDeterministicSamples(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, HalfLife: 170 * time.Millisecond, QAlphaMeV: 10.75, Daughter: "284Nh", CitationLink: "https://www.nndc.bnl.gov/ensnds/288/Mc/adopted.pdf"},
		"284Nh": {Symbol: "Nh", Z: 113, A: 284, HalfLife: 1 * time.Second, Daughter: "280Rg", CitationLink: "data/research.seed.json"},
		"280Rg": {Symbol: "Rg", Z: 111, A: 280, HalfLife: 1 * time.Second, CitationLink: "data/research.seed.json"},
	}

	summary, err := DecaySimulationSummary("288Mc", catalog, SimulationOptions{Samples: 5, Seed: 20260429})
	if err != nil {
		t.Fatalf("DecaySimulationSummary returned error: %v", err)
	}
	if summary.StartID != "288Mc" || summary.SampleCount != 5 || summary.Seed != 20260429 {
		t.Fatalf("summary metadata = start %q samples %d seed %d", summary.StartID, summary.SampleCount, summary.Seed)
	}
	if summary.Method != "exponential_decay_fixed_seed" {
		t.Fatalf("method = %q, want exponential_decay_fixed_seed", summary.Method)
	}
	if got, want := len(summary.Records), 3; got != want {
		t.Fatalf("record count = %d, want %d", got, want)
	}
	first := summary.Records[0]
	if first.ID != "288Mc" || first.SourcePath != "https://www.nndc.bnl.gov/ensnds/288/Mc/adopted.pdf" {
		t.Fatalf("first record identity/source = %+v", first)
	}
	if !almostEqual(first.DecayConstantPerSecond, 4.077336, 0.000001) {
		t.Fatalf("lambda = %.9f, want about 4.077336", first.DecayConstantPerSecond)
	}
	if !almostEqual(first.MeanLifeSeconds, 0.245258, 0.000001) {
		t.Fatalf("mean life = %.9f, want about 0.245258", first.MeanLifeSeconds)
	}
	if first.MonteCarloP05Seconds <= 0 || first.MonteCarloMedianSeconds <= 0 || first.MonteCarloP95Seconds <= 0 {
		t.Fatalf("Monte Carlo quantiles must be positive: %+v", first)
	}
	if first.MonteCarloP05Seconds > first.MonteCarloMedianSeconds || first.MonteCarloMedianSeconds > first.MonteCarloP95Seconds {
		t.Fatalf("Monte Carlo quantiles are not ordered: %+v", first)
	}

	repeat, err := DecaySimulationSummary("288Mc", catalog, SimulationOptions{Samples: 5, Seed: 20260429})
	if err != nil {
		t.Fatalf("repeat DecaySimulationSummary returned error: %v", err)
	}
	if !reflect.DeepEqual(summary, repeat) {
		t.Fatalf("fixed-seed summary is not deterministic\nfirst:  %+v\nsecond: %+v", summary, repeat)
	}
}

func TestDecaySimulationSummaryRejectsInvalidSimulationInputs(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, HalfLife: 170 * time.Millisecond, CitationLink: "seed"},
	}
	for _, tc := range []struct {
		name    string
		catalog Catalog
		options SimulationOptions
		want    string
	}{
		{name: "zero samples", catalog: catalog, options: SimulationOptions{Samples: 0, Seed: 1}, want: "simulation sample count must be positive"},
		{name: "zero seed", catalog: catalog, options: SimulationOptions{Samples: 1}, want: "simulation seed must be non-zero"},
		{name: "missing half life", catalog: Catalog{"288Mc": {Symbol: "Mc", Z: 115, A: 288, CitationLink: "seed"}}, options: SimulationOptions{Samples: 1, Seed: 1}, want: "288Mc half-life must be positive for decay simulation"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecaySimulationSummary("288Mc", tc.catalog, tc.options)
			if err == nil {
				t.Fatal("DecaySimulationSummary returned nil error")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %q, want %q", err, tc.want)
			}
		})
	}
}

func almostEqual(got, want, tolerance float64) bool {
	if got > want {
		return got-want <= tolerance
	}
	return want-got <= tolerance
}
