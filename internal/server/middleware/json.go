package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bliuchak/heroes/internal/storage"
)

type errInvalidJSON struct {
	Message string `json:"message"`
}

func IsJSONValid(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			jsonInvalid := errInvalidJSON{Message: "unable to read body from request"}
			bytes, _ := json.Marshal(jsonInvalid)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(bytes)
			return
		}
		defer r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewReader(reqBody))

		var hero storage.Hero
		err = json.Unmarshal(reqBody, &hero)
		if err != nil {
			jsonInvalid := errInvalidJSON{Message: "unable to unmarshall body to structure"}
			bytes, _ := json.Marshal(jsonInvalid)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(bytes)
			return
		}

		if hero.IsValid() {
			next.ServeHTTP(w, r)
		} else {
			jsonInvalid := errInvalidJSON{Message: "data in json are not valid"}
			bytes, _ := json.Marshal(jsonInvalid)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(bytes)
			return
		}
	})
}
