package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	sample "github.com/friendsofgo/gopher-api/cmd/sample-data"
	gopher "github.com/friendsofgo/gopher-api/pkg"
	"github.com/friendsofgo/gopher-api/pkg/storage/inmem"
)

func TestFetchGophers(t *testing.T) {
	req, err := http.NewRequest("GET", "/gophers", nil)
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}

	repo := inmem.NewGopherRepository(sample.Gophers)
	s := New(repo)

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

	var got []*gopher.Gopher
	err = json.Unmarshal(b, &got)
	if err != nil {
		t.Fatalf("could not unmarshall response %v", err)
	}

	expected := len(sample.Gophers)

	if len(got) != expected {
		t.Errorf("expected %d gophers, got: %d gopher", sample.Gophers, got)
	}
}
