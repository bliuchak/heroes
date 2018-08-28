package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bliuchak/heroes/internal/storage"
	"github.com/stretchr/testify/mock"
)

type mockStorager struct {
	mock.Mock
}

func (ms *mockStorager) Status() (string, error) {
	return "", nil
}

func (ms *mockStorager) GetHeroes() ([]storage.Hero, error) {
	return []storage.Hero{}, nil
}

func (ms *mockStorager) GetHero(name string) (storage.Hero, error) {
	return storage.Hero{
		ID:   "1",
		Name: "Batman",
	}, nil
}

func (ms *mockStorager) CreateHero(id, name string) error {
	return nil
}

func (ms *mockStorager) DeleteHero(id string) error {
	return nil
}

func TestGetHeroHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/hero/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	s := new(mockStorager)
	hh := HeroHandler{}
	hh.SetStorage(s)
	handler := http.HandlerFunc(hh.GetHeroHandler)

	// Our handlers satisfy http.HeroHandler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	log.Println(rr.Code, rr.Body.String())

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"id":"1","name":"Batman"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
