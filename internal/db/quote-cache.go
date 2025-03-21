package db

import (
	"errors"
	"os"
)

interface Cache {
}

type QuoteCache struct {
	Quotes *[]Quote
}


func initializeCache() {

}

func 

func createQuoteCache() error {

	csvPath := os.Getenv("QUOTE_CSV_PATH")
	if csvPath == "" {
		return errors.New("QUOTE_CSV_PATH not found")
	}

	// read file
	// sort by author
	// tranform into array of Quote objects


}
