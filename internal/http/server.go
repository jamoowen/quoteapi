package http

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	quoteapi "github.com/jamoowen/quoteapi/internal"
	"github.com/jamoowen/quoteapi/internal/auth"
	"github.com/jamoowen/quoteapi/internal/email"
)

type Server struct {
}

func (s *Server) registerRoutes(mux *http.ServeMux, h Handler) {
	mux.HandleFunc("GET /authenticate", h.handleAuthenticate)
	mux.HandleFunc("POST /authenticate", h.handleAuthenticate)

	mux.Handle("GET /quote/random", h.authMiddleware(http.HandlerFunc(h.randomQuoteHandler)))
	mux.Handle("GET /quote/author", h.authMiddleware(http.HandlerFunc(h.getQuotesForAuthorHandler)))
	mux.Handle("POST /quote/author", h.authMiddleware(http.HandlerFunc(h.addQuote)))
}

func (s *Server) StartServer(quoteService quoteapi.QuoteService, smtpService email.MailService, authService auth.AuthService, logger *log.Logger) error {
	mux := http.NewServeMux()
	handler := Handler{
		quoteService: quoteService,
		authService:  authService,
		mailer:       smtpService,
		logger:       logger,
	}
	s.registerRoutes(mux, handler)
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
