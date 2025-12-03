package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mashurimansur/goCMS/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserUseCase is a mock implementation of userusecase.UseCase
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Register(ctx context.Context, u *user.User, password string) error {
	args := m.Called(ctx, u, password)
	return args.Error(0)
}

func (m *MockUserUseCase) Login(ctx context.Context, email, password string) (string, *user.User, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Get(1).(*user.User), args.Error(2)
}

func (m *MockUserUseCase) GetProfile(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserUseCase) UpdateProfile(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserUseCase) ListUsers(ctx context.Context, limit, offset int) ([]*user.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *MockUserUseCase) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	// Mock auth middleware for registration (not needed for register but needed for Register method signature)
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	reqBody := registerRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Username: "testuser",
	}
	body, _ := json.Marshal(reqBody)

	mockUseCase.On("Register", mock.Anything, mock.AnythingOfType("*user.User"), reqBody.Password).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUserHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	reqBody := loginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	expectedUser := &user.User{ID: "user-id", Email: reqBody.Email}
	mockUseCase.On("Login", mock.Anything, reqBody.Email, reqBody.Password).Return("access-token", expectedUser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUserHandler_GetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	expectedUser := &user.User{ID: "user-123", Email: "test@example.com", FullName: "Test User"}
	mockUseCase.On("GetProfile", mock.Anything, "user-123").Return(expectedUser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/user-123", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUserHandler_GetProfile_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	mockUseCase.On("GetProfile", mock.Anything, "user-123").Return((*user.User)(nil), assert.AnError)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/user-123", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	reqBody := user.User{
		FullName: "Updated User",
		Email:    "updated@example.com",
	}
	body, _ := json.Marshal(reqBody)

	mockUseCase.On("UpdateProfile", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/admin/users/user-123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUserHandler_UpdateProfile_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/admin/users/user-123", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_UpdateProfile_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	reqBody := user.User{
		FullName: "Updated User",
		Email:    "updated@example.com",
	}
	body, _ := json.Marshal(reqBody)

	mockUseCase.On("UpdateProfile", mock.Anything, mock.AnythingOfType("*user.User")).Return(assert.AnError)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/admin/users/user-123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUserHandler_ListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	users := []*user.User{
		{ID: "user-1", Email: "user1@example.com"},
		{ID: "user-2", Email: "user2@example.com"},
	}
	mockUseCase.On("ListUsers", mock.Anything, 10, 0).Return(users, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUserHandler_ListUsers_WithPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	users := []*user.User{
		{ID: "user-1", Email: "user1@example.com"},
	}
	mockUseCase.On("ListUsers", mock.Anything, 5, 10).Return(users, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/?limit=5&offset=10", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUserHandler_ListUsers_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	mockUseCase.On("ListUsers", mock.Anything, 10, 0).Return(([]*user.User)(nil), assert.AnError)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUserHandler_DeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	mockUseCase.On("DeleteUser", mock.Anything, "user-123").Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/users/user-123", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUserHandler_DeleteUser_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	mockUseCase.On("DeleteUser", mock.Anything, "user-123").Return(assert.AnError)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/users/user-123", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUserHandler_Register_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	reqBody := registerRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	mockUseCase.On("Register", mock.Anything, mock.AnythingOfType("*user.User"), reqBody.Password).Return(assert.AnError)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUserHandler_Register_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString("invalid json"))
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Login_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	reqBody := loginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	mockUseCase.On("Login", mock.Anything, reqBody.Email, reqBody.Password).Return("", (*user.User)(nil), assert.AnError)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUserHandler_Login_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockUserUseCase)
	handler := NewUserHandler(mockUseCase)

	router := gin.New()
	authMiddleware := func(c *gin.Context) { c.Next() }
	handler.Register(router.Group("/api/v1"), authMiddleware)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString("invalid json"))
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
