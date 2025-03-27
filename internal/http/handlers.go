package http

import (
	"fmt"
	"net/http"

	quoteapi "github.com/jamoowen/quoteapi/internal"
)

type Handler struct {
	QuoteService quoteapi.QuoteService
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Req received=> url:%v, method:%v \n", r.URL, r.Method)
	w.Write([]byte("hello there "))
}

// func (h *Handler) randomQuoteHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Request received... ", r)
// 	w.Header().Set("Content-Type", "application/json")
// 	randomQuote, err := h.QuoteService.GetRandomQuote()
// 	if err != nil {

// 	}
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(randomQuote)
// }
