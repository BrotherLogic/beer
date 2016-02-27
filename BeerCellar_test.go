package main

import "testing"

func TestGetVersion(t *testing.T) {
	mine := NewBeerCellar()
	version := mine.GetVersion()
	if len(version) == 0 {
		t.Errorf("Version %q is not legit\n", version)
	}
}

func TestMain(t *testing.T) {
	main()
}

func TestVersion(t *testing.T) {
	RunVersion(true, NewBeerCellar())
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
