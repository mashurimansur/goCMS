package database

import (
	"context"
	"crypto/tls"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Config supplies database connectivity details. All fields are optional for now.
type Config struct {
	Driver            string
	DSN               string
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetime   time.Duration
	ConnectionTimeout time.Duration
}

// Connection wraps the sql.DB pointer to keep infrastructure details isolated.
type Connection struct {
	DB *sql.DB
}

// NewConnection establishes a SQL connection when configuration is provided.
// When driver or DSN are empty the function returns a no-op connection so the
// rest of the application can continue to run without a database.
func NewConnection(ctx context.Context, cfg Config) (*Connection, error) {
	if cfg.Driver == "" || cfg.DSN == "" {
		return &Connection{}, nil
	}

	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "gateway01.ap-southeast-1.prod.aws.tidbcloud.com",
	})

	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, err
	}

	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}

	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}

	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	if cfg.ConnectionTimeout <= 0 {
		cfg.ConnectionTimeout = 5 * time.Second
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, cfg.ConnectionTimeout)
	defer cancel()

	if err := db.PingContext(timeoutCtx); err != nil {
		return nil, err
	}

	return &Connection{DB: db}, nil
}

// Close shuts down the sql.DB when it was initialized.
func (c *Connection) Close() error {
	if c == nil || c.DB == nil {
		return nil
	}

	return c.DB.Close()
}
