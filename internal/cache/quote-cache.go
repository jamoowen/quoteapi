package cache

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math/rand"
	"os"

	quoteapi "github.com/jamoowen/quoteapi/internal"
	"github.com/jamoowen/quoteapi/internal/utils"
)

// GetRandomQuote() (Quote, error)
// GetQuotesForAuthor(author string) ([]Quote, error)
// InsertQuote(Quote) error

type QuoteCache struct {
	Quotes []quoteapi.Quote
}

func (c *QuoteCache) GetRandomQuote() (quoteapi.Quote, error) {
	if len(c.Quotes) == 0 {
		return quoteapi.Quote{}, errors.New("no quotes available")
	}
	randomIndex := rand.Intn(len(c.Quotes))
	return c.Quotes[randomIndex], nil
}

func (c *QuoteCache) GetQuotesForAuthor(author string) ([]quoteapi.Quote, error) {
	quotes := make([]quoteapi.Quote, 0, 10)
	for _, quote := range c.Quotes {
		if utils.LooselyCompareTwoStrings(quote.Author, author) {
			quotes = append(quotes, quote)
		}
	}
	return quotes, nil
}

func NewQuoteCache(csvPath string) (*QuoteCache, error) {
	if csvPath == "" {
		return nil, errors.New("missing path to csv")
	}
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv: %v", err)
	}
	slice := make([]quoteapi.Quote, 0, len(records))
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) != 2 {
			fmt.Printf("malformed record: %v", record)
			continue
		}
		var author string = record[0]
		var message string = record[1]
		slice = append(slice, quoteapi.Quote{Author: author, Message: message})
	}

	return &QuoteCache{Quotes: slice}, nil
}
