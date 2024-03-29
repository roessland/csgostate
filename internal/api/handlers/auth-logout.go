package handlers

import (
	"github.com/markbates/goth/gothic"
	"github.com/roessland/csgostate/internal/server"
	"net/http"
)

func GetAuthLogout(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		sess, _ := app.SessionStore.New(r)
		sess.Values = map[interface{}]interface{}{}
		app.SessionStore.Save(r, w, sess)
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}

}
