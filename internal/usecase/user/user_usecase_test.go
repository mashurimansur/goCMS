package user

import (
	"context"
	"testing"
	"time"

	"github.com/mashurimansur/goCMS/internal/domain/user"
	"github.com/mashurimansur/goCMS/internal/utils/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*user.User), args.Error(1)
}

type MockTokenMaker struct {
	mock.Mock
}

func (m *MockTokenMaker) CreateToken(username string, duration time.Duration) (string, *token.Payload, error) {
	args := m.Called(username, duration)
	return args.String(0), args.Get(1).(*token.Payload), args.Error(2)
}

func (m *MockTokenMaker) VerifyToken(tokenStr string) (*token.Payload, error) {
	args := m.Called(tokenStr)
	return args.Get(0).(*token.Payload), args.Error(1)
}

func TestUserUseCase_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockMaker := new(MockTokenMaker)
	uc := NewUserUseCase(mockRepo, mockMaker, time.Hour)

	u := &user.User{
		FullName: "Test User",
		Email:    "test@example.com",
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(arg *user.User) bool {
		return arg.Email == u.Email && arg.Role == "user" && arg.Status == "active" && arg.PasswordHash != ""
	})).Return(nil)

	err := uc.Register(context.Background(), u, "password123")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockMaker := new(MockTokenMaker)
	uc := NewUserUseCase(mockRepo, mockMaker, time.Hour)

	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	u := &user.User{
		ID:           "user-id",
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	mockRepo.On("GetByEmail", mock.Anything, email).Return(u, nil)
	mockMaker.On("CreateToken", u.ID, time.Hour).Return("access_token", &token.Payload{}, nil)

	token, user, err := uc.Login(context.Background(), email, password)
	assert.NoError(t, err)
	assert.Equal(t, "access_token", token)
	assert.Equal(t, u, user)
	mockRepo.AssertExpectations(t)
	mockMaker.AssertExpectations(t)
}

func TestUserUseCase_GetProfile(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockMaker := new(MockTokenMaker)
	uc := NewUserUseCase(mockRepo, mockMaker, time.Hour)

	userID := "user-id"
	u := &user.User{
		ID:       userID,
		FullName: "Test User",
		Email:    "test@example.com",
	}

	mockRepo.On("GetByID", mock.Anything, userID).Return(u, nil)

	result, err := uc.GetProfile(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, u, result)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_UpdateProfile(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockMaker := new(MockTokenMaker)
	uc := NewUserUseCase(mockRepo, mockMaker, time.Hour)

	u := &user.User{
		ID:       "user-id",
		FullName: "Updated User",
		Email:    "updated@example.com",
	}

	mockRepo.On("Update", mock.Anything, u).Return(nil)

	err := uc.UpdateProfile(context.Background(), u)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_ListUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockMaker := new(MockTokenMaker)
	uc := NewUserUseCase(mockRepo, mockMaker, time.Hour)

	users := []*user.User{
		{ID: "user-1", Email: "user1@example.com"},
		{ID: "user-2", Email: "user2@example.com"},
	}

	mockRepo.On("List", mock.Anything, 10, 0).Return(users, nil)

	result, err := uc.ListUsers(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, users, result)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockMaker := new(MockTokenMaker)
	uc := NewUserUseCase(mockRepo, mockMaker, time.Hour)

	userID := "user-id"
	mockRepo.On("Delete", mock.Anything, userID).Return(nil)

	err := uc.DeleteUser(context.Background(), userID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
