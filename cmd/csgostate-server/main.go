package main

import (
	"embed"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/roessland/csgostate/csgostate"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

//go:embed static/*
//go:embed index.html
var static embed.FS
var players PlayerRepo

func main() {
	// Store the state of every player that has ever sent a gamestate to us
	players = NewPlayerRepo()

	// Listen to gamestate integration push messages
	listener := csgostate.NewListener()

	// Serve player states as an API
	go ServeAPI(listener.HandlerFunc)

	for state := range listener.Updates {
		//fmt.Printf("%v\n", state)
		players.Update(&state)
	}
}

func ServeAPI(pushHandler http.HandlerFunc) {
	router := mux.NewRouter()

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	router.HandleFunc("/api/push", func(w http.ResponseWriter, r *http.Request) {
		pushHandler.ServeHTTP(w, r)
	})

	router.HandleFunc("/api/players", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		err := json.NewEncoder(w).Encode(players)
		if err != nil {
			log.Println(err)
		}
	})

	oauthCfg := &oauth2.Config{
		ClientID:     "",
		ClientSecret: "https://localhost:3528/steamlogin/return",
		Endpoint:     oauth2.Endpoint{
			AuthURL:   "https://steamcommunity.com/openid",
			TokenURL:  "",
			AuthStyle: oauth2.AuthStyleAutoDetect,
		},
		RedirectURL:  "",
		Scopes:       nil,
	}

	router.HandleFunc("/steamlogin/return", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("this is steam login return"))
		// To verify the user, make a call from your backend to
		//https://steamcommunity.com/openid/login copying every query string parameter
		//from that response with one exception: replace &openid.mode=id_res
		//with &openid.mode=check_authentication.
	})

	router.HandleFunc("/steamlogin/redirect", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, oauthCfg.AuthCodeURL("banana"), http.StatusFound)
	})

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

	log.Fatal(srv.ListenAndServe())
}
