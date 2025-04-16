package cache

import (
	"time"
)

func CleanupOtpCache(c *OtpCache, cleanupDelay time.Duration, otpValiditySeconds int64) {
	for {
		time.Sleep(cleanupDelay)
		now := time.Now().Unix()
		c.mu.Lock()
		for email, otp := range c.cache {
			if now-otp.CreatedTime > otpValiditySeconds {
				delete(c.cache, email)
			}
		}
	}
}
