package sessions

import (
	"github.com/gorilla/sessions"
)

const sessionCookieName = "_CSGOSS_SESSION"
const steamNickNameSessionKey = "SteamNick"
const steamAvatarURLCookieKey = "SteamAvatarURL"
const steamIDCookieKey = "SteamID"

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
	return sess.GetString(steamIDCookieKey)
}

func (sess *Session) SetSteamID(steamID string) {
	sess.Values[steamIDCookieKey] = steamID
}

func (sess *Session) IsLoggedIn() bool {
	return sess.SteamID() != ""
}

func (sess *Session) NickName() string {
	return sess.GetString(steamNickNameSessionKey)
}

func (sess *Session) SetNickName(nickName string) {
	sess.Values[steamNickNameSessionKey] = nickName
}

func (sess *Session) AvatarURL() string {
	return sess.GetString(steamAvatarURLCookieKey)
}

func (sess *Session) SetAvatarURL(avatarURL string) {
	sess.Values[steamAvatarURLCookieKey] = avatarURL
}