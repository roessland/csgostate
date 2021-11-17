package teamevents

import "github.com/roessland/csgostate/cmd/csgostate-server/playerevents"

type pollingEmitter struct {
	playerEvents *playerevents.EventRepo
	teamEvents   *EventRepo
}

func NewPollingEmitter(playerEvents *playerevents.EventRepo, teamEvents *EventRepo) {
	pe := &pollingEmitter{}
	pe.playerEvents = playerEvents
	pe.teamEvents = teamEvents
}

func (pe *pollingEmitter) registerEventHandlers() {

}

func (pe *pollingEmitter) poll() {

}