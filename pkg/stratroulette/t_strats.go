package stratroulette

import (
	"github.com/roessland/csgostate/pkg/csgostate"
	"github.com/roessland/csgostate/pkg/maps"
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
			maps.Basalt:  true,
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
		Name:   "Rush A",
		DescNO: "Rush A!",
		DescEN: "",
		ApplicableMaps: map[maps.Map]bool{
			maps.Basalt: true,
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
		Team:   csgostate.PlayerTeamT,
		Name:   "Default",
		DescNO: "Default strat for T, få mid-kontroll, få noen picks.",
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
