package main

import "strconv"
import "strings"
import "time"

type Beer struct {
	id         int
	drink_date string
}

func NewBeer(line string) (Beer, error) {
	elems := strings.Split(line, "~")
	bid, _ := strconv.Atoi(elems[0])

	// Ensure that the date parses correctly
	_, err := time.Parse("02/01/06", elems[1])

	beer := Beer{id: bid, drink_date: elems[1]}
	return beer, err
}
