package http

import (
	"encoding/json"
	"io"
	"net/http"

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
	// need to get author out of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		badRequestError(w, "Could not parse request body")
		return
	}
	var newQuote quoteapi.Quote
	err = json.Unmarshal(body, &newQuote)
	if err != nil {
		badRequestError(w, "Malformed Json. Expected author & message object")
		return
	}
	// now we need to
	// a) append to list (cache)
	// b) write to csv
	w.WriteHeader(http.StatusCreated)
}
