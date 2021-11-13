package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/roessland/csgostate/csgostate"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/steam"
	bolt "go.etcd.io/bbolt"
)

func main() {
	app, err := NewApp(NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	go app.ServeAPI()

	for state := range app.StateListener.Updates {
		fmt.Printf("%v\n", state)
		app.PlayerRepo.Update(&state)
	}
}

//go:embed static/*
//go:embed index.html
var static embed.FS

//go:embed templates/*
var templates embed.FS

type App struct {
	Config          Config
	SteamHTTPClient *http.Client
	SessionStore    *SessionStore
	DB              *bolt.DB
	UserRepo        UserRepo
	PlayerRepo      PlayerRepo
	StateListener   *csgostate.Listener
}

func NewApp(config Config) (*App, error) {
	var err error
	app := &App{}
	app.Config = config
	app.SessionStore = NewSessionStore([]byte(app.Config.SessionSecret))
	app.DB, err = bolt.Open("csgostate.db", 0666, nil)
	if err != nil {
		return nil, err
	}
	app.PlayerRepo = NewPlayerRepo()
	app.UserRepo, err = NewDBUserRepo(app.DB, app.Config.PushTokenSecret)
	if err != nil {
		return nil, err
	}
	app.StateListener = csgostate.NewListener()
	return app, nil
}

func (app *App) Close() {
	app.DB.Close()
}

func (app *App) ServeAPI() {
	goth.UseProviders(
		steam.New(app.Config.SteamKey, app.Config.URL+"auth/callback"),
	)
	gothic.GetProviderName = func(r *http.Request) (string, error) {
		return "steam", nil
	}

	router := mux.NewRouter()

	router.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r)
	})

	router.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		// Ask Steam API for user details
		oauthUser, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to fetch user details from Steam: %s", err), 500)
			return
		}

		// Create a new session
		sess, err := app.SessionStore.New(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to create session: %s", err), 500)
			return
		}
		sess.SetNickName(oauthUser.NickName)
		sess.SetAvatarURL(oauthUser.AvatarURL)
		sess.SetSteamID(oauthUser.UserID)

		user, err := app.UserRepo.GetBySteamID(oauthUser.UserID)
		if err != nil {
			log.Println("SteamID ", oauthUser.UserID)
			http.Error(w, fmt.Sprintf("unable to get user: %s", err), 500)
			return
		}
		if user == nil {
			user = &User{
				SteamID:   oauthUser.UserID,
				NickName:  oauthUser.NickName,
				AvatarURL: oauthUser.AvatarURL,
			}
			err := app.UserRepo.Create(user)
			if err != nil {
				http.Error(w, fmt.Sprintf("unable to create user: %s", err), 500)
				return
			}
		}

		err = app.SessionStore.Save(r, w, sess)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to save session cookie: %s", err), 500)
			return
		}

		// Redirect to front page
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusFound)
	})

	router.HandleFunc("/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		sess, _ := app.SessionStore.New(r)
		sess.Values = map[interface{}]interface{}{}
		app.SessionStore.Save(r, w, sess)
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			log.Print(err)
		}
	})

	router.HandleFunc("/api/push", func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
		}
		fmt.Println("\n\nxxxxx\n\n", string(buf))

		r.Body = io.NopCloser(bytes.NewReader(buf))
		app.StateListener.HandlerFunc(w, r)
	})

	router.HandleFunc("/api/players", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		err := json.NewEncoder(w).Encode(app.PlayerRepo.GetAll())
		if err != nil {
			log.Println(err)
		}
	})

	router.HandleFunc("/gamestate_integration_csgostate.cfg", func(w http.ResponseWriter, r *http.Request) {
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
			Sess *Session
			User *User
		}{
			Sess: sess,
			User: user,
		})
		if err != nil {
			_, _ = fmt.Fprintln(w, "unable to execute template: ", err)
			return
		}
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// Catch-all for remaining requests. Must be last.
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		http.FileServer(http.FS(static)).ServeHTTP(w, r)
	})

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:3528",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("listening to %s", srv.Addr)
	log.Print(srv.ListenAndServe())
}
