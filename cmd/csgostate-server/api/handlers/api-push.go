package handlers

import (
	"bytes"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func GetApiPush(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
		}

		r.Body = io.NopCloser(bytes.NewReader(buf))
		app.StateListener.HandlerFunc(w, r)
	}
}
