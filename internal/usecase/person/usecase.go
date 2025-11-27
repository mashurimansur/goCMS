package person

import (
	"context"
	"errors"

	domain "github.com/mashurimansur/goCMS/internal/domain/person"
)

// UseCase provides person-related business logic.
type UseCase interface {
	GetDefaultPerson(ctx context.Context) (domain.Person, error)
}

type service struct {
	repo domain.Repository
}

// New wires a person use case with its repository dependency.
func New(repo domain.Repository) UseCase {
	return &service{repo: repo}
}

// GetDefaultPerson fetches the default person data from the repository.
func (uc *service) GetDefaultPerson(ctx context.Context) (domain.Person, error) {
	if uc.repo == nil {
		return domain.Person{}, errors.New("person repository is not configured")
	}

	return uc.repo.GetDefault(ctx)
}
