package handlers

import (
	"fmt"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"html/template"
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

		err = tmpl.Execute(w, sess)
		if err != nil {
			_, _ = fmt.Fprintln(w, "unable to execute template: ", err)
			return
		}
	}
}
