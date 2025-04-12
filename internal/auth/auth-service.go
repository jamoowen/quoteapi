package auth

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/jamoowen/quoteapi/internal/cache"
)

type AuthService interface {
	GenerateAPIKey(email string) (string, error)
	GenerateOtp(email string) (string, error)
	Authenticate(apiKey string) (bool, error)
}

type PersistedKey struct {
	id           int
	email        string
	apiKey       string
	requestCount int
	expiration   int
	createdAt    int
}

type Auth struct {
	otpService cache.OtpService
	db         *sql.DB
}

func (a *Auth) NewAuthService(db *sql.DB) Auth {
	otpCache := cache.NewOtpCache()
	return Auth{
		otpService: otpCache,
		db:         db,
	}
}

func (a *Auth) GenerateOtp(email string, secondsTilExpiration int64) (string, error) {
	maxDigits := 6
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))),
	)
	if err != nil {
		return "", fmt.Errorf("Failed to generate otp: %w\n", err)
	}
	pin := fmt.Sprintf("%0*d", maxDigits, bi)

	return pin, nil
}

func (a *Auth) GenerateApiKey(email string) (string, error) {
	return "", nil
}

// gen new key for email => create new key and assign to cache
// cache:
