package teamevents

import (
	"fmt"
	"github.com/roessland/csgostate/cmd/csgostate-server/playerevents"
	"github.com/roessland/csgostate/csgostate"
	"sync"
)

var ta *teamAssembler

func init() {
	ta = NewTeamAssembler()
}

func NewTeamAssembler() *teamAssembler {
	ta := &teamAssembler{}
	ta.teams = make(map[int]*teamInfo)
	ta.teamForPlayer = make(map[string]int)
	ta.registerEventHandlers()
	return ta
}

func (ta *teamAssembler) registerEventHandlers() {
	playerevents.Spectating.Register(ta.handleSpectating)
}

type teamAssembler struct {
	sync.Mutex
	lastTeamId    int
	teams         map[int]*teamInfo
	teamForPlayer map[string]int
}

func (ta *teamAssembler) handleSpectating(payload playerevents.SpectatingPayload) {
	ta.Lock()
	defer ta.Unlock()

	playerId := payload.CurrState.Provider.SteamID
	spectatedId := payload.CurrState.Player.SteamID
	playerTeamId := ta.teamForPlayer[playerId]
	spectatedTeamId := ta.teamForPlayer[spectatedId]
	if playerTeamId == 0 && spectatedTeamId == 0 {
		// Nobody has a team.
		// Create a team, put both players in it.
		teamId := ta.createTeam()
		ta.teams[teamId].addPlayer(playerId)
		ta.teams[teamId].addPlayer(spectatedId)
		ta.teamForPlayer[playerId] = teamId
		ta.teamForPlayer[spectatedId] = teamId
	} else if playerTeamId == 0 && spectatedTeamId != 0 {
		// Spectated has a team but player does not.
		// Put player on spectated's team.
		ta.teams[spectatedTeamId].addPlayer(playerId)
		ta.teamForPlayer[playerId] = spectatedTeamId
	} else if playerTeamId != 0 && spectatedTeamId == 0 {
		// Player has a team but spectated does not.
		// Put spectated on player's team.
		ta.teams[playerTeamId].addPlayer(spectatedId)
		ta.teamForPlayer[spectatedId] = playerTeamId
	} else if playerTeamId == spectatedTeamId {
		// No need to do anything, they are already on same team.
	} else if playerTeamId != spectatedTeamId {
		// Expected player and spectated to be on same team,
		// but they are already on two separate teams.
		// Delete spectated's team and add them to players's team.
		ta.deleteTeam(spectatedTeamId)
		ta.teams[playerTeamId].addPlayer(spectatedId)
		ta.teamForPlayer[spectatedId] = playerTeamId
	} else {
		panic("should never happen")
	}
}

func (ta *teamAssembler) createTeam() int {
	ta.Lock()
	defer ta.Unlock()

	ta.lastTeamId++
	ta.teams[ta.lastTeamId] = newTeam()
	return ta.lastTeamId
}

func (ta *teamAssembler) deleteTeam(teamID int) {
	ta.Lock()
	defer ta.Unlock()

	delete(ta.teams, teamID)
}

type teamInfo struct {
	players map[string]*playerInfo
}

func newTeam() *teamInfo {
	team := teamInfo{}
	team.players = make(map[string]*playerInfo)
	return &team
}

func (ti *teamInfo) addPlayer(steamID string) {
	ti.players[steamID] = newPlayer()
	fmt.Println(len(ti.players))
}

type playerInfo struct {
	lastState *csgostate.State
}

func newPlayer() *playerInfo {
	player := playerInfo{}
	return &player
}
