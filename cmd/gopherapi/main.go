package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/friendsofgo/gopherapi/cmd/sample-data"
	gopher "github.com/friendsofgo/gopherapi/pkg"
	"github.com/friendsofgo/gopherapi/pkg/adding"
	"github.com/friendsofgo/gopherapi/pkg/fetching"
	"github.com/friendsofgo/gopherapi/pkg/log/logrus"
	"github.com/friendsofgo/gopherapi/pkg/modifying"
	"github.com/friendsofgo/gopherapi/pkg/removing"
	"github.com/friendsofgo/gopherapi/pkg/server"
	"github.com/friendsofgo/gopherapi/pkg/storage/cockroach"
	"github.com/friendsofgo/gopherapi/pkg/storage/inmem"
	"github.com/friendsofgo/gopherapi/pkg/tracer"
	_ "github.com/joho/godotenv/autoload"
	"github.com/openzipkin/zipkin-go"
)

func main() {

	var (
		hostName, _     = os.Hostname()
		defaultServerID = fmt.Sprintf("%s-%s", os.Getenv("GOPHERAPI_NAME"), hostName)
		defaultHost     = os.Getenv("GOPHERAPI_SERVER_HOST")
		defaultPort, _  = strconv.Atoi(os.Getenv("GOPHERAPI_SERVER_PORT"))

		zipkinURL = os.Getenv("ZIPKIN_ENDPOINT")
	)

	host := flag.String("host", defaultHost, "define host of the server")
	port := flag.Int("port", defaultPort, "define port of the server")
	serverID := flag.String("server-id", defaultServerID, "define server identifier")
	withData := flag.Bool("withData", false, "initialize the api with some gophers")
	withTrace := flag.Bool("withTrace", false, "initialize the api with tracing")
	cockroachDB := flag.Bool("cockroach", false, "initialize the api using CrockoachDB as db engine")
	flag.Parse()

	var gophers map[string]gopher.Gopher
	if *withData {
		gophers = sample.Gophers
	}

	logger := logrus.NewLogger()
	trc := tracer.NewNoopTracer()
	if *withTrace {
		var err error
		trc, err = tracer.NewTracer(*serverID, zipkinURL)
		if err != nil {
			log.Fatal(err)
		}
	}

	repo := inmem.NewRepository(gophers, trc)
	if *cockroachDB {
		repo = newCockroachRepository(trc)
	}
	fetchingService := fetching.NewService(repo, logger)
	addingService := adding.NewService(repo)
	modifyingService := modifying.NewService(repo)
	removingService := removing.NewService(repo)

	httpAddr := fmt.Sprintf("%s:%d", *host, *port)

	s := server.New(
		*serverID,
		trc,
		fetchingService,
		addingService,
		modifyingService,
		removingService,
	)

	fmt.Println("The gopher server is on tap now:", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, s.Router()))
}

func newCockroachRepository(trc *zipkin.Tracer) gopher.Repository {
	cockroachAddr := os.Getenv("COCKROACH_ADDR")
	cockroachDBName := os.Getenv("COCKROACH_DB")

	cockroachConn, err := cockroach.NewConn(cockroachAddr, cockroachDBName)
	if err != nil {
		log.Fatal(err)
	}
	return cockroach.NewRepository(cockroachConn, trc)
}
