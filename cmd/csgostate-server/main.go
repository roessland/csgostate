package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/roessland/csgostate/internal/api"
	"github.com/roessland/csgostate/internal/metrics"
	"github.com/roessland/csgostate/internal/playerevents"
	"github.com/roessland/csgostate/internal/server"
	"github.com/roessland/csgostate/pkg/csgostate"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	app, err := server.NewApp(server.NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	// stateGrep(app, `"name": "weapon_flashbang",`)

	registerEventHandlers(app)

	// debugEventHandlers(app)

	go metrics.Serve(app.Log)

	api.ServeAPI(app)
}

func writeArchive(app *server.App, year, month int) {
	fileName := fmt.Sprintf("archive-states-%d-%d.json", year, month)

	stat, err := os.Stat(fileName)
	if !errors.Is(err, fs.ErrNotExist) && err != nil {
		app.Log.Warnw("stat archive", "err", err.Error())
		return
	}
	if stat != nil && stat.Size() > 0 {
		app.Log.Warnw("archive already exists")
		return
	}

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		app.Log.Warnw("opening archive", "err", err.Error())
		return
	}
	defer f.Close()

	err = app.StateRepo.ArchiveMonth(2021, month, f)
	if err != nil {
		app.Log.Errorw("writing archive", "err", err.Error())
		return
	}

	app.Log.Info(fmt.Sprintf("wrote archive for %d-%d", year, month))
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

		if isHoldingGrenade(payload.PrevState) {
			app.Discord.Post(payload.PrevState.Player.Name + " ble drept med en granat i hånden. Aiaiai, for en noob.")
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

		app.Log.Sync()
	}
}

// stateGrep searches DB for a string.
func stateGrep(app *server.App, str string) {
	states, err := app.StateRepo.GetAll()
	if err != nil {
		panic(err)
	}
	for i := 1; i < len(states); i++ {
		if strings.Contains(string(states[i].RawJson), str) {
			fmt.Println(string(states[i].RawJson))
		}
	}
	os.Exit(0)
}

func isReloading(state *csgostate.State) bool {
	for _, weap := range *state.Player.Weapons {
		if weap.State == csgostate.WeaponStateReloading {
			return true
		}
	}
	return false
}

func isHoldingGrenade(state *csgostate.State) bool {
	for _, weap := range *state.Player.Weapons {
		if weap.State != csgostate.WeaponStateActive {
			continue
		}
		if strings.HasSuffix(weap.Name, "flashbang") || strings.HasSuffix(weap.Name, "grenade") {
			return true
		}
	}
	return false
}
