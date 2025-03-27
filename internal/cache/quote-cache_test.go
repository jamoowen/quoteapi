package cache

import (
	"os"
	"testing"

	quoteapi "github.com/jamoowen/quoteapi/internal"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading .env file")
	}
}

// Test for QuoteCache initialization
func TestQuoteCacheInitialization(t *testing.T) {
	csvPath := os.Getenv("QUOTE_CSV_PATH")
	if csvPath == "" {
		t.Fatalf("QUOTE_CSV_PATH env var not set")
	}
	quoteCache, err := NewQuoteCache(csvPath)
	if err != nil {
		t.Fatalf("Quote cache initialization failed: %v", err)
	}
	if len(quoteCache.Quotes) < 100 {
		t.Errorf("Expected at least 100 quotes, got %d", len(quoteCache.Quotes))
	}
}

// Test for GetRandomQuote()
func TestGetRandomQuote(t *testing.T) {
	quoteCache := getQuoteCache()
	randomQuote, err := quoteCache.GetRandomQuote()
	if err != nil || randomQuote.Author == "" {
		t.Fatalf("Failed to fetch a random quote")
	}
}

func getQuoteCache() *QuoteCache {
	return &QuoteCache{
		Quotes: []quoteapi.Quote{
			{Author: "Shakespear", Message: "To be or not to be"},
			{Author: "Trump", Message: "Theyre eating the dogs"},
			{Author: "Conor Mcgregor", Message: "Salaam aleykum alaaada"},
			{Author: "Wanderlei Silva", Message: "I wanna now"},
		},
	}
}
