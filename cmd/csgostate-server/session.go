package main

import (
	"github.com/gorilla/sessions"
	"net/http"
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

type SessionStore struct {
	basicStore sessions.Store
}

func NewSessionStore(keyPairs ...[]byte) *SessionStore {
	sessionStore := SessionStore{
		basicStore: sessions.NewCookieStore(keyPairs...),
	}
	return &sessionStore
}

/*
sess.SetNickName(user.NickName)
sess.SetAvatarURL(user.AvatarURL)
sess.SetSteamID(user.UserID)
*/

func (store *SessionStore) New(r *http.Request) (*Session, error) {
	sess, err := store.basicStore.New(r, SessionCookieName)
	return (*Session)(sess), err
}

func (store *SessionStore) Get(r *http.Request) (*Session, error) {
	sess, err := store.basicStore.Get(r, SessionCookieName)
	return (*Session)(sess), err
}

func (store *SessionStore) Save(r *http.Request, w http.ResponseWriter, s *Session) error {
	err := store.basicStore.Save(r, w, (*sessions.Session)(s))
	return err
}
