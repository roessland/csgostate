package api

import (
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/steam"
	"github.com/roessland/csgostate/cmd/csgostate-server/api/handlers"
	"github.com/roessland/csgostate/cmd/csgostate-server/api/middleware"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"net/http"
	"time"
)

func ServeAPI(app *server.App) {
	goth.UseProviders(
		steam.New(app.Config.SteamKey, app.Config.URL+"auth/callback"),
	)
	gothic.GetProviderName = func(r *http.Request) (string, error) {
		return "steam", nil
	}

	router := mux.NewRouter()

	router.Use(middleware.NewRequestIDMiddleware(app))

	router.Use(middleware.NewRequestResponseLoggingMiddleware(app))

	router.HandleFunc("/auth/login", handlers.GetAuthLogin()).
		Methods(http.MethodGet)

	router.HandleFunc("/auth/callback", handlers.GetAuthCallback(app)).
		Methods(http.MethodGet)

	router.HandleFunc("/auth/logout", handlers.GetAuthLogout(app)).
		Methods(http.MethodGet)

	router.HandleFunc("/api/health", handlers.GetApiHealth()).
		Methods(http.MethodGet)

	router.HandleFunc("/api/push", handlers.GetApiPush(app)).
		Methods(http.MethodGet)

	router.HandleFunc("/api/players", handlers.GetApiPlayers(app)).
		Methods(http.MethodGet)

	router.HandleFunc("/gamestate_integration_csgostate.cfg", handlers.GetGamestateCfg(app)).
		Methods(http.MethodGet)

	router.HandleFunc("/", handlers.GetIndex(app)).
		Methods(http.MethodGet)

	// Catch-all for remaining requests. Must be last.
	router.PathPrefix("/").HandlerFunc(handlers.GetStatic()).
		Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:3528",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	app.Log.Infof("listening to %s", srv.Addr)
	app.Log.Infow("server closed", "err", srv.ListenAndServe())
}
