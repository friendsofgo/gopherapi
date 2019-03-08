package server

import (
	"encoding/json"
	"net/http"

	gopher "github.com/friendsofgo/gopherapi/pkg"

	"github.com/gorilla/mux"
)

type api struct {
	router     http.Handler
	repository gopher.GopherRepository
}

// Server representation of gopher server
type Server interface {
	Router() http.Handler
	FetchGophers(w http.ResponseWriter, r *http.Request)
	FetchGopher(w http.ResponseWriter, r *http.Request)
}

// New initialize the server
func New(repo gopher.GopherRepository) Server {
	a := &api{repository: repo}

	r := mux.NewRouter()
	r.HandleFunc("/gophers", a.FetchGophers).Methods(http.MethodGet)
	r.HandleFunc("/gophers/{ID:[a-zA-Z0-9_]+}", a.FetchGopher).Methods(http.MethodGet)

	a.router = r
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}

// FetchGophers return a list of all gophers
func (a *api) FetchGophers(w http.ResponseWriter, r *http.Request) {
	gophers, _ := a.repository.FetchGophers()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gophers)
}

// FetchGopher return a gopher by ID
func (a *api) FetchGopher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gopher, err := a.repository.FetchGopherByID(vars["ID"])
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusNotFound) // We use not found for simplicity
		json.NewEncoder(w).Encode("Gopher Not found")
		return
	}

	json.NewEncoder(w).Encode(gopher)
}
