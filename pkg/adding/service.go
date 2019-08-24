package adding

import (
	"context"

	gopher "github.com/friendsofgo/gopherapi/pkg"
)

// Service provides adding operations.
type Service interface {
	AddGopher(ctx context.Context, ID, name, image string, age int) error
}

type service struct {
	repository gopher.Repository
}

// NewService creates an adding service with the necessary dependencies
func NewService(repository gopher.Repository) Service {
	return &service{repository}
}

// AddGopher adds the given gopher to storage
func (s *service) AddGopher(ctx context.Context, ID, name, image string, age int) error {
	g := gopher.New(ID, name, image, age)
	return s.repository.CreateGopher(ctx, g)
}
