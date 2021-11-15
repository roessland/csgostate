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
	return "appeared"
}

func (e *appeared) Register(handler func(payload AppearedPayload)) {
	e.handlers = append(e.handlers, handler)
}

func (e *appeared) Trigger(payload AppearedPayload) {
	for _, handler := range e.handlers {
		go handler(payload)
	}
}
