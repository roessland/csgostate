package api

import (
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/steam"
	"github.com/roessland/csgostate/cmd/csgostate-server/api/handlers"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"log"
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

	router.HandleFunc("/auth/login", handlers.AuthLogin())
	router.HandleFunc("/auth/callback", handlers.AuthCallback(app))
	router.HandleFunc("/auth/logout", handlers.AuthLogout(app))
	router.HandleFunc("/api/health", handlers.ApiHealth())
	router.HandleFunc("/api/push", handlers.ApiPush(app))
	router.HandleFunc("/api/players", handlers.ApiPlayers(app))
	router.HandleFunc("/gamestate_integration_csgostate.cfg", handlers.GamestateCfg(app))
	router.HandleFunc("/", handlers.Index(app))

	// Catch-all for remaining requests. Must be last.
	router.PathPrefix("/").HandlerFunc(handlers.Static())

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:3528",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("listening to %s", srv.Addr)
	log.Print(srv.ListenAndServe())
}

