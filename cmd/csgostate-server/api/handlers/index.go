package handlers

import (
	"fmt"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"github.com/roessland/csgostate/cmd/csgostate-server/sessions"
	"html/template"
	"log"
	"net/http"
)

func Index(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := app.SessionStore.Get(r)
		if err != nil {
			_, _ = fmt.Fprintln(w, "unable to get session")
			return
		}

		tmpl, err := template.ParseFS(templates, "templates/index.tmpl.html")
		if err != nil {
			_, _ = fmt.Fprintln(w, "unable to load template: ", err)
			return
		}

		var lastStateRawJson string
		if sess.IsLoggedIn() {
			lastState, err := app.StateRepo.GetLatest(sess.SteamID())
			if err != nil {
				log.Print("error getting latest state for steamid ", sess.SteamID(), "err ", err)
			} else {
				if lastState != nil && lastState.RawJson != nil {
					lastStateRawJson = string(lastState.RawJson)
				}
			}
		}

		err = tmpl.Execute(w, struct {
			Sess      *sessions.Session
			LastState string
		}{
			Sess: sess,
			LastState: lastStateRawJson,
		})
		if err != nil {
			_, _ = fmt.Fprintln(w, "unable to execute template: ", err)
			return
		}
	}
}
