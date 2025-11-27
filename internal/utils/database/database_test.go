package database

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewConnection_EmptyConfig(t *testing.T) {
	conn, err := NewConnection(context.Background(), Config{})
	if err != nil {
		t.Fatalf("expected nil error for empty config, got %v", err)
	}
	if conn == nil || conn.DB != nil {
		t.Fatalf("expected no-op connection with nil DB")
	}
}

func TestNewConnection_WithSQLMock(t *testing.T) {
	dsn := "conn_success"
	db, mock, err := sqlmock.NewWithDSN(dsn, sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectPing()

	cfg := Config{
		Driver:            "sqlmock",
		DSN:               dsn,
		MaxOpenConns:      2,
		MaxIdleConns:      1,
		ConnMaxLifetime:   time.Second,
		ConnectionTimeout: time.Second,
	}

	conn, err := NewConnection(context.Background(), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn.DB == nil {
		t.Fatalf("expected DB to be initialized")
	}

	mock.ExpectClose()
	if err := conn.Close(); err != nil {
		t.Fatalf("close returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestNewConnection_PingFailure(t *testing.T) {
	dsn := "conn_failure"
	db, mock, err := sqlmock.NewWithDSN(dsn, sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectPing().WillReturnError(errors.New("ping failed"))

	if _, err := NewConnection(context.Background(), Config{
		Driver: "sqlmock",
		DSN:    dsn,
	}); err == nil {
		t.Fatalf("expected error on ping failure")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestConnection_CloseNilSafe(t *testing.T) {
	var conn *Connection
	if err := conn.Close(); err != nil {
		t.Fatalf("expected nil error for nil receiver, got %v", err)
	}

	empty := &Connection{}
	if err := empty.Close(); err != nil {
		t.Fatalf("expected nil error for empty connection, got %v", err)
	}
}
