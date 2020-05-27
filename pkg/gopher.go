package gopher

import (
	"context"
	"time"
)

// Gopher defines the properties of a gopher to be listed
type Gopher struct {
	ID        string     `json:"ID"`
	Name      string     `json:"name,omitempty"`
	Image     string     `json:"image,omitempty"`
	Age       int        `json:"age,omitempty"`
	CreatedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-"`
}

// New creates a gopher
func New(ID, name, image string, age int) *Gopher {
	return &Gopher{
		ID:    ID,
		Name:  name,
		Image: image,
		Age:   age,
	}
}

//Repository provides access to the gopher storage
type Repository interface {
	// CreateGopher saves a given gopher
	CreateGopher(ctx context.Context, gopher *Gopher) error
	// FetchGophers return all gophers saved in storage
	FetchGophers(ctx context.Context) ([]Gopher, error)
	// DeleteGopher remove gopher with given ID
	DeleteGopher(ctx context.Context, ID string) error
	// UpdateGopher modify gopher with given ID and given new data
	UpdateGopher(ctx context.Context, ID string, gopher Gopher) error
	// FetchGopherByID returns the gopher with given ID
	FetchGopherByID(ctx context.Context, ID string) (*Gopher, error)
}
