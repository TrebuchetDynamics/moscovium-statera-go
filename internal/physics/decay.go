package physics

import (
	"errors"
	"fmt"
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
	if i.A < i.Z {
		return fmt.Errorf("%s has mass number below atomic number", i.ID())
	}
	if i.CitationLink == "" {
		return fmt.Errorf("%s missing citation link", i.ID())
	}
	return nil
}

func DecayChain(start string, catalog Catalog) ([]Isotope, error) {
	if start == "" {
		return nil, errors.New("start isotope is required")
	}

	visited := make(map[string]bool)
	chain := make([]Isotope, 0, 8)
	if err := traverse(start, catalog, visited, &chain); err != nil {
		return nil, err
	}
	return chain, nil
}

func traverse(id string, catalog Catalog, visited map[string]bool, chain *[]Isotope) error {
	if visited[id] {
		return fmt.Errorf("cycle detected at %s", id)
	}

	isotope, ok := catalog[id]
	if !ok {
		return fmt.Errorf("isotope %s not found in catalog", id)
	}
	if err := isotope.Validate(); err != nil {
		return err
	}

	visited[id] = true
	*chain = append(*chain, isotope)
	if isotope.Daughter == "" {
		return nil
	}
	return traverse(isotope.Daughter, catalog, visited, chain)
}
