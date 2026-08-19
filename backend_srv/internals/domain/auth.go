package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, username string, password string, full_name string, email string) error
	GetUser(ctx context.Context, username string) (User, error)
	CreateVerifyEmail(ctx context.Context, username string, email string, secret_code string) (InfoEmailVerification, error)
	VerifyEmail(ctx context.Context, email string) error
	VerifySecretCode(ctx context.Context, username string, secret_code string) (string, error)
}

type User struct {
	Username          string     `json:"username"`
	HashedPassword    string     `json:"password"`
	CreatedAt         *time.Time `json:"created_at"`
	FullName          string     `json:"full_name"`
	Email             string     `json:"email"`
	PasswordChangedAt time.Time  `json:"password_changed_at"`
	IsEmailVerified   bool       `json:"is_email_verified"`
	Role              string     `json:"role"`
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
