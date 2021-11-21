package handlers

import (
	"encoding/json"
	"github.com/roessland/csgostate/internal/server"
	"net/http"
)

func GetApiHealth(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		w.WriteHeader(400)
		return

		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			app.Log.Errorw("error encoding body", "err", err)
			http.Error(w, "error encoding body", http.StatusInternalServerError)
		}
	}
}
