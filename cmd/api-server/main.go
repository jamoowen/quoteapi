package main

import (
	"log"

	"github.com/jamoowen/quoteapi/internal/server"
)

func main() {
	//  set up db here
	// create controller with db connection? (pointer to the array)
	// pass controller to the router?
	log.Fatal(server.NewServer())

}
