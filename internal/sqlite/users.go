package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jamoowen/quoteapi/internal/problems"
)

type User struct {
	Email              string
	HashedApiKey       string
	RequestCount       int64
	CreatedAtTimestamp int64
	LastUsedAt         time.Time
}

type NotFoundError struct {
}

type UsersStorage interface {
	GetUserByEmail(email string, ctx context.Context) (User, error)
	GetUserByKey(hashedKey string, ctx context.Context) (User, error)
	UpsertKeyForUser(email, hashedKey string, ctx context.Context) error
	IncrementRequestCountForUser(hashedKey string, ctx context.Context) error
}

type Users struct {
	db *sql.DB
}

func NewUsersStorage(db *sql.DB) *Users {
	return &Users{db}
}

func (u *Users) GetUserByEmail(email string, ctx context.Context) (User, error) {
	var user User
	row := u.db.QueryRowContext(ctx, `
		SELECT email, hashed_api_key, request_count, created_at_timestamp, last_used_at 
		FROM users
		WHERE email = ?
		`, email)
	err := row.Scan(
		&user.Email,
		&user.HashedApiKey,
		&user.RequestCount,
		&user.CreatedAtTimestamp,
		&user.LastUsedAt,
	)
	if err == sql.ErrNoRows {
		return User{}, problems.NewNotFoundError(fmt.Sprintf("No user found for this email (%v), email", email))
	}
	if err != nil {
		return User{}, fmt.Errorf("Failed to fetch user (%v): %v", email, err)
	}
	return user, nil
}

func (u *Users) GetUserByKey(hashedKey string, ctx context.Context) (User, error) {
	var user User
	row := u.db.QueryRowContext(ctx, `
		SELECT email, hashed_api_key, request_count, created_at_timestamp, last_used_at 
		FROM users
		WHERE email = ?
		`, hashedKey)
	err := row.Scan(
		&user.Email,
		&user.HashedApiKey,
		&user.RequestCount,
		&user.CreatedAtTimestamp,
		&user.LastUsedAt,
	)
	if err == sql.ErrNoRows {
		return User{}, problems.NewNotFoundError(fmt.Sprintf("No user found for this key (%v), email", hashedKey))

	}
	if err != nil {
		return User{}, fmt.Errorf("Failed to fetch user by key (%v): %v", hashedKey, err)
	}
	return user, nil
}

func (u *Users) UpsertKeyForUser(email, hashedKey string, ctx context.Context) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := u.db.ExecContext(ctx, `
        INSERT INTO users (email, hashed_api_key, last_used_at )
        VALUES (?, ?, ?)
        ON CONFLICT(email) DO UPDATE SET
        hashed_api_key = ?
    `, email, hashedKey, now, hashedKey)
	if err != nil {
		return fmt.Errorf("Failed to upsert API key for user (%v): %v", email, err)
	}
	return nil
}

func (u *Users) IncrementRequestCountForUser(hashedKey string, ctx context.Context) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	result, err := u.db.ExecContext(ctx, `
        UPDATE users 
        SET request_count = request_count + 1,
            last_used_at = ?
        WHERE hashed_api_key = ?
    `, now, hashedKey)
	if err != nil {
		return fmt.Errorf("Failed to increment request count for key (%v): %v", hashedKey, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return problems.NewNotFoundError(fmt.Sprintf("failed to increment key usage for this key (%v)", hashedKey))

	}
	return nil
}
