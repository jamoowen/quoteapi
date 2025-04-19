package auth

import (
	"context"
	"testing"

	"github.com/jamoowen/quoteapi/internal/sqlite"
)

func TestGenerateOtp(t *testing.T) {
	auth := Auth{}
	auth.otpSecondsValid = 10
	_, err := auth.generatePinForOtp()
	if err != nil {
		t.Fatalf(":Failed to generate a random pin for otp: %v", err.Error())
	}
}

func TestGenerateApiKey(t *testing.T) {
	a := Auth{}
	apiKey, hashedApiKey := a.generateApiKey()
	if apiKey == "" || hashedApiKey == "" {
		t.Fatal("Failed to generate api keys")
	}
}

func TestAuthenticateApiKey(t *testing.T) {
	a := Auth{}
	apiKey, hashedApiKey := a.generateApiKey()
	mockUsers := MockUsers{HashedApiKey: hashedApiKey}
	a.usersStorage = &mockUsers

	var c context.Context
	authorized, err := a.AuthenticateApiKey(apiKey, c)
	if err != nil {
		t.Fatal("Failed to authenticate api key")
	}
	if !authorized {
		t.Fatal("Auth returned false")
	}
}

type MockUsers struct {
	HashedApiKey string
	RequestCount int64
}

func (u *MockUsers) GetUserByEmail(email string, ctx context.Context) (sqlite.User, error) {
	user := sqlite.User{Email: email, HashedApiKey: u.HashedApiKey, RequestCount: u.RequestCount}
	return user, nil
}

func (u *MockUsers) GetUserByKey(hashedKey string, ctx context.Context) (sqlite.User, error) {
	return sqlite.User{}, nil
}

func (u *MockUsers) UpsertKeyForUser(email, hashedKey string, ctx context.Context) error {
	return nil
}

func (u *MockUsers) IncrementRequestCountForUser(hashedKey string, ctx context.Context) error {
	return nil
}
