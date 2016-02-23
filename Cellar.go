package main

import "bufio"
import "log"
import "os"

type Cellar struct {
	name     string
	contents []Beer
}

func NewCellar(cname string) Cellar {
	return Cellar{name: cname, contents: make([]Beer, 0)}
}

func (cellar *Cellar) AddBeer(beer Beer) {
	cellar.contents = append(cellar.contents, beer)
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
		beer := NewBeer(scanner.Text())
		cellar.AddBeer(beer)
	}

	return &cellar
}
