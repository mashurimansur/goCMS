package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	domain "github.com/mashurimansur/goCMS/internal/domain/person"
)

type stubPersonUseCase struct {
	person domain.Person
	err    error
}

func (s stubPersonUseCase) GetDefaultPerson(ctx context.Context) (domain.Person, error) {
	return s.person, s.err
}

func TestGetDefaultPerson_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	handler := NewPersonHandler(stubPersonUseCase{person: domain.Person{Name: "Bob"}})
	api := router.Group("/api/v1")
	handler.Register(api)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/person", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var payload domain.Person
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload.Name != "Bob" {
		t.Fatalf("unexpected person name: %q", payload.Name)
	}
}

func TestGetDefaultPerson_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	handler := NewPersonHandler(stubPersonUseCase{err: context.DeadlineExceeded})
	api := router.Group("/api/v1")
	handler.Register(api)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/person", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload["error"] == "" {
		t.Fatalf("expected error message in response")
	}
}
