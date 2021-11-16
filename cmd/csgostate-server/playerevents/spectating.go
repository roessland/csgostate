package playerevents

import (
	"github.com/roessland/csgostate/csgostate"
)

var Spectating spectating

type SpectatingPayload struct {
	PrevState *csgostate.State
	CurrState *csgostate.State
}

type spectating struct {
	handlers []func(payload SpectatingPayload)
}

func (e *spectating) String() string {
	return "spectating"
}

func (e *spectating) Register(handler func(payload SpectatingPayload)) {
	e.handlers = append(e.handlers, handler)
}

func (e *spectating) Trigger(payload SpectatingPayload) {
	for _, handler := range e.handlers {
		go handler(payload)
	}
}


func extractSpectating(prevState, currState *csgostate.State) error {
	if prevState == nil {
		// Should never happen
		return nil
	}

	if currState.Player == nil {
		// Not spectating anyone.
		return nil
	}

	if currState.Player.SteamID == currState.Provider.SteamID {
		// Playing, not spectating.
		return nil
	}
	if prevState.Player == nil || currState.Player.SteamID != prevState.Player.SteamID {
		Spectating.Trigger(SpectatingPayload{
			PrevState: prevState,
			CurrState: currState,
		})
	}
	return nil
}