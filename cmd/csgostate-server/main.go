package main

import (
	"github.com/roessland/csgostate/cmd/csgostate-server/api"
	"github.com/roessland/csgostate/cmd/csgostate-server/maps"
	"github.com/roessland/csgostate/cmd/csgostate-server/playerevents"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"github.com/roessland/csgostate/cmd/csgostate-server/stratroulette"
	"github.com/roessland/csgostate/cmd/csgostate-server/teamevents"
	_ "github.com/roessland/csgostate/cmd/csgostate-server/teamevents"
	"github.com/roessland/csgostate/csgostate"
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
	app.PlayerEvents.Spectating.Register(func(payload playerevents.SpectatingPayload) {
		app.Log.Infow("",
			"event", app.PlayerEvents.Spectating.String(),
			"auth_nick", payload.CurrState.Auth.Nick,
			"player", payload.CurrState.Player.Name)
	})

	app.PlayerEvents.Died.Register(func(payload playerevents.DiedPayload) {
		app.Log.Infow("",
			"event", app.PlayerEvents.Died.String(),
			"nick", payload.PrevState.Player.Name,
			"time", payload.CurrState.Provider.Timestamp)

		if isReloading(payload.PrevState) {
			app.Discord.Post("lol! Ikke reload når du peeker, lmao for en noob du er " + payload.PrevState.Player.Name)
		}
	})

	app.PlayerEvents.Spawned.Register(func(payload playerevents.SpawnedPayload) {
		app.Log.Infow("",
			"event", app.PlayerEvents.Spawned.String(),
			"nick", payload.CurrState.Player.Name,
			"time", payload.CurrState.Provider.Timestamp)
	})

	app.PlayerEvents.Appeared.Register(func(payload playerevents.AppearedPayload) {
		app.Log.Infow("",
			"event", app.PlayerEvents.Appeared.String(),
			"nick", payload.CurrState.Player.Name,
			"time", payload.CurrState.Provider.Timestamp)
	})

	app.TeamEvents.Created.Register(func(payload teamevents.CreatedPayload) {
		app.Log.Infow("",
			"event", app.TeamEvents.Created.String(),
			"team_id", payload.TeamID)
	})

	app.TeamEvents.PlayerJoined.Register(func(payload teamevents.PlayerJoinedPayload) {
		app.Log.Infow("",
			"event", app.TeamEvents.PlayerJoined.String(),
			"team_id", payload.TeamID,
			"player_id", payload.PlayerID)
	})

	app.TeamEvents.RoundPhaseChanged.Register(func(payload teamevents.RoundPhaseChangedPayload) {
		app.Log.Infow("",
			"event", app.TeamEvents.RoundPhaseChanged.String(),
			"from_phase", payload.From,
			"to_phase", payload.To)

		if payload.To == csgostate.RoundPhaseFreezetime {
			m := maps.FromString(payload.CurrState.Map.Name)
			strat := stratroulette.GetRandom(m, payload.CurrState.Player.Team, false)
			if strat != nil {
				app.Discord.Post("YOOO BOIIS HERE IS THE PLAN: " + strat.DescNO)
			}
		}

		if payload.To == csgostate.RoundPhaseLive {
			//fmt.Println("Lets gooooo!")
		}
	})
}

// debugEventHandlers feeds all states in database to the player events extractor.
func debugEventHandlers(app *server.App) {
	states, err := app.StateRepo.GetAll()
	if err != nil {
		panic(err)
	}
	for i := 1; i < len(states); i++ {
		sleepMillis := (states[i].Provider.Timestamp - states[i-1].Provider.Timestamp) * 20
		if sleepMillis > 100 {
			sleepMillis = 100
		}
		if sleepMillis < 0 {
			sleepMillis = 0
		}
		if states[i].Player.Activity == csgostate.PlayerActivityMenu {
			sleepMillis = 0
		}
		time.Sleep(time.Duration(sleepMillis) * time.Millisecond)

		err = app.PlayerEventsExtractor.Feed(&states[i])
		if err != nil {
			panic(err)
		}

		app.TeamsRepo.Feed(&states[i])
		app.Log.Sync()
	}
}

func isReloading(state *csgostate.State) bool {
	for _, weap := range *state.Player.Weapons {
		if weap.State == csgostate.WeaponStateReloading {
			return true
		}
	}
	return false
}