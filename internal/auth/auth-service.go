package auth

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/jamoowen/quoteapi/internal/cache"
	"github.com/jamoowen/quoteapi/internal/problems"
)

type AuthService interface {
	GenerateApiKey(email string) (string, error)
	GenerateOtp(email string) (string, error)
	AuthenticateOtp(email, otp string) (OTPStatus, error)
	AuthenticateApiKey(apiKey string) (bool, error)
}

type OTPStatus int

const (
	OTPValid OTPStatus = iota
	OTPInvalid
	OTPExpired
	OTPUserNotFound
	OTPError
)

type PersistedKey struct {
	id           int
	email        string
	apiKey       string
	requestCount int
	expiration   int
	createdAt    int
}

type Auth struct {
	otpService      cache.OtpService
	otpSecondsValid int64
	db              *sql.DB
}

func NewAuthService(db *sql.DB, otpSecondsValid int64) *Auth {
	otpCache := cache.NewOtpCache()
	return &Auth{
		otpService:      otpCache,
		otpSecondsValid: otpSecondsValid,
		db:              db,
	}
}

func (a *Auth) GenerateOtp(email string) (string, error) {
	newPin, err := a.generatePinForOtp()
	if err != nil {
		return "", err
	}
	a.otpService.CacheOtp(email, newPin)
	return newPin, nil
}

func (a *Auth) generatePinForOtp() (string, error) {
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

func (a *Auth) AuthenticateOtp(email, pin string) (OTPStatus, error) {
	otp, err := a.otpService.GetOtp(email)
	if errors.Is(err, &problems.NotFoundError{}) {
		return OTPUserNotFound, nil
	}
	if err != nil {
		return OTPError, err
	}
	if otp.Pin != pin {
		return OTPInvalid, nil
	}
	pinExpiration := time.Now().Unix() - a.otpSecondsValid
	if otp.CreatedTime < pinExpiration {
		return OTPExpired, nil
	}
	a.otpService.InvalidateOtp(email)
	return OTPValid, nil
}

// todo
func (a *Auth) GenerateApiKey(email string) (string, error) {
	return "SUPER SECRTET API KEY", nil
}

// todo
func (a *Auth) AuthenticateApiKey(apiKey string) (bool, error) {
	return true, nil
}

// gen new key for email => create new key and assign to cache
// cache:
