package stratroulette

import (
	"github.com/roessland/csgostate/cmd/csgostate-server/maps"
	"github.com/roessland/csgostate/csgostate"
	"math/rand"
)

type Tag string

const TagSerious Tag = "serious"

type Strat struct {
	Team           csgostate.PlayerTeam
	Name           string
	DescNO         string
	DescEN         string
	ApplicableMaps map[maps.Map]bool
	Tags           map[Tag]bool
}

var Strats []Strat

func addStrat(strat Strat) {
	Strats = append(Strats, strat)
}

func init() {
	Strats = make([]Strat, 0)
}

func GetRandom(m maps.Map, team csgostate.PlayerTeam, onlySerious bool) *Strat {
	var matchingStrats []Strat
	for _, strat := range Strats {
		if team != strat.Team {
			continue
		}
		if onlySerious && !strat.Tags[TagSerious] {
			continue
		}
		if !strat.ApplicableMaps[m] {
			continue
		}
		matchingStrats = append(matchingStrats, strat)
	}
	if len(matchingStrats) == 0 {
		return nil
	}

	i := rand.Intn(len(matchingStrats))
	return &matchingStrats[i]
}
