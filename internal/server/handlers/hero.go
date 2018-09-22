package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bliuchak/heroes/internal/storage"
	"github.com/gorilla/mux"
)

// HeroHandler contains hero handler data
// extend common handler
type HeroHandler struct {
	CommonHandler
}

// GetHeroesHandler handler to get all heroes
func (hh *HeroHandler) GetHeroesHandler(w http.ResponseWriter, r *http.Request) {
	hs, err := hh.Storage.GetHeroes()
	if err != nil {
		hh.Logger.Error().Err(err).Msg("Unable to get heroes")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(hs) > 0 {
		data, err := json.Marshal(hs)
		if err != nil {
			hh.Logger.Error().Err(err).Msg("Unable to marshall data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

// GetHeroHandler handler to get single hero
func (hh *HeroHandler) GetHeroHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	h, err := hh.Storage.GetHero(v["id"])
	if err != nil {
		switch err.(type) {
		case *storage.ErrHeroNotExist:
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			hh.Logger.Error().Err(err).Msg("Unable to ger var from url")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	data, err := hh.JSON.Marshal(h)
	if err != nil {
		hh.Logger.Error().Err(err).Msg("Unable to marshall data")
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// CreateHeroHandler handler to create a new hero
func (hh *HeroHandler) CreateHeroHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		hh.Logger.Error().Err(err).Msg("Unable to read body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var hero storage.Hero
	err = json.Unmarshal(b, &hero)
	if err != nil {
		hh.Logger.Error().Err(err).Msg("Unable to unmarshall data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = hh.Storage.CreateHero(hero.ID, hero.Name)
	if err != nil {
		hh.Logger.Error().Err(err).Msg("Unable to send create hero request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteHeroHandler handler to delete hero
func (hh *HeroHandler) DeleteHeroHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	err := hh.Storage.DeleteHero(v["id"])
	if err != nil {
		switch err.(type) {
		case *storage.ErrNothingToDelete:
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			hh.Logger.Error().Err(err).Msg("Unable to ger var from url")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
