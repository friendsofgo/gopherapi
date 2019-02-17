package gopher

// Gopher defines the properties of a gopher to be listed
type Gopher struct {
	ID    string `json:"ID"`
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
	Age   int    `json:"age,omitempty"`
}

// Repository provides access to the gopher storage
type GopherRepository interface {
	// CreateGopher saves a given gopher
	CreateGopher(g *Gopher) error
	// FetchGophers return all gophers saved in storage
	FetchGophers() ([]*Gopher, error)
	// DeleteGopher remove gopher with given ID
	DeleteGopher(ID string) error
	// UpdateGopher modify gopher with given ID and given new data
	UpdateGopher(ID string, g *Gopher) error
	// FetchGopherByID returns the gopher with given ID
	FetchGopherByID(ID string) (*Gopher, error)
}
