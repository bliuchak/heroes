package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	jsmocks "github.com/bliuchak/heroes/internal/server/json/mocks"
	"github.com/bliuchak/heroes/internal/storage"
	stmocks "github.com/bliuchak/heroes/internal/storage/mocks"
	. "github.com/stretchr/testify/mock"
)

// TestifyMockCall is representation of method call on mock structure
type TestifyMockCall struct {
	Method   string
	Call     []interface{}
	Response []interface{}
	Times    int
}

type expected struct {
	code int
}

func TestHeroHandler_GetHeroHandler(t *testing.T) {
	tests := []struct {
		name     string
		storage  []TestifyMockCall
		json     []TestifyMockCall
		expected expected
	}{
		{
			name: "should return storage.ErrHeroNotExist on hh.Storage.GetHero",
			storage: []TestifyMockCall{
				{
					Method: "GetHero",
					Call: []interface{}{
						AnythingOfType("string"),
					},
					Response: []interface{}{
						storage.Hero{},
						storage.NewErrHeroNotExist("dummy"),
					},
				},
			},
			json: []TestifyMockCall{
				{
					Method: "Marshal",
					Call: []interface{}{
						AnythingOfType("storage.Hero"),
					},
					Response: []interface{}{
						[]byte{},
						nil,
					},
				},
			},
			expected: expected{
				code: http.StatusNotFound,
			},
		},
		{
			name: "should return default error on hh.Storage.GetHero",
			storage: []TestifyMockCall{
				{
					Method: "GetHero",
					Call: []interface{}{
						AnythingOfType("string"),
					},
					Response: []interface{}{
						storage.Hero{},
						errors.New("dummy"),
					},
				},
			},
			json: []TestifyMockCall{
				{
					Method: "Marshal",
					Call: []interface{}{
						AnythingOfType("storage.Hero"),
					},
					Response: []interface{}{
						[]byte{},
						nil,
					},
				},
			},
			expected: expected{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "should return error on Marshal json response",
			storage: []TestifyMockCall{
				{
					Method: "GetHero",
					Call: []interface{}{
						AnythingOfType("string"),
					},
					Response: []interface{}{
						storage.Hero{},
						nil,
					},
				},
			},
			json: []TestifyMockCall{
				{
					Method: "Marshal",
					Call: []interface{}{
						AnythingOfType("storage.Hero"),
					},
					Response: []interface{}{
						[]byte{},
						errors.New("marshaler error"),
					},
				},
			},
			expected: expected{
				code: http.StatusNotImplemented,
			},
		},
		{
			name: "should return hero json",
			storage: []TestifyMockCall{
				{
					Method: "GetHero",
					Call: []interface{}{
						AnythingOfType("string"),
					},
					Response: []interface{}{
						storage.Hero{},
						nil,
					},
				},
			},
			json: []TestifyMockCall{
				{
					Method: "Marshal",
					Call: []interface{}{
						AnythingOfType("storage.Hero"),
					},
					Response: []interface{}{
						[]byte{},
						nil,
					},
				},
			},
			expected: expected{
				code: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/hero/1", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			s := new(stmocks.Storager)
			for _, mockCall := range tt.storage {
				s.On(mockCall.Method, mockCall.Call...).Return(mockCall.Response...)
			}

			j := new(jsmocks.Marshaler)
			for _, mockCall := range tt.json {
				j.On(mockCall.Method, mockCall.Call...).Return(mockCall.Response...)
			}

			hh := HeroHandler{}
			hh.SetStorage(s)
			hh.SetJSON(j)

			handler := http.HandlerFunc(hh.GetHeroHandler)

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expected.code {
				t.Errorf("handler returned unexpected response code: got %v want %v",
					rr.Code, tt.expected.code)
			}
		})
	}
}

func TestHeroHandler_DeleteHeroHandler(t *testing.T) {
	tests := []struct {
		name     string
		storage  []TestifyMockCall
		expected expected
	}{
		{
			name: "should return HeroNotExists error",
			storage: []TestifyMockCall{
				{
					Method: "DeleteHero",
					Call: []interface{}{
						AnythingOfType("string"),
					},
					Response: []interface{}{
						storage.NewErrNothingToDelete("dummy"),
					},
				},
			},
			expected: expected{
				code: http.StatusNotFound,
			},
		},
		{
			name: "should return internal server error",
			storage: []TestifyMockCall{
				{
					Method: "DeleteHero",
					Call: []interface{}{
						AnythingOfType("string"),
					},
					Response: []interface{}{
						errors.New("dummy"),
					},
				},
			},
			expected: expected{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "should successfully delete hero",
			storage: []TestifyMockCall{
				{
					Method: "DeleteHero",
					Call: []interface{}{
						AnythingOfType("string"),
					},
					Response: []interface{}{
						nil,
					},
				},
			},
			expected: expected{
				code: http.StatusNoContent,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, "/hero/1", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			s := new(stmocks.Storager)
			for _, mockCall := range tt.storage {
				s.On(mockCall.Method, mockCall.Call...).Return(mockCall.Response...)
			}

			hh := HeroHandler{}
			hh.SetStorage(s)

			handler := http.HandlerFunc(hh.DeleteHeroHandler)

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expected.code {
				t.Errorf("handler returned unexpected response code: got %v want %v",
					rr.Code, tt.expected.code)
			}
		})
	}
}
