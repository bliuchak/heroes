package handlers

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	jsmocks "github.com/bliuchak/heroes/internal/server/json/mocks"
	"github.com/bliuchak/heroes/internal/storage"
	stmocks "github.com/bliuchak/heroes/internal/storage/mocks"
	. "github.com/stretchr/testify/mock"
)

func TestGetHeroHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/hero/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	s := new(stmocks.Storager)
	s.On("GetHero", AnythingOfType("string")).Return(storage.Hero{}, nil)

	j := new(jsmocks.Marshaler)
	j.On("Marshal", AnythingOfType("storage.Hero")).Return([]byte{}, errors.New("marshaler error"))

	hh := HeroHandler{}
	hh.SetStorage(s)
	hh.SetJSON(j)

	handler := http.HandlerFunc(hh.GetHeroHandler)

	// Our handlers satisfy http.HeroHandler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	log.Println(rr.Code, rr.Body.String())

	// Check the status code is what we expect.
	// if status := rr.Code; status != http.StatusOK {
	// 	t.Errorf("handler returned wrong status code: got %v want %v",
	// 		status, http.StatusOK)
	// }

	// // Check the response body is what we expect.
	// expected := `{"id":"1","name":"Batman"}`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}
