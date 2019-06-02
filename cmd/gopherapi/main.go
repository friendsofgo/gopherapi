package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/friendsofgo/gopherapi/pkg/adding"
	"github.com/friendsofgo/gopherapi/pkg/fetching"

	"github.com/friendsofgo/gopherapi/pkg/storage/inmem"

	sample "github.com/friendsofgo/gopherapi/cmd/sample-data"
	gopher "github.com/friendsofgo/gopherapi/pkg"

	"github.com/friendsofgo/gopherapi/pkg/server"
)

func main() {
	withData := flag.Bool("withData", false, "initialize the api with some gophers")
	flag.Parse()

	var gophers map[string]gopher.Gopher
	if *withData {
		gophers = sample.Gophers
	}

	repo := inmem.NewRepository(gophers)
	fetching := fetching.NewService(repo)
	adding := adding.NewService(repo)

	s := server.New(fetching, adding)

	fmt.Println("The gopher server is on tap now: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", s.Router()))
}
