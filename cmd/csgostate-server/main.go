package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/steam"

	"github.com/roessland/csgostate/csgostate"
)



//go:embed static/*
//go:embed index.html
var static embed.FS

//go:embed templates/*
var templates embed.FS

type App struct {
	Config          Config
	SteamHTTPClient *http.Client
	SessionStore    *SessionStore
	PlayerRepo      InMemoryPlayerRepo
}

type Config struct {
	SessionSecret string
	URL string
	SteamKey string
}

func (config Config) Verify() {
	if len(config.SessionSecret) < 10 {
		panic("missing or too short SESSION_SECRET environment variable")
	}

	if config.URL == "" {
		panic("Must specify CSGOSS_URL")
	}

	if len(config.SteamKey) == 0 {
		panic("you must set STEAM_KEY env to get profile info")
	}
}

func NewConfig() Config {
	config := Config{}
	config.URL = "http://localhost:3528/"
	config.SessionSecret = os.Getenv("SESSION_SECRET")
	config.SteamKey = os.Getenv("STEAM_KEY")
	config.Verify()
	return config
}

func NewApp(config Config) *App {
	app := &App{}
	app.Config = config
	app.SessionStore = NewSessionStore([]byte(app.Config.SessionSecret))
	app.PlayerRepo = NewPlayerRepo()
	return app
}

func main() {
	app := NewApp(NewConfig())

	// Listen to gamestate integration push messages
	listener := csgostate.NewListener()

	// Serve player states as an API
	go ServeAPI(app, listener.HandlerFunc)

	for state := range listener.Updates {
		fmt.Printf("%v\n", state)
		app.PlayerRepo.Update(&state)
	}
}

func ServeAPI(app *App, pushHandler http.HandlerFunc) {
	goth.UseProviders(
		steam.New(os.Getenv("STEAM_KEY"), app.Config.URL + "auth/callback"),
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
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			_, _ = fmt.Fprintln(w, "unable to fetch user details from Steam: ", err)
			return
		}

		// Create a new session
		sess, err := app.SessionStore.New(r)
		if err != nil {
			_, _ = fmt.Fprintln(w, "unable to create session")
			return
		}

		sess.SetNickName(user.NickName)
		sess.SetAvatarURL(user.AvatarURL)
		sess.SetSteamID(user.UserID)

		err = app.SessionStore.Save(r, w, sess)
		if err != nil {
			_, _ = fmt.Fprintln(w, "unable to save session cookie: ", err)
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
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	router.HandleFunc("/api/push", func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
		}
		fmt.Println("\n\nxxxxx\n\n", string(buf))

		r.Body = io.NopCloser(bytes.NewReader(buf))
		pushHandler.ServeHTTP(w, r)
	})

	router.HandleFunc("/api/players", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		err := json.NewEncoder(w).Encode(app.PlayerRepo.GetAll())
		if err != nil {
			log.Println(err)
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
		Handler: router,
		Addr:    "127.0.0.1:3528",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("listening to %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
