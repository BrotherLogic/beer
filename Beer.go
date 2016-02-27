package main

import "strconv"
import "strings"
import "time"

type Beer struct {
	id         int
	drink_date string
}

func (beer1 Beer) IsAfter(beer2 Beer) bool {
	time1, _ := time.Parse("02/01/06", beer1.drink_date)
	time2, _ := time.Parse("02/01/06", beer2.drink_date)

	return time1.Before(time2)
}

func NewBeer(line string) (Beer, error) {
	elems := strings.Split(line, "~")
	bid, _ := strconv.Atoi(elems[0])

	// Ensure that the date parses correctly
	_, err := time.Parse("02/01/06", elems[1])

	beer := Beer{id: bid, drink_date: elems[1]}
	return beer, err
}
