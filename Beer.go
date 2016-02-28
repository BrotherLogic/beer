package main

import "strconv"
import "strings"
import "time"

// Beer Central definition of a beer used in the system.
type Beer struct {
	id        int
	drinkDate string
}

// IsAfter Should beer1 be drunk after beer2
func (beer1 Beer) IsAfter(beer2 Beer) bool {
	time1, _ := time.Parse("02/01/06", beer1.drinkDate)
	time2, _ := time.Parse("02/01/06", beer2.drinkDate)

	return time1.Before(time2)
}

// NewBeer Builds a new beer
func NewBeer(line string) (Beer, error) {
	elems := strings.Split(line, "~")
	bid, _ := strconv.Atoi(elems[0])

	// Ensure that the date parses correctly
	_, err := time.Parse("02/01/06", elems[1])

	beer := Beer{id: bid, drinkDate: elems[1]}
	return beer, err
}
