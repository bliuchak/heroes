package server

import (
	"net/http"

	"github.com/bliuchak/heroes/internal/server/json"

	"github.com/bliuchak/heroes/internal/server/handlers"
	"github.com/bliuchak/heroes/internal/server/middleware"
)

// SetRoutes setter for basic routes
func (s *Server) SetRoutes() {
	heroHandler := handlers.HeroHandler{}
	heroHandler.SetLogger(s.Logger)
	heroHandler.SetStorage(s.Storage)
	heroHandler.SetJSON(new(json.JSON))

	statusHandler := handlers.StatusHandler{}

	s.Router.HandleFunc("/status", statusHandler.GetStatusHandler).Methods(http.MethodGet)
	s.Router.HandleFunc("/heroes", heroHandler.GetHeroesHandler).Methods(http.MethodGet)
	s.Router.HandleFunc("/hero/{id:[0-9]+}", heroHandler.GetHeroHandler).Methods(http.MethodGet)
	s.Router.HandleFunc("/hero", middleware.IsJSONValid(heroHandler.CreateHeroHandler)).Methods(http.MethodPost)
	s.Router.HandleFunc("/hero/{id:[0-9]+}", heroHandler.DeleteHeroHandler).Methods(http.MethodDelete)
}
