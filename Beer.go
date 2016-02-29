package main

import "errors"
import "strconv"
import "strings"
import "time"

// Beer Central definition of a beer used in the system.
type Beer struct {
	id        int
	drinkDate string
	size      string
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

	if len(elems) != 3 {
		return Beer{}, errors.New("Line is misspecified: " + line)
	}

	bid, _ := strconv.Atoi(elems[0])
	size := elems[2]

	// Ensure that the date parses correctly
	_, err := time.Parse("02/01/06", elems[1])

	// Ensure that the beer size is set
	if err == nil && size != "bomber" && size != "small" {
		err = errors.New(size + " is not a valid beer size")
	}

	beer := Beer{id: bid, drinkDate: elems[1], size: size}
	return beer, err
}
