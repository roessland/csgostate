package stratroulette

import (
	"github.com/roessland/csgostate/cmd/csgostate-server/maps"
	"github.com/roessland/csgostate/csgostate"
)

func init() {

	addStrat(Strat{
		Team:   csgostate.PlayerTeamT,
		Name:   "Rush B",
		DescNO: "Rush B! No stop! Cyka blyat!",
		DescEN: "",
		ApplicableMaps: map[maps.Map]bool{
			maps.Inferno: true,
			maps.Dust2:   true,
			maps.Vertigo: true,
			maps.Ancient: true,
		},
		Tags: map[Tag]bool{
			TagSerious: true,
		},
	})

	addStrat(Strat{
		Team:   csgostate.PlayerTeamT,
		Name:   "Rush short A",
		DescNO: "Rush short A!",
		DescEN: "",
		ApplicableMaps: map[maps.Map]bool{
			maps.Inferno: true,
			maps.Dust2:   true,
			maps.Cache:   true,
		},
		Tags: map[Tag]bool{
			TagSerious: true,
		},
	})


	addStrat(Strat{
		Team:   csgostate.PlayerTeamT,
		Name:   "Kitchen",
		DescNO: "Camp på kjøkkenet frem til det skjer noe, eller rush A ved 0:50 hvis ingen pusher",
		DescEN: "",
		ApplicableMaps: map[maps.Map]bool{
			maps.Inferno: true,
		},
		Tags: map[Tag]bool{
			TagSerious: false,
		},
	})

	addStrat(Strat{
		Team:   csgostate.PlayerTeamT,
		Name:   "Full vent dive",
		DescNO: "Alle løper rett ned vents til B",
		DescEN: "",
		ApplicableMaps: map[maps.Map]bool{
			maps.Nuke: true,
		},
		Tags: map[Tag]bool{
			TagSerious: true,
		},
	})

	addStrat(Strat{
		Team:   csgostate.PlayerTeamT,
		Name:   "Vent/ramp split",
		DescNO: "Beste spawn løper rett ned i vents og planter B. Resten rusher Ramp til B.",
		DescEN: "",
		ApplicableMaps: map[maps.Map]bool{
			maps.Nuke: true,
		},
		Tags: map[Tag]bool{
			TagSerious: true,
		},
	})

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
}
