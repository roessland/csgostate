# CS:GO State

Counter-Strike: Global Offensive Game State Integration client 
and library written in Go.

The primary functionality is to listen to gamestate transitions,
and share these with an API. 

You could for example make an IFTTT trigger based on CS:GO events,
forward health/cash/weapon information to the in-game leader,
pick a suitable strategy based on the entire team's weapon/cash information,
switch your smart lighting to a red color scheme when you have low HP,
and countless other near-useless scenarios.

## Game State Integration (GSI) information

* [Valve Developer Community Wiki page on CS:GO GSI](https://developer.valvesoftware.com/wiki/Counter-Strike:_Global_Offensive_Game_State_Integration#Game_State_Components)
* [In-depth explanation by Bkid on Reddit](https://www.reddit.com/r/GlobalOffensive/comments/cjhcpy/game_state_integration_a_very_large_and_indepth/)
* [go-csgi](https://github.com/dank/go-csgsi): Another library written in Go
* [csgo-gsi-events](https://github.com/tsuriga/csgo-gsi-events/blob/master/src/csgo-event-emitter.js): Emit certain sequences of messages as more useful events.