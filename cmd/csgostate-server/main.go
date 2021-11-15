package main

import (
	"fmt"
	"github.com/roessland/csgostate/cmd/csgostate-server/api"
	"github.com/roessland/csgostate/cmd/csgostate-server/playerevents"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"log"
)

func main() {
	app, err := server.NewApp(server.NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	registerEventHandlers(app)

	// debugEventHandlers(app)

	api.ServeAPI(app)
}

func registerEventHandlers(app *server.App) {
	playerevents.Died.Register(func(payload playerevents.DiedPayload) {
		app.Log.Infow("",
			"event", playerevents.Died.String(),
			"nick", payload.PrevState.Player.Name)
	})

	playerevents.DiedReloading.Register(func(payload playerevents.DiedReloadingPayload) {
		app.Log.Infow("",
			"event", playerevents.DiedReloading.String(),
			"nick", payload.PrevState.Player.Name)
	})

	playerevents.Appeared.Register(func(payload playerevents.AppearedPayload) {
		app.Log.Infow("",
			"event", playerevents.Appeared.String(),
			"nick", payload.CurrState.Player.Name)
	})
}

// debugEventHandlers feeds all states in database to the player events extractor.
func debugEventHandlers(app *server.App) {
	users, _ := app.UserRepo.GetAll()
	for _, user := range users {
		fmt.Println(user.NickName)
		states, err := app.StateRepo.GetAllForPlayer(user.SteamID)
		if err != nil {
			panic(err)
		}
		for i, _ := range states {
			err = app.PlayerEventsExtractor.Feed(&states[i])
			if err != nil {
				panic(err)
			}
		}
	}
}
