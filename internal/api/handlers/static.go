package handlers

import (
	"embed"
	"net/http"
)

//go:embed static/*
var static embed.FS

func GetStatic() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		http.FileServer(http.FS(static)).ServeHTTP(w, r)
	}
}
