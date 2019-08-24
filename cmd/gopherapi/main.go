package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/friendsofgo/gopherapi/cmd/sample-data"
	gopher "github.com/friendsofgo/gopherapi/pkg"
	"github.com/friendsofgo/gopherapi/pkg/adding"
	"github.com/friendsofgo/gopherapi/pkg/fetching"
	"github.com/friendsofgo/gopherapi/pkg/log"
	"github.com/friendsofgo/gopherapi/pkg/log/hooks"
	"github.com/friendsofgo/gopherapi/pkg/modifying"
	"github.com/friendsofgo/gopherapi/pkg/removing"
	"github.com/friendsofgo/gopherapi/pkg/server"
	"github.com/friendsofgo/gopherapi/pkg/storage/inmem"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	var (
		hostName, _       = os.Hostname()
		defaultServerName = fmt.Sprintf("%s-%s", os.Getenv("GOPHERAPI_NAME"), hostName)
		defaultHost       = os.Getenv("GOPHERAPI_SERVER_HOST")
		defaultPort, _    = strconv.Atoi(os.Getenv("GOPHERAPI_SERVER_PORT"))
	)

	host := flag.String("host", defaultHost, "define host of the server")
	port := flag.Int("port", defaultPort, "define port of the server")
	serverName := flag.String("server-name", defaultServerName, "define name of the server")
	withData := flag.Bool("withData", false, "initialize the api with some gophers")
	flag.Parse()

	var gophers map[string]gopher.Gopher
	if *withData {
		gophers = sample.Gophers
	}

	logger := log.NewLogger(hooks.NewServerInformationHook())

	repo := inmem.NewRepository(gophers)
	fetchingService := fetching.NewService(repo, logger)
	addingService := adding.NewService(repo)
	modifyingService := modifying.NewService(repo)
	removingService := removing.NewService(repo)

	httpAddr := fmt.Sprintf("%s:%d", *host, *port)

	s := server.New(
		*serverName,
		httpAddr,
		fetchingService,
		addingService,
		modifyingService,
		removingService,
	)

	fmt.Println("The gopher server is on tap now:", httpAddr)
	logger.Fatal(http.ListenAndServe(httpAddr, s.Router()))
}
