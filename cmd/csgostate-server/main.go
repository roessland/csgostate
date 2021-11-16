package main

import (
	"fmt"
	"github.com/roessland/csgostate/cmd/csgostate-server/api"
	"github.com/roessland/csgostate/cmd/csgostate-server/playerevents"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"log"
	"time"
)

func main() {
	app, err := server.NewApp(server.NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	registerEventHandlers(app)

	//debugEventHandlers(app)

	api.ServeAPI(app)
}

func registerEventHandlers(app *server.App) {
	playerevents.Spectating.Register(func(payload playerevents.SpectatingPayload) {
		app.Log.Infow("",
			"event", playerevents.Spectating.String(),
			"auth_nick", payload.CurrState.Auth.Nick,
			"player", payload.CurrState.Player.Name)
	})

	playerevents.Died.Register(func(payload playerevents.DiedPayload) {
		app.Log.Infow("",
			"event", playerevents.Died.String(),
			"nick", payload.PrevState.Player.Name,
			"time", payload.CurrState.Provider.Timestamp)
	})

	playerevents.Spawned.Register(func(payload playerevents.SpawnedPayload) {
		app.Log.Infow("",
			"event", playerevents.Spawned.String(),
			"nick", payload.CurrState.Player.Name,
			"time", payload.CurrState.Provider.Timestamp)
	})

	playerevents.Appeared.Register(func(payload playerevents.AppearedPayload) {
		app.Log.Infow("",
			"event", playerevents.Appeared.String(),
			"nick", payload.CurrState.Player.Name,
			"time", payload.CurrState.Provider.Timestamp)
	})
}

// debugEventHandlers feeds all states in database to the player events extractor.
func debugEventHandlers(app *server.App) {
	states, err := app.StateRepo.GetAll()
	if err != nil {
		panic(err)
	}
	fmt.Print(len(states))
	for i, _ := range states {
		//s := states[i]
		//if s.Provider.Timestamp < 1636919279 || 1636919287 < s.Provider.Timestamp {
		//continue
		//}
		//fmt.Println(string(s.RawJson))
		time.Sleep(time.Millisecond * 30)

		err = app.PlayerEventsExtractor.Feed(&states[i])
		if err != nil {
			panic(err)
		}
	}
}
