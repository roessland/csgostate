package handlers

import (
	"encoding/json"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"log"
	"net/http"
)

func ApiPlayers(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		err := json.NewEncoder(w).Encode(app.PlayerRepo.GetAll())
		if err != nil {
			log.Println(err)
		}
	}
}