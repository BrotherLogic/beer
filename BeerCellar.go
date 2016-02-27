package main

import "flag"
import "fmt"
import "strconv"

type BeerCellar struct {
	version string
	bcellar []Cellar
}

func (cellar BeerCellar) GetNumberOfCellars() int {
	return len(cellar.bcellar)
}

func (cellar BeerCellar) AddBeerToCellar(beer Beer) Cellar {
	best_cellar := -1
	best_score := -1

	for i, v := range cellar.bcellar {
		insert_count := v.ComputeInsertCost(beer)

		if insert_count > 0 && (insert_count < best_score || best_score < 0) {
			best_score = insert_count
			best_cellar = i
		}
	}

	cellar.bcellar[best_cellar].AddBeer(beer)

	return cellar.bcellar[best_cellar]
}

func (cellar BeerCellar) GetEmptyCellarCount() int {
	count := 0
	for _, v := range cellar.bcellar {
		if v.Size() == 0 {
			count++
		}
	}
	return count
}

func NewBeerCellar() *BeerCellar {
	bc := BeerCellar{
		version: "0.1",
		bcellar: make([]Cellar, 0),
	}

	for i := 1; i < 9; i++ {
		bc.bcellar = append(bc.bcellar, NewCellar("cellar"+strconv.Itoa(i)))
	}

	return &bc
}

func (bc *BeerCellar) GetVersion() string {
	return bc.version
}

func RunVersion(version bool, cellar *BeerCellar) {
	if version {
		fmt.Printf("BeerCellar: %q\n", cellar.GetVersion())
	}
}

func main() {
	fmt.Printf("HERE\n")

	var version bool
	flag.BoolVar(&version, "version", false, "Prints version")
	flag.Parse()

	cellar := NewBeerCellar()
	RunVersion(version, cellar)
}
