package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	sample "github.com/friendsofgo/gopherapi/cmd/sample-data"
	"github.com/friendsofgo/gopherapi/internal/container"
	gopher "github.com/friendsofgo/gopherapi/pkg"
)

func main() {
	withData := flag.Bool("withData", false, "initialize the api with some gophers")
	flag.Parse()

	var gophers map[string]gopher.Gopher
	if *withData {
		gophers = sample.Gophers
	}
	s := container.InitializeServer(gophers)

	fmt.Println("The gopher server is on tap now: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", s.Router()))
}
