package playerevents

import (
	"github.com/roessland/csgostate/pkg/csgostate"
)

type DiedPayload struct {
	PrevState *csgostate.State
	CurrState *csgostate.State
}

type died struct {
	handlers []func(payload DiedPayload)
}

func (e *died) String() string {
	return "player_died"
}

func (e *died) Register(handler func(payload DiedPayload)) {
	e.handlers = append(e.handlers, handler)
}

func (e *died) Trigger(payload DiedPayload) {
	for _, handler := range e.handlers {
		go handler(payload)
	}
}


func (e *died) extractFromStateDiff(prevState, currState *csgostate.State) error {
	if prevState == nil {
		// This cannot be the first event
		return nil
	}

	if prevState.Player == nil {
		// Was not playing
		return nil
	}

	if prevState.Player.State == nil {
		// Was not playing
		return nil
	}

	if prevState.Player.SteamID != prevState.Provider.SteamID {
		// Was spectating someone
		return nil
	}

	if prevState.Player.State.Health == 0 {
		// Was already dead
		return nil
	}

	if currState.Player != nil && currState.Player.Activity == csgostate.PlayerActivityMenu {
		// Left the game
	}

	if prevState.Player.State.Health > 0 && (
		currState.Player == nil ||
			currState.Player.SteamID != currState.Provider.SteamID ||
			currState.Player.State == nil ||
			currState.Player.State.Health == 0) {
		e.Trigger(DiedPayload{
			PrevState: prevState,
			CurrState: currState,
		})
	}
	return nil
}