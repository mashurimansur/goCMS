package sql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNew_WithNilDB(t *testing.T) {
	if _, err := New(nil); err == nil {
		t.Fatalf("expected error when db is nil")
	}
}

func TestGetDefault_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo, err := New(db)
	if err != nil {
		t.Fatalf("unexpected error from New: %v", err)
	}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
	mock.ExpectQuery("SELECT name FROM persons LIMIT 1").WillReturnRows(rows)

	person, err := repo.GetDefault(context.Background())
	if err != nil {
		t.Fatalf("unexpected error fetching default person: %v", err)
	}

	if person.Name != "Alice" {
		t.Fatalf("expected name Alice, got %q", person.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetDefault_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo, err := New(db)
	if err != nil {
		t.Fatalf("unexpected error from New: %v", err)
	}

	mock.ExpectQuery("SELECT name FROM persons LIMIT 1").WillReturnError(sql.ErrNoRows)

	if _, err := repo.GetDefault(context.Background()); err == nil {
		t.Fatalf("expected not found error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
