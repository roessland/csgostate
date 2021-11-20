package playerevents

import (
	"github.com/roessland/csgostate/pkg/csgostate"
)

type SpectatingPayload struct {
	PrevState *csgostate.State
	CurrState *csgostate.State
}

type spectating struct {
	handlers []func(payload SpectatingPayload)
}

func (e *spectating) String() string {
	return "player_spectating"
}

func (e *spectating) Register(handler func(payload SpectatingPayload)) {
	e.handlers = append(e.handlers, handler)
}

func (e *spectating) Trigger(payload SpectatingPayload) {
	for _, handler := range e.handlers {
		go handler(payload)
	}
}


func (e *spectating) extractFromStateDiff(prevState, currState *csgostate.State) error {
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
		e.Trigger(SpectatingPayload{
			PrevState: prevState,
			CurrState: currState,
		})
	}
	return nil
}