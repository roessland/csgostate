package handlers

import (
	"github.com/markbates/goth/gothic"
	"net/http"
)

func AuthLogin() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r)
	}
}
