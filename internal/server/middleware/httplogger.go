package middleware

import (
	"net/http"
)

// HttpLogger ...
func (md *Middleware) HttpLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		md.Logger.Info().Str("method", r.Method).Str("url", r.RequestURI).Msg("Middleware call")
		next.ServeHTTP(w, r)
	})
}
