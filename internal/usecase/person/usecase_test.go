package person

import (
	"context"
	"testing"

	domain "github.com/mashurimansur/goCMS/internal/domain/person"
)

type stubRepository struct {
	person domain.Person
	err    error
}

func (s stubRepository) GetDefault(ctx context.Context) (domain.Person, error) {
	return s.person, s.err
}

func TestGetDefaultPerson_NilRepository(t *testing.T) {
	uc := New(nil)

	_, err := uc.GetDefaultPerson(context.Background())
	if err == nil {
		t.Fatalf("expected error when repository is nil")
	}
}

func TestGetDefaultPerson_Success(t *testing.T) {
	expected := domain.Person{Name: "Jane"}
	uc := New(stubRepository{person: expected})

	person, err := uc.GetDefaultPerson(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if person != expected {
		t.Fatalf("unexpected person returned: %+v", person)
	}
}
