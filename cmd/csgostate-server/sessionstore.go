package main

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type SessionStore struct {
	basicStore sessions.Store
}

func NewSessionStore(keyPairs ...[]byte) *SessionStore {
	sessionStore := SessionStore{
		basicStore: sessions.NewCookieStore(keyPairs...),
	}
	return &sessionStore
}

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
