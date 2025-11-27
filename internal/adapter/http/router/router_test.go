package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/mashurimansur/goCMS/internal/adapter/http/handler"
	domain "github.com/mashurimansur/goCMS/internal/domain/person"
)

func TestNewGinEngine_WithHandler(t *testing.T) {
	previousMode := gin.Mode()
	defer gin.SetMode(previousMode)

	gin.SetMode(gin.DebugMode)

	personHandler := handler.NewPersonHandler(stubPersonUseCase{person: domain.Person{Name: "Router"}})
	engine := NewGinEngine(Options{
		Mode:          gin.TestMode,
		PersonHandler: personHandler,
	})

	if mode := gin.Mode(); mode != gin.TestMode {
		t.Fatalf("expected gin mode %s, got %s", gin.TestMode, mode)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/person", nil)
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected handler to respond 200, got %d", rec.Code)
	}
}

func TestNewGinEngine_HealthAndMissingHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := NewGinEngine(Options{})

	// Person route should be missing without handler.
	missing := httptest.NewRecorder()
	engine.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/api/v1/admin/person", nil))
	if missing.Code != http.StatusNotFound {
		t.Fatalf("expected 404 when no handler registered, got %d", missing.Code)
	}

	// Healthz route always available.
	health := httptest.NewRecorder()
	engine.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/health", nil))
	if health.Code != http.StatusOK {
		t.Fatalf("expected 200 on /health, got %d", health.Code)
	}
}

type stubPersonUseCase struct {
	person domain.Person
	err    error
}

func (s stubPersonUseCase) GetDefaultPerson(ctx context.Context) (domain.Person, error) {
	return s.person, s.err
}
