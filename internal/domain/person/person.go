package person

import "context"

// Person models the data returned by the API and stored in the domain layer.
type Person struct {
	Name string `json:"name"`
}

// Repository abstracts the data source that stores person information.
type Repository interface {
	GetDefault(ctx context.Context) (Person, error)
}
