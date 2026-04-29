package physics

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type NuclideID struct {
	ID     string
	A      int
	Symbol string
	Z      int
}

var elementNameBySymbol = map[string]string{
	"H":  "Hydrogen",
	"He": "Helium",
	"Li": "Lithium",
	"Be": "Beryllium",
	"B":  "Boron",
	"C":  "Carbon",
	"N":  "Nitrogen",
	"O":  "Oxygen",
	"F":  "Fluorine",
	"Ne": "Neon",
	"Na": "Sodium",
	"Mg": "Magnesium",
	"Al": "Aluminium",
	"Si": "Silicon",
	"P":  "Phosphorus",
	"S":  "Sulfur",
	"Cl": "Chlorine",
	"Ar": "Argon",
	"K":  "Potassium",
	"Ca": "Calcium",
	"Sc": "Scandium",
	"Ti": "Titanium",
	"V":  "Vanadium",
	"Cr": "Chromium",
	"Mn": "Manganese",
	"Fe": "Iron",
	"Co": "Cobalt",
	"Ni": "Nickel",
	"Cu": "Copper",
	"Zn": "Zinc",
	"Ga": "Gallium",
	"Ge": "Germanium",
	"As": "Arsenic",
	"Se": "Selenium",
	"Br": "Bromine",
	"Kr": "Krypton",
	"Rb": "Rubidium",
	"Sr": "Strontium",
	"Y":  "Yttrium",
	"Zr": "Zirconium",
	"Nb": "Niobium",
	"Mo": "Molybdenum",
	"Tc": "Technetium",
	"Ru": "Ruthenium",
	"Rh": "Rhodium",
	"Pd": "Palladium",
	"Ag": "Silver",
	"Cd": "Cadmium",
	"In": "Indium",
	"Sn": "Tin",
	"Sb": "Antimony",
	"Te": "Tellurium",
	"I":  "Iodine",
	"Xe": "Xenon",
	"Cs": "Caesium",
	"Ba": "Barium",
	"La": "Lanthanum",
	"Ce": "Cerium",
	"Pr": "Praseodymium",
	"Nd": "Neodymium",
	"Pm": "Promethium",
	"Sm": "Samarium",
	"Eu": "Europium",
	"Gd": "Gadolinium",
	"Tb": "Terbium",
	"Dy": "Dysprosium",
	"Ho": "Holmium",
	"Er": "Erbium",
	"Tm": "Thulium",
	"Yb": "Ytterbium",
	"Lu": "Lutetium",
	"Hf": "Hafnium",
	"Ta": "Tantalum",
	"W":  "Tungsten",
	"Re": "Rhenium",
	"Os": "Osmium",
	"Ir": "Iridium",
	"Pt": "Platinum",
	"Au": "Gold",
	"Hg": "Mercury",
	"Tl": "Thallium",
	"Pb": "Lead",
	"Bi": "Bismuth",
	"Po": "Polonium",
	"At": "Astatine",
	"Rn": "Radon",
	"Fr": "Francium",
	"Ra": "Radium",
	"Ac": "Actinium",
	"Th": "Thorium",
	"Pa": "Protactinium",
	"U":  "Uranium",
	"Np": "Neptunium",
	"Pu": "Plutonium",
	"Am": "Americium",
	"Cm": "Curium",
	"Bk": "Berkelium",
	"Cf": "Californium",
	"Es": "Einsteinium",
	"Fm": "Fermium",
	"Md": "Mendelevium",
	"No": "Nobelium",
	"Lr": "Lawrencium",
	"Rf": "Rutherfordium",
	"Db": "Dubnium",
	"Sg": "Seaborgium",
	"Bh": "Bohrium",
	"Hs": "Hassium",
	"Mt": "Meitnerium",
	"Ds": "Darmstadtium",
	"Rg": "Roentgenium",
	"Cn": "Copernicium",
	"Nh": "Nihonium",
	"Fl": "Flerovium",
	"Mc": "Moscovium",
	"Lv": "Livermorium",
	"Ts": "Tennessine",
	"Og": "Oganesson",
}

