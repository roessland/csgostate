package middleware

import (
	"github.com/gorilla/mux"
	"github.com/roessland/csgostate/internal/server"
	"net/http"
)

func NewRequireAdminMiddleware(app *server.App) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, _ := app.SessionStore.New(r)
			steamID := sess.SteamID()
			if !app.Config.IsAdmin(steamID) {
				http.Error(w, "admin permission required to access this page", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
