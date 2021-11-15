package playerevents

import (
	"github.com/roessland/csgostate/csgostate"
)

var DiedReloading diedReloading

func init() {
	Died.Register(func(payload DiedPayload) {
		for _, weap := range *payload.PrevState.Player.Weapons {
			if weap.State == csgostate.WeaponStateReloading {
				DiedReloading.Trigger(DiedReloadingPayload{
					PrevState: payload.PrevState,
					CurrState: payload.CurrState,
				})
			}
		}
	})
}

type DiedReloadingPayload struct {
	PrevState *csgostate.State
	CurrState *csgostate.State
}

type diedReloading struct {
	handlers []func(payload DiedReloadingPayload)
}

func (e *diedReloading) String() string {
	return "died-reloading"
}

func (e *diedReloading) Register(handler func(payload DiedReloadingPayload)) {
	e.handlers = append(e.handlers, handler)
}

func (e *diedReloading) Trigger(payload DiedReloadingPayload) {
	for _, handler := range e.handlers {
		go handler(payload)
	}
}
