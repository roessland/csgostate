package maps

type Map string

const Nil Map = ""
const Unknown Map = "unknown"
const Dust2 Map = "de_dust2"
const Inferno Map = "de_inferno"
const Mirage Map = "de_mirage"
const Vertigo Map = "de_vertigo"
const Cache Map = "de_cache"
const Nuke Map = "de_nuke"
const Ancient Map = "de_ancient"
const Insertion2 Map = "cs_insertion2"
const Basalt Map = "de_basalt"

func FromString(name string) Map {
	switch Map(name) {
	case Dust2:
		return Dust2
	case Inferno:
		return Inferno
	case Mirage:
		return Mirage
	case Vertigo:
		return Vertigo
	case Cache:
		return Cache
	case Nuke:
		return Nuke
	case Ancient:
		return Ancient
	case Insertion2:
		return Insertion2
	case Basalt:
		return Basalt
	default:
		return Unknown
	}
}
