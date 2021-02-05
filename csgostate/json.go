package csgostate

type State struct {
	Provider Provider `json:"provider"`
	Player   Player   `json:"player"`
	Round    Round    `json:"round"`
	Map Map `json:"map"`
	RawJson  []byte
}

type Map struct {
	Mode string `json:"mode"`
	Name string `json:"name"`
	Phase string `json:"phase"`
}

/*
	"provider": {
		"name": "Counter-Strike: Global Offensive",
		"appid": 730,
		"version": 13779,
		"steamid": "76561197993200126",
		"timestamp": 1612019048
	},
*/
type Provider struct {
	SteamID   string `json:"steamid"`
	Timestamp int    `json:"timestamp"`
}

type Player struct {
	SteamID  string        `json:"steamid"`
	Clan     string        `json:"clan"`
	Name     string        `json:"name"`
	Team     string        `json:"team"`
	Activity string        `json:"activity"`
	State    PlayerState   `json:"state"`
	Weapons  PlayerWeapons `json:"weapons"`
}

/*
	"state": {
		"health": 100,
		"armor": 100,
		"helmet": false,
		"flashed": 0,
		"smoked": 0,
		"burning": 0,
		"money": 150,
		"round_kills": 0,
		"round_killhs": 0,
		"equip_value": 850
	},
*/
type PlayerState struct {
	Health      int  `json:"health"`
	Armor       int  `json:"armor"`
	Helmet      bool `json:"helmet""`
	Flashed     int  `json:"flashed"`
	Smoked      int  `json:"smoked"`
	Burning     int  `json:"burning"`
	Money       int  `json:"money"`
	RoundKills  int  `json:"round_kills"`
	RoundKillHS int  `json:"round_killhs"`
	EquipValue  int  `json:"equip_value"`
}

/*
"weapons": {
			"weapon_0": {
				"name": "weapon_knife_t",
				"paintkit": "default",
				"type": "Knife",
				"state": "holstered"
			},
			"weapon_1": {
				"name": "weapon_glock",
				"paintkit": "aq_glock18_flames_blue",
				"type": "Pistol",
				"ammo_clip": 20,
				"ammo_clip_max": 20,
				"ammo_reserve": 120,
				"state": "holstered"
			},
			"weapon_2": {
				"name": "weapon_c4",
				"paintkit": "default",
				"type": "C4",
				"state": "active"
			}
		}
*/
type PlayerWeapons map[string]Weapon

type Weapon struct {
	Name        string `json:"name"`
	Paintkit    string `json:"paintkit"`
	Type        string `json:"type"`
	AmmoClip    int    `json:"ammo_clip"`
	AmmoClipMax int    `json:"ammo_clip_max"`
	AmmoReserve int    `json:"ammo_reserve"`
	// "holstered" or "active"
	State string `json:"state"`
}

/*
	"round": {
		"phase": "live",
		"bomb": "planted"
	},
*/
// Bomb is "" before bomb is planted, then becomes "planted" when planted.
type Round struct {
	// over -> freezetime -> live -> freezetime
	Phase string `json:"phase"`
	// "" or "planted"
	Bomb string `json:"bomb"`
}

/*
First event: weapons change (no more c4)
Second event: bomb planted in game state
{{76561197993200126 1612023326} {playing {100 100 false 0 0 0 150 0 0 850} map[weapon_0:{weapon_knife_t default Knife 0 0 0 holstered} weapon_1:{weapon_glock aq_glock18_flames_blue Pistol 20 20 120 active} weapon_2:{weapon_c4 default C4 0 0 0 holstered}]} {live }}
{{76561197993200126 1612023340} {playing {100 100 false 0 0 0 150 0 0 850} map[weapon_0:{weapon_knife_t default Knife 0 0 0 holstered} weapon_1:{weapon_glock aq_glock18_flames_blue Pistol 20 20 120 holstered} weapon_2:{weapon_c4 default C4 0 0 0 active}]} {live }}
{{76561197993200126 1612023348} {playing {100 100 false 0 0 0 450 0 0 850} map[weapon_0:{weapon_knife_t default Knife 0 0 0 holstered} weapon_1:{weapon_glock aq_glock18_flames_blue Pistol 20 20 120 active}]} {live }}
{{76561197993200126 1612023349} {playing {100 100 false 0 0 0 450 0 0 850} map[weapon_0:{weapon_knife_t default Knife 0 0 0 holstered} weapon_1:{weapon_glock aq_glock18_flames_blue Pistol 20 20 120 active}]} {live planted}}
{{76561197993200126 1612023351} {playing {100 100 false 0 0 0 450 0 0 850} map[weapon_0:{weapon_knife_t default Knife 0 0 0 active} weapon_1:{weapon_glock aq_glock18_flames_blue Pistol 20 20 120 holstered}]} {live planted}}
{{76561197993200126 1612023362} {playing {84 91 false 0 0 0 450 0 0 850} map[weapon_0:{weapon_knife_t default Knife 0 0 0 active} weapon_1:{weapon_glock aq_glock18_flames_blue Pistol 20 20 120 holstered}]} {live planted}}
{{76561197993200126 1612023363} {playing {60 91 false 0 0 0 450 0 0 850} map[weapon_0:{weapon_knife_t default Knife 0 0 0 active} weapon_1:{weapon_glock aq_glock18_flames_blue Pistol 20 20 120 holstered}]} {live planted}}
{{76561197993200126 1612023363} {playing {36 91 false 0 0 0 450 0 0 850} map[weapon_0:{weapon_knife_t default Knife 0 0 0 active} weapon_1:{weapon_glock aq_glock18_flames_blue Pistol 20 20 120 holstered}]} {live planted}}
{{76561197993200126 1612023363} {playing {12 91 false 0 0 0 450 0 0 850} map[weapon_0:{weapon_knife_t default Knife 0 0 0 active} weapon_1:{weapon_glock aq_glock18_flames_blue Pistol 20 20 120 holstered}]} {live planted}}

*/

/*
apon_c4 default C4 0 0 0 holstered}]} {live }}
{{76561197993200126 1612023901} {playing {0 0 false 0 0 0 10250 2 0 6250} map[]} {live }}

*/

/*
{
	"provider": {
		"name": "Counter-Strike: Global Offensive",
		"appid": 730,
		"version": 13779,
		"steamid": "76561197993200126",
		"timestamp": 1612019048
	},
	"player": {
		"steamid": "76561197993200126",
		"clan": "VUKUKARSK",
		"name": "Andy",
		"observer_slot": 6,
		"team": "T",
		"activity": "playing",
		"state": {
			"health": 100,
			"armor": 100,
			"helmet": false,
			"flashed": 0,
			"smoked": 0,
			"burning": 0,
			"money": 150,
			"round_kills": 0,
			"round_killhs": 0,
			"equip_value": 850
		},
		"match_stats": {
			"kills": 0,
			"assists": 0,
			"deaths": 0,
			"mvps": 0,
			"score": 0
		},
		"weapons": {
			"weapon_0": {
				"name": "weapon_knife_t",
				"paintkit": "default",
				"type": "Knife",
				"state": "holstered"
			},
			"weapon_1": {
				"name": "weapon_glock",
				"paintkit": "aq_glock18_flames_blue",
				"type": "Pistol",
				"ammo_clip": 20,
				"ammo_clip_max": 20,
				"ammo_reserve": 120,
				"state": "holstered"
			},
			"weapon_2": {
				"name": "weapon_c4",
				"paintkit": "default",
				"type": "C4",
				"state": "active"
			}
		}
	},
*/
