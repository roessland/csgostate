package main

import (
	"github.com/roessland/csgostate/csgostate"
	"log"
)

type PlayerRepo map[string]*Player

func NewPlayerRepo() PlayerRepo {
	return make(map[string]*Player)
}

func (pr PlayerRepo) GetOrCreatePlayer(steamID string) *Player {
	_, ok := pr[steamID]
	if !ok {
		pr[steamID] = NewPlayer()
	}
	return pr[steamID]
}

func (pr PlayerRepo) Update(state *csgostate.State) {
	steamID := state.Provider.SteamID
	if steamID == "" {
		log.Print("empty steamID in update")
	}
	p := pr.GetOrCreatePlayer(steamID)
	p.Update(state)
}

type Player struct {
	LatestState csgostate.State
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Update(state *csgostate.State) {
	p.LatestState = *state
}
