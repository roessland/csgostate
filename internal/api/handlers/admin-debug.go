package handlers

import (
	"github.com/roessland/csgostate/internal/repos/userrepo"
	"github.com/roessland/csgostate/internal/server"
	"github.com/roessland/csgostate/internal/sessions"
	"html/template"
	"net/http"
)

func GetAdminDebug(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := app.SessionStore.Get(r)
		if err != nil {
			app.Log.Errorw("unable to get session", "err", err)
			http.Error(w, "unable to get session", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFS(templates, "templates/admin-debug.tmpl.html")
		if err != nil {
			app.Log.Errorw("error loading template", "err", err)
			http.Error(w, "couldn't render page ", http.StatusInternalServerError)
			return
		}

		allUsers, err := app.UserRepo.GetAll()
		if err != nil {
			panic(err)
		}

		var lastStateRawJson string
		if sess.IsLoggedIn() {
			lastState, err := app.StateRepo.GetLatestForPlayer(sess.SteamID())
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
			AllUsers  []userrepo.User
		}{
			Sess:      sess,
			LastState: lastStateRawJson,
			AllUsers:  allUsers,
		})
		if err != nil {
			app.Log.Errorw("error executing template", "err", err)
			http.Error(w, "couldn't render page", http.StatusInternalServerError)
			return
		}
	}
}
