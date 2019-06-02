package fetching

import gopher "github.com/friendsofgo/gopherapi/pkg"

// Service provides fetching operations.
type Service interface {
	FetchGophers() ([]gopher.Gopher, error)
	FetchGopherByID(ID string) (*gopher.Gopher, error)
}

type service struct {
	repository gopher.Repository
}

// NewService creates a fetching service with the necessary dependencies
func NewService(repository gopher.Repository) Service {
	return &service{repository}
}

// FetchGophers returns all gophers
func (s *service) FetchGophers() ([]gopher.Gopher, error) {
	return s.repository.FetchGophers()
}

// FetchGopherByID returns a gopher
func (s *service) FetchGopherByID(ID string) (*gopher.Gopher, error) {
	return s.repository.FetchGopherByID(ID)
}
