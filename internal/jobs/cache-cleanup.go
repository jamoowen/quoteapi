package jobs

import (
	"time"

	"github.com/jamoowen/quoteapi/internal/auth"
	"github.com/jamoowen/quoteapi/internal/http"
)

func CacheCleanupJob(interval time.Duration, otpCacheExpirationSeconds int64, authService auth.AuthService, server *http.Server) {
	for {
		time.Sleep(interval)
		authService.CleanupCache(otpCacheExpirationSeconds)
		server.CleanupMiddlewareCache()
	}
}
