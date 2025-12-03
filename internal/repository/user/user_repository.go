package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mashurimansur/goCMS/internal/domain/user"
)

// UserRepository implements user.Repository for MySQL.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new MySQL user repository.
func NewUserRepository(db *sql.DB) user.Repository {
	return &UserRepository{db: db}
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}

	query := `
		INSERT INTO users (
			id, full_name, username, email, phone, password_hash, avatar_url, role, status, email_verified, phone_verified, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.FullName, u.Username, u.Email, u.Phone, u.PasswordHash, u.AvatarURL, u.Role, u.Status, u.EmailVerified, u.PhoneVerified, u.CreatedAt, u.UpdatedAt,
	)
	return err
}

// GetByEmail retrieves a user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, full_name, username, email, phone, password_hash, avatar_url, role, status, last_login, email_verified, phone_verified, created_at, updated_at
		FROM users
		WHERE email = ?
	`
	return r.scanUser(ctx, query, email)
}

// GetByID retrieves a user by ID.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	query := `
		SELECT id, full_name, username, email, phone, password_hash, avatar_url, role, status, last_login, email_verified, phone_verified, created_at, updated_at
		FROM users
		WHERE id = ?
	`
	return r.scanUser(ctx, query, id)
}

// GetByUsername retrieves a user by username.
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	query := `
		SELECT id, full_name, username, email, phone, password_hash, avatar_url, role, status, last_login, email_verified, phone_verified, created_at, updated_at
		FROM users
		WHERE username = ?
	`
	return r.scanUser(ctx, query, username)
}

func (r *UserRepository) scanUser(ctx context.Context, query string, args ...interface{}) (*user.User, error) {
	u := &user.User{}
	var lastLogin sql.NullTime
	var avatarURL sql.NullString
	var phone sql.NullString
	var username sql.NullString

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&u.ID, &u.FullName, &username, &u.Email, &phone, &u.PasswordHash, &avatarURL, &u.Role, &u.Status, &lastLogin, &u.EmailVerified, &u.PhoneVerified, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil if not found, or custom error
		}
		return nil, err
	}

	if lastLogin.Valid {
		u.LastLogin = lastLogin.Time
	}
	if avatarURL.Valid {
		u.AvatarURL = avatarURL.String
	}
	if phone.Valid {
		u.Phone = phone.String
	}
	if username.Valid {
		u.Username = username.String
	}

	return u, nil
}

// Update updates an existing user.
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	u.UpdatedAt = time.Now()
	query := `
		UPDATE users
		SET full_name = ?, username = ?, email = ?, phone = ?, avatar_url = ?, role = ?, status = ?, email_verified = ?, phone_verified = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		u.FullName, u.Username, u.Email, u.Phone, u.AvatarURL, u.Role, u.Status, u.EmailVerified, u.PhoneVerified, u.UpdatedAt, u.ID,
	)
	return err
}

// Delete deletes a user by ID.
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves a list of users with pagination.
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	query := `
		SELECT id, full_name, username, email, phone, password_hash, avatar_url, role, status, last_login, email_verified, phone_verified, created_at, updated_at
		FROM users
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		u := &user.User{}
		var lastLogin sql.NullTime
		var avatarURL sql.NullString
		var phone sql.NullString
		var username sql.NullString

		err := rows.Scan(
			&u.ID, &u.FullName, &username, &u.Email, &phone, &u.PasswordHash, &avatarURL, &u.Role, &u.Status, &lastLogin, &u.EmailVerified, &u.PhoneVerified, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if lastLogin.Valid {
			u.LastLogin = lastLogin.Time
		}
		if avatarURL.Valid {
			u.AvatarURL = avatarURL.String
		}
		if phone.Valid {
			u.Phone = phone.String
		}
		if username.Valid {
			u.Username = username.String
		}

		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
