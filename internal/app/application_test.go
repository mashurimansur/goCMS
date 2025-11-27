package app

import (
	"context"
	"net"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"

	"github.com/mashurimansur/goCMS/internal/utils/config"
	"github.com/mashurimansur/goCMS/internal/utils/database"
)

func TestBuildPersonRepository(t *testing.T) {
	if _, err := buildPersonRepository(nil); err == nil {
		t.Fatalf("expected error when connection is nil")
	}

	if _, err := buildPersonRepository(&database.Connection{}); err == nil {
		t.Fatalf("expected error when db is nil")
	}

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo, err := buildPersonRepository(&database.Connection{DB: db})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo == nil {
		t.Fatalf("expected repository instance")
	}
}

func TestApplicationNew_WithMissingDBConfig(t *testing.T) {
	_, err := New(context.Background(), config.AppConfig{})
	if err == nil {
		t.Fatalf("expected error when database config is missing")
	}
}

func TestApplicationNew_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dsn := "app_success"
	db, mock, err := sqlmock.NewWithDSN(dsn, sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectPing()
	mock.ExpectClose()

	cfg := config.AppConfig{
		HTTPAddr: ":8081",
		GinMode:  gin.TestMode,
		Database: database.Config{
			Driver: "sqlmock",
			DSN:    dsn,
		},
	}

	application, err := New(context.Background(), cfg)
	if err != nil {
		t.Fatalf("unexpected error from New: %v", err)
	}
	if application.engine == nil {
		t.Fatalf("expected engine to be initialized")
	}

	if err := application.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestApplicationRun_DefaultAddress(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Occupy the default port so Run fails immediately instead of blocking.
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Skipf("unable to listen on :8080 to test default address: %v", err)
	}
	defer ln.Close()

	app := &Application{engine: gin.New()}
	if err := app.Run(); err == nil {
		t.Fatalf("expected error when default port is unavailable")
	}
}

func TestApplicationRun_RequiresEngine(t *testing.T) {
	app := &Application{}
	if err := app.Run(); err == nil {
		t.Fatalf("expected error when engine is nil")
	}
}

func TestApplicationClose_NilSafe(t *testing.T) {
	app := &Application{}
	if err := app.Close(); err != nil {
		t.Fatalf("expected nil error on nil db connection, got %v", err)
	}
}
