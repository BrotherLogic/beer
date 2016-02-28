package main

import "flag"
import "fmt"
import "strconv"

type BeerCellar struct {
	version string
	bcellar []Cellar
}

func (cellar BeerCellar) PrintCellar(printer Printer) {
	for i, v := range cellar.bcellar {
		if i > 0 {
			printer.Println("--------------")
		}
		v.PrintCellar(printer)
	}
}

func (cellar BeerCellar) GetNumberOfCellars() int {
	return len(cellar.bcellar)
}

func (cellar BeerCellar) Save() {
	for _, v := range cellar.bcellar {
		v.Save()
	}
}

func (cellar BeerCellar) Size() int {
	size := 0
	for _, v := range cellar.bcellar {
		size += v.Size()
	}
	return size
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

func LoadBeerCellar() *BeerCellar {

	bc := BeerCellar{
		version: "0.1",
		bcellar: make([]Cellar, 0),
	}

	for i := 1; i < 9; i++ {
		bc.bcellar = append(bc.bcellar, *BuildCellar("cellar" + strconv.Itoa(i)))
	}

	return &bc
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

func (bc *BeerCellar) AddBeer(id string, date string) *Cellar {
	id_num, _ := strconv.Atoi(id)
	if id_num >= 0 {
		cellar := bc.AddBeerToCellar(Beer{id: id_num, drink_date: date})
		return &cellar
	}

	return nil
}

func RunVersion(version bool, cellar *BeerCellar) {
	if version {
		fmt.Printf("BeerCellar: %q\n", cellar.GetVersion())
	}
}

func RunAddBeer(add_beer bool, id string, date string, cellar *BeerCellar) {
	if add_beer {
		box := cellar.AddBeer(id, date)
		print := &StdOutPrint{}
		box.PrintCellar(print)
	}
}

func main() {
	fmt.Printf("HERE\n")

	var version bool
	flag.BoolVar(&version, "version", false, "Prints version")

	var add_beer bool
	var beerid string
	var drink_date string
	flag.BoolVar(&add_beer, "add", false, "Adds a beer")
	flag.StringVar(&beerid, "id", "-1", "ID of the beer")
	flag.StringVar(&drink_date, "date", "", "Date to be drunk by")

	flag.Parse()
	cellar := LoadBeerCellar()
	RunVersion(version, cellar)
	RunAddBeer(add_beer, beerid, drink_date, cellar)
}
