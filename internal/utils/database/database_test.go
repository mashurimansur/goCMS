package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"
)

// --- Fake drivers to avoid real DB connections ---

// successDriver will always open a connection whose Ping succeeds.
type successDriver struct{}

func (d *successDriver) Open(name string) (driver.Conn, error) {
	return &successConn{}, nil
}

type successConn struct{}

func (c *successConn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("not implemented")
}

func (c *successConn) Close() error { return nil }

func (c *successConn) Begin() (driver.Tx, error) {
	return nil, errors.New("not implemented")
}

// Implement driver.Pinger so db.PingContext succeeds.
func (c *successConn) Ping(ctx context.Context) error {
	return nil
}

// pingFailDriver returns a connection whose Ping always fails.
type pingFailDriver struct{}

func (d *pingFailDriver) Open(name string) (driver.Conn, error) {
	return &pingFailConn{}, nil
}

type pingFailConn struct{}

func (c *pingFailConn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("not implemented")
}

func (c *pingFailConn) Close() error { return nil }

func (c *pingFailConn) Begin() (driver.Tx, error) {
	return nil, errors.New("not implemented")
}

func (c *pingFailConn) Ping(ctx context.Context) error {
	return errors.New("ping failed")
}

// Register fake drivers for tests.
func init() {
	sql.Register("successdb", &successDriver{})
	sql.Register("pingfaildb", &pingFailDriver{})
}

// --- Tests ---

func TestNewConnection_NoDriverReturnsNoOpConnection(t *testing.T) {
	ctx := context.Background()
	cfg := Config{} // Driver empty

	conn, err := NewConnection(ctx, cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if conn == nil {
		t.Fatalf("expected non-nil Connection, got nil")
	}

	if conn.DB != nil {
		t.Fatalf("expected nil DB on no-op connection, got %v", conn.DB)
	}
}

func TestNewConnection_OpenError(t *testing.T) {
	ctx := context.Background()
	cfg := Config{
		Driver: "unknown-driver", // not registered -> sql.Open should fail
	}

	conn, err := NewConnection(ctx, cfg)
	if err == nil {
		t.Fatalf("expected error when opening with unknown driver, got nil")
	}

	if conn != nil {
		t.Fatalf("expected nil Connection on error, got %+v", conn)
	}
}

func TestNewConnection_Success(t *testing.T) {
	ctx := context.Background()
	cfg := Config{
		Driver:          "successdb",
		Username:        "user",
		Password:        "pass",
		Protocol:        "tcp",
		Address:         "localhost",
		Port:            "3306",
		DatabaseName:    "testdb",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Minute,
		// ConnectionTimeout > 0 so it uses provided value, not default
		ConnectionTimeout: 2 * time.Second,
	}

	conn, err := NewConnection(ctx, cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if conn == nil || conn.DB == nil {
		t.Fatalf("expected valid DB connection, got %+v", conn)
	}

	// Ensure Close works when DB is initialized.
	if err := conn.Close(); err != nil {
		t.Fatalf("expected no error closing DB, got %v", err)
	}
}

func TestNewConnection_PingErrorAndDefaultTimeout(t *testing.T) {
	ctx := context.Background()
	cfg := Config{
		Driver:       "pingfaildb",
		Username:     "user",
		Password:     "pass",
		Protocol:     "tcp",
		Address:      "localhost",
		Port:         "3306",
		DatabaseName: "testdb",
		// ConnectionTimeout <= 0 triggers default (5 * time.Second)
		ConnectionTimeout: 0,
	}

	conn, err := NewConnection(ctx, cfg)
	if err == nil {
		t.Fatalf("expected error from PingContext, got nil")
	}

	if conn != nil {
		t.Fatalf("expected nil Connection on ping error, got %+v", conn)
	}
}

func TestConnectionClose_NoOp(t *testing.T) {
	// Case 1: nil receiver
	var conn *Connection
	if err := conn.Close(); err != nil {
		t.Fatalf("expected nil error when closing nil connection, got %v", err)
	}

	// Case 2: non-nil Connection but nil DB
	conn = &Connection{}
	if err := conn.Close(); err != nil {
		t.Fatalf("expected nil error when closing connection with nil DB, got %v", err)
	}
}
