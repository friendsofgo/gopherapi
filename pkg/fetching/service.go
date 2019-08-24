package fetching

import (
	"context"

	gopher "github.com/friendsofgo/gopherapi/pkg"
)

// Service provides fetching operations.
type Service interface {
	FetchGophers(ctx context.Context) ([]gopher.Gopher, error)
	FetchGopherByID(ctx context.Context, ID string) (*gopher.Gopher, error)
}

type service struct {
	repository gopher.Repository
}

// NewService creates a fetching service with the necessary dependencies
func NewService(repository gopher.Repository) Service {
	return &service{repository}
}

// FetchGophers returns all gophers
func (s *service) FetchGophers(ctx context.Context) ([]gopher.Gopher, error) {
	return s.repository.FetchGophers(ctx)
}

// FetchGopherByID returns a gopher
func (s *service) FetchGopherByID(ctx context.Context, ID string) (*gopher.Gopher, error) {
	return s.repository.FetchGopherByID(ctx, ID)
}
