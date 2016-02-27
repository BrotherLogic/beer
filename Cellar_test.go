package main

import "math"
import "testing"

func TestBuildCellar(t *testing.T) {
	cellar := BuildCellar("testdata/simplecellar.cellar")
	if cellar.Size() != 3 {
		t.Errorf("Cellar has %d entries, should have 3\n", cellar.Size())
	}
}

func TestBuildCellarFileOpenFail(t *testing.T) {
	cellar := BuildCellar("testdata/makdeupfilename.cellar")
	if cellar != nil {
		t.Errorf("Cellar has been built with a non-existant file\n")
	}
}

func TestComputeEmptyCellarCost(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer, _ := NewBeer("1234~01/01/16")

	ic := cellar.ComputeInsertCost(beer)

	if ic != math.MaxInt16 {
		t.Errorf("Insert cost should be really high here; in fact it's %d\n", ic)
	}
}

func TestComputeBackCost(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, _ := NewBeer("1234~01/01/16")
	beer2, _ := NewBeer("1234~01/02/16")

	cellar.AddBeer(beer1)
	ic := cellar.ComputeInsertCost(beer2)

	if ic != 1 {
		t.Errorf("Insert costs should have been one, in fact it's %d\n", ic)
	}
}

func TestComputeMiddleCost(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, _ := NewBeer("1234~01/01/16")
	beer2, _ := NewBeer("1234~01/02/16")
	beer3, _ := NewBeer("1234~01/03/16")

	cellar.AddBeer(beer1)
	cellar.AddBeer(beer3)
	ic := cellar.ComputeInsertCost(beer2)

	if ic != 1 {
		t.Errorf("Insert costs should have been one, in fact it's %d\n", ic)
	}
}
