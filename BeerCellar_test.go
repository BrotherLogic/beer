package main

import "flag"
import "log"
import "testing"

func TestAddToCellars(t *testing.T) {
	mine1 := NewBeerCellar("testaddcellar")
	mine1.AddBeer("1234", "01/01/16", "bomber")
	mine1.AddBeer("1234", "01/01.15", "bomber")

	if mine1.GetEmptyCellarCount() != 7 {
		t.Errorf("Cellar is not balanced: %v\n", mine1)
	}
}

func TestSyncAndSave(t *testing.T) {
     mine := NewBeerCellar("testsync")
     mine.AddBeer("1234", "01/01/16", "bomber")
     mine.AddBeer("791939", "01/02/16", "bomber")
     mine.AddBeer("1234", "01/03/16", "bomber")
     mine.syncTime = "28/02/16"

     mine.Sync(blankVenuePageFetcher{}, venuePageConverter{})
     mine.Save();

     mine2, err := LoadBeerCellar("testsync")
     if err != nil {
     	t.Errorf("Error loading cellar: %v", err)
     }

     if mine2.Size() != 2 {
     	t.Errorf("Sync remove has failed(%v): %v\n", mine2.Size(), mine2)
     }

     if mine2.syncTime == "28/02/16" {
     	t.Errorf("Sync time was not updated")
     }
}

func TestRemoveFromCellar(t *testing.T) {
     mine := NewBeerCellar("testremovefromcellar")
     mine.AddBeer("1234", "01/01/16", "bomber")
     mine.AddBeer("1235", "01/02/16", "bomber")
     mine.AddBeer("1234", "01/03/16", "bomber")

     mine.RemoveBeer(1234)

     if len(mine.bcellar[0].contents) != 2 {
     	t.Errorf("Beer has not been removed: %v\n", mine)
     }
}

func TestSaveAndReload(t *testing.T) {
	mine1, _ := LoadOrNewBeerCellar("cellar1")
	mine1.AddBeer("1234", "01/01/16", "bomber")
	mine1.Save()

	mine2, _ := LoadOrNewBeerCellar("cellar1")
	if mine2.Size() == 0 {
		t.Errorf("Cellar is not being reloaded correctly\n")
	}
}

func TestForClobbering(t *testing.T) {
	mine1 := NewBeerCellar("cellar1")
	mine2 := NewBeerCellar("cellar2")

	mine1.AddBeer("1234", "01/01/16", "bomber")
	mine1.Save()

	t1, _ := LoadBeerCellar("cellar1")
	if t1.Size() != 1 {
		t.Errorf("Cellar is missized on first load: %v\n", t1.Size())
	}

	mine2.Save()
	t2, _ := LoadBeerCellar("cellar1")
	if t2.Size() != 1 {
		t.Errorf("Cellar is missized on second load: %v\n", t2.Size())
	}
}

func TestGetSyncTime(t *testing.T) {
	mine := NewBeerCellar("testing")
	mine.syncTime = "01/01/16"
	mine.Save()

	mine2, _ := LoadBeerCellar("testing")
	if mine2.syncTime != "01/01/16" {
		t.Errorf("BeerCellar is not saving the sync date\n")
	}
}

func TestGetVersion(t *testing.T) {
	mine := NewBeerCellar("test")
	version := mine.GetVersion()
	if len(version) == 0 {
		t.Errorf("Version %q is not legit\n", version)
	}
}

func TestAddBeer(t *testing.T) {
	mine := NewBeerCellar("test")
	mine.AddBeer("1234", "01/01/16", "bomber")
	mine.Save()

	mine2, _ := LoadBeerCellar("test")
	if mine2.Size() != mine.Size() && mine2.Size() == 1 {
		t.Errorf("Size on reload is incorrect %v vs %v\n", mine.Size(), mine2.Size())
	}
}

func TestAddNoBeer(t *testing.T) {
	mine := NewBeerCellar("test")
	mine.AddBeer("-1", "01/01/16", "bomber")
	mine.Save()

	mine2, _ := LoadBeerCellar("test")
	if mine2.Size() != 0 {
		t.Errorf("Error on adding no beer - %v but %v\n", mine.Size(), mine2.Size())
	}
}

func TestMain(t *testing.T) {
	main()
}

