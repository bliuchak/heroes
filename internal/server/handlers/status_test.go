package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	stmocks "github.com/bliuchak/heroes/internal/storage/mocks"
)

func TestStatusHandler_GetStatusHandler(t *testing.T) {
	tests := []struct {
		name      string
		reader    io.Reader
		storage   []TestifyMockCall
		marshaler func(v interface{}) ([]byte, error)
		expected  expected
	}{
		{
			name: "should return response that storage is in OK state",
			storage: []TestifyMockCall{
				{
					Method: "Status",
					Response: []interface{}{
						"PONG",
						nil,
					},
				},
			},
			expected: expected{
				code: http.StatusOK,
			},
		},
		{
			name: "should return error from method sh.Storage.Status()",
			storage: []TestifyMockCall{
				{
					Method: "Status",
					Response: []interface{}{
						"",
						errors.New("dummy error"),
					},
				},
			},
			expected: expected{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "should return error on Marshaler",
			storage: []TestifyMockCall{
				{
					Method: "Status",
					Response: []interface{}{
						"PONG",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sh := StatusHandler{}
			rr := httptest.NewRecorder()

			s := new(stmocks.Storager)
			for _, mockCall := range tt.storage {
				s.On(mockCall.Method, mockCall.Call...).Return(mockCall.Response...)
			}

			sh.SetStorage(s)

			if tt.marshaler != nil {
				sh.Marshaler = tt.marshaler
			}

			r := httptest.NewRequest("GET", "/status", nil)
			sh.GetStatusHandler(rr, r)

			if rr.Code != tt.expected.code {
				t.Errorf("handler returned unexpected response code: got %v want %v",
					rr.Code, tt.expected.code)
			}
		})
	}
}
