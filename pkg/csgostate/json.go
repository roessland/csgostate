package csgostate

type State struct {
	Added      interface{} `json:"added,omitempty"`
	Auth       *Auth       `json:"auth,omitempty"`
	Map        *Map        `json:"map,omitempty"`
	Player     *Player     `json:"player,omitempty"`
	Previously interface{} `json:"previously,omitempty"`
	Provider   *Provider   `json:"provider,omitempty"`
	Round      *Round      `json:"round,omitempty"`
	RawJson    []byte      `json:"-"`
}

// Auth provides a way of authenticating gamestate messages.
// If auth token matches SteamID in backend we know who sent the message.
type Auth struct {
	Nick  string `json:"nick,omitempty"`
	Token string `json:"token,omitempty"`
}

type Map struct {
	CurrentSpectators     int               `json:"current_spectators"`
	Mode                  string            `json:"mode"`
	Name                  string            `json:"name"`
	NumMatchesToWinSeries int               `json:"num_matches_to_win_series"`
	Phase                 string            `json:"phase"`
	Round                 int               `json:"round"`
	RoundWins             map[string]string `json:"round_wins,omitempty"`
	SouvenirsTotal        int               `json:"souvenirs_total"`
	TeamCT                Team              `json:"team_ct"`
	TeamT                 Team              `json:"team_t"`
}

type Player struct {
	Activity     PlayerActivity    `json:"activity,omitempty"`
	Clan         string            `json:"clan,omitempty"`
	MatchStats   *PlayerMatchStats `json:"match_stats,omitempty"`
	Name         string            `json:"name,omitempty"`
	ObserverSlot *int              `json:"observer_slot,omitempty"`
	State        *PlayerState      `json:"state,omitempty"`
	SteamID      string            `json:"steamid,omitempty"`
	Team         PlayerTeam        `json:"team,omitempty"`
	Weapons      *PlayerWeapons    `json:"weapons,omitempty"`
}

const PlayerNameUnconnected = "unconnected"

type PlayerActivity string

const PlayerActivityPlaying PlayerActivity = "playing"
const PlayerActivityMenu PlayerActivity = "menu"

type PlayerTeam string

const PlayerTeamNil PlayerTeam = ""
const PlayerTeamCT PlayerTeam = "CT"
const PlayerTeamT PlayerTeam = "T"

func PlayerTeamFromString(name string) PlayerTeam {
	switch PlayerTeam(name) {
	case PlayerTeamCT:
		return PlayerTeamCT
	case PlayerTeamT:
		return PlayerTeamT
	default:
		return PlayerTeamNil
	}
}


type Provider struct {
	AppID     int    `json:"appid"`
	Name      string `json:"name"`
	SteamID   string `json:"steamid"`
	Timestamp int    `json:"timestamp"`
	Version   int    `json:"version"`
}

type Team struct {
	ConsecutiveRoundLosses int    `json:"consecutive_round_losses"`
	MatchesWonThisSeries   int    `json:"matches_won_this_series"`
	Name                   string `json:"name,omitempty"`
	Score                  int    `json:"score"`
	TimeoutsRemaining      int    `json:"timeouts_remaining"`
}

type PlayerState struct {
	Armor       int   `json:"armor"`
	Burning     int   `json:"burning"`
	Defusekit   *bool `json:"defusekit,omitempty"`
	EquipValue  int   `json:"equip_value"`
	Flashed     int   `json:"flashed"`
	Health      int   `json:"health"`
	Helmet      bool  `json:"helmet"`
	Money       int   `json:"money"`
	RoundKillHS int   `json:"round_killhs"`
	RoundKills  int   `json:"round_kills"`
	Smoked      int   `json:"smoked"`
}

type PlayerWeapons map[string]Weapon

type Weapon struct {
	AmmoClip    *int        `json:"ammo_clip,omitempty"`
	AmmoClipMax *int        `json:"ammo_clip_max,omitempty"`
	AmmoReserve *int        `json:"ammo_reserve,omitempty"`
	Name        string      `json:"name"`
	Paintkit    string      `json:"paintkit"`
	State       WeaponState `json:"state"`
	Type        string      `json:"type"`
}

// If a weapon is held and active is is either Active or Reloading.
type WeaponState string

const WeaponStateNil WeaponState = ""
const WeaponStateHolstered WeaponState = "holstered"
const WeaponStateActive WeaponState = "active"
const WeaponStateReloading WeaponState = "reloading"

type PlayerMatchStats struct {
	Assists int `json:"assists"`
	Deaths  int `json:"deaths"`
	Kills   int `json:"kills"`
	MVPs    int `json:"mvps"`
	Score   int `json:"score"`
}

type Round struct {
	// Bomb is "" before bomb is planted, then becomes "planted" when planted, then "exploded"
	Bomb string `json:"bomb,omitempty"`
	// Phase is over -> freezetime -> live -> freezetime -> over
	Phase   RoundPhase   `json:"phase"`
	WinTeam RoundWinTeam `json:"win_team,omitempty"` // Missing, T or CT
}

type RoundPhase string

const RoundPhaseOver RoundPhase = "over"
const RoundPhaseFreezetime RoundPhase = "freezetime"
const RoundPhaseLive RoundPhase = "live"

type RoundWinTeam string

const RoundWinTeamCT RoundWinTeam = "CT"
const RoundWinTeamT RoundWinTeam = "T"
