package main

import (
	"github.com/gorilla/sessions"
)

const SessionCookieName = "_CSGOSS_SESSION"
const SteamNickNameSessionKey = "SteamNick"
const SteamAvatarURLCookieKey = "SteamAvatarURL"
const SteamIDCookieKey = "SteamID"

type Session sessions.Session

func (sess *Session) GetString(key string) string {
	val := sess.Values[key]
	switch v := val.(type) {
	case string:
		return v
	default:
		return ""
	}
}

func (sess *Session) SteamID() string {
	return sess.GetString(SteamIDCookieKey)
}

func (sess *Session) SetSteamID(steamID string) {
	sess.Values[SteamIDCookieKey] = steamID
}

func (sess *Session) IsLoggedIn() bool {
	return sess.SteamID() != ""
}

func (sess *Session) NickName() string {
	return sess.GetString(SteamNickNameSessionKey)
}

func (sess *Session) SetNickName(nickName string) {
	sess.Values[SteamNickNameSessionKey] = nickName
}

func (sess *Session) AvatarURL() string {
	return sess.GetString(SteamAvatarURLCookieKey)
}

func (sess *Session) SetAvatarURL(avatarURL string) {
	sess.Values[SteamAvatarURLCookieKey] = avatarURL
}