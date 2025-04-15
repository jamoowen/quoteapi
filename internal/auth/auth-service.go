package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jamoowen/quoteapi/internal/cache"
	"github.com/jamoowen/quoteapi/internal/problems"
	"github.com/jamoowen/quoteapi/internal/sqlite"
)

type AuthService interface {
	CreateNewApiKeyForUser(email string, ctx context.Context) (string, error)
	GenerateOtp(email string) (string, error)
	AuthenticateOtp(email, otp string) (OTPStatus, error)
	AuthenticateApiKey(apiKey string, ctx context.Context) (bool, error)
}

type OTPStatus int

const (
	OTPValid OTPStatus = iota
	OTPInvalid
	OTPExpired
	OTPUserNotFound
	OTPError
)

type Auth struct {
	otpService      cache.OtpService
	otpSecondsValid int64
	usersStorage    sqlite.UsersStorage
	apiKeySecret    string
}

func NewAuthService(db *sql.DB, otpSecondsValid int64) *Auth {
	usersStorage := sqlite.NewUsersStorage(db)
	otpCache := cache.NewOtpCache()
	return &Auth{
		otpService:      otpCache,
		otpSecondsValid: otpSecondsValid,
		usersStorage:    usersStorage,
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
		return "", fmt.Errorf("failed to generate otp: %w", err)
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

func (a *Auth) generateApiKey() (string, string) {
	// is determinism a word???
	aLittleBitOfDeterminism := strconv.FormatInt(time.Now().UnixMicro(), 12)
	apiKey := uuid.New().String() + aLittleBitOfDeterminism
	hashedApiKey := a.hashString(apiKey)
	return apiKey, hashedApiKey
}

func (a *Auth) CreateNewApiKeyForUser(email string, ctx context.Context) (string, error) {
	apiKey, hashedApiKey := a.generateApiKey()
	err := a.usersStorage.UpsertKeyForUser(email, hashedApiKey, ctx)
	if err != nil {
		return "", err
	}
	return apiKey, nil
}

func (a *Auth) AuthenticateApiKey(apiKey string, ctx context.Context) (bool, error) {
	hashedApiKey := a.hashString(apiKey)
	_, err := a.usersStorage.GetUserByKey(hashedApiKey, ctx)
	if errors.Is(err, &problems.NotFoundError{}) {
		return false, nil
	}
	return true, err
}

func (a *Auth) hashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashBytes := h.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

// gen new key for email => create new key and assign to cache
// cache:
