package main

import (
	"fmt"
	"github.com/roessland/csgostate/cmd/csgostate-server/api"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"github.com/roessland/csgostate/csgostate"
	"log"
)

func main() {
	app, err := server.NewApp(server.NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	// For debugging JSON structure against recorded data.
	users, _ := app.UserRepo.GetAll()
	for _, user := range users {
		fmt.Println(user.NickName)
		app.Log.Info(app.StateRepo.DebugJsonForPlayer(user.SteamID))
	}

	// For debugging JSON structure against recorded data.
	for _, user := range users {
		states, _ := app.StateRepo.GetAllForPlayer(user.SteamID)
		for i := 1; i < len(states); i++ {
			prevState := states[i-1]
			state := states[i]

			if prevState.Player == nil || state.Player == nil {
				continue
			}

			if prevState.Player.State == nil || state.Player.State == nil {
				continue
			}

			if prevState.Provider == nil || state.Provider == nil {
				continue
			}

			if prevState.Player.SteamID != prevState.Provider.SteamID || state.Player.SteamID != state.Provider.SteamID {
				continue
			}

			// Find active weapon in prevState
			var activeWeapon csgostate.Weapon
			for _, weap := range *prevState.Player.Weapons {
				if weap.State == csgostate.WeaponStateActive || weap.State == csgostate.WeaponStateReloading {
					activeWeapon = weap
				}
			}

			if prevState.Player.State.Health > 0 && state.Player.State.Health == 0 {
				fmt.Println(state.Player.Name,
					" died on ", state.Map.Name,
					" gamemode ", state.Map.Mode,
					activeWeapon.State, activeWeapon.Name,
					" score ", prevState.Map.TeamCT.Score, "-", prevState.Map.TeamT.Score)
			}
		}
	}

	api.ServeAPI(app)
}
