package sql

import (
	"context"
	"database/sql"
	"errors"

	domain "github.com/mashurimansur/goCMS/internal/domain/person"
)

// Repository reads person data from a SQL database.
type Repository struct {
	db *sql.DB
}

// New creates a SQL-backed person repository.
func New(db *sql.DB) (domain.Repository, error) {
	if db == nil {
		return nil, errors.New("sql repository requires a non-nil db")
	}

	return &Repository{db: db}, nil
}

// GetDefault fetches the first person row. Table schema is intentionally minimal.
func (r *Repository) GetDefault(ctx context.Context) (domain.Person, error) {
	const query = `SELECT name FROM persons LIMIT 1`

	var person domain.Person
	err := r.db.QueryRowContext(ctx, query).Scan(&person.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return person, errors.New("person not found")
		}
		return person, err
	}

	return person, nil
}
