package physics

import (
	"strings"
	"testing"
)

func TestParseNuclideIDExtractsMassSymbolAndAtomicNumber(t *testing.T) {
	cases := []struct {
		id     string
		mass   int
		symbol string
		z      int
	}{
		{id: "288Mc", mass: 288, symbol: "Mc", z: 115},
		{id: "284Nh", mass: 284, symbol: "Nh", z: 113},
		{id: "264Lr", mass: 264, symbol: "Lr", z: 103},
		{id: "4He", mass: 4, symbol: "He", z: 2},
	}

	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			nuclide, err := ParseNuclideID(tc.id)
			if err != nil {
				t.Fatalf("ParseNuclideID(%q) returned error: %v", tc.id, err)
			}
			if nuclide.ID != tc.id || nuclide.A != tc.mass || nuclide.Symbol != tc.symbol || nuclide.Z != tc.z {
				t.Fatalf("ParseNuclideID(%q) = %+v, want ID=%q A=%d Symbol=%q Z=%d", tc.id, nuclide, tc.id, tc.mass, tc.symbol, tc.z)
			}
		})
	}
}

func TestParseNuclideIDRejectsMalformedOrUnknownIDs(t *testing.T) {
	cases := []struct {
		id   string
		want string
	}{
		{id: "", want: "nuclide ID is required"},
		{id: "Mc288", want: "must start with a mass number"},
		{id: "288", want: "must include an element symbol"},
		{id: "0Mc", want: "mass number must be positive"},
		{id: "288mc", want: "unknown element symbol"},
		{id: "288Xx", want: "unknown element symbol"},
		{id: "288Mc+", want: "invalid element symbol"},
	}

	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			_, err := ParseNuclideID(tc.id)
			if err == nil {
				t.Fatalf("ParseNuclideID(%q) returned nil error", tc.id)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("ParseNuclideID(%q) error = %q, want substring %q", tc.id, err.Error(), tc.want)
			}
		})
	}
}

func TestAtomicNumberForSymbolCoversCurrentAlphaChainSymbols(t *testing.T) {
	cases := map[string]int{
		"Mc": 115,
		"Nh": 113,
		"Rg": 111,
		"Mt": 109,
		"Bh": 107,
		"Db": 105,
		"Lr": 103,
	}

	for symbol, wantZ := range cases {
		gotZ, ok := AtomicNumberForSymbol(symbol)
		if !ok {
			t.Fatalf("AtomicNumberForSymbol(%q) not found", symbol)
		}
		if gotZ != wantZ {
			t.Fatalf("AtomicNumberForSymbol(%q) = %d, want %d", symbol, gotZ, wantZ)
		}
	}
}

func TestAlphaDaughterIDComputesExpectedMassAndAtomicNumber(t *testing.T) {
	cases := []struct {
		parent   string
		daughter string
	}{
		{parent: "288Mc", daughter: "284Nh"},
		{parent: "290Mc", daughter: "286Nh"},
		{parent: "284Nh", daughter: "280Rg"},
	}

	for _, tc := range cases {
		t.Run(tc.parent+"->"+tc.daughter, func(t *testing.T) {
			if err := ValidateAlphaDaughterID(tc.parent, tc.daughter); err != nil {
				t.Fatalf("ValidateAlphaDaughterID(%q, %q) returned error: %v", tc.parent, tc.daughter, err)
			}
		})
	}
}

func TestAlphaDaughterIDRejectsNonAlphaTransitions(t *testing.T) {
	cases := []struct {
		parent   string
		daughter string
		want     string
	}{
		{parent: "288Mc", daughter: "285Nh", want: "A=285, want 284"},
		{parent: "288Mc", daughter: "284Fl", want: "Z=114, want 113"},
		{parent: "288Mc", daughter: "284Xx", want: "unknown element symbol"},
	}

	for _, tc := range cases {
		t.Run(tc.parent+"->"+tc.daughter, func(t *testing.T) {
			err := ValidateAlphaDaughterID(tc.parent, tc.daughter)
			if err == nil {
				t.Fatalf("ValidateAlphaDaughterID(%q, %q) returned nil error", tc.parent, tc.daughter)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("ValidateAlphaDaughterID(%q, %q) error = %q, want substring %q", tc.parent, tc.daughter, err.Error(), tc.want)
			}
		})
	}
}
