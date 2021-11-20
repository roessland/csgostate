package handlers

import (
	"github.com/markbates/goth/gothic"
	"net/http"
)

func GetAuthLogin() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r)
	}
}
