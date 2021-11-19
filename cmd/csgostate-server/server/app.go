package server

import (
	"github.com/roessland/csgostate/cmd/csgostate-server/discord"
	"github.com/roessland/csgostate/cmd/csgostate-server/logger"
	"github.com/roessland/csgostate/cmd/csgostate-server/playerevents"
	"github.com/roessland/csgostate/cmd/csgostate-server/repos/playerrepo"
	"github.com/roessland/csgostate/cmd/csgostate-server/repos/staterepo"
	"github.com/roessland/csgostate/cmd/csgostate-server/repos/userrepo"
	"github.com/roessland/csgostate/cmd/csgostate-server/sessions"
	"github.com/roessland/csgostate/csgostate"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

type App struct {
	Config                Config
	Log                   logger.Logger
	SteamHTTPClient       *http.Client
	SessionStore          *sessions.SessionStore
	DB                    *bolt.DB
	Discord               *discord.Client
	UserRepo              userrepo.UserRepo
	PlayerRepo            playerrepo.PlayerRepo
	StateRepo             staterepo.StateRepo
	StateListener         *csgostate.Listener
	PlayerEvents          *playerevents.EventRepo
	PlayerEventsExtractor *playerevents.Extractor
}

func NewApp(config Config) (*App, error) {
	var err error

	app := &App{}

	app.Config = config

	app.Log, err = logger.NewLogger()
	if err != nil {
		return nil, err
	}

	app.SessionStore = sessions.NewSessionStore([]byte(app.Config.SessionSecret))

	app.DB, err = bolt.Open("csgostate.db", 0666, nil)
	if err != nil {
		return nil, err
	}

	app.Discord = discord.NewClient(app.Config.DiscordWebhookURL, app.Log)

	app.PlayerRepo = playerrepo.NewInMemoryPlayerRepo()

	app.UserRepo, err = userrepo.NewDBUserRepo(app.DB, app.Config.PushTokenSecret)
	if err != nil {
		return nil, err
	}

	app.StateRepo, err = staterepo.NewDBStateRepo(app.DB)
	if err != nil {
		return nil, err
	}

	app.StateListener = csgostate.NewListener()
	app.PlayerEvents = playerevents.NewRepo()
	app.PlayerEventsExtractor = playerevents.NewExtractor(app.PlayerEvents)

	return app, nil
}

func (app *App) Close() {
	_ = app.DB.Close()
	_ = app.Log.Sync()
}
