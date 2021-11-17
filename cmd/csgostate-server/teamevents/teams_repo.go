package teamevents

import (
	"github.com/roessland/csgostate/cmd/csgostate-server/playerevents"
	"github.com/roessland/csgostate/csgostate"
	"sort"
	"strings"
	"sync"
)

type TeamsRepo struct {
	teamAssember *TeamAssembler
	teams        map[int]*teamInfo
	playerTeam   map[string]int
	teamsMutex   sync.Mutex
	teamEvents   *EventRepo
	playerEvents *playerevents.EventRepo
}

func NewTeamsRepo(teamAssembler *TeamAssembler, teamEvents *EventRepo, playerEvents *playerevents.EventRepo) *TeamsRepo {
	tr := &TeamsRepo{}
	tr.teamAssember = teamAssembler
	tr.teams = make(map[int]*teamInfo)
	tr.playerTeam = make(map[string]int)
	tr.teamEvents = teamEvents
	tr.playerEvents = playerEvents
	tr.registerEventHandlers()
	return tr
}

func (tr *TeamsRepo) GetTeamForPlayer(playerID string) *teamInfo {
	teamID := tr.getTeamForPlayer(playerID)
	return tr.teams[teamID]
}

func (tr *TeamsRepo) registerEventHandlers() {
	tr.teamEvents.Created.Register(func(payload CreatedPayload) {
		tr.teamsMutex.Lock()
		defer tr.teamsMutex.Unlock()
		tr.createTeamIfNotExists(payload.TeamID)
	})

	tr.teamEvents.PlayerJoined.Register(func(payload PlayerJoinedPayload) {
		tr.teamsMutex.Lock()
		defer tr.teamsMutex.Unlock()
		tr.createTeamIfNotExists(payload.TeamID)
		tr.addPlayerToTeam(payload.TeamID, payload.PlayerID, payload.PlayerNick)
	})
}

func (tr *TeamsRepo) Feed(state *csgostate.State) {
	if state.Player == nil {
		return
	}
	if state.Provider.SteamID != state.Player.SteamID {
		// We could store state of spectated player, but it's
		// easier to only store self-reported states.
		return
	}
	playerID := state.Player.SteamID
	playerTeam := tr.getTeamForPlayer(playerID)
	if playerTeam == 0 {
		// Player does't have a team yet
		return
	}

	tr.teams[playerTeam].playersMutex.Lock()
	playerInfo := tr.teams[playerTeam].players[playerID]
	playerInfo.lastState = state
	playerInfo.nick = state.Player.Name
	tr.teams[playerTeam].playersMutex.Unlock()

	tr.teams[playerTeam].lastPhaseMutex.Lock()
	tr.teams[playerTeam].updateLastRoundPhase(state, tr.teamEvents)
	tr.teams[playerTeam].lastPhaseMutex.Unlock()
}

func (tr *TeamsRepo) addPlayerToTeam(teamID int, playerID string, playerNick string) {
	tr.createTeamIfNotExists(teamID)
	tr.teams[teamID].playersMutex.Lock()
	defer tr.teams[teamID].playersMutex.Unlock()
	tr.teams[teamID].addPlayer(playerID, playerNick)
	tr.playerTeam[playerID] = teamID
}

func (tr *TeamsRepo) createTeamIfNotExists(teamID int) {
	if tr.teams[teamID] == nil {
		tr.teams[teamID] = &teamInfo{
			players: make(map[string]*playerInfo),
		}
	}
}

func (tr *TeamsRepo) getTeamForPlayer(playerID string) int {
	return tr.playerTeam[playerID]
}

type teamInfo struct {
	playersMutex   sync.Mutex
	players        map[string]*playerInfo
	lastPhaseMutex sync.Mutex
	lastRoundPhase csgostate.RoundPhase
}

func (ti *teamInfo) addPlayer(id, nick string) {
	ti.players[id] = newPlayer(id, nick)
}

func (ti *teamInfo) updateLastRoundPhase(state *csgostate.State, teamEvents *EventRepo) {
	if state.Round == nil {
		return
	}
	currRoundPhase := state.Round.Phase
	lastRoundPhase := ti.lastRoundPhase
	if lastRoundPhase != currRoundPhase {
		teamEvents.RoundPhaseChanged.Trigger(RoundPhaseChangedPayload{
			From:      lastRoundPhase,
			To:        currRoundPhase,
			CurrState: state,
		})
	}
	ti.lastRoundPhase = currRoundPhase
}

func (ti *teamInfo) String() string {
	if ti == nil {
		return "<nil>"
	}
	var names []string
	for playerID, playerInfo := range ti.players {
		if playerInfo.nick != "" {
			names = append(names, playerInfo.nick)
		} else {
			names = append(names, playerID)
		}
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

type playerInfo struct {
	id        string
	nick      string
	lastState *csgostate.State
}

func newPlayer(id, nick string) *playerInfo {
	player := playerInfo{}
	player.id = id
	player.nick = nick
	return &player
}
