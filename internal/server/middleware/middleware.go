package middleware

import "github.com/rs/zerolog"

type Middleware struct {
	Logger zerolog.Logger
}

func (m *Middleware) SetLogger(logger zerolog.Logger) {
	m.Logger = logger
}
