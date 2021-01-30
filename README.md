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