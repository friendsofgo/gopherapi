package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/friendsofgo/gopherapi/pkg/adding"
	"github.com/friendsofgo/gopherapi/pkg/fetching"
	"github.com/friendsofgo/gopherapi/pkg/modifying"
	"github.com/friendsofgo/gopherapi/pkg/removing"

	"github.com/gorilla/mux"
)

// server all server necessary dependencies
type server struct {
	serverName string
	httpAddr   string

	router    http.Handler
	fetching  fetching.Service
	adding    adding.Service
	modifying modifying.Service
	removing  removing.Service
}

// Server representation of gopher server
type Server interface {
	Router() http.Handler
	FetchGophers(w http.ResponseWriter, r *http.Request)
	FetchGopher(w http.ResponseWriter, r *http.Request)
	AddGopher(w http.ResponseWriter, r *http.Request)
	ModifyGopher(w http.ResponseWriter, r *http.Request)
	RemoveGopher(w http.ResponseWriter, r *http.Request)
}

// New initialize the server
func New(
	serverName string,
	httpAddr string,
	fS fetching.Service,
	aS adding.Service,
	mS modifying.Service,
	rS removing.Service) Server {
	a := &server{
		serverName: serverName,
		httpAddr:   httpAddr,
		fetching:   fS,
		adding:     aS,
		modifying:  mS,
		removing:   rS}
	router(a)

	return a
}

func (s *server) createContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, contextKeyServerName, s.serverName)
	ctx = context.WithValue(ctx, contextKeyHttpAddr, s.httpAddr)

	return ctx
}

func router(s *server) {
	r := mux.NewRouter()
	r.HandleFunc("/gophers", s.FetchGophers).Methods(http.MethodGet)
	r.HandleFunc("/gophers/{ID:[a-zA-Z0-9_]+}", s.FetchGopher).Methods(http.MethodGet)
	r.HandleFunc("/gophers", s.AddGopher).Methods(http.MethodPost)
	r.HandleFunc("/gophers/{ID:[a-zA-Z0-9_]+}", s.ModifyGopher).Methods(http.MethodPut)
	r.HandleFunc("/gophers/{ID:[a-zA-Z0-9_]+}", s.RemoveGopher).Methods(http.MethodDelete)

	s.router = r
}

func (s *server) Router() http.Handler {
	return s.router
}

// FetchGophers return a list of all gophers
func (s *server) FetchGophers(w http.ResponseWriter, r *http.Request) {
	ctx := s.createContext(r.Context())

	gophers, _ := s.fetching.FetchGophers(ctx)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(gophers)

}

// FetchGopher return a gopher by ID
func (s *server) FetchGopher(w http.ResponseWriter, r *http.Request) {
	ctx := s.createContext(r.Context())

	vars := mux.Vars(r)
	gopher, err := s.fetching.FetchGopherByID(ctx, vars["ID"])
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusNotFound) // We use not found for simplicity
		_ = json.NewEncoder(w).Encode("Gopher Not found")
		return
	}

	_ = json.NewEncoder(w).Encode(gopher)

}

type addGopherRequest struct {
	ID    string `json:"ID"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Age   int    `json:"age"`
}

// AddGopher save a gopher
func (s *server) AddGopher(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var g addGopherRequest
	err := decoder.Decode(&g)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("Error unmarshalling request body")
		return
	}
	ctx := s.createContext(r.Context())
	if err := s.adding.AddGopher(ctx, g.ID, g.Name, g.Image, g.Age); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("Can't create a gopher")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type modifyGopherRequest struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Age   int    `json:"age"`
}

// ModifyGopher modify gopher data
func (s *server) ModifyGopher(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var g addGopherRequest
	err := decoder.Decode(&g)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("Error unmarshalling request body")
		return
	}
	vars := mux.Vars(r)
	ctx := s.createContext(r.Context())

	if err := s.modifying.ModifyGopher(ctx, vars["ID"], g.Name, g.Image, g.Age); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("Can't modify a gopher")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveGopher remove a gopher
func (s *server) RemoveGopher(w http.ResponseWriter, r *http.Request) {
	ctx := s.createContext(r.Context())
	vars := mux.Vars(r)

	_ = s.removing.RemoveGopher(ctx, vars["ID"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

}
