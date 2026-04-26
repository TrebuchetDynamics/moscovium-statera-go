package main

import (
	"fmt"
	"log"
	"time"

	"github.com/TrebuchetDynamics/moscovium-statera-go/internal/physics"
)

func main() {
	catalog := physics.Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, HalfLife: 170 * time.Millisecond, QAlphaMeV: 10.75, Daughter: "284Nh", CitationLink: "https://www.nndc.bnl.gov/ensnds/288/Mc/adopted.pdf"},
		"284Nh": {Symbol: "Nh", Z: 113, A: 284, Daughter: "280Rg", CitationLink: "data/research.seed.json"},
		"280Rg": {Symbol: "Rg", Z: 111, A: 280, Daughter: "276Mt", CitationLink: "data/research.seed.json"},
		"276Mt": {Symbol: "Mt", Z: 109, A: 276, Daughter: "272Bh", CitationLink: "data/research.seed.json"},
		"272Bh": {Symbol: "Bh", Z: 107, A: 272, Daughter: "268Db", CitationLink: "data/research.seed.json"},
		"268Db": {Symbol: "Db", Z: 105, A: 268, Daughter: "264Lr", CitationLink: "data/research.seed.json"},
		"264Lr": {Symbol: "Lr", Z: 103, A: 264, CitationLink: "data/research.seed.json"},
	}

	chain, err := physics.DecayChain("288Mc", catalog)
	if err != nil {
		log.Fatal(err)
	}

	for index, isotope := range chain {
		if index > 0 {
			fmt.Print(" -> ")
		}
		fmt.Print(isotope.ID())
	}
	fmt.Println()
}
