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
	"golang.org/x/time/rate"
)

const (
	API_RATE_LIMIT_SECONDS        = 60
	IP_ADDRESS_RATE_LIMIT_SECONDS = 10
)

type Server struct {
	ipRateLimiter     ipAddressRateLimiter
	apiKeyRateLimiter apiKeyRateLimiter
}

// middleware looks ugly surely theres a cleaner way to write this
func (s *Server) registerRoutes(mux *http.ServeMux, h Handler) {
	mux.Handle("GET /authenticate",
		s.ipRateLimiter.limit(
			http.HandlerFunc(h.handleAuthenticate),
		),
	)
	mux.Handle("POST /authenticate",
		s.ipRateLimiter.limit(
			http.HandlerFunc(h.handleAuthenticate),
		),
	)

	// key protected

	mux.Handle("GET /quote/random",
		h.authMiddleware(
			s.apiKeyRateLimiter.limit(
				s.ipRateLimiter.limit(
					http.HandlerFunc(h.randomQuoteHandler),
				),
			),
		),
	)

	mux.Handle("GET /quote/author",
		h.authMiddleware(
			s.apiKeyRateLimiter.limit(
				s.ipRateLimiter.limit(
					http.HandlerFunc(h.getQuotesForAuthorHandler),
				),
			),
		),
	)

	mux.Handle("POST /quote/author",
		h.authMiddleware(
			s.apiKeyRateLimiter.limit(
				s.ipRateLimiter.limit(
					http.HandlerFunc(h.addQuote),
				),
			),
		),
	)
}

func (s *Server) StartServer(quoteService quoteapi.QuoteService, smtpService email.MailService, authService auth.AuthService, logger *log.Logger) error {
	mux := http.NewServeMux()
	limiter := rate.NewLimiter(1, 3)
	handler := Handler{
		quoteService: quoteService,
		authService:  authService,
		mailer:       smtpService,
		limiter:      limiter,
		logger:       logger,
	}

	ipAddressUsageMap := make(map[string]int64)
	s.ipRateLimiter = ipAddressRateLimiter{ipAddresses: ipAddressUsageMap}

	apiKeyUsageMap := make(map[string]int64)
	s.apiKeyRateLimiter = apiKeyRateLimiter{apiKeys: apiKeyUsageMap}

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
