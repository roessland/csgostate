package handlers

import (
	"fmt"
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
			_, _ = fmt.Fprintln(w, "unable to get session")
			return
		}

		user, err := app.UserRepo.GetBySteamID(sess.SteamID())
		if err != nil {
			_, _ = fmt.Fprintln(w, "error retrieving your user")
			return
		}
		if user == nil {
			_, _ = fmt.Fprintln(w, "you must be logged in to view cfg")
			return
		}

		tmpl, err := template.ParseFS(templates, "templates/gamestate_integration_csgostate.tmpl.cfg")
		if err != nil {
			_, _ = fmt.Fprintln(w, "unable to load template: ", err)
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
			_, _ = fmt.Fprintln(w, "unable to execute template: ", err)
			return
		}
	}
}
