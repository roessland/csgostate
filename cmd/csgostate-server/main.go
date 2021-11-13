package main

import (
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
		//fmt.Printf("%v\n", state)
		app.PlayerRepo.Update(&state)
		err := app.StateRepo.Push(&state)
		if err != nil {
			log.Println("cannot push state to repo: ", err)
		}
	}
}