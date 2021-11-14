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

	// For debugging JSON structure against recorded data.
	//users, _ := app.UserRepo.GetAll()
	//for _, user := range users {
	//	fmt.Println(user.NickName)
	//	app.Log.Info(app.StateRepo.DebugJsonForPlayer(user.SteamID))
	//}

	api.ServeAPI(app)
}
