package http

import (
	"net"
	"net/http"
	"sync"
	"time"
)

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			http.Error(w, "X-API-KEY header missing", http.StatusForbidden)
			return
		}
		// check db for user
		authorized, err := h.authService.AuthenticateApiKey(apiKey, r.Context())
		if err != nil {
			h.logger.Println("Failed to authenticate api key: ", err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		if !authorized {
			http.Error(w, "Invalid api key", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type ApiKeyRateLimiter struct {
	requiredIntervalSeconds int64
	apiKeys                 map[string]int64
	mu                      sync.RWMutex
}

type IpAddressRateLimiter struct {
	requiredIntervalSeconds int64
	ipAddresses             map[string]int64
	mu                      sync.RWMutex
}

func (l *ApiKeyRateLimiter) limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Unix()
		apiKey := r.Header.Get("X-API-KEY")
		// no key then we cant infer usage
		if apiKey == "" {
			next.ServeHTTP(w, r)
		}
		l.mu.Lock()
		defer l.mu.Unlock()
		lastUsed, ok := l.apiKeys[apiKey]
		l.apiKeys[apiKey] = now
		if !ok {
			next.ServeHTTP(w, r)
		} else if now-lastUsed > l.requiredIntervalSeconds {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Too many requests. API keys are limited to 1 request per min", http.StatusTooManyRequests)
		}
	})
}

func (l *IpAddressRateLimiter) limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Unix()
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// no key then we cant infer usage
		if ip == "" {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			next.ServeHTTP(w, r)
		}
		l.mu.Lock()
		defer l.mu.Unlock()
		lastUsed, ok := l.ipAddresses[ip]
		l.ipAddresses[ip] = now
		if ok && now-lastUsed < l.requiredIntervalSeconds {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (s *Server) CleanupMiddlewareCache() {
	now := time.Now().Unix()
	s.ApiKeyRateLimiter.mu.Lock()
	for key, lastUsed := range s.ApiKeyRateLimiter.apiKeys {
		if now-lastUsed > s.ApiKeyRateLimiter.requiredIntervalSeconds {
			delete(s.ApiKeyRateLimiter.apiKeys, key)
		}
	}
	s.ApiKeyRateLimiter.mu.Unlock()
	s.IpAddressRateLimiter.mu.Lock()
	for key, lastUsed := range s.IpAddressRateLimiter.ipAddresses {
		if now-lastUsed > s.IpAddressRateLimiter.requiredIntervalSeconds {
			delete(s.IpAddressRateLimiter.ipAddresses, key)
		}
	}
	s.IpAddressRateLimiter.mu.Unlock()
}
