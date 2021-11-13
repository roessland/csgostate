package main

import (
	"fmt"
	"github.com/roessland/csgostate/cmd/csgostate-server/api"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"log"
)

func main() {
	app, err := server.NewApp(server.NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	go api.ServeAPI(app)

	for state := range app.StateListener.Updates {
		fmt.Printf("%v\n", state)
		app.PlayerRepo.Update(&state)
	}
}