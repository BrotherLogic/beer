package main

import "testing"

func TestGetVersion(t *testing.T) {
	mine := NewBeerCellar()
	version := mine.GetVersion()
	if len(version) == 0 {
		t.Errorf("Version %q is not legit\n", version)
	}
}

func TestAddBeer(t *testing.T) {
	mine := NewBeerCellar()
	mine.AddBeer("1234", "01/01/16")
	mine.Save()

	mine2 := LoadBeerCellar()
	if mine2.Size() != mine.Size() && mine2.Size() == 1 {
		t.Errorf("Size on reload is incorrect %v vs %v\n", mine.Size(), mine2.Size())
	}
}

func TestAddNoBeer(t *testing.T) {
	mine := NewBeerCellar()
	mine.AddBeer("-1", "01/01/16")
	mine.Save()

	mine2 := LoadBeerCellar()
	if mine2.Size() != 0 {
		t.Errorf("Error on adding no beer - %v but %v\n", mine.Size(), mine2.Size())
	}
}

func TestMain(t *testing.T) {
	main()
}

func TestRunVersion(t *testing.T) {
	RunVersion(true, NewBeerCellar())
}

func TestRunAddBeer(t *testing.T) {
	RunAddBeer(true, "1234", "01/02/16", NewBeerCellar())
}

func TestGetNumberOfCellars(t *testing.T) {
	bc := NewBeerCellar()
	if bc.GetNumberOfCellars() != 8 {
		t.Errorf("Wrong number of cellars: %d\n", bc.GetNumberOfCellars())
	}
}

func TestAddToCellar(t *testing.T) {
	cellar := NewBeerCellar()
	beer1, _ := NewBeer("1234~01/01/16")
	beer2, _ := NewBeer("1234~01/02/16")
	beer3, _ := NewBeer("1234~01/03/16")

	cellar.AddBeerToCellar(beer1)
	cellar.AddBeerToCellar(beer2)
	cellar.AddBeerToCellar(beer3)

	if cellar.GetEmptyCellarCount() != 7 {
		t.Errorf("Too many cellars are not empty %d\n", cellar.GetEmptyCellarCount())
	}
}

func TestPrintOutBeerCellar(t *testing.T) {
	cellar := NewBeerCellar()
	beer1, _ := NewBeer("1234~01/01/16")
	beer2, _ := NewBeer("1234~01/02/16")
	beer3, _ := NewBeer("1234~01/03/16")

	cellar.AddBeerToCellar(beer1)
	cellar.AddBeerToCellar(beer2)
	cellar.AddBeerToCellar(beer3)

	linecounter := &LineCounterPrinter{}
	cellar.PrintCellar(linecounter)

	if linecounter.GetCount() != 3+8+7 {
		t.Errorf("Print cellar has printed the wrong number of lines: %v\n", linecounter.GetCount())
	}
}

func TestSaveBeerCellar(t *testing.T) {
	cellar := NewBeerCellar()
	beer1, _ := NewBeer("1234~01/01/16")
	beer2, _ := NewBeer("1234~01/02/16")
	beer3, _ := NewBeer("1234~01/03/16")

	cellar.AddBeerToCellar(beer1)
	cellar.AddBeerToCellar(beer2)
	cellar.AddBeerToCellar(beer3)

	cellar.Save()

	cellar2 := LoadBeerCellar()
	if cellar2.Size() != cellar.Size() {
		t.Errorf("Mismatched sizes %v and %v\n", cellar.Size(), cellar2.Size())
	}
}
