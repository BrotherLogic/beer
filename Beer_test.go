package main

import "testing"

func TestNewBeer(t *testing.T) {
	beer, _ := NewBeer("123~01/02/16")

	if beer.id != 123 {
		t.Errorf("Beer id %d is not 123\n", beer.id)
	}

	if beer.drink_date != "01/02/16" {
		t.Errorf("Date %q is not 01/02/16\n", beer.drink_date)
	}
}

func TestBadDate(t *testing.T) {
	_, err := NewBeer("123~01/15/16")

	if err == nil {
		t.Errorf("Beer with bad date has been parsed correctly\n")
	}
}
