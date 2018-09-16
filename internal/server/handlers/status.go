package handlers

import (
	"net/http"
)

// StatusHandler contains status handler data
// extend common handler
type StatusHandler struct {
	CommonHandler
}

// GetStatusHandler handle for application status endpoint
func (sh *StatusHandler) GetStatusHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
