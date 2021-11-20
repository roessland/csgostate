package stratroulette

import (
	"github.com/roessland/csgostate/pkg/csgostate"
	"github.com/roessland/csgostate/pkg/maps"
)

func init() {

	addStrat(Strat{
		Team:   csgostate.PlayerTeamCT,
		Name:   "Camp taket i et minutt",
		DescNO: "All spawner ved apartments, går opp stigen til taket. Camper der til det er 40 sekund igjen, deretter rusher nærmeste hostage.",
		DescEN: "",
		ApplicableMaps: map[maps.Map]bool{
			maps.Insertion2: true,
		},
		Tags: map[Tag]bool{
			TagSerious: false,
		},
	})

	addStrat(Strat{
		Team:   csgostate.PlayerTeamCT,
		Name:   "Default",
		DescNO: "Default strat for CT",
		DescEN: "",
		ApplicableMaps: map[maps.Map]bool{
			maps.Basalt:     true,
			maps.Dust2:      true,
			maps.Inferno:    true,
			maps.Vertigo:    true,
			maps.Cache:      true,
			maps.Nuke:       true,
			maps.Ancient:    true,
			maps.Insertion2: true,
			maps.Mirage:     true,
		},
		Tags: map[Tag]bool{
			TagSerious: true,
		},
	})
}
