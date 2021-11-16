package handlers

import (
	"encoding/json"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"github.com/roessland/csgostate/csgostate"
	"io/ioutil"
	"net/http"
)

func PostApiPush(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			app.Log.Errorw("error reading body", "err", err)
			http.Error(w, "error reading body", http.StatusInternalServerError)
			return
		}

		// Decode body
		var state csgostate.State
		state.RawJson = body
		err = json.Unmarshal(body, &state)
		if err != nil {
			app.Log.Errorw("error unmarshalling body", "err", err)
			http.Error(w, "error parsing body", http.StatusInternalServerError)
			return
		}

		// Save event to DB
		err = app.StateRepo.Push(&state)
		if err != nil {
			app.Log.Errorw("error storing state", "err", err)
			http.Error(w, "error storing state", http.StatusInternalServerError)
			return
		}

		// Feed to player events extractor
		err = app.PlayerEventsExtractor.Feed(&state)
		if err != nil {
			app.Log.Errorw("error feeding PlayerEventsExtractor", "err", err)
			http.Error(w, "error parsing state", http.StatusInternalServerError)
			return
		}

		// Update PlayerRepo
		app.PlayerRepo.Update(&state)

		// Success
		w.Header().Set("Content-Type", "text/html")
	}
}
