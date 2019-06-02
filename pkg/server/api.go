package server

import (
	"encoding/json"
	"net/http"

	"github.com/friendsofgo/gopherapi/pkg/adding"
	"github.com/friendsofgo/gopherapi/pkg/fetching"
	"github.com/friendsofgo/gopherapi/pkg/modifying"

	"github.com/gorilla/mux"
)

type api struct {
	router    http.Handler
	fetching  fetching.Service
	adding    adding.Service
	modifying modifying.Service
}

// Server representation of gopher server
type Server interface {
	Router() http.Handler
	FetchGophers(w http.ResponseWriter, r *http.Request)
	FetchGopher(w http.ResponseWriter, r *http.Request)
	AddGopher(w http.ResponseWriter, r *http.Request)
	ModifyGopher(w http.ResponseWriter, r *http.Request)
}

// New initialize the server
func New(fS fetching.Service, aS adding.Service, mS modifying.Service) Server {
	a := &api{fetching: fS, adding: aS, modifying: mS}

	r := mux.NewRouter()
	r.HandleFunc("/gophers", a.FetchGophers).Methods(http.MethodGet)
	r.HandleFunc("/gophers/{ID:[a-zA-Z0-9_]+}", a.FetchGopher).Methods(http.MethodGet)
	r.HandleFunc("/gophers", a.AddGopher).Methods(http.MethodPost)
	r.HandleFunc("/gophers/{ID:[a-zA-Z0-9_]+}", a.ModifyGopher).Methods(http.MethodPut)

	a.router = r
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}

// FetchGophers return a list of all gophers
func (a *api) FetchGophers(w http.ResponseWriter, r *http.Request) {
	gophers, _ := a.fetching.FetchGophers()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gophers)
}

// FetchGopher return a gopher by ID
func (a *api) FetchGopher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gopher, err := a.fetching.FetchGopherByID(vars["ID"])
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusNotFound) // We use not found for simplicity
		json.NewEncoder(w).Encode("Gopher Not found")
		return
	}

	json.NewEncoder(w).Encode(gopher)
}

type addGopherRequest struct {
	ID    string `json:"ID"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Age   int    `json:"age"`
}

// AddGopher save a gopher
func (a *api) AddGopher(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var g addGopherRequest
	err := decoder.Decode(&g)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error unmarshalling request body")
		return
	}

	if err := a.adding.AddGopher(g.ID, g.Name, g.Image, g.Age); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Can't create a gopher")
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
func (a *api) ModifyGopher(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var g addGopherRequest
	err := decoder.Decode(&g)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error unmarshalling request body")
		return
	}
	vars := mux.Vars(r)
	if err := a.modifying.ModifyGopher(vars["ID"], g.Name, g.Image, g.Age); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Can't modify a gopher")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
