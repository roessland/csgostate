package teamevents

import (
	"github.com/roessland/csgostate/cmd/csgostate-server/playerevents"
	"sync"
)

type TeamAssembler struct {
	sync.Mutex
	playerEvents  *playerevents.EventRepo
	teamEvents    *EventRepo
	lastTeamId    int
	teamForPlayer map[string]int
}

func NewTeamAssembler(playerEvents *playerevents.EventRepo, teamEvents *EventRepo) *TeamAssembler {
	ta := &TeamAssembler{}
	ta.teamForPlayer = make(map[string]int)
	ta.playerEvents = playerEvents
	ta.teamEvents = teamEvents
	ta.registerEventHandlers()
	return ta
}

func (ta *TeamAssembler) GetTeamForPlayer(playerID string) int {
	return ta.teamForPlayer[playerID]
}

func (ta *TeamAssembler) registerEventHandlers() {
	ta.playerEvents.Spectating.Register(ta.handleSpectating)
}

func (ta *TeamAssembler) handleSpectating(payload playerevents.SpectatingPayload) {
	ta.Lock()
	defer ta.Unlock()

	playerId := payload.CurrState.Provider.SteamID
	spectatedId := payload.CurrState.Player.SteamID
	playerTeamId := ta.teamForPlayer[playerId]
	spectatedTeamId := ta.teamForPlayer[spectatedId]
	spectatedNick := payload.CurrState.Player.Name

	if playerTeamId == 0 && spectatedTeamId == 0 {
		// Nobody has a team.
		// Create a team, put both players in it.
		teamId := ta.createTeam()
		ta.teamEvents.Created.Trigger(CreatedPayload{TeamID: ta.lastTeamId})
		ta.teamForPlayer[playerId] = teamId
		ta.teamForPlayer[spectatedId] = teamId
		ta.teamEvents.PlayerJoined.Trigger(PlayerJoinedPayload{TeamID: teamId, PlayerID: playerId})
		ta.teamEvents.PlayerJoined.Trigger(PlayerJoinedPayload{TeamID: teamId, PlayerID: spectatedId, PlayerNick: spectatedNick})
	} else if playerTeamId == 0 && spectatedTeamId != 0 {
		// Spectated has a team but player does not.
		// Put player on spectated's team.
		ta.teamForPlayer[playerId] = spectatedTeamId
		ta.teamEvents.PlayerJoined.Trigger(PlayerJoinedPayload{TeamID: spectatedTeamId, PlayerID: playerId})
	} else if playerTeamId != 0 && spectatedTeamId == 0 {
		// Player has a team but spectated does not.
		// Put spectated on player's team.
		ta.teamForPlayer[spectatedId] = playerTeamId
		ta.teamEvents.PlayerJoined.Trigger(PlayerJoinedPayload{TeamID: playerTeamId, PlayerID: spectatedId, PlayerNick: spectatedNick})
	} else if playerTeamId == spectatedTeamId {
		// No need to do anything, they are already on same team.
	} else if playerTeamId != spectatedTeamId {
		// Expected player and spectated to be on same team,
		// but they are already on two separate teams.
		// Move spectated from existing team to player's team.
		ta.teamForPlayer[spectatedId] = playerTeamId
		ta.teamEvents.PlayerJoined.Trigger(PlayerJoinedPayload{TeamID: playerTeamId, PlayerID: spectatedId, PlayerNick: spectatedNick})
	} else {
		panic("should never happen")
	}
}

func (ta *TeamAssembler) createTeam() int {
	ta.lastTeamId++
	return ta.lastTeamId
}
