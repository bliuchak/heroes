package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestHeroHandler_GetHeroesHandler(t *testing.T) {
	tests := []struct {
		name      string
		reader    io.Reader
		storage   []TestifyMockCall
		marshaler func(v interface{}) ([]byte, error)
		expected  expected
	}{
		{
			name: "should return error hh.Storage.GetHeroes",
			storage: []TestifyMockCall{
				{
					Method: "GetHeroes",
					Response: []interface{}{
						[]storage.Hero{
							storage.Hero{ID: "1", Name: "Batman"},
							storage.Hero{ID: "2", Name: "Superman"},
						},
						errors.New("get heroes error"),
					},
				},
			},
			expected: expected{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "should return error hh.Marshal",
			storage: []TestifyMockCall{
				{
					Method: "GetHeroes",
					Response: []interface{}{
						[]storage.Hero{
							storage.Hero{ID: "1", Name: "Batman"},
							storage.Hero{ID: "2", Name: "Superman"},
						},
						nil,
					},
				},
			},
			marshaler: func(v interface{}) ([]byte, error) {
				return []byte{}, errors.New("marshal error")
			},
			expected: expected{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "should return heroes",
			storage: []TestifyMockCall{
				{
					Method: "GetHeroes",
					Response: []interface{}{
						[]storage.Hero{
							storage.Hero{ID: "1", Name: "Batman"},
							storage.Hero{ID: "2", Name: "Superman"},
						},
						nil,
					},
				},
			},
			expected: expected{
				code: http.StatusOK,
			},
		},
		{
			name: "should return no heroes (empty array)",
			storage: []TestifyMockCall{
				{
					Method: "GetHeroes",
					Response: []interface{}{
						[]storage.Hero{},
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
			hh := HeroHandler{}
			rr := httptest.NewRecorder()

			s := new(stmocks.Storager)
			for _, mockCall := range tt.storage {
				s.On(mockCall.Method, mockCall.Call...).Return(mockCall.Response...)
			}

			hh.SetStorage(s)

			if tt.marshaler != nil {
				hh.Marshaler = tt.marshaler
			}

			r := httptest.NewRequest("GET", "/heroes", nil)
			hh.GetHeroesHandler(rr, r)

			if rr.Code != tt.expected.code {
				t.Errorf("handler returned unexpected response code: got %v want %v",
					rr.Code, tt.expected.code)
			}
		})
	}
}

func TestHeroHandler_CreateHeroHandler(t *testing.T) {
	tests := []struct {
		name        string
		reader      io.Reader
		storage     []TestifyMockCall
		unmarshaler func(data []byte, v interface{}) error
		expected    expected
	}{
		{
			name:   "should return error from ioutil.ReadAll",
			reader: errReader(0),
			expected: expected{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "should return error from hh.Unmarshal",
			unmarshaler: func(data []byte, v interface{}) error {
				return errors.New("unmarshal error")
			},
			expected: expected{
				code: http.StatusInternalServerError,
			},
		},
		{
			name:   "should return error hh.Storage.CreateHero",
			reader: strings.NewReader(`{"id":"1","name":"Batman"}`),
			storage: []TestifyMockCall{
				{
					Method: "CreateHero",
					Call: []interface{}{
						AnythingOfType("string"),
						AnythingOfType("string"),
					},
					Response: []interface{}{
						errors.New("create hero error"),
					},
				},
			},
			expected: expected{
				code: http.StatusInternalServerError,
			},
		},
		{
			name:   "should create hero",
			reader: strings.NewReader(`{"id":"1","name":"Batman"}`),
			storage: []TestifyMockCall{
				{
					Method: "CreateHero",
					Call: []interface{}{
						AnythingOfType("string"),
						AnythingOfType("string"),
					},
					Response: []interface{}{
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
			hh := HeroHandler{}
			rr := httptest.NewRecorder()

			reader := io.Reader(nil)
			if tt.reader != nil {
				reader = tt.reader
			}

			s := new(stmocks.Storager)
			for _, mockCall := range tt.storage {
				s.On(mockCall.Method, mockCall.Call...).Return(mockCall.Response...)
			}

			hh.SetStorage(s)

			r := httptest.NewRequest("POST", "/hero", reader)
			hh.CreateHeroHandler(rr, r)

			if tt.unmarshaler != nil {
				hh.Unmarshaler = tt.unmarshaler
			}

			if rr.Code != tt.expected.code {
				t.Errorf("handler returned unexpected response code: got %v want %v",
					rr.Code, tt.expected.code)
			}
		})
	}
}

func TestHeroHandler_GetHeroHandler(t *testing.T) {
	tests := []struct {
		name      string
		storage   []TestifyMockCall
		marshaler func(v interface{}) ([]byte, error)
		expected  expected
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
			marshaler: func(v interface{}) ([]byte, error) {
				return []byte{}, errors.New("marshal error")
			},
			expected: expected{
				code: http.StatusInternalServerError,
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

			hh := HeroHandler{}
			hh.SetStorage(s)

			// set custom marshaler to trigger error
			if tt.marshaler != nil {
				hh.Marshaler = tt.marshaler
			}

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
