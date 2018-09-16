package middleware

import (
	"net/http"
)

// HTTPLogger middleware to log http request
func (md *Middleware) HTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		md.Logger.Info().Str("method", r.Method).Str("url", r.RequestURI).Msg("Middleware call")
		next.ServeHTTP(w, r)
	})
}
