package handlers

import (
	"fmt"
	"github.com/markbates/goth/gothic"
	"github.com/roessland/csgostate/cmd/csgostate-server/repos/userrepo"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"log"
	"net/http"
)

func GetAuthCallback(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
			user = &userrepo.User{
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
	}
}
