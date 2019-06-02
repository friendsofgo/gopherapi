package removing

import (
	gopher "github.com/friendsofgo/gopherapi/pkg"
)

// Service provides removing operations.
type Service interface {
	RemoveGopher(ID string) error
}

type service struct {
	repository gopher.Repository
}

// NewService creates a removing service with the necessary dependencies
func NewService(repository gopher.Repository) Service {
	return &service{repository}
}

// RemoveGopher remove gopher from the storage
func (s *service) RemoveGopher(ID string) error {
	return s.repository.DeleteGopher(ID)
}
