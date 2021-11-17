package playerevents

import (
	"github.com/roessland/csgostate/csgostate"
)

type SpawnedPayload struct {
	PrevState *csgostate.State
	CurrState *csgostate.State
}

type spawned struct {
	handlers []func(payload SpawnedPayload)
}

func (e *spawned) String() string {
	return "player_spawned"
}

func (e *spawned) Register(handler func(payload SpawnedPayload)) {
	e.handlers = append(e.handlers, handler)
}

func (e *spawned) Trigger(payload SpawnedPayload) {
	for _, handler := range e.handlers {
		go handler(payload)
	}
}

func (e *spawned) extractFromStateDiff(prevState, currState *csgostate.State) error {
	if prevState == nil {
		// Should never happen
		return nil
	}

	if currState.Player == nil || currState.Player.Activity == csgostate.PlayerActivityMenu {
		// Not playing
		return nil
	}

	if currState.Provider.SteamID != currState.Player.SteamID {
		// Spectating
		return nil
	}

	if currState.Player.Name == csgostate.PlayerNameUnconnected {
		// Joined server just now. Or nick is "unconnected"...
		return nil
	}

	if prevState.Player == nil ||
		prevState.Player.State == nil ||
		prevState.Player.SteamID != currState.Provider.SteamID ||
		(prevState.Player != nil && prevState.Player.State.Health == 0 && currState.Player.State.Health > 0) {
		// Not already playing or not already alive.
		e.Trigger(SpawnedPayload{
			PrevState: prevState,
			CurrState: currState,
		})
	}

	return nil
}