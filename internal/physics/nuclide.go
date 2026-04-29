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

type Element struct {
	Z      int
	Symbol string
	Name   string
}

var periodicTable = []Element{
	{Z: 1, Symbol: "H", Name: "Hydrogen"},
	{Z: 2, Symbol: "He", Name: "Helium"},
	{Z: 3, Symbol: "Li", Name: "Lithium"},
	{Z: 4, Symbol: "Be", Name: "Beryllium"},
	{Z: 5, Symbol: "B", Name: "Boron"},
	{Z: 6, Symbol: "C", Name: "Carbon"},
	{Z: 7, Symbol: "N", Name: "Nitrogen"},
	{Z: 8, Symbol: "O", Name: "Oxygen"},
	{Z: 9, Symbol: "F", Name: "Fluorine"},
	{Z: 10, Symbol: "Ne", Name: "Neon"},
	{Z: 11, Symbol: "Na", Name: "Sodium"},
	{Z: 12, Symbol: "Mg", Name: "Magnesium"},
	{Z: 13, Symbol: "Al", Name: "Aluminium"},
	{Z: 14, Symbol: "Si", Name: "Silicon"},
	{Z: 15, Symbol: "P", Name: "Phosphorus"},
	{Z: 16, Symbol: "S", Name: "Sulfur"},
	{Z: 17, Symbol: "Cl", Name: "Chlorine"},
	{Z: 18, Symbol: "Ar", Name: "Argon"},
	{Z: 19, Symbol: "K", Name: "Potassium"},
	{Z: 20, Symbol: "Ca", Name: "Calcium"},
	{Z: 21, Symbol: "Sc", Name: "Scandium"},
	{Z: 22, Symbol: "Ti", Name: "Titanium"},
	{Z: 23, Symbol: "V", Name: "Vanadium"},
	{Z: 24, Symbol: "Cr", Name: "Chromium"},
	{Z: 25, Symbol: "Mn", Name: "Manganese"},
	{Z: 26, Symbol: "Fe", Name: "Iron"},
	{Z: 27, Symbol: "Co", Name: "Cobalt"},
	{Z: 28, Symbol: "Ni", Name: "Nickel"},
	{Z: 29, Symbol: "Cu", Name: "Copper"},
	{Z: 30, Symbol: "Zn", Name: "Zinc"},
	{Z: 31, Symbol: "Ga", Name: "Gallium"},
	{Z: 32, Symbol: "Ge", Name: "Germanium"},
	{Z: 33, Symbol: "As", Name: "Arsenic"},
	{Z: 34, Symbol: "Se", Name: "Selenium"},
	{Z: 35, Symbol: "Br", Name: "Bromine"},
	{Z: 36, Symbol: "Kr", Name: "Krypton"},
	{Z: 37, Symbol: "Rb", Name: "Rubidium"},
	{Z: 38, Symbol: "Sr", Name: "Strontium"},
	{Z: 39, Symbol: "Y", Name: "Yttrium"},
	{Z: 40, Symbol: "Zr", Name: "Zirconium"},
	{Z: 41, Symbol: "Nb", Name: "Niobium"},
	{Z: 42, Symbol: "Mo", Name: "Molybdenum"},
	{Z: 43, Symbol: "Tc", Name: "Technetium"},
	{Z: 44, Symbol: "Ru", Name: "Ruthenium"},
	{Z: 45, Symbol: "Rh", Name: "Rhodium"},
	{Z: 46, Symbol: "Pd", Name: "Palladium"},
	{Z: 47, Symbol: "Ag", Name: "Silver"},
	{Z: 48, Symbol: "Cd", Name: "Cadmium"},
	{Z: 49, Symbol: "In", Name: "Indium"},
	{Z: 50, Symbol: "Sn", Name: "Tin"},
	{Z: 51, Symbol: "Sb", Name: "Antimony"},
	{Z: 52, Symbol: "Te", Name: "Tellurium"},
	{Z: 53, Symbol: "I", Name: "Iodine"},
	{Z: 54, Symbol: "Xe", Name: "Xenon"},
	{Z: 55, Symbol: "Cs", Name: "Caesium"},
	{Z: 56, Symbol: "Ba", Name: "Barium"},
	{Z: 57, Symbol: "La", Name: "Lanthanum"},
	{Z: 58, Symbol: "Ce", Name: "Cerium"},
	{Z: 59, Symbol: "Pr", Name: "Praseodymium"},
	{Z: 60, Symbol: "Nd", Name: "Neodymium"},
	{Z: 61, Symbol: "Pm", Name: "Promethium"},
	{Z: 62, Symbol: "Sm", Name: "Samarium"},
	{Z: 63, Symbol: "Eu", Name: "Europium"},
	{Z: 64, Symbol: "Gd", Name: "Gadolinium"},
	{Z: 65, Symbol: "Tb", Name: "Terbium"},
	{Z: 66, Symbol: "Dy", Name: "Dysprosium"},
	{Z: 67, Symbol: "Ho", Name: "Holmium"},
	{Z: 68, Symbol: "Er", Name: "Erbium"},
	{Z: 69, Symbol: "Tm", Name: "Thulium"},
	{Z: 70, Symbol: "Yb", Name: "Ytterbium"},
	{Z: 71, Symbol: "Lu", Name: "Lutetium"},
	{Z: 72, Symbol: "Hf", Name: "Hafnium"},
	{Z: 73, Symbol: "Ta", Name: "Tantalum"},
	{Z: 74, Symbol: "W", Name: "Tungsten"},
	{Z: 75, Symbol: "Re", Name: "Rhenium"},
	{Z: 76, Symbol: "Os", Name: "Osmium"},
	{Z: 77, Symbol: "Ir", Name: "Iridium"},
	{Z: 78, Symbol: "Pt", Name: "Platinum"},
	{Z: 79, Symbol: "Au", Name: "Gold"},
	{Z: 80, Symbol: "Hg", Name: "Mercury"},
	{Z: 81, Symbol: "Tl", Name: "Thallium"},
	{Z: 82, Symbol: "Pb", Name: "Lead"},
	{Z: 83, Symbol: "Bi", Name: "Bismuth"},
	{Z: 84, Symbol: "Po", Name: "Polonium"},
	{Z: 85, Symbol: "At", Name: "Astatine"},
	{Z: 86, Symbol: "Rn", Name: "Radon"},
	{Z: 87, Symbol: "Fr", Name: "Francium"},
	{Z: 88, Symbol: "Ra", Name: "Radium"},
	{Z: 89, Symbol: "Ac", Name: "Actinium"},
	{Z: 90, Symbol: "Th", Name: "Thorium"},
	{Z: 91, Symbol: "Pa", Name: "Protactinium"},
	{Z: 92, Symbol: "U", Name: "Uranium"},
	{Z: 93, Symbol: "Np", Name: "Neptunium"},
	{Z: 94, Symbol: "Pu", Name: "Plutonium"},
	{Z: 95, Symbol: "Am", Name: "Americium"},
	{Z: 96, Symbol: "Cm", Name: "Curium"},
	{Z: 97, Symbol: "Bk", Name: "Berkelium"},
	{Z: 98, Symbol: "Cf", Name: "Californium"},
	{Z: 99, Symbol: "Es", Name: "Einsteinium"},
	{Z: 100, Symbol: "Fm", Name: "Fermium"},
	{Z: 101, Symbol: "Md", Name: "Mendelevium"},
	{Z: 102, Symbol: "No", Name: "Nobelium"},
	{Z: 103, Symbol: "Lr", Name: "Lawrencium"},
	{Z: 104, Symbol: "Rf", Name: "Rutherfordium"},
	{Z: 105, Symbol: "Db", Name: "Dubnium"},
	{Z: 106, Symbol: "Sg", Name: "Seaborgium"},
	{Z: 107, Symbol: "Bh", Name: "Bohrium"},
	{Z: 108, Symbol: "Hs", Name: "Hassium"},
	{Z: 109, Symbol: "Mt", Name: "Meitnerium"},
	{Z: 110, Symbol: "Ds", Name: "Darmstadtium"},
	{Z: 111, Symbol: "Rg", Name: "Roentgenium"},
	{Z: 112, Symbol: "Cn", Name: "Copernicium"},
	{Z: 113, Symbol: "Nh", Name: "Nihonium"},
	{Z: 114, Symbol: "Fl", Name: "Flerovium"},
	{Z: 115, Symbol: "Mc", Name: "Moscovium"},
	{Z: 116, Symbol: "Lv", Name: "Livermorium"},
	{Z: 117, Symbol: "Ts", Name: "Tennessine"},
	{Z: 118, Symbol: "Og", Name: "Oganesson"},
}

var atomicNumberBySymbol = func() map[string]int {
	lookup := make(map[string]int, len(periodicTable))
	for _, element := range periodicTable {
		lookup[element.Symbol] = element.Z
	}
	return lookup
}()

var elementNameBySymbol = func() map[string]string {
	lookup := make(map[string]string, len(periodicTable))
	for _, element := range periodicTable {
		lookup[element.Symbol] = element.Name
	}
	return lookup
}()

var symbolByAtomicNumber = func() map[int]string {
	lookup := make(map[int]string, len(periodicTable))
	for _, element := range periodicTable {
		lookup[element.Z] = element.Symbol
	}
	return lookup
}()

func Elements() []Element {
	elements := make([]Element, len(periodicTable))
	copy(elements, periodicTable)
	return elements
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
	symbol, ok := symbolByAtomicNumber[z]
	return symbol, ok
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
