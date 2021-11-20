package handlers

import (
	"encoding/json"
	"github.com/roessland/csgostate/internal/server"
	"net/http"
)

func GetApiPlayers(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		err := json.NewEncoder(w).Encode(app.PlayerRepo.GetAll())
		if err != nil {
			app.Log.Errorw("error encoding list of players", "err", err)
			http.Error(w, "error encoding body", http.StatusInternalServerError)
			return
		}
	}
}
