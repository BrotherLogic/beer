package main

import "log"
import "math"
import "testing"

type LineCounterPrinter struct {
	lcount int
}

func (lineCounter *LineCounterPrinter) Println(output string) {
	lineCounter.lcount = lineCounter.lcount + 1
}

func (lineCounter LineCounterPrinter) GetCount() int {
	return lineCounter.lcount
}

func TestCellarFree(t *testing.T) {
	mine1 := NewCellar("freetest1")
	mine2 := NewCellar("freetest2")
	mine3 := NewCellar("freetest3")

	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")
	beer3, _ := NewBeer("1235~01/03/16~small")

	mine1.AddBeer(beer1)
	mine1.AddBeer(beer2)
	mine2.AddBeer(beer3)

	large1, small1 := mine1.GetFreeSlots()
	if large1 != 20-2 || small1 != 0 {
		t.Errorf("Problem computing free slots of bomber cellar: %v, %v, %v", large1, small1, mine1)
	}
	large2, small2 := mine2.GetFreeSlots()
	if large2 != 0 || small2 != 30-1 {
		t.Errorf("Problem computing free slots of bomber cellar: %v, %v, %v", large2, small2, mine2)
	}
	large3, small3 := mine3.GetFreeSlots()
	if large3 != 20 || small3 != 30 {
		t.Errorf("Problem computing free slots of bomber cellar: %v, %v, %v", large3, small3, mine3)
	}
}

func TestCellarDiffSimple(t *testing.T) {
	mine1 := NewCellar("difftest1")
	mine2 := NewCellar("difftest2")

	cacheBeer(1234, "Tester1")
	cacheBeer(1235, "Tester2")

	beer1, _ := NewBeer("1234~01/01/16~small")
	beer2, _ := NewBeer("1234~01/01/16~small")

	mine1.AddBeer(beer1)
	mine2.AddBeer(beer2)

	diffs := mine1.Diff(mine2)

	if len(diffs) != 0 {
		t.Errorf("Diff error : %v vs %v gives %v", mine1, mine2, diffs)
	}
}

func TestCellarDiff(t *testing.T) {
	mine1 := NewCellar("difftest1")
	mine2 := NewCellar("difftest2")

	cacheBeer(1234, "Tester1")
	cacheBeer(1235, "Tester2")

	beer1, _ := NewBeer("1234~01/01/16~small")
	beer2, _ := NewBeer("1234~01/02/16~small")
	beer3, _ := NewBeer("1235~01/03/16~small")

	mine1.AddBeer(beer1)
	mine1.AddBeer(beer2)
	mine2.AddBeer(beer1)
	mine2.AddBeer(beer3)

	diffs := mine1.Diff(mine2)

	if len(diffs) != 2 {
		t.Errorf("Not enough lines in the diff")
	} else {
		if diffs[0] != "+ "+beer2.Name() {
			t.Errorf("Second line of diff is wrong (should be +) %v", diffs)
		}

		if diffs[1] != "- "+beer3.Name() {
			t.Errorf("Third line of diff is wrong (should be -) %v", diffs)
		}
	}

	rdiffs := mine2.Diff(mine1)

	if len(rdiffs) != 2 {
		t.Errorf("Not enough lines in the diff")
	} else {
		if rdiffs[0] != "- "+beer2.Name() {
			t.Errorf("Second line of diff is wrong (should be +) %v", rdiffs)
		}

		if rdiffs[1] != "+ "+beer3.Name() {
			t.Errorf("Third line of diff is wrong (should be -) %v", rdiffs)
		}
	}

}

func TestRemoveBeer(t *testing.T) {
	mine := NewCellar("testremovecostcellar")
	beer1, _ := NewBeer("1234~01/01/16~small")
	beer2, _ := NewBeer("1235~01/02/16~small")
	beer3, _ := NewBeer("1236~01/03/16~small")
	mine.AddBeer(beer1)
	mine.AddBeer(beer2)
	mine.AddBeer(beer3)

	mine.Remove(1235)

	if len(mine.contents) != 2 {
		t.Errorf("Beer has not been removed: %v\n", mine)
	}
}

func TestRemoveCost(t *testing.T) {
	mine := NewCellar("testremovecostcellar")
	beer1, _ := NewBeer("1234~01/01/16~small")
	beer2, _ := NewBeer("1235~01/02/16~small")
	beer3, _ := NewBeer("1236~01/03/16~small")
	mine.AddBeer(beer1)
	mine.AddBeer(beer2)
	mine.AddBeer(beer3)

	if mine.GetRemoveCost(1237) >= 0 {
		t.Errorf("Removing non cellared beer is not less than zero: %v\n", mine.GetRemoveCost(1237))
	}

	if mine.GetRemoveCost(1234) != 0 {
		t.Errorf("Remove cost is wrong (0): %v\n", mine.GetRemoveCost(1234))
	}
	if mine.GetRemoveCost(1235) != 1 {
		t.Errorf("Remove cost is wrong (1): %v\n", mine.GetRemoveCost(1234))
	}
	if mine.GetRemoveCost(1236) != 2 {
		t.Errorf("Remove cost is wrong (2): %v\n", mine.GetRemoveCost(1234))
	}

}

