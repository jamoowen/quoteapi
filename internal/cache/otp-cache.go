package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/jamoowen/quoteapi/internal/problems"
)

type Otp struct {
	Pin         string
	CreatedTime int64
}

type OtpService interface {
	CacheOtp(email, pin string) error
	GetOtp(email string) (Otp, error)
	InvalidateOtp(email string) error
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
	c.cache[email] = Otp{Pin: pin, CreatedTime: time.Now().Unix()}
	return nil
}

func (c *OtpCache) GetOtp(email string) (Otp, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	otp, ok := c.cache[email]
	if ok == false {
		return Otp{}, problems.NewNotFoundError(fmt.Sprintf("No OTP found for email (%v)", email))
	}
	return otp, nil
}

func (c *OtpCache) InvalidateOtp(email string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, email)
	return nil
}

func (c *OtpCache) CleanupCache(expirationUnixTimestamp int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for email, otp := range c.cache {
		if otp.CreatedTime < expirationUnixTimestamp {
			delete(c.cache, email)
		}
	}
	return nil
}
