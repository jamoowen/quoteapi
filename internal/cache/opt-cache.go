package cache

import (
	"sync"
	"time"
)

type Otp struct {
	pin         string
	createdTime int64
}

type OtpService interface {
	CacheOtp(email, pin string) error
	GetOtp(email string) Otp
	CleanupCache(expirationTime int64) error
}

type OtpCache struct {
	cache map[string]Otp
	mu    sync.RWMutex
}

func NewOtpCache() *OtpCache {
	return &OtpCache{cache: make(map[string]Otp)}
}

func (c *OtpCache) CacheOtp(email, pin string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[email] = Otp{pin: pin, createdTime: time.Now().Unix()}
	return nil
}

func (c *OtpCache) GetOtp(email string) Otp {
	c.mu.RLock()
	defer c.mu.Unlock()
	return c.cache[email]
}

func (c *OtpCache) CleanupCache(expirationUnixTimestamp int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for email, otp := range c.cache {
		if otp.createdTime < expirationUnixTimestamp {
			delete(c.cache, email)
		}
	}
	return nil
}
