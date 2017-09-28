package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/beerserver/proto"
	pbdi "github.com/brotherlogic/discovery/proto"
	"github.com/brotherlogic/goserver/utils"
)

func getIP(servername string) (string, int) {
	conn, _ := grpc.Dial(utils.RegistryIP+":"+strconv.Itoa(utils.RegistryPort), grpc.WithInsecure())
	defer conn.Close()

	registry := pbdi.NewDiscoveryServiceClient(conn)
	entry := pbdi.RegistryEntry{Name: servername}
	r, err := registry.Discover(context.Background(), &entry)
	if err != nil {
		return "", -1
	}
	return r.Ip, int(r.Port)
}

func main() {
	addFlags := flag.NewFlagSet("AddBeer", flag.ExitOnError)
	var addSize = addFlags.String("size", "bomber", "Size of the beer")
	var id = addFlags.Int64("id", -1, "Id of the beer")
	var date = addFlags.String("date", "", "Date of drinking")

	getFlags := flag.NewFlagSet("GetBeer", flag.ExitOnError)
	var size = getFlags.String("size", "bomber", "Size of the beer")

	getCellarFlags := flag.NewFlagSet("GetCellar", flag.ExitOnError)
	var cellar = getCellarFlags.Int("cellar", 1, "The number of the cellar")

	removeFlags := flag.NewFlagSet("Remove", flag.ExitOnError)
	var rID = removeFlags.Int64("id", -1, "The id to be removed")

	ip, port := getIP("beerserver")
	conn, _ := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()

	client := pb.NewBeerCellarServiceClient(conn)

	if len(os.Args) <= 1 {
		fmt.Printf("Commands: get\n")
	} else {
		switch os.Args[1] {
		case "get":
			if err := getFlags.Parse(os.Args[2:]); err == nil {
				beer, err := client.GetBeer(context.Background(), &pb.Beer{Size: *size})
				if err == nil {
					fmt.Printf("%v\n", beer)
				}
				log.Printf("%v", err)
			}
		case "add":
			if err := addFlags.Parse(os.Args[2:]); err == nil {

				ts, err := time.Parse("02/01/06", *date)
				if err != nil {
					log.Fatalf("Error parsing date: %v", err)
				}
				_, err = client.AddBeer(context.Background(), &pb.Beer{Size: *addSize, Id: *id, DrinkDate: ts.Unix()})
				if err != nil {
					log.Fatalf("Error adding beer: %v", err)
				}
			}
		case "cellar":
			if err := getCellarFlags.Parse(os.Args[2:]); err == nil {
				cellar, err := client.GetCellar(context.Background(), &pb.Cellar{Name: "cellar" + strconv.Itoa(*cellar)})
				if err != nil {
					log.Fatalf("Error getting cellar: %v", err)
				}
				for i, beer := range cellar.Beers {
					fmt.Printf("%v. %v\n", i+1, beer)
				}
			}
		case "remove":
			if err := removeFlags.Parse(os.Args[2:]); err == nil {
				beer, err := client.RemoveBeer(context.Background(), &pb.Beer{Id: *rID})
				if err != nil {
					log.Fatalf("Error removing beer: %v", err)
				}
				fmt.Printf("REMOVED %v\n", beer)
			}
		case "drunk":
			list, err := client.GetDrunk(context.Background(), &pb.Empty{})
			if err != nil {
				log.Fatalf("Error getting drunk beer: %v", err)
			}

			for i, beer := range list.Beers {
				fmt.Printf("%v. %v\n", i+1, beer)
			}
		}
	}
}
