package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/friendsofgo/gopherapi/pkg/adding"
	"github.com/friendsofgo/gopherapi/pkg/fetching"
	"github.com/friendsofgo/gopherapi/pkg/modifying"

	sample "github.com/friendsofgo/gopherapi/cmd/sample-data"
	gopher "github.com/friendsofgo/gopherapi/pkg"
	"github.com/friendsofgo/gopherapi/pkg/storage/inmem"
)

func TestFetchGophers(t *testing.T) {
	req, err := http.NewRequest("GET", "/gophers", nil)
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}

	s := buildServer()

	rec := httptest.NewRecorder()

	s.FetchGophers(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %d, got: %d", http.StatusOK, res.StatusCode)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	var got []gopher.Gopher
	err = json.Unmarshal(b, &got)
	if err != nil {
		t.Fatalf("could not unmarshall response %v", err)
	}

	expected := len(sample.Gophers)

	if len(got) != expected {
		t.Errorf("expected %v gophers, got: %v gopher", sample.Gophers, got)
	}
}

func TestFetchGopher(t *testing.T) {

	testData := []struct {
		name   string
		g      *gopher.Gopher
		status int
		err    string
	}{
		{name: "gopher found", g: gopherSample(), status: http.StatusOK},
		{name: "gopher not found", g: &gopher.Gopher{ID: "123"}, status: http.StatusNotFound, err: "Gopher Not found"},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			uri := fmt.Sprintf("/gophers/%s", tt.g.ID)
			req, err := http.NewRequest("GET", uri, nil)
			if err != nil {
				t.Fatalf("could not created request: %v", err)
			}

			s := buildServer()

			rec := httptest.NewRecorder()
			s.Router().ServeHTTP(rec, req)

			res := rec.Result()

			defer res.Body.Close()
			if tt.status != res.StatusCode {
				t.Errorf("expected %d, got: %d", tt.status, res.StatusCode)
			}
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if tt.err == "" {
				var got *gopher.Gopher
				err = json.Unmarshal(b, &got)
				if err != nil {
					t.Fatalf("could not unmarshall response %v", err)
				}

				if *got != *tt.g {
					t.Fatalf("expected %v, got: %v", tt.g, got)
				}
			}
		})
	}

}

func TestAddGopher(t *testing.T) {
	bodyJSON := []byte(`{
        "ID": "01DCBP0R0MSNZY975ZQF1DCQCH",
        "name": "Eustaqio",
        "image": "https://storage.googleapis.com/gopherizeme.appspot.com/gophers/f73f25d73c06cc81c482821391a85c4b7dd34ba5.png",
        "age": 99
    }`)
	req, err := http.NewRequest("POST", "/gophers", bytes.NewBuffer(bodyJSON))
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}
	s := buildServer()
	rec := httptest.NewRecorder()

	s.AddGopher(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected %d, got: %d", http.StatusCreated, res.StatusCode)
	}
}

func gopherSample() *gopher.Gopher {
	return &gopher.Gopher{
		ID:    "01D3XZ3ZHCP3KG9VT4FGAD8KDR",
		Name:  "Jenny",
		Age:   18,
		Image: "https://storage.googleapis.com/gopherizeme.appspot.com/gophers/0ceb2c10fc0c30575c18ff1defa1ffd41501bc62.png",
	}
}

func buildServer() Server {
	repo := inmem.NewRepository(sample.Gophers)
	fetching := fetching.NewService(repo)
	adding := adding.NewService(repo)
	modifying := modifying.NewService(repo)

	return New(fetching, adding, modifying)
}
