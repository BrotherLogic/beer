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
