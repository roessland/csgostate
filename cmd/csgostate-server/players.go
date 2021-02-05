package main

import (
	"github.com/roessland/csgostate/csgostate"
	"log"
)

type PlayerRepo interface {
	GetAll() []*Player
	GetOrCreatePlayer(steamID string) *Player
	Update(state *csgostate.State)
}

// verify that it satisfies interface
var _ PlayerRepo = InMemoryPlayerRepo{}

type InMemoryPlayerRepo map[string]*Player

func NewPlayerRepo() InMemoryPlayerRepo {
	return make(map[string]*Player)
}

func (pr InMemoryPlayerRepo) GetAll() []*Player {
	var players []*Player
	for _, p := range pr {
		players = append(players, p)
	}
	return players
}

func (pr InMemoryPlayerRepo) GetOrCreatePlayer(steamID string) *Player {
	_, ok := pr[steamID]
	if !ok {
		pr[steamID] = NewPlayer()
	}
	return pr[steamID]
}

func (pr InMemoryPlayerRepo) Update(state *csgostate.State) {
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
