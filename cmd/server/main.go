package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/jamoowen/quoteapi/internal/cache"
	"github.com/jamoowen/quoteapi/internal/http"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	//  set up db here
	// create controller with db connection? (pointer to the array)
	// pass controller to the router?

	// load env vars
	// initialize cache
	//
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Startup error: %v", err.Error())
	}
	csvPath := path.Join(wd, "/data/quotes.csv")
	fmt.Println("csvPath: ", csvPath)
	cache, err := cache.NewQuoteCache(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{}
	log.Fatal(server.StartServer(cache))

}
