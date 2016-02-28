package main

import "testing"

func TestNewBeer(t *testing.T) {
	beer, _ := NewBeer("123~01/02/16")

	if beer.id != 123 {
		t.Errorf("Beer id %d is not 123\n", beer.id)
	}

	if beer.drinkDate != "01/02/16" {
		t.Errorf("Date %q is not 01/02/16\n", beer.drinkDate)
	}
}

func TestBadDate(t *testing.T) {
	_, err := NewBeer("123~01/15/16")

	if err == nil {
		t.Errorf("Beer with bad date has been parsed correctly\n")
	}
}

func TestBeerAfter(t *testing.T) {
	beer1, _ := NewBeer("123~01/04/16")
	beer2, _ := NewBeer("123~01/03/16")

	if beer1.IsAfter(beer2) {
		t.Errorf("%v is described as being after %v\n", beer1, beer2)
	}

	if !beer2.IsAfter(beer1) {
		t.Errorf("%v is described as not being after %v\n", beer2, beer1)
	}

	if beer1.IsAfter(beer1) {
		t.Errorf("Beer is described as being after itself\n")
	}
}
