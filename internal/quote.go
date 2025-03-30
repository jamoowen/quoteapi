package quoteapi

type Quote struct {
	Author  string `json:"author"`
	Message string `json:"message"`
}

type QuoteService interface {
	GetRandomQuote() (Quote, error)
	GetQuotesForAuthor(author string) ([]Quote, error)
	// InsertQuote(Quote) error
}
