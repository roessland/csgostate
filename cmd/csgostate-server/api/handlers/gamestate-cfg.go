package handlers

import (
	"github.com/roessland/csgostate/cmd/csgostate-server/repos/userrepo"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"github.com/roessland/csgostate/cmd/csgostate-server/sessions"
	"net/http"
	"text/template"
)

func GetGamestateCfg(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := app.SessionStore.Get(r)
		if err != nil {
			app.Log.Errorw("unable to get session", "err", err)
			http.Error(w, "unable to get session", http.StatusInternalServerError)
			return
		}

		user, err := app.UserRepo.GetBySteamID(sess.SteamID())
		if err != nil {
			app.Log.Errorw("error retrieving user", "err", err)
			http.Error(w, "couldn't retrieve your user", http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "couldn't retrieve your user", http.StatusForbidden)
			return
		}

		tmpl, err := template.ParseFS(templates, "templates/gamestate_integration_csgostate.tmpl.cfg")
		if err != nil {
			app.Log.Errorw("unable to load template", "err", err)
			http.Error(w, "couldn't render your gamestate.cfg", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, struct {
			Sess *sessions.Session
			User *userrepo.User
		}{
			Sess: sess,
			User: user,
		})
		if err != nil {
			app.Log.Errorw("unable to execute template", "err", err)
			http.Error(w, "couldn't render your gamestate.cfg", http.StatusInternalServerError)
			return
		}
	}
}
