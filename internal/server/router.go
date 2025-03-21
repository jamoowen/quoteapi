package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func NewServer() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/quote/random", randomQuoteHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("Server starting on: %s\n", s.Addr)
	return s.ListenAndServe()

}

func randomQuoteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received... ", r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	author := "Thomas Edison"
	message := "Genius is one percent inspiration and ninety-nine percent perspiration."
	response := map[string]string{
		"author":  author,
		"message": message,
	}
	json.NewEncoder(w).Encode(response)
}
