package main

import (
	"github.com/gorilla/mux"
	"github.com/roessland/csgostate/pkg/csgostate"
	"log"
	"net/http"
	"os"
	"time"
)


func main() {
	// Listen to gamestate integration push messages
	listener := csgostate.NewListener()

	// Serve player states as an API
	go ServeAPI(listener.HandlerFunc)

	for state := range listener.Updates {
		os.Stdout.Write(state.RawJson)
		os.Stdout.Write([]byte("\n\n"))
	}
}

func ServeAPI(pushHandler http.HandlerFunc) {
	router := mux.NewRouter()

	router.HandleFunc("/api/push", func(w http.ResponseWriter, r *http.Request) {
		pushHandler.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:3528",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
