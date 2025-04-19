package cache

import (
	"testing"

	quoteapi "github.com/jamoowen/quoteapi/internal"
)

// Test for QuoteCache initialization
func TestQuoteCacheInitialization(t *testing.T) {
	quoteCache, err := NewQuoteCache("../../data/quotes.csv")
	if err != nil {
		t.Fatalf("Quote cache initialization failed: %v", err.Error())
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

func TestInsertQuote(t *testing.T) {
	// Defining the columns of the table
	tests := []struct {
		name           string
		quote          quoteapi.Quote
		want           error
		expectedIndex  int
		expectedLength int
	}{
		// the table itself
		{"Should insert the quote", quoteapi.Quote{Author: "Julius Caesar", Message: "YOLO"}, nil, 1, 5},
		{"Should append the quote", quoteapi.Quote{Author: "zulius Caesar", Message: "YOLO"}, nil, 4, 5},
		{"Should prepend the quote", quoteapi.Quote{Author: "aulius Caesar", Message: "YOLO"}, nil, 0, 5},
		{"Should insert at beginning of authors quotes if same author", quoteapi.Quote{Author: "Trump", Message: "mOLO"}, nil, 2, 5},
		{"Should ignore if a duplicate", quoteapi.Quote{Author: "Trump", Message: "Theyre eating the dogs"}, nil, 2, 4},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quoteCache := getQuoteCache()
			err := quoteCache.AddNewQuote(tt.quote)
			if len(quoteCache.Quotes) != tt.expectedLength {
				t.Errorf("Expected length: %v, got: %v", tt.expectedLength, len(quoteCache.Quotes))
			}

			if err != tt.want {
				t.Errorf("got %t, want %t", err, tt.want)
			}
			if tt.want == nil {
				quoteAtExpectedInsertIndex := quoteCache.Quotes[tt.expectedIndex]
				if quoteAtExpectedInsertIndex.Message != tt.quote.Message {
					t.Errorf("Expected quote (%v) to be inserted at ined %v. Found: %v", tt.quote, tt.expectedIndex, quoteAtExpectedInsertIndex)
				}
			}
		})
	}
}

func getQuoteCache() *QuoteCache {
	return &QuoteCache{
		Quotes: []quoteapi.Quote{
			{Author: "Conor Mcgregor", Message: "Salaam aleykum alaaada"},
			{Author: "Shakespear", Message: "To be or not to be"},
			{Author: "Trump", Message: "Theyre eating the dogs"},
			{Author: "Wanderlei Silva", Message: "I wanna now"},
		},
	}
}