func TestMidInsert(t *testing.T) {
	mine := NewCellar("testinsertcellar")
	beer1, _ := NewBeer("1234~01/01/16~small")
	beer2, _ := NewBeer("1235~01/02/16~small")
	beer3, _ := NewBeer("1236~01/03/16~small")
	beer4, _ := NewBeer("1237~01/04/16~small")

	mine.AddBeer(beer1)
	mine.AddBeer(beer2)
	mine.AddBeer(beer4)

	log.Printf("%v\n", mine)

	mine.AddBeer(beer3)

	log.Printf("%v\n", mine)

	if mine.CountBeersInCellar(1237) != 1 {
		t.Errorf("Problem with cellar: %v\n", mine)
	}
}

func TestProblemsOfClobbering(t *testing.T) {
	mine1 := NewCellar("testprobcellar")
	beer1, _ := NewBeer("938229~01/03/18~small")
	mine1.AddBeer(beer1)
	log.Printf("%v\n", mine1)
	if mine1.CountBeersInCellar(938229) != 1 {
		t.Errorf("Problem with cellar 1. (%v) %v\n", mine1.CountBeersInCellar(938229), mine1)
	}
	beer2, _ := NewBeer("938229~01/03/17~small")
	mine1.AddBeer(beer2)
	log.Printf("%v\n", mine1)
	if mine1.CountBeersInCellar(938229) != 2 {
		t.Errorf("Problem with cellar 2. (%v) %v\n", mine1.CountBeersInCellar(938229), mine1)
	}
	beer3, _ := NewBeer("938229~01/03/16~small")
	mine1.AddBeer(beer3)
	log.Printf("%v\n", mine1)
	if mine1.CountBeersInCellar(938229) != 3 {
		t.Errorf("Problem with cellar 3. (%v) %v\n", mine1.CountBeersInCellar(938229), mine1)
	}
	beer4, _ := NewBeer("768356~01/09/18~small")
	mine1.AddBeer(beer4)
	log.Printf("%v\n", mine1)
	beer5, _ := NewBeer("768356~01/09/18~small")
	mine1.AddBeer(beer5)
	log.Printf("%v\n", mine1)
	beer6, _ := NewBeer("768356~01/09/17~small")
	mine1.AddBeer(beer6)
	log.Printf("%v\n", mine1)
	beer7, _ := NewBeer("938229~01/03/17~small")
	mine1.AddBeer(beer7)
	log.Printf("%v\n", mine1)
	if mine1.CountBeersInCellar(938229) != 4 {
		t.Errorf("Problem with cellar 4. (%v) %v\n", mine1.CountBeersInCellar(938229), mine1)
	}
	beer8, _ := NewBeer("552346~01/03/17~small")
	mine1.AddBeer(beer8)
	log.Printf("%v\n", mine1)

	if mine1.Size() != 8 {
		t.Errorf("Not the right number of beers, 8 but %v\n", mine1.Size())
	}
}

func TestMergeCellar(t *testing.T) {
	cellar1 := NewCellar("cellar1")
	cellar2 := NewCellar("cellar2")
	cellar3 := NewCellar("cellar3")

	beer1, _ := NewBeer("1234~10/01/15~bomber")
	cellar1.AddBeer(beer1)
	beer2, _ := NewBeer("1235~01/01/15~bomber")
	cellar2.AddBeer(beer2)
	beer3, _ := NewBeer("1236~05/01/15~bomber")
	cellar3.AddBeer(beer3)

	merged := MergeCellars("bomber", cellar1, cellar2, cellar3)
	merged2 := MergeCellars("small", cellar1, cellar2, cellar3)

	if merged[0].id != 1235 {
		t.Errorf("Merged list is ordered incorrectly %v\n", merged[0])
	}

	if len(merged2) != 0 {
		t.Errorf("Merged small list is non-empty: %v\n", merged2)
	}
}

func TestPrint(t *testing.T) {
	printer := StdOutPrint{}

	//Line below should not fail
	printer.Println("Made up")
}

func TestBuildCellar(t *testing.T) {
	cellar := BuildCellar("testdata/simplecellar.cellar")
	if cellar.Size() != 3 {
		t.Errorf("Cellar has %d entries, should have 3\n", cellar.Size())
	}
}

func TestBuildCellarWithBadSize(t *testing.T) {
	cellar := BuildCellar("testdata/badsize.cellar")
	if cellar.Size() != 2 {
		t.Errorf("Cellar has %d entries, should have 2\n", cellar.Size())
	}
}

func TestBuildCellarFileOpenFail(t *testing.T) {
	cellar := BuildCellar("testdata/makdeupfilename.cellar")
	if cellar != nil {
		t.Errorf("Cellar has been built with a non-existant file\n")
	}
}

