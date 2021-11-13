package handlers

import (
	"bytes"
	"fmt"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func ApiPush(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
		}
		fmt.Println("\n\nxxxxx\n\n", string(buf))

		r.Body = io.NopCloser(bytes.NewReader(buf))
		app.StateListener.HandlerFunc(w, r)
	}
}
