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
