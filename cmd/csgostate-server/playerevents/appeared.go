package playerevents

import (
	"github.com/roessland/csgostate/csgostate"
)

var Appeared appeared

type AppearedPayload struct {
	PrevState *csgostate.State
	CurrState *csgostate.State
}

type appeared struct {
	handlers []func(payload AppearedPayload)
}

func (e *appeared) String() string {
	return "player_appeared"
}

func (e *appeared) Register(handler func(payload AppearedPayload)) {
	e.handlers = append(e.handlers, handler)
}

func (e *appeared) Trigger(payload AppearedPayload) {
	for _, handler := range e.handlers {
		go handler(payload)
	}
}

func extractAppeared(prevState, currState *csgostate.State) error {
	lastEventTime := 0
	if prevState != nil {
		lastEventTime = prevState.Provider.Timestamp
	}
	secondsSincePrevState := currState.Provider.Timestamp - lastEventTime
	if secondsSincePrevState > 120 {
		Appeared.Trigger(AppearedPayload{
			PrevState: prevState,
			CurrState: currState,
		})
	}
	return nil
}