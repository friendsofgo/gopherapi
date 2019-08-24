package fetching

import (
	"context"

	gopher "github.com/friendsofgo/gopherapi/pkg"
	"github.com/friendsofgo/gopherapi/pkg/log"
)

// Service provides fetching operations.
type Service interface {
	FetchGophers(ctx context.Context) ([]gopher.Gopher, error)
	FetchGopherByID(ctx context.Context, ID string) *gopher.Gopher
}

type service struct {
	repository gopher.Repository
	logger     *log.Logger
}

// NewService creates a fetching service with the necessary dependencies
func NewService(repository gopher.Repository, logger *log.Logger) Service {
	return &service{repository, logger}
}

// FetchGophers returns all gophers
func (s *service) FetchGophers(ctx context.Context) ([]gopher.Gopher, error) {
	return s.repository.FetchGophers(ctx)
}

// FetchGopherByID returns a gopher
func (s *service) FetchGopherByID(ctx context.Context, ID string) *gopher.Gopher {
	g, err := s.repository.FetchGopherByID(ctx, ID)

	// This error can be any error type of our repository
	if err != nil {
		s.logger.UnexpectedError(ctx, err)
		return nil
	}

	return g
}
