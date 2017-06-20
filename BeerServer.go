package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/keystore/client"
	"google.golang.org/grpc"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	cellar *BeerCellar
}

// DoRegister Registers this server
func (s *Server) DoRegister(server *grpc.Server) {
	//Nothing to register
}

// ReportHealth Determines if the server is healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Mote promotes this server
func (s *Server) Mote(master bool) error {
	return nil
}

//Init builds a server
func Init() Server {
	cellar, _ := LoadOrNewBeerCellar("prod", ".beer")
	s := Server{&goserver.GoServer{}, cellar}
	s.Register = &s
	return s
}

func runSaveUntappd(command string, flags *flag.FlagSet, key string, secret string, cellar *BeerCellar) {
	if command == "untappd" {
		if flags.Parsed() {
			cellar.SetUntappd(key, secret)
			untappdKey = key
			untappdSecret = secret
		} else {
			flags.PrintDefaults()
		}
	} else {
		if cellar.untappdKey == "" {
			//Set from environment variables
			untappdKey = os.Getenv("CLIENTID")
			untappdSecret = os.Getenv("CLIENTSECRET")
		} else {
			untappdKey = cellar.untappdKey
			untappdSecret = cellar.untappdSecret
		}
	}
}

func run() {

	//Turn off logging
	//log.SetFlags(0)
	//log.SetOutput(ioutil.Discard)

	var beerid string
	var drinkDate string
	var size string
	var days string
	var years string
	var count string
	addBeerFlags := flag.NewFlagSet("add", flag.ContinueOnError)
	addBeerFlags.SetOutput(ioutil.Discard)
	addBeerFlags.StringVar(&beerid, "id", "", "ID of the beer")
	addBeerFlags.StringVar(&drinkDate, "date", "", "Date to be drunk by")
	addBeerFlags.StringVar(&size, "size", "", "Size of bottle")
	addBeerFlags.StringVar(&days, "days", "", "Number of separate days")
	addBeerFlags.StringVar(&years, "years", "", "Number of separate years")
	addBeerFlags.StringVar(&count, "count", "", "Number of bottles")

	var numBombers int
	var numSmall int
	listBeerFlags := flag.NewFlagSet("list", flag.ContinueOnError)
	listBeerFlags.SetOutput(ioutil.Discard)
	listBeerFlags.IntVar(&numBombers, "bombers", 2, "Number of bombers to list from the cellar")
	listBeerFlags.IntVar(&numSmall, "small", 4, "Number of small beers to list from the cellar")

	var key string
	var secret string
	saveUntappdFlags := flag.NewFlagSet("untappd", flag.ContinueOnError)
	saveUntappdFlags.SetOutput(ioutil.Discard)
	saveUntappdFlags.StringVar(&key, "key", "", "Key for untappd")
	saveUntappdFlags.StringVar(&secret, "secret", "", "Secret for untappd")

	var search string
	searchFlags := flag.NewFlagSet("search", flag.ContinueOnError)
	searchFlags.SetOutput(ioutil.Discard)
	searchFlags.StringVar(&search, "string", "", "String to search for")

	var removeID int
	removeFlags := flag.NewFlagSet("remove", flag.ContinueOnError)
	removeFlags.SetOutput(ioutil.Discard)
	removeFlags.IntVar(&removeID, "id", 0, "The ID of the beer to be removed")

	cellarName := "prod"
	dirName := "/Users/simon/.beer/"

	addBeerFlags.Parse(os.Args[2:])
	listBeerFlags.Parse(os.Args[2:])
	saveUntappdFlags.Parse(os.Args[2:])
	searchFlags.Parse(os.Args[2:])
	cellar, _ := LoadOrNewBeerCellar(cellarName, dirName)
	LoadCache(dirName + "prod_cache")

	runSaveUntappd(os.Args[1], saveUntappdFlags, key, secret, cellar)
	runVersion(os.Args[1], cellar)
	runAddBeer(os.Args[1], addBeerFlags, beerid, drinkDate, size, days, years, count, cellar)
	runPrintCellar(os.Args[1], cellar)
	runListBeers(os.Args[1], listBeerFlags, numBombers, numSmall, cellar)
	runSearch(os.Args[1], searchFlags, search)
	runRemoveBeer(os.Args[1], removeFlags, removeID, cellar)
	runSync(os.Args[1], cellar)

	cellar.Save()
	SaveCache("prod_cache")

	//Print the free slots
	largeFree, smallFree := cellar.GetFreeSlots()
	fmt.Printf("\nThere are %v free bomber slots\n", largeFree)
	fmt.Printf("There are %v free small slots\n", smallFree)
}

func main() {
	server := Init()
	server.GoServer.KSclient = *keystoreclient.GetClient()
	server.PrepServer()
	server.RegisterServer("beer", false)
	server.Serve()
}