func TestRunVersion(t *testing.T) {
	runVersion("version", NewBeerCellar("test"))
}

func TestRunSaveUntappd(t *testing.T) {
	runSaveUntappd("untappd", flag.NewFlagSet("blah", flag.ExitOnError), "testkey", "testsecret", NewBeerCellar("untappdtest"))
}

func TestRunAddBeer(t *testing.T) {
	runAddBeer("add", flag.NewFlagSet("mock", flag.ExitOnError), "1234", "01/02/16", "bomber", NewBeerCellar("test"))
}

func TestRunPrintCellar(t *testing.T) {
	runPrintCellar("print", NewBeerCellar("test"))
}

func TestRunListBeers(t *testing.T) {
	runListBeers("list", flag.NewFlagSet("mock", flag.ExitOnError), 5, 5, NewBeerCellar("test"))
}

func TestMin(t *testing.T) {
	if Min(3, 2) == 3 || Min(2, 3) == 3 {
		t.Errorf("Min is not returning the min\n")
	}
}

func TestGetNumberOfCellars(t *testing.T) {
	bc := NewBeerCellar("test")
	if bc.GetNumberOfCellars() != 8 {
		t.Errorf("Wrong number of cellars: %d\n", bc.GetNumberOfCellars())
	}
}

func TestAddToCellar(t *testing.T) {
	cellar := NewBeerCellar("test")
	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")
	beer3, _ := NewBeer("1234~01/03/16~bomber")

	cellar.AddBeerToCellar(beer1)
	cellar.AddBeerToCellar(beer2)
	cellar.AddBeerToCellar(beer3)

	if cellar.GetEmptyCellarCount() != 7 {
		t.Errorf("Too many cellars are not empty %d\n", cellar.GetEmptyCellarCount())
	}
}

func TestPrintOutBeerCellar(t *testing.T) {
	cellar := NewBeerCellar("test")
	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")
	beer3, _ := NewBeer("1234~01/03/16~bomber")

	cellar.AddBeerToCellar(beer1)
	cellar.AddBeerToCellar(beer2)
	cellar.AddBeerToCellar(beer3)

	linecounter := &LineCounterPrinter{}
	cellar.PrintCellar(linecounter)

	if linecounter.GetCount() != 3+8+7 {
		t.Errorf("Print cellar has printed the wrong number of lines: %v\n", linecounter.GetCount())
	}
}

func TestSaveBeerCellar(t *testing.T) {
	cellar := NewBeerCellar("test")
	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")
	beer3, err := NewBeer("1234~01/03/16~bomber")

	if err != nil {
		t.Errorf("Parse fail %v\n", err)
	}

	cellar.AddBeerToCellar(beer1)
	cellar.AddBeerToCellar(beer2)
	cellar.AddBeerToCellar(beer3)

	cellar.Save()

	cellar2, _ := LoadBeerCellar("test")
	if cellar2.Size() != cellar.Size() {
		t.Errorf("Mismatched sizes %v and %v\n", cellar.Size(), cellar2.Size())
	}
}

func TestSaveBadBeerCellar(t *testing.T) {
	cellar := NewBeerCellar("madeupdirectory/blah")
	cellar.Save()
}

func TestLoadBadBeerCellar(t *testing.T) {
	_, err := LoadBeerCellar("madeupdirectory/blah")
	if err == nil {
		t.Errorf("No Error on opening bad cellar\n")
	}
}

func TestPrintBeers(t *testing.T) {
	log.Printf("Starting Here\n")
	cellar := NewBeerCellar("test")
	beer1, _ := NewBeer("1234~12/05/12~bomber")
	beer2, _ := NewBeer("1235~12/05/12~small")
	cellar.AddBeerToCellar(beer1)
	cellar.AddBeerToCellar(beer2)
	cellar.PrintBeers(5, 5)
}

func TestGetUntappdKeys(t *testing.T) {
	cellar := NewBeerCellar("uttest")
	cellar.SetUntappd("testkey", "testsecret")
	cellar.Save()

	cellar2, _ := LoadOrNewBeerCellar("uttest")

	key, secret := cellar2.GetUntappd()
	if key != "testkey" || secret != "testsecret" {
		t.Errorf("Keys have come back incorrectly: %v and %v\n", key, secret)
	}
}
