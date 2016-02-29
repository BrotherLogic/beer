package main

import "flag"
import "fmt"
import "strconv"

// BeerCellar the overall beer cellar
type BeerCellar struct {
	version string
	name    string
	bcellar []Cellar
}

// PrintCellar prints out the cellar
func (cellar BeerCellar) PrintCellar(printer Printer) {
	for i, v := range cellar.bcellar {
		if i > 0 {
			printer.Println("--------------")
		}
		v.PrintCellar(printer)
	}
}

// GetNumberOfCellars gets the number of cellar boxes in the cellar
func (cellar BeerCellar) GetNumberOfCellars() int {
	return len(cellar.bcellar)
}

// Save stores a cellar to disk
func (cellar BeerCellar) Save() {
	for _, v := range cellar.bcellar {
		v.Save()
	}
}

// Size gets the size of the cellar
func (cellar BeerCellar) Size() int {
	size := 0
	for _, v := range cellar.bcellar {
		size += v.Size()
	}
	return size
}

// AddBeerToCellar Adds a beer to the cellar
func (cellar BeerCellar) AddBeerToCellar(beer Beer) Cellar {
	bestCellar := -1
	bestScore := -1

	for i, v := range cellar.bcellar {
		insertCount := v.ComputeInsertCost(beer)

		if insertCount > 0 && (insertCount < bestScore || bestScore < 0) {
			bestScore = insertCount
			bestCellar = i
		}
	}

	cellar.bcellar[bestCellar].AddBeer(beer)

	return cellar.bcellar[bestCellar]
}

// GetEmptyCellarCount Gets the number of empty cellars
func (cellar BeerCellar) GetEmptyCellarCount() int {
	count := 0
	for _, v := range cellar.bcellar {
		if v.Size() == 0 {
			count++
		}
	}
	return count
}

// LoadBeerCellar loads a set of beer cellar files
func LoadBeerCellar(name string) *BeerCellar {

	bc := BeerCellar{
		version: "0.1",
		name:    name,
		bcellar: make([]Cellar, 0),
	}

	for i := 1; i < 9; i++ {
		bc.bcellar = append(bc.bcellar, *BuildCellar(name + strconv.Itoa(i) + ".cellar"))
	}

	return &bc
}

// NewBeerCellar creates new beer cellar
func NewBeerCellar(name string) *BeerCellar {
	bc := BeerCellar{
		version: "0.1",
		name:    name,
		bcellar: make([]Cellar, 0),
	}

	for i := 1; i < 9; i++ {
		bc.bcellar = append(bc.bcellar, NewCellar(name+strconv.Itoa(i)+".cellar"))
	}

	return &bc
}

// GetVersion gets the version of the cellar code
func (cellar *BeerCellar) GetVersion() string {
	return cellar.version
}

// AddBeer adds the beer to the cellar
func (cellar *BeerCellar) AddBeer(id string, date string, size string) *Cellar {
	idNum, _ := strconv.Atoi(id)
	if idNum >= 0 {
		cellarBox := cellar.AddBeerToCellar(Beer{id: idNum, drinkDate: date, size: size})
		return &cellarBox
	}

	return nil
}

func runVersion(version bool, cellar *BeerCellar) {
	if version {
		fmt.Printf("BeerCellar: %q\n", cellar.GetVersion())
	}
}

func runAddBeer(addBeer bool, id string, date string, size string, cellar *BeerCellar) {
	if addBeer {
		box := cellar.AddBeer(id, date, size)
		print := &StdOutPrint{}
		box.PrintCellar(print)
	}
}

func main() {
	var version bool
	flag.BoolVar(&version, "version", false, "Prints version")

	var addBeer bool
	var beerid string
	var drinkDate string
	var size string
	flag.BoolVar(&addBeer, "add", false, "Adds a beer")
	flag.StringVar(&beerid, "id", "-1", "ID of the beer")
	flag.StringVar(&drinkDate, "date", "", "Date to be drunk by")
	flag.StringVar(&size, "size", "", "Size of bottle")

	cellarName := "prod"

	flag.Parse()
	cellar := NewBeerCellar(cellarName)
	runVersion(version, cellar)
	runAddBeer(addBeer, beerid, drinkDate, size, cellar)
}
