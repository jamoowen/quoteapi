package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Quote struct {
	Author  string `json:"author"`
	Message string `json:"message"`
}

func NewServer() *http.Server {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/quote", randomQuoteHandler)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get timeouts from environment or use defaults
	readTimeout := 10 * time.Second
	writeTimeout := 10 * time.Second

	if rt := os.Getenv("READ_TIMEOUT"); rt != "" {
		if duration, err := time.ParseDuration(rt); err == nil {
			readTimeout = duration
		}
	}

	if wt := os.Getenv("WRITE_TIMEOUT"); wt != "" {
		if duration, err := time.ParseDuration(wt); err == nil {
			writeTimeout = duration
		}
	}

	s := &http.Server{
		Addr:           ":" + port, // Added the colon here
		Handler:        mux,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	return s
}

func StartServer() error {
	s := NewServer()
	fmt.Printf("Server starting on %s\n", s.Addr)
	return s.ListenAndServe()
}

func randomQuoteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received... ", r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	author := "Thomas Edison"
	message := "Genius is one percent inspiration and ninety-nine percent perspiration."
	response := Quote{Author: author, Message: message}
	json.NewEncoder(w).Encode(response)
}
