package main

import "strconv"
import "strings"

type Beer struct {
	id         int
	drink_date string
}

func NewBeer(line string) Beer {
	elems := strings.Split(line, "~")
	bid, _ := strconv.Atoi(elems[0])
	beer := Beer{id: bid, drink_date: elems[1]}
	return beer
}
