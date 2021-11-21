# CS:GO State

[![Tests](https://github.com/roessland/csgostate/actions/workflows/tests.yml/badge.svg)](https://github.com/roessland/csgostate/actions/workflows/tests.yml) [![Production deploy](https://github.com/roessland/csgostate/actions/workflows/production-deploy.yml/badge.svg)](https://github.com/roessland/csgostate/actions/workflows/production-deploy.yml)

Live demo:
https://csgostate.roessland.com/

Counter-Strike: Global Offensive Game State Integration client and library
written in Go.

The primary functionality is to listen to gamestate transitions, and share
these with an API.

You could for example make an IFTTT trigger based on CS:GO events, forward
health/cash/weapon information to the in-game leader, pick a suitable strategy
based on the entire team's weapon/cash information, switch your smart lighting
to a red color scheme when you have low HP, and countless other near-useless
scenarios.

## Components

### csgostate

Library to parse incoming JSON. Can also help with installing the gamestate.cfg
file automatically.

### cmd/csgostate-printer

Acts as a gamestate target API. Records all incoming data to a JSON file.

### cmd/csgostate-pusher

Replays a recorded JSON file in real time, pushing each line to an API.

### cmd/csgostate-server

Web server and gamestate events listener. Steam login. Generate the
gamestate.cfg file with a per-user token so incoming messages can be connected
to a registered user. User registration happens automatically on first login. I
intend to add some kind of team functionality so that events from team members
are shown on the same screen. For example strategies.

## Game State Integration (GSI) information

* [Valve Developer Community Wiki page on CS:GO GSI](https://developer.valvesoftware.com/wiki/Counter-Strike:_Global_Offensive_Game_State_Integration#Game_State_Components)
* [In-depth explanation by Bkid on Reddit](https://www.reddit.com/r/GlobalOffensive/comments/cjhcpy/game_state_integration_a_very_large_and_indepth/)
* [go-csgi](https://github.com/dank/go-csgsi): Another library written in Go
* [csgo-gsi-events](https://github.com/tsuriga/csgo-gsi-events/blob/master/src/csgo-event-emitter.js):
  Emit certain sequences of messages as more useful events.
