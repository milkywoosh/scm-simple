package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"scm-simple-luke.com/dir/internals"
)

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, username string, password string, full_name string, email string) error
	GetUser(ctx context.Context, username string) (UserRow, error)
	CreateVerifyEmail(ctx context.Context, username string, email string, secret_code string) (InfoEmailVerification, error)
	VerifyEmail(ctx context.Context, email string) error
	VerifySecretCode(ctx context.Context, username string, secret_code string) (string, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
}

type User struct {
	Username          string               `json:"username"`
	HashedPassword    internals.NullString `json:"password"`
	CreatedAt         internals.NullTime   `json:"created_at"`
	FullName          string               `json:"full_name"`
	Email             string               `json:"email"`
	PasswordChangedAt internals.NullTime   `json:"password_changed_at"`
	IsEmailVerified   bool                 `json:"is_email_verified"`
	Role              internals.NullString `json:"role"`
}

type UserRow struct {
	Username       internals.NullString `json:"username"`
	HashedPassword internals.NullString `json:"password"`
	CreatedAt      internals.NullTime   `json:"created_at"`
	Email          internals.NullString `json:"email"`
	Fullname       internals.NullString `json:"full_name"`
}

type Authentication interface {
	UserRepository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type InfoEmailVerification struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	SecretCode string `json:"secret_code"`
}
