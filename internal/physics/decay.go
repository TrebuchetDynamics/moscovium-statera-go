package physics

import (
	"errors"
	"fmt"
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
