package handlers

import (
	"github.com/bliuchak/heroes/internal/storage"
	"github.com/rs/zerolog"
)

type CommonHandler struct {
	Storage storage.Storager
	Logger  zerolog.Logger
}

func (ch *CommonHandler) SetLogger(logger zerolog.Logger) {
	ch.Logger = logger
}

func (ch *CommonHandler) SetStorage(st storage.Storager) {
	ch.Storage = st
}
