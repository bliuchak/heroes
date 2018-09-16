package middleware

import "github.com/rs/zerolog"

// Middleware contains dependencies for middleware structure
type Middleware struct {
	Logger zerolog.Logger
}

// SetLogger sets logger
func (m *Middleware) SetLogger(logger zerolog.Logger) {
	m.Logger = logger
}
