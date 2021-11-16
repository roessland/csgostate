package main

import (
	"fmt"
	"github.com/roessland/csgostate/cmd/csgostate-server/api"
	"github.com/roessland/csgostate/cmd/csgostate-server/playerevents"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"github.com/roessland/csgostate/csgostate"
	"log"
	"sort"
	"time"
)

func main() {
	app, err := server.NewApp(server.NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	registerEventHandlers(app)

	migrateDb(app)

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

// debugEventHandlers feeds all states in database to the player events extractor.
func migrateDb(app *server.App) {
	var allStates []csgostate.State
	users, _ := app.UserRepo.GetAll()
	for _, user := range users {
		fmt.Println(user.NickName)
		states, err := app.StateRepo.GetAllForPlayer(user.SteamID)
		if err != nil {
			panic(err)
		}
		allStates = append(allStates, states...)
	}

	sort.Slice(allStates, func(i, j int) bool {
		return allStates[i].Provider.Timestamp < allStates[j].Provider.Timestamp
	})

	fmt.Println("inserting", len(allStates), "allstates")

	for i, _ := range allStates {
		err := app.StateRepo.PushMigrate(&allStates[i])
		if err != nil {
			panic(err)
		}
	}
}
