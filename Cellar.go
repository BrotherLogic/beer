package main

import "bufio"
import "fmt"
import "log"
import "math"
import "os"
import "strconv"

// Printer prints to various places
type Printer interface {
	// Println prints a line
	Println(string)
}

// StdOutPrint a printer which prints to standard out
type StdOutPrint struct{}

// Println prints a line to standard out
func (stdprinter *StdOutPrint) Println(output string) {
	fmt.Printf("%v\n", output)
}

// Cellar a single cellar box
type Cellar struct {
	name     string
	contents []Beer
}

// NewCellar builds a new cellar
func NewCellar(cname string) Cellar {
	return Cellar{name: cname, contents: make([]Beer, 0)}
}

// Save saves the cellar file
func (cellar *Cellar) Save() {
	f, err := os.Create(cellar.name)

	if err != nil {
		log.Printf("Error opening file %v\n", err)
	}

	defer f.Close()

	for _, v := range cellar.contents {
		fmt.Fprintf(f, "%v~%v~%v\n", v.id, v.drinkDate, v.size)
	}
}

// PrintCellar prints the contents of the cellar using the Printer
func (cellar *Cellar) PrintCellar(out Printer) {
	out.Println(cellar.name)

	for _, v := range cellar.contents {
		out.Println("BeerName " + strconv.Itoa(v.id))
	}
}

// GetNext Removes the next beer from the cellar
func (cellar *Cellar) GetNext() Beer {
	beer := cellar.contents[0]
	cellar.contents = cellar.contents[1:]
	return beer
}

func (cellar *Cellar) getInsertPoint(beer Beer) int {
	insertPoint := -1
	for i := 0; i < len(cellar.contents); i++ {
		if beer.IsAfter(cellar.contents[i]) {
			insertPoint = i
			break
		}
	}

	if insertPoint == -1 {
		insertPoint = len(cellar.contents)
	}
	return insertPoint
}

// ComputeInsertCost Determines the cost of inserting beer into the cellar
func (cellar *Cellar) ComputeInsertCost(beer Beer) int {
	//Insert cost of an empty cellar should be high
	if len(cellar.contents) == 0 {
		return int(math.MaxInt16)
	}

	//Don't mix sizes
	if cellar.contents[0].size != beer.size {
		return -1
	}

	// Ensure that cellars don't overflow
	if cellar.contents[0].size == "small" && len(cellar.contents) >= 30 {
		return -1
	} else if cellar.contents[0].size == "bomber" && len(cellar.contents) >= 20 {
		return -1
	}

	insertPoint := cellar.getInsertPoint(beer)

	return insertPoint
}

// AddBeer adds a beer to the cellar
func (cellar *Cellar) AddBeer(beer Beer) {
	insertPoint := cellar.getInsertPoint(beer)
	before := cellar.contents[:insertPoint]
	after := cellar.contents[insertPoint:]
	cellar.contents = append(before, beer)
	cellar.contents = append(cellar.contents, after...)
}

// Size Determines the size of the cellar
func (cellar Cellar) Size() int {
	return len(cellar.contents)
}

// BuildCellar Constructs the cellar for the given fileName
func BuildCellar(fileName string) *Cellar {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Error opening %q\n", fileName)
		return nil
	}

	defer file.Close()

	cellar := NewCellar(fileName)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		beer, err := NewBeer(line)

		if err == nil {
			cellar.AddBeer(beer)
		} else {
			log.Printf("Unable to parse beer: %v -> %v\n", line, err)
		}
	}

	return &cellar
}
