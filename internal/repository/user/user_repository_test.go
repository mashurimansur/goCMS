package user

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mashurimansur/goCMS/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	u := &user.User{
		FullName:     "Test User",
		Username:     "testuser",
		Email:        "test@example.com",
		Phone:        "1234567890",
		PasswordHash: "hash",
		Role:         "user",
		Status:       "active",
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users")).
		WithArgs(sqlmock.AnyArg(), u.FullName, u.Username, u.Email, u.Phone, u.PasswordHash, u.AvatarURL, u.Role, u.Status, u.EmailVerified, u.PhoneVerified, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), u)
	assert.NoError(t, err)
	assert.NotEmpty(t, u.ID)
	assert.NotZero(t, u.CreatedAt)
	assert.NotZero(t, u.UpdatedAt)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	email := "test@example.com"
	rows := sqlmock.NewRows([]string{"id", "full_name", "username", "email", "phone", "password_hash", "avatar_url", "role", "status", "last_login", "email_verified", "phone_verified", "created_at", "updated_at"}).
		AddRow("uuid", "Test User", "testuser", email, "1234567890", "hash", "avatar.jpg", "user", "active", time.Now(), true, true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, full_name, username, email")).
		WithArgs(email).
		WillReturnRows(rows)

	u, err := repo.GetByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, email, u.Email)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, full_name, username, email")).
		WithArgs("unknown@example.com").
		WillReturnError(sql.ErrNoRows)

	u, err := repo.GetByEmail(context.Background(), "unknown@example.com")
	assert.NoError(t, err)
	assert.Nil(t, u)
}

func TestUserRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	userID := "test-uuid"
	rows := sqlmock.NewRows([]string{"id", "full_name", "username", "email", "phone", "password_hash", "avatar_url", "role", "status", "last_login", "email_verified", "phone_verified", "created_at", "updated_at"}).
		AddRow(userID, "Test User", "testuser", "test@example.com", "1234567890", "hash", "avatar.jpg", "user", "active", time.Now(), true, true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, full_name, username, email")).
		WithArgs(userID).
		WillReturnRows(rows)

	u, err := repo.GetByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, userID, u.ID)
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	username := "testuser"
	rows := sqlmock.NewRows([]string{"id", "full_name", "username", "email", "phone", "password_hash", "avatar_url", "role", "status", "last_login", "email_verified", "phone_verified", "created_at", "updated_at"}).
		AddRow("uuid", "Test User", username, "test@example.com", "1234567890", "hash", "avatar.jpg", "user", "active", time.Now(), true, true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, full_name, username, email")).
		WithArgs(username).
		WillReturnRows(rows)

	u, err := repo.GetByUsername(context.Background(), username)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, username, u.Username)
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	u := &user.User{
		ID:       "uuid",
		FullName: "Updated User",
		Username: "updateduser",
		Email:    "updated@example.com",
		Phone:    "9876543210",
		Role:     "admin",
		Status:   "active",
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users")).
		WithArgs(u.FullName, u.Username, u.Email, u.Phone, u.AvatarURL, u.Role, u.Status, u.EmailVerified, u.PhoneVerified, sqlmock.AnyArg(), u.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(context.Background(), u)
	assert.NoError(t, err)
	assert.NotZero(t, u.UpdatedAt)
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	userID := "test-uuid"

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = ?")).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), userID)
	assert.NoError(t, err)
}

func TestUserRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "full_name", "username", "email", "phone", "password_hash", "avatar_url", "role", "status", "last_login", "email_verified", "phone_verified", "created_at", "updated_at"}).
		AddRow("uuid1", "User 1", "user1", "user1@example.com", "1234567890", "hash", "avatar.jpg", "user", "active", time.Now(), true, true, time.Now(), time.Now()).
		AddRow("uuid2", "User 2", "user2", "user2@example.com", "9876543210", "hash", "avatar2.jpg", "user", "active", time.Now(), true, true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, full_name, username, email")).
		WithArgs(10, 0).
		WillReturnRows(rows)

	users, err := repo.List(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
}

func TestUserRepository_List_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "full_name", "username", "email", "phone", "password_hash", "avatar_url", "role", "status", "last_login", "email_verified", "phone_verified", "created_at", "updated_at"})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, full_name, username, email")).
		WithArgs(10, 0).
		WillReturnRows(rows)

	users, err := repo.List(context.Background(), 10, 0)
	assert.NoError(t, err)
	// List returns nil when there are no results, which is fine
	if users != nil {
		assert.Len(t, users, 0)
	}
}
