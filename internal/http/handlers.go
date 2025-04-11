package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	quoteapi "github.com/jamoowen/quoteapi/internal"
)

type Handler struct {
	QuoteService quoteapi.QuoteService
}

func (h *Handler) randomQuoteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	randomQuote, err := h.QuoteService.GetRandomQuote()
	if err != nil {
		internalServerError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(randomQuote)
}

// /quote/:quthor
func (h *Handler) getQuotesForAuthorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// need to get author out of the request
	author := r.URL.Query().Get("name")
	if author == "" {
		badRequestError(w, "author name must be provided")

	}
	randomQuote, err := h.QuoteService.GetQuotesForAuthor(author)
	if err != nil {
		internalServerError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	if len(randomQuote) == 0 {
		w.Write([]byte("No quotes found for this homie"))
		return
	}
	json.NewEncoder(w).Encode(randomQuote)
}

// POST /quote/add
func (h *Handler) addQuote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Limit request body size to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		if err == http.ErrBodyReadAfterClose {
			badRequestError(w, "Request body too large")
			return
		}
		badRequestError(w, "Could not parse request body")
		return
	}
	var newQuote quoteapi.Quote
	err = json.Unmarshal(body, &newQuote)
	if err != nil {
		badRequestError(w, "Malformed JSON. Expected author & message object")
		return
	}
	// Validate required fields
	if newQuote.Author == "" || newQuote.Message == "" {
		badRequestError(w, "Author and message are required")
		return
	}
	// Validate message length (100 words)
	words := strings.Fields(newQuote.Message)
	if len(words) > 100 {
		badRequestError(w, "Message cannot exceed 100 words")
		return
	}

	// now we need to
	// a) append to list (cache)
	// b) write to csv
	w.WriteHeader(http.StatusCreated)
}
