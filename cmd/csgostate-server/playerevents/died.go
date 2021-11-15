package playerevents

import (
	"github.com/roessland/csgostate/csgostate"
)

var Died died

type DiedPayload struct {
	PrevState *csgostate.State
	CurrState *csgostate.State
}

type died struct {
	handlers []func(payload DiedPayload)
}

func (e *died) String() string {
	return "died"
}

func (e *died) Register(handler func(payload DiedPayload)) {
	e.handlers = append(e.handlers, handler)
}

func (e *died) Trigger(payload DiedPayload) {
	for _, handler := range e.handlers {
		go handler(payload)
	}
}
