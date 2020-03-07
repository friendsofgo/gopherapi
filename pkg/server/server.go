package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"

	"github.com/friendsofgo/gopherapi/pkg/adding"
	"github.com/friendsofgo/gopherapi/pkg/fetching"
	"github.com/friendsofgo/gopherapi/pkg/modifying"
	"github.com/friendsofgo/gopherapi/pkg/removing"

	"github.com/gorilla/mux"
)

// server all server necessary dependencies
type server struct {
	serverID string

	tracer *zipkin.Tracer

	router http.Handler

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
	serverID string,
	tracer *zipkin.Tracer,
	fS fetching.Service,
	aS adding.Service,
	mS modifying.Service,
	rS removing.Service,
) Server {
	a := &server{
		serverID:  serverID,
		tracer:    tracer,
		fetching:  fS,
		adding:    aS,
		modifying: mS,
		removing:  rS}
	router(a)

	return a
}

func router(s *server) {
	r := mux.NewRouter()

	r.Use(
		zipkinhttp.NewServerMiddleware(s.tracer),
		newServerMiddleware(s.serverID),
	)

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
	gophers, _ := s.fetching.FetchGophers(r.Context())

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(gophers)

}

// FetchGopher return a gopher by ID
func (s *server) FetchGopher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gopher := s.fetching.FetchGopherByID(r.Context(), vars["ID"])
	w.Header().Set("Content-Type", "application/json")
	if gopher == nil {
		w.WriteHeader(http.StatusNotFound)
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
	if err := s.adding.AddGopher(r.Context(), g.ID, g.Name, g.Image, g.Age); err != nil {
		log.Println(err)
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
	if err := s.modifying.ModifyGopher(r.Context(), vars["ID"], g.Name, g.Image, g.Age); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("Can't modify a gopher")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveGopher remove a gopher
func (s *server) RemoveGopher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_ = s.removing.RemoveGopher(r.Context(), vars["ID"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

}
