package user

import (
	"context"
	"time"
)

// User models the user data.
type User struct {
	ID            string    `json:"id"`
	FullName      string    `json:"full_name"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	PasswordHash  string    `json:"-"`
	AvatarURL     string    `json:"avatar_url"`
	Role          string    `json:"role"`
	Status        string    `json:"status"`
	LastLogin     time.Time `json:"last_login"`
	EmailVerified bool      `json:"email_verified"`
	PhoneVerified bool      `json:"phone_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Repository abstracts the data source that stores user information.
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
}
