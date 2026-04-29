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

func TestElementsCoversExactly118UniqueChemicalElements(t *testing.T) {
	elements := Elements()
	if len(elements) != 118 {
		t.Fatalf("Elements() returned %d entries, want 118", len(elements))
	}

	seenSymbols := make(map[string]bool, len(elements))
	seenAtomicNumbers := make(map[int]bool, len(elements))
	for _, element := range elements {
		if element.Z < 1 || element.Z > 118 {
			t.Fatalf("Elements() contains %s with Z=%d outside 1..118", element.Symbol, element.Z)
		}
		if element.Symbol == "" {
			t.Fatalf("Elements() contains empty symbol at Z=%d", element.Z)
		}
		if element.Name == "" {
			t.Fatalf("Elements() contains empty name for symbol %s", element.Symbol)
		}
		if seenSymbols[element.Symbol] {
			t.Fatalf("Elements() contains duplicate symbol %s", element.Symbol)
		}
		if seenAtomicNumbers[element.Z] {
			t.Fatalf("Elements() contains duplicate atomic number %d", element.Z)
		}
		seenSymbols[element.Symbol] = true
		seenAtomicNumbers[element.Z] = true
	}

	spotChecks := map[string]struct {
		z    int
		name string
	}{
		"Mc": {z: 115, name: "Moscovium"},
		"Nh": {z: 113, name: "Nihonium"},
	}
	for symbol, want := range spotChecks {
		matched := false
		for _, element := range elements {
			if element.Symbol == symbol {
				matched = true
				if element.Z != want.z || element.Name != want.name {
					t.Fatalf("Elements() %s entry = Z=%d Name=%q, want Z=%d Name=%q", symbol, element.Z, element.Name, want.z, want.name)
				}
			}
		}
		if !matched {
			t.Fatalf("Elements() missing symbol %s", symbol)
		}
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

func TestElementNameForSymbolCoversCurrentSeedSymbols(t *testing.T) {
	cases := map[string]string{
		"Mc": "Moscovium",
		"Nh": "Nihonium",
		"He": "Helium",
	}

	for symbol, wantName := range cases {
		gotName, ok := ElementNameForSymbol(symbol)
		if !ok {
			t.Fatalf("ElementNameForSymbol(%q) not found", symbol)
		}
		if gotName != wantName {
			t.Fatalf("ElementNameForSymbol(%q) = %q, want %q", symbol, gotName, wantName)
		}
	}
}

func TestSymbolForAtomicNumberCoversCurrentAlphaChainSymbols(t *testing.T) {
	cases := map[int]string{
		115: "Mc",
		113: "Nh",
		111: "Rg",
		109: "Mt",
		107: "Bh",
		105: "Db",
		103: "Lr",
	}

	for z, wantSymbol := range cases {
		gotSymbol, ok := SymbolForAtomicNumber(z)
		if !ok {
			t.Fatalf("SymbolForAtomicNumber(%d) not found", z)
		}
		if gotSymbol != wantSymbol {
			t.Fatalf("SymbolForAtomicNumber(%d) = %q, want %q", z, gotSymbol, wantSymbol)
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
		{parent: "288Mc", daughter: "284Fl", want: "Z=114, want 113 for expected alpha daughter 284Nh"},
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
