package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/beerserver/proto"
	pbdi "github.com/brotherlogic/discovery/proto"
)

func getIP(servername string) (string, int) {
	conn, _ := grpc.Dial("192.168.86.64:50055", grpc.WithInsecure())
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

	getFlags := flag.NewFlagSet("GetBeer", flag.ExitOnError)
	var size = getFlags.String("string", "bomber", "Size of the beer")

	if len(os.Args) <= 1 {
		fmt.Printf("Commands: get\n")
	} else {
		switch os.Args[1] {
		case "get":
			if err := getFlags.Parse(os.Args[2:]); err == nil {
				ip, port := getIP("beerserver")
				conn, _ := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
				defer conn.Close()

				client := pb.NewBeerCellarServiceClient(conn)
				beer, err := client.GetBeer(context.Background(), &pb.Beer{Size: *size})
				if err != nil {
					log.Fatalf("Error getting beer: %v", err)
				}
				fmt.Printf("%v\n", beer)
			}
		}
	}
}
