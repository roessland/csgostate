package maps

type Map string

const Nil Map = ""
const Unknown Map = "unknown"
const Dust2 Map = "de_dust2"
const Inferno Map = "de_inferno"
const Vertigo Map = "de_vertigo"
const Cache Map = "de_cache"
const Nuke Map = "de_nuke"

func FromString(name string) Map {
	switch Map(name) {
	case Dust2:
		return Dust2
	case Inferno:
		return Inferno
	case Vertigo:
		return Vertigo
	case Cache:
		return Cache
	case Nuke:
		return Nuke
	default:
		return Unknown
	}
}
