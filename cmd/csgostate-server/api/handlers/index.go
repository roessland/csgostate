package handlers

import (
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"github.com/roessland/csgostate/cmd/csgostate-server/sessions"
	"html/template"
	"net/http"
)

func GetIndex(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := app.SessionStore.Get(r)
		if err != nil {
			app.Log.Errorw("unable to get session", "err", err)
			http.Error(w, "unable to get session", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFS(templates, "templates/index.tmpl.html")
		if err != nil {
			app.Log.Errorw("error loading template", "err", err)
			http.Error(w, "couldn't render page ", http.StatusInternalServerError)
			return
		}

		var lastStateRawJson string
		if sess.IsLoggedIn() {
			lastState, err := app.StateRepo.GetLatest(sess.SteamID())
			if err != nil {
				app.Log.Errorw("error getting latest state for steamid",
					"steamid", sess.SteamID(),
					"err", err)
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
			Sess:      sess,
			LastState: lastStateRawJson,
		})
		if err != nil {
			app.Log.Errorw("error executing template", "err", err)
			http.Error(w, "couldn't render page", http.StatusInternalServerError)
			return
		}
	}
}
