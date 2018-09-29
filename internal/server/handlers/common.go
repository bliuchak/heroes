package handlers

import (
	"encoding/json"

	"github.com/bliuchak/heroes/internal/storage"
	"github.com/rs/zerolog"
)

// CommonHandler strucute for basic dependencies for rest handlers
type CommonHandler struct {
	Storage     storage.Storager
	Logger      zerolog.Logger
	Marshaler   func(v interface{}) ([]byte, error)
	Unmarshaler func(data []byte, v interface{}) error
}

// SetLogger sets logger
func (ch *CommonHandler) SetLogger(logger zerolog.Logger) {
	ch.Logger = logger
}

// SetStorage sets storage
func (ch *CommonHandler) SetStorage(st storage.Storager) {
	ch.Storage = st
}

// Marshal will marshal the provided value with the Marshaller defined on ch.
// If un-set, json.Marshal will be used.
func (ch *CommonHandler) Marshal(v interface{}) ([]byte, error) {
	if ch.Marshaler == nil {
		return json.Marshal(v)
	}
	return ch.Marshaler(v)
}

// Unmarshal will unmarshal the provided value with the Unmarshaler defined on ch.
// If un-set, json.Unmarshal will be used.
func (ch *CommonHandler) Unmarshal(data []byte, v interface{}) error {
	if ch.Unmarshaler == nil {
		return json.Unmarshal(data, &v)
	}
	return ch.Unmarshaler(data, &v)
}
