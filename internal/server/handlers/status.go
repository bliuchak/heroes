package handlers

import (
	"net/http"
)

type StatusHandler struct {
	CommonHandler
}

func (sh *StatusHandler) GetStatusHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
