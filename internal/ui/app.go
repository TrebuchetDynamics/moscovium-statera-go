package ui

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

func DefaultSpec() Spec {
	return Spec{
		Title:  "Moscovium Statera Go",
		Width:  1100,
		Height: 720,
		Modules: []ModuleSpec{
			{
				Name:        "Education",
				Description: "Nuclide chart and source-verification sandbox for evaluated Moscovium data.",
			},
			{
				Name:        "Research",
				Description: "Citation browser for DOI-backed nuclear data and evaluated records.",
			},
			{
				Name:        "Design",
				Description: "Theoretical isotope-lattice modeling with explicit uncertainty boundaries.",
			},
		},
	}
}
