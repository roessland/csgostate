package server

import (
	"github.com/roessland/csgostate/cmd/csgostate-server/repos/playerrepo"
	"github.com/roessland/csgostate/cmd/csgostate-server/repos/userrepo"
	"github.com/roessland/csgostate/cmd/csgostate-server/sessions"
	"github.com/roessland/csgostate/csgostate"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

type App struct {
	Config          Config
	SteamHTTPClient *http.Client
	SessionStore    *sessions.SessionStore
	DB              *bolt.DB
	UserRepo        userrepo.UserRepo
	PlayerRepo      playerrepo.PlayerRepo
	StateListener   *csgostate.Listener
}

func NewApp(config Config) (*App, error) {
	var err error
	app := &App{}
	app.Config = config
	app.SessionStore = sessions.NewSessionStore([]byte(app.Config.SessionSecret))
	app.DB, err = bolt.Open("csgostate.db", 0666, nil)
	if err != nil {
		return nil, err
	}
	app.PlayerRepo = playerrepo.NewPlayerRepo()
	app.UserRepo, err = userrepo.NewDBUserRepo(app.DB, app.Config.PushTokenSecret)
	if err != nil {
		return nil, err
	}
	app.StateListener = csgostate.NewListener()
	return app, nil
}

func (app *App) Close() {
	app.DB.Close()
}
