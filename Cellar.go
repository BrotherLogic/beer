package main

import "bufio"
import "fmt"
import "log"
import "math"
import "os"
import "strconv"

type Printer interface {
	Println(string)
}
type StdOutPrint struct{}

func (stdprinter *StdOutPrint) Println(output string) {
	fmt.Printf("%v\n", output)
}

type Cellar struct {
	name     string
	contents []Beer
}

func NewCellar(cname string) Cellar {
	return Cellar{name: cname, contents: make([]Beer, 0)}
}

func (cellar *Cellar) PrintCellar(out Printer) {
	out.Println(cellar.name)

	for _, v := range cellar.contents {
		out.Println("BeerName " + strconv.Itoa(v.id))
	}
}

func (cellar *Cellar) GetNext() Beer {
	beer := cellar.contents[0]
	cellar.contents = cellar.contents[1:]
	return beer
}

func (cellar *Cellar) getInsertPoint(beer Beer) int {
	insert_point := -1
	for i := 0; i < len(cellar.contents); i++ {
		if beer.IsAfter(cellar.contents[i]) {
			insert_point = i
			break
		}
	}

	if insert_point == -1 {
		insert_point = len(cellar.contents)
	}
	return insert_point
}

func (cellar *Cellar) ComputeInsertCost(beer Beer) int {
	//Insert cost of an empty cellar should be high
	if len(cellar.contents) == 0 {
		return int(math.MaxInt16)
	}

	insert_point := cellar.getInsertPoint(beer)

	return insert_point
}

func (cellar *Cellar) AddBeer(beer Beer) {
	insert_point := cellar.getInsertPoint(beer)
	before := cellar.contents[:insert_point]
	after := cellar.contents[insert_point:]
	cellar.contents = append(before, beer)
	cellar.contents = append(cellar.contents, after...)
}

func (cellar Cellar) Size() int {
	return len(cellar.contents)
}

func BuildCellar(file_name string) *Cellar {
	file, err := os.Open(file_name)
	if err != nil {
		log.Printf("Error opening %q\n", file_name)
		return nil
	}

	defer file.Close()

	cellar := NewCellar(file_name)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		beer, err := NewBeer(scanner.Text())

		if err == nil {
			cellar.AddBeer(beer)
		}
	}

	return &cellar
}
