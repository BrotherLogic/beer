package main

import "flag"
import "fmt"

type BeerCellar struct {
	version string
}

func NewBeerCellar() *BeerCellar {
	bc := BeerCellar{
		version: "0.1",
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