var atomicNumberBySymbol = map[string]int{
	"H":  1,
	"He": 2,
	"Li": 3,
	"Be": 4,
	"B":  5,
	"C":  6,
	"N":  7,
	"O":  8,
	"F":  9,
	"Ne": 10,
	"Na": 11,
	"Mg": 12,
	"Al": 13,
	"Si": 14,
	"P":  15,
	"S":  16,
	"Cl": 17,
	"Ar": 18,
	"K":  19,
	"Ca": 20,
	"Sc": 21,
	"Ti": 22,
	"V":  23,
	"Cr": 24,
	"Mn": 25,
	"Fe": 26,
	"Co": 27,
	"Ni": 28,
	"Cu": 29,
	"Zn": 30,
	"Ga": 31,
	"Ge": 32,
	"As": 33,
	"Se": 34,
	"Br": 35,
	"Kr": 36,
	"Rb": 37,
	"Sr": 38,
	"Y":  39,
	"Zr": 40,
	"Nb": 41,
	"Mo": 42,
	"Tc": 43,
	"Ru": 44,
	"Rh": 45,
	"Pd": 46,
	"Ag": 47,
	"Cd": 48,
	"In": 49,
	"Sn": 50,
	"Sb": 51,
	"Te": 52,
	"I":  53,
	"Xe": 54,
	"Cs": 55,
	"Ba": 56,
	"La": 57,
	"Ce": 58,
	"Pr": 59,
	"Nd": 60,
	"Pm": 61,
	"Sm": 62,
	"Eu": 63,
	"Gd": 64,
	"Tb": 65,
	"Dy": 66,
	"Ho": 67,
	"Er": 68,
	"Tm": 69,
	"Yb": 70,
	"Lu": 71,
	"Hf": 72,
	"Ta": 73,
	"W":  74,
	"Re": 75,
	"Os": 76,
	"Ir": 77,
	"Pt": 78,
	"Au": 79,
	"Hg": 80,
	"Tl": 81,
	"Pb": 82,
	"Bi": 83,
	"Po": 84,
	"At": 85,
	"Rn": 86,
	"Fr": 87,
	"Ra": 88,
	"Ac": 89,
	"Th": 90,
	"Pa": 91,
	"U":  92,
	"Np": 93,
	"Pu": 94,
	"Am": 95,
	"Cm": 96,
	"Bk": 97,
	"Cf": 98,
	"Es": 99,
	"Fm": 100,
	"Md": 101,
	"No": 102,
	"Lr": 103,
	"Rf": 104,
	"Db": 105,
	"Sg": 106,
	"Bh": 107,
	"Hs": 108,
	"Mt": 109,
	"Ds": 110,
	"Rg": 111,
	"Cn": 112,
	"Nh": 113,
	"Fl": 114,
	"Mc": 115,
	"Lv": 116,
	"Ts": 117,
	"Og": 118,
}

func AtomicNumberForSymbol(symbol string) (int, bool) {
	z, ok := atomicNumberBySymbol[symbol]
	return z, ok
}

func ElementNameForSymbol(symbol string) (string, bool) {
	name, ok := elementNameBySymbol[symbol]
	return name, ok
}

func SymbolForAtomicNumber(z int) (string, bool) {
	for symbol, atomicNumber := range atomicNumberBySymbol {
		if atomicNumber == z {
			return symbol, true
		}
	}
	return "", false
}

func ParseNuclideID(id string) (NuclideID, error) {
	trimmed := strings.TrimSpace(id)
	if trimmed == "" {
		return NuclideID{}, fmt.Errorf("nuclide ID is required")
	}

	massEnd := 0
	for massEnd < len(trimmed) {
		r := rune(trimmed[massEnd])
		if !unicode.IsDigit(r) {
			break
		}
		massEnd++
	}
	if massEnd == 0 {
		return NuclideID{}, fmt.Errorf("nuclide ID %q must start with a mass number", id)
	}
	if massEnd == len(trimmed) {
		return NuclideID{}, fmt.Errorf("nuclide ID %q must include an element symbol", id)
	}

	mass, err := strconv.Atoi(trimmed[:massEnd])
	if err != nil || mass <= 0 {
		return NuclideID{}, fmt.Errorf("nuclide ID %q mass number must be positive", id)
	}

	symbol := trimmed[massEnd:]
	for _, r := range symbol {
		if !unicode.IsLetter(r) {
			return NuclideID{}, fmt.Errorf("nuclide ID %q has invalid element symbol %q", id, symbol)
		}
	}
	z, ok := AtomicNumberForSymbol(symbol)
	if !ok {
		return NuclideID{}, fmt.Errorf("nuclide ID %q has unknown element symbol %q", id, symbol)
	}

	return NuclideID{ID: trimmed, A: mass, Symbol: symbol, Z: z}, nil
}

func ValidateAlphaDaughterID(parentID, daughterID string) error {
	parent, err := ParseNuclideID(parentID)
	if err != nil {
		return err
	}
	daughter, err := ParseNuclideID(daughterID)
	if err != nil {
		return err
	}

	expectedA := parent.A - 4
	if daughter.A != expectedA {
		return fmt.Errorf("%s alpha daughter %s has A=%d, want %d", parent.ID, daughter.ID, daughter.A, expectedA)
	}
	expectedZ := parent.Z - 2
	if daughter.Z != expectedZ {
		expectedSymbol, ok := SymbolForAtomicNumber(expectedZ)
		if !ok {
			return fmt.Errorf("%s alpha daughter %s has Z=%d, want %d", parent.ID, daughter.ID, daughter.Z, expectedZ)
		}
		expectedID := fmt.Sprintf("%d%s", expectedA, expectedSymbol)
		return fmt.Errorf("%s alpha daughter %s has Z=%d, want %d for expected alpha daughter %s", parent.ID, daughter.ID, daughter.Z, expectedZ, expectedID)
	}
	return nil
}
