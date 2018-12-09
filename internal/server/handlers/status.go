package handlers

import (
	"net/http"
)

// StatusHandler contains status handler data
// extend common handler
type StatusHandler struct {
	CommonHandler
}

// StatusResponse displays status of required app dependencies (e.g. storage)
type StatusResponse struct {
	Redis string `json:"redis"`
}

// GetStatusHandler handle for application status endpoint
func (sh *StatusHandler) GetStatusHandler(w http.ResponseWriter, _ *http.Request) {
	status, err := sh.Storage.Status()
	if err != nil {
		sh.Logger.Error().Err(err).Msg("Unable to get storage status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := StatusResponse{
		Redis: status,
	}

	data, err := sh.Marshal(resp)
	if err != nil {
		sh.Logger.Error().Err(err).Msg("Unable to unmarshall data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
