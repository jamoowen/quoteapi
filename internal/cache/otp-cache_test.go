package cache

import (
	"testing"
	"time"
)

func TestCacheOtp(t *testing.T) {
	c := NewOtpCache()
	pin := "testOtp"
	email := "example@gmail.com"
	err := c.CacheOtp(email, pin)
	if err != nil {
		t.Fatalf("Failewd to cache otp: %v", err.Error())
	}
	if len(c.cache) != 1 {
		t.Errorf("Expected 1 item in the cache")
	}
	otp, ok := c.cache[email]
	if ok == false || otp.Pin != pin {
		t.Errorf("Failed to cache pin")
	}
}

func TestCleanupCache(t *testing.T) {
	c := NewOtpCache()
	pin := "testOtp"
	email := "example@gmail.com"
	err := c.CacheOtp(email, pin)
	if err != nil {
		t.Fatalf("Failewd to cache otp: %v", err.Error())
	}
	if len(c.cache) != 1 {
		t.Errorf("Expected 1 item in the cache")
	}
	time.Sleep(5 * time.Second)
	otpSecondsValid := 4
	c.CleanupCache(int64(otpSecondsValid))
	otp, ok := c.cache[email]
	if ok == true {
		t.Errorf("Expected cleanup to delete expired keys: %v", otp)
	}
	if len(c.cache) != 0 {
		t.Errorf("Expected cache to be empty")
	}
}
