package cache

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"sync"

	quoteapi "github.com/jamoowen/quoteapi/internal"
	"github.com/jamoowen/quoteapi/internal/csv"
	"github.com/jamoowen/quoteapi/internal/utils"
)

// GetRandomQuote() (Quote, error)
// GetQuotesForAuthor(author string) ([]Quote, error)
// InsertQuote(Quote) error

type QuoteCache struct {
	mu     sync.RWMutex
	Quotes []quoteapi.Quote
}

func (c *QuoteCache) GetRandomQuote() (quoteapi.Quote, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.Quotes) == 0 {
		return quoteapi.Quote{}, fmt.Errorf("no quotes available")
	}
	randomIndex := rand.Intn(len(c.Quotes))
	return c.Quotes[randomIndex], nil
}

func (c *QuoteCache) GetQuotesForAuthor(author string) ([]quoteapi.Quote, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	quotes := make([]quoteapi.Quote, 0, 10)
	for _, quote := range c.Quotes {
		if utils.LooselyCompareTwoStrings(quote.Author, author) {
			quotes = append(quotes, quote)
		}
	}
	return quotes, nil
}

func (c *QuoteCache) AddNewQuote(newQuote quoteapi.Quote) error {
	// 1 -> lock the array
	c.mu.Lock()
	defer c.mu.Unlock()
	// find the place to insert (retaining ordering)
	// iterate until
	for i, quote := range c.Quotes {
		// if its an author we havent seen before insert alphabetically
		newAuthor := strings.ToLower(newQuote.Author)
		currAuthor := strings.ToLower(quote.Author)
		if newAuthor < currAuthor {
			c.Quotes = slices.Insert(c.Quotes, i, newQuote)
			break
		}
		// if we have seen the author insert only if the Message is ordered
		if newAuthor == currAuthor {
			// ignore if a duplicate
			if newQuote.Message == quote.Message {
				break
			}
			c.Quotes = slices.Insert(c.Quotes, i, newQuote)
			break
		}
		// append if its the last ordered
		if i == len(c.Quotes)-1 {
			c.Quotes = append(c.Quotes, newQuote)
		}
	}
	// now we want to overwrite our csv file
	return nil
}

func NewQuoteCache(csvPath string) (*QuoteCache, error) {
	records, err := csv.ReadCsv(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize quote cache: %w", err)
	}
	quotes := make([]quoteapi.Quote, 0, len(records))
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) != 2 {
			fmt.Printf("malformed record: %v", record)
			continue
		}
		author := record[0]
		message := record[1]
		quotes = append(quotes, quoteapi.Quote{Author: author, Message: message})
	}

	slices.SortFunc(quotes, func(a, b quoteapi.Quote) int {
		if n := strings.Compare(strings.ToLower(a.Author), strings.ToLower(b.Author)); n != 0 {
			return n
		}
		return strings.Compare(strings.ToLower(a.Message), strings.ToLower(b.Message))
	})
	return &QuoteCache{Quotes: quotes}, nil
}
