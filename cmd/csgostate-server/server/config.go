package server

import (
	"os"
	"strings"
)

type Config struct {
	// SessionSecret is used to encrypt session cookies.
	SessionSecret string

	// URL is the public URL of this application. For example `https://csgostate.roessland.com/`.
	URL string

	// SteamKey identifies this application to Steam API
	SteamKey string

	// PushTokenSecret is used along with SteamID to generate push tokens that users provide as
	// authentication. Very hard to rotate this secret since the tokens
	// are distributed in gamestate.cfg files on every users computer.
	PushTokenSecret string

	// Admins is a list of the SteamIDs of admin users.
	Admins []string
}

func NewConfig() Config {
	config := Config{}
	config.URL = os.Getenv("CSGOSS_URL")
	config.SessionSecret = os.Getenv("SESSION_SECRET")
	config.SteamKey = os.Getenv("STEAM_KEY")
	config.PushTokenSecret = os.Getenv("PUSH_TOKEN_SECRET")
	config.Admins = strings.Split(os.Getenv("ADMINS"), ",")
	config.Verify()
	return config
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

	if len(config.PushTokenSecret) < 20 {
		panic("missing or too short PUSH_TOKEN_SECRET environment variable")
	}
}

func (config Config) IsAdmin(steamID string) bool {
	for _, adminSteamID := range config.Admins {
		if adminSteamID == steamID && adminSteamID != "" {
			return true
		}
	}
	return false
}