package teamevents

import "github.com/roessland/csgostate/csgostate"

type RoundPhaseChangedPayload struct {
	From      csgostate.RoundPhase
	To        csgostate.RoundPhase
	CurrState *csgostate.State
}

type roundPhaseChanged struct {
	handlers []func(payload RoundPhaseChangedPayload)
}

func (e *roundPhaseChanged) String() string {
	return "round_phase_changed"
}

func (e *roundPhaseChanged) Register(handler func(payload RoundPhaseChangedPayload)) {
	e.handlers = append(e.handlers, handler)
}

func (e *roundPhaseChanged) Trigger(payload RoundPhaseChangedPayload) {
	for _, handler := range e.handlers {
		go handler(payload)
	}
}
