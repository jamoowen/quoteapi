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
	IpRateLimiter     ipAddressRateLimiter
	ApiKeyRateLimiter apiKeyRateLimiter
	handler           Handler
}

// middleware looks ugly surely theres a cleaner way to write this
func (s *Server) registerRoutes(mux *http.ServeMux, h Handler) {
	mux.Handle("GET /authenticate",
		s.IpRateLimiter.limit(
			http.HandlerFunc(h.handleAuthenticate),
		),
	)
	mux.Handle("POST /authenticate",
		s.IpRateLimiter.limit(
			http.HandlerFunc(h.handleAuthenticate),
		),
	)

	// key protected

	mux.Handle("GET /quote/random",
		h.authMiddleware(
			s.ApiKeyRateLimiter.limit(
				s.IpRateLimiter.limit(
					http.HandlerFunc(h.randomQuoteHandler),
				),
			),
		),
	)

	mux.Handle("GET /quote/author",
		h.authMiddleware(
			s.ApiKeyRateLimiter.limit(
				s.IpRateLimiter.limit(
					http.HandlerFunc(h.getQuotesForAuthorHandler),
				),
			),
		),
	)

	mux.Handle("POST /quote/author",
		h.authMiddleware(
			s.ApiKeyRateLimiter.limit(
				s.IpRateLimiter.limit(
					http.HandlerFunc(h.addQuote),
				),
			),
		),
	)
}

func NewServer(quoteService quoteapi.QuoteService, smtpService email.MailService, authService auth.AuthService, ipRateLimitSeconds, apiKeyRateLimitSeconds int64, logger *log.Logger) (*Server, error) {
	s := Server{}
	ipAddressUsageMap := make(map[string]int64)
	s.IpRateLimiter = ipAddressRateLimiter{ipAddresses: ipAddressUsageMap, requiredIntervalSeconds: ipRateLimitSeconds}

	apiKeyUsageMap := make(map[string]int64)
	s.ApiKeyRateLimiter = apiKeyRateLimiter{apiKeys: apiKeyUsageMap, requiredIntervalSeconds: apiKeyRateLimitSeconds}

	handler := Handler{
		quoteService: quoteService,
		authService:  authService,
		mailer:       smtpService,
		logger:       logger,
	}
	s.handler = handler
	return &s, nil
}

func (s *Server) StartServer() error {
	mux := http.NewServeMux()
	s.registerRoutes(mux, s.handler)

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