func TestCellarOverflowSmall(t *testing.T) {
	cellar := NewCellar("testing_cellar")

	//Should be able to add 30 small bottles before insert cost is < 0
	for i := 0; i < 30; i++ {
		beer1, _ := NewBeer("1234~01/01/15~small")
		cost := cellar.ComputeInsertCost(beer1)
		if cost < 0 {
			t.Errorf("Inserting %v into %v is too costly\n", beer1, cellar)
		}
		cellar.AddBeer(beer1)
	}

	beer2, _ := NewBeer("1234~01/01/15~small")
	cost := cellar.ComputeInsertCost(beer2)
	if cost >= 0 {
		t.Errorf("Inserting %v into %v is not prohibited (%v)\n", beer2, cellar, cost)
	}
}

func TestCellarOverflowBomber(t *testing.T) {
	cellar := NewCellar("testing_cellar")

	//Should be able to add 30 small bottles before insert cost is < 0
	for i := 0; i < 20; i++ {
		beer1, _ := NewBeer("1234~01/01/15~bomber")
		cost := cellar.ComputeInsertCost(beer1)
		if cost < 0 {
			t.Errorf("Inserting %v into %v is too costly\n", beer1, cellar)
		}
		cellar.AddBeer(beer1)
	}

	beer2, _ := NewBeer("1234~01/01/15~bomber")
	cost := cellar.ComputeInsertCost(beer2)
	if cost >= 0 {
		t.Errorf("Inserting %v into %v is not prohibited (%v)\n", beer2, cellar, cost)
	}
}

func TestComputeMixSizes(t *testing.T) {
	cellar := NewCellar("testing_cellar")
	beer1, err1 := NewBeer("1234~01/01/16~bomber")
	beer2, err2 := NewBeer("1235~01/01/16~small")

	if err1 != nil || err2 != nil {
		t.Errorf("Parse issue %v,%v\n", err1, err2)
	}

	cellar.AddBeer(beer1)

	cost := cellar.ComputeInsertCost(beer2)
	if cost >= 0 {
		t.Errorf("Beer sizes have been mixed: %v\n", cost)
	}
}

func TestComputeEmptyCellarCost(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer, _ := NewBeer("1234~01/01/16")

	ic := cellar.ComputeInsertCost(beer)

	if ic != math.MaxInt16 {
		t.Errorf("Insert cost should be really high here; in fact it's %d\n", ic)
	}
}

func TestComputeBackCost(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")

	cellar.AddBeer(beer1)
	ic := cellar.ComputeInsertCost(beer2)

	if ic != 1 {
		t.Errorf("Insert costs should have been one, in fact it's %d\n", ic)
	}
}

func TestComputeMiddleCost(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")
	beer3, _ := NewBeer("1234~01/03/16~bomber")

	cellar.AddBeer(beer1)
	cellar.AddBeer(beer3)
	ic := cellar.ComputeInsertCost(beer2)

	if ic != 1 {
		t.Errorf("Insert costs should have been one, in fact it's %d\n", ic)
	}
}

func TestCellarInsert(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")
	beer3, _ := NewBeer("1234~01/03/16~bomber")

	cellar.AddBeer(beer1)
	cellar.AddBeer(beer2)
	cellar.AddBeer(beer3)

	beert1 := cellar.GetNext()
	beert2 := cellar.GetNext()
	beert3 := cellar.GetNext()

	if beert1 != beer1 {
		t.Errorf("Beer1 has been inserted in the wrong position\n")
	}
	if beert2 != beer2 {
		t.Errorf("Beer2 has been inserted in the wrong position\n")
	}
	if beert3 != beer3 {
		t.Errorf("Beer3 has been inserted in the wrong position\n")
	}
}

func TestPrintOutCellar(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, _ := NewBeer("1234~01/01/16")
	beer2, _ := NewBeer("1234~01/02/16")
	beer3, _ := NewBeer("1234~01/03/16")

	cellar.AddBeer(beer1)
	cellar.AddBeer(beer2)
	cellar.AddBeer(beer3)

	linecounter := &LineCounterPrinter{}
	cellar.PrintCellar(linecounter)

	if linecounter.GetCount() != 4 {
		t.Errorf("Print cellar has printed the wrong number of lines: %v\n", linecounter.GetCount())
	}
}

func TestSaveCellar(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, err1 := NewBeer("1234~01/01/16~bomber")
	beer2, err2 := NewBeer("1234~01/02/16~bomber")
	beer3, err3 := NewBeer("1234~01/03/16~bomber")

	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("Problems loading beers: %v, %v, %v\n", err1, err2, err3)
	}

	cellar.AddBeer(beer1)
	cellar.AddBeer(beer2)
	cellar.AddBeer(beer3)

	log.Printf("Pre save size: %v\n", cellar.Size())

	cellar.Save()

	cellar2 := BuildCellar("test_cellar")
	if cellar2.Size() != 3 {
		t.Errorf("Reloading cellar is not the right size: %v\n", cellar2.Size())
	}
}

func TestSaveBadCellar(t *testing.T) {
	cellar := NewCellar("madeupdirectory/blah/blah/cellar")
	cellar.Save()
}
