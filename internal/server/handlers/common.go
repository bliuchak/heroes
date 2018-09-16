package handlers

import (
	"github.com/bliuchak/heroes/internal/storage"
	"github.com/rs/zerolog"
)

// CommonHandler strucute for basic dependencies for rest handlers
type CommonHandler struct {
	Storage storage.Storager
	Logger  zerolog.Logger
}

// SetLogger sets logger
func (ch *CommonHandler) SetLogger(logger zerolog.Logger) {
	ch.Logger = logger
}

// SetStorage sets storage
func (ch *CommonHandler) SetStorage(st storage.Storager) {
	ch.Storage = st
}
