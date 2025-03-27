package http

import (
	"fmt"
	"net/http"
	"os"
	"time"

	quoteapi "github.com/jamoowen/quoteapi/internal"
)

type Server struct {
}

func (s *Server) StartServer(quoteService quoteapi.QuoteService) error {
	mux := http.NewServeMux()
	handler := Handler{QuoteService: quoteService}
	mux.Handle("/", &handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("Server starting on: http://localhost%s\n", server.Addr)
	return server.ListenAndServe()
}
