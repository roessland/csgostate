package playerevents

import (
	"github.com/pkg/errors"
	"github.com/roessland/csgostate/csgostate"
)

type Extractor struct {
	players      map[string]playerState
	numEventsFed int
}

type playerState struct {
	prevState, currState *csgostate.State
}

func NewExtractor() *Extractor {
	e := &Extractor{}
	e.players = make(map[string]playerState)
	return e
}

func (e *Extractor) Feed(state *csgostate.State) error {
	e.numEventsFed++

	if state == nil {
		return errors.New("state was nil")
	}
	if state.Provider == nil {

		return errors.New("provider was nil")
	}
	if state.Provider.SteamID == "" {
		return errors.New("provider steamid was empty")
	}
	providerSteamID := state.Provider.SteamID
	newPlayerState := playerState{
		prevState: e.players[providerSteamID].currState,
		currState: state,
	}
	e.players[providerSteamID] = newPlayerState
	err := extractAll(newPlayerState.prevState, newPlayerState.currState)
	if err != nil {
		return errors.Wrap(err, "error during events extraction")
	}
	return nil
}

func extractAll(prevState, currState *csgostate.State) error {
	err := checkProvidersMatch(prevState, currState)
	if err != nil {
		return errors.Wrap(err, "provider mismatch")
	}

	err = extractAppeared(prevState, currState)
	if err != nil {
		return errors.Wrap(err, "error extracting appeared event")
	}

	err = extractSpawned(prevState, currState)
	if err != nil {
		return errors.Wrap(err, "error extracting spawned event")
	}

	err = extractSpectating(prevState, currState)
	if err != nil {
		return errors.Wrap(err, "error extracting spectating event")
	}

	err = extractDied(prevState, currState)
	if err != nil {
		return errors.Wrap(err, "error extracting died event")
	}

	return nil
}

func checkProvidersMatch(prevState, currState *csgostate.State) error {
	if prevState == nil {
		return nil
	}
	if prevState.Provider == nil || currState.Provider == nil {
		return errors.New("missing provider in either event")
	}

	if prevState.Provider.SteamID != currState.Provider.SteamID {
		return errors.New("got events from different providers")
	}
	return nil
}
