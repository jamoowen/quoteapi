package main

import (
	"log"
	"os"

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

	cache, err := cache.NewQuoteCache(os.Getenv("QUOTE_CSV_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{}
	log.Fatal(server.StartServer(cache))

}
