package teamevents

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
	return "team_spawned"
}

func (e *spawned) Register(handler func(payload SpawnedPayload)) {
	e.handlers = append(e.handlers, handler)
}

func (e *spawned) Trigger(payload SpawnedPayload) {
	for _, handler := range e.handlers {
		go handler(payload)
	}
}