package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetApiHealth() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			log.Print(err)
		}
	}
}
