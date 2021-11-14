package handlers

import (
	"fmt"
	"github.com/markbates/goth/gothic"
	"github.com/roessland/csgostate/cmd/csgostate-server/repos/userrepo"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"net/http"
)

func GetAuthCallback(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ask Steam API for user details
		oauthUser, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			app.Log.Errorw("error completing steam auth",
				"err", err)
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
			app.Log.Errorw("error getting user",
				"steamid", oauthUser.UserID,
				"nickname", oauthUser.NickName,
				"err", err)
			http.Error(w, "error during auth callback", http.StatusInternalServerError)
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
				app.Log.Errorw("error creating user",
					"steamid", oauthUser.UserID,
					"nickname", oauthUser.NickName,
					"err", err)
				http.Error(w, "couldn't create a user", http.StatusInternalServerError)
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
