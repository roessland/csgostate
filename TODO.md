## Todos and memos

For my personal use, but might be useful for someone else too.

Enabling GSI:
Must put a cfg file in:
`\Steam\SteamApps\common\Counter-Strike Global Offensive\csgo\cfg\`

`CfgDir string` option.

Windows: You can find this by reading the registry value.
`HKEY_CURRENT_USER\Software\Valve\Steam`,
then add the rest of the install path.

`SteamPath Key Value + "\SteamApps\common\Counter-Strike Global Offensive\csgo\cfg\"
`

The file should be named `gamestate_integration_YourServiceName.cfg` and
the filename should be unique since it could be
overwritten by other gamestate clients.

For MVP I can just create the file manually.
Later on I should provide a feature to create the file
if it does not exist, and maybe a config file specifying the cfg directory.

The contents of this file depends on whether you are spectating or playing.
Since this library focuses on playing, only the player config is added.
Adding extra spectator fields you don't have access to will result in an
empty or missing field (which isn't really a problem).
```
"Config name goes here, can be anything"
{
    "uri" "http://127.0.0.1:3528"
    "timeout" "5.0"
    "buffer"  "0.1"
    "throttle" "0.1"
    "heartbeat" "30.0"
    "auth"
    {
      "token" "Q79v5tcxVQ8u"
    }
    "data"
    {
      "provider"                "1"
      "player_id"               "1"
      "player_state"            "1"
      "map"                     "1"
      "map_round_wins"          "1"
      "player_match_stats"      "1"
      "player_weapons"          "1"
      "round"                   "1"
      "allgrenades"             "1"
      "allplayers_id"           "1"
      "allplayers_match_stats"  "1"
      "allplayers_position"     "1"
      "allplayers_state"        "1"
      "allplayers_weapons"      "1"
      "bomb"                    "1"
      "phase_countdowns"        "1"
      "player_position"         "1"
    }
}
```

The auth section is added to every request,
and can be used to avoid people sending messages on
behalf of other players. Since this MVP assumes players
are trusted I'll just add a random constant value here.

## Endpoint settings

These are important. In particular, since I want a live-updated
dashboard:
* uri: Use a HTTPS address with SSL. Steam will validate the cert.
  But since we have a middleman client to filter out useless info,
  localhost:3528 should be fine.
* timeout: Client will consider a message timed out if there is no 
  response after this amount of time. The implications of this
  is that the API should buffer incoming messages and immediately
  return 200 OK to Steam, then process messages async.
* buffer: How live are we? Default of 0.1 sec should be fine.
  This clusters messages occurring in a short time interval to save bandwidth
  and connections.
* throttle: Don't send another message this amount of time
  after getting 200 OK from API. Default of 1.0 sec is probably too high,
  consider decreasing to 0.2-0.4 sec.
* heartbeat: Even if no game state change occurs, send a heartbeat.
  This is probably useful. Around 3-5 seconds should be fine.
  
## Filtering data

We can do some filtering in the config file,
by eliminating useless components,  
so that messages of that type are never sent to begin with.

## Event triggers
Terrorist player drops C4 -> Can start bomb-timer.
Restart it whenever someone drops C4, since they may have picked it up again.
Remove counter whenever someone picks up C4.
There is a bomb planted event but it could be delayed???
Hypothesis: Drop C4 event is accurate while bomb planted event is inaccurate.

## Messages

When launching game
```
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612016710
        },
        "player": {
                "steamid": "76561197993200126",
                "name": "Andy",
                "activity": "menu"
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}

```

When connecting to a game
```
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612016988
        },
        "player": {
                "steamid": "76561197993200126",
                "name": "unconnected",
                "activity": "playing",
                "state": {
                        "health": 0,
                        "armor": 0,
                        "helmet": false,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 0
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 0,
                        "mvps": 0,
                        "score": 0
                },
                "weapons": {

                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "warmup",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "previously": {
                "player": {
                        "name": "Andy",
                        "activity": "menu"
                }
        },
        "added": {
                "player": {
                        "state": true,
                        "match_stats": true,
                        "weapons": true
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612016990
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "activity": "playing",
                "state": {
                        "health": 0,
                        "armor": 0,
                        "helmet": false,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 0
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 0,
                        "mvps": 0,
                        "score": 0
                },
                "weapons": {

                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "warmup",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "previously": {
                "player": {
                        "name": "unconnected"
                }
        },
        "added": {
                "player": {
                        "clan": true
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612016990
        },
        "player": {
                "steamid": "76561198034276409",
                "name": "intse",
                "observer_slot": 1,
                "team": "CT",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
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
                                "name": "weapon_knife",
                                "paintkit": "default",
                                "type": "Knife",
                                "state": "holstered"
                        },
                        "weapon_1": {
                                "name": "weapon_usp_silencer",
                                "paintkit": "cu_usp_progressiv",
                                "type": "Pistol",
                                "ammo_clip": 12,
                                "ammo_clip_max": 12,
                                "ammo_reserve": 24,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "steamid": "76561197993200126",
                        "clan": "VUKUKARSK",
                        "name": "Andy",
                        "state": {
                                "health": 0,
                                "armor": 0,
                                "helmet": false,
                                "equip_value": 0
                        }
                },
                "map": {
                        "phase": "warmup"
                }
        },
        "added": {
                "player": {
                        "weapons": {
                                "weapon_0": true,
                                "weapon_1": true
                        },
                        "observer_slot": true,
                        "team": true
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612016992
        },
        "player": {
                "steamid": "76561198034276409",
                "name": "intse",
                "observer_slot": 1,
                "team": "CT",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 3250
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
                                "name": "weapon_knife",
                                "paintkit": "default",
                                "type": "Knife",
                                "state": "holstered"
                        },
                        "weapon_1": {
                                "name": "weapon_usp_silencer",
                                "paintkit": "cu_usp_progressiv",
                                "type": "Pistol",
                                "ammo_clip": 12,
                                "ammo_clip_max": 12,
                                "ammo_reserve": 24,
                                "state": "holstered"
                        },
                        "weapon_2": {
                                "name": "weapon_famas",
                                "paintkit": "am_nuclear_skulls2_famas",
                                "type": "Rifle",
                                "ammo_clip": 25,
                                "ammo_clip_max": 25,
                                "ammo_reserve": 90,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "state": {
                                "equip_value": 1200
                        },
                        "weapons": {
                                "weapon_1": {
                                        "state": "active"
                                }
                        }
                }
        },
        "added": {
                "player": {
                        "weapons": {
                                "weapon_2": true
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612016994
        },
        "player": {
                "steamid": "76561198034276409",
                "name": "intse",
                "observer_slot": 1,
                "team": "CT",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 3550
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
                                "name": "weapon_knife",
                                "paintkit": "default",
                                "type": "Knife",
                                "state": "holstered"
                        },
                        "weapon_1": {
                                "name": "weapon_fiveseven",
                                "paintkit": "default",
                                "type": "Pistol",
                                "ammo_clip": 20,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 100,
                                "state": "holstered"
                        },
                        "weapon_2": {
                                "name": "weapon_famas",
                                "paintkit": "am_nuclear_skulls2_famas",
                                "type": "Rifle",
                                "ammo_clip": 25,
                                "ammo_clip_max": 25,
                                "ammo_reserve": 90,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "state": {
                                "equip_value": 3250
                        },
                        "weapons": {
                                "weapon_1": {
                                        "name": "weapon_usp_silencer",
                                        "paintkit": "cu_usp_progressiv",
                                        "ammo_clip": 12,
                                        "ammo_clip_max": 12,
                                        "ammo_reserve": 24
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612016994
        },
        "player": {
                "steamid": "76561198034276409",
                "name": "intse",
                "observer_slot": 1,
                "team": "CT",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 3550
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
                                "name": "weapon_knife",
                                "paintkit": "default",
                                "type": "Knife",
                                "state": "holstered"
                        },
                        "weapon_1": {
                                "name": "weapon_fiveseven",
                                "paintkit": "default",
                                "type": "Pistol",
                                "ammo_clip": 20,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 100,
                                "state": "active"
                        },
                        "weapon_2": {
                                "name": "weapon_famas",
                                "paintkit": "am_nuclear_skulls2_famas",
                                "type": "Rifle",
                                "ammo_clip": 25,
                                "ammo_clip_max": 25,
                                "ammo_reserve": 90,
                                "state": "holstered"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "state": "holstered"
                                },
                                "weapon_2": {
                                        "state": "active"
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612016994
        },
        "player": {
                "steamid": "76561198034276409",
                "name": "intse",
                "observer_slot": 1,
                "team": "CT",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 3550
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
                                "name": "weapon_knife",
                                "paintkit": "default",
                                "type": "Knife",
                                "state": "holstered"
                        },
                        "weapon_1": {
                                "name": "weapon_fiveseven",
                                "paintkit": "default",
                                "type": "Pistol",
                                "ammo_clip": 20,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 100,
                                "state": "holstered"
                        },
                        "weapon_2": {
                                "name": "weapon_famas",
                                "paintkit": "am_nuclear_skulls2_famas",
                                "type": "Rifle",
                                "ammo_clip": 25,
                                "ammo_clip_max": 25,
                                "ammo_reserve": 90,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "state": "active"
                                },
                                "weapon_2": {
                                        "state": "holstered"
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612016999
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "observer_slot": 0,
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
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
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "steamid": "76561198034276409",
                        "name": "intse",
                        "observer_slot": 1,
                        "team": "CT",
                        "state": {
                                "money": 0,
                                "equip_value": 3550
                        },
                        "weapons": {
                                "weapon_0": {
                                        "name": "weapon_knife"
                                },
                                "weapon_1": {
                                        "name": "weapon_fiveseven",
                                        "paintkit": "default",
                                        "ammo_reserve": 100,
                                        "state": "holstered"
                                },
                                "weapon_2": {
                                        "name": "weapon_famas",
                                        "paintkit": "am_nuclear_skulls2_famas",
                                        "type": "Rifle",
                                        "ammo_clip": 25,
                                        "ammo_clip_max": 25,
                                        "ammo_reserve": 90,
                                        "state": "active"
                                }
                        }
                }
        },
        "added": {
                "player": {
                        "clan": true
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}

```


When idling for 5 seconds in deathmatch
```
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017061
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 1,
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
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017075
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 1,
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
                                "ammo_clip": 19,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 20
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
```


When shooting with Glock and emptying a magazine,

We get a LOT of superfluous information here,
that should be filtered out.

```
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017137
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017139
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 19,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 20
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017140
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 18,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 19
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017140
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 17,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 18
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017141
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 15,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 17
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017141
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 13,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 15
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017141
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 12,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 13
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017141
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 11,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 12
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017142
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 10,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 11
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017142
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 8,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 10
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017142
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 7,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 8
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017143
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 5,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 7
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017143
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 4,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 5
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017143
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 2,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 4
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017143
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 1,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 2
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017143
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "ammo_clip": 0,
                                "ammo_clip_max": 20,
                                "ammo_reserve": 120,
                                "state": "reloading"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 1,
                                        "state": "active"
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017145
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "state": "reloading"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "ammo_clip": 0
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017146
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 34463,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 2,
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
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "state": "reloading"
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}

```


When switching from AK to SMG, then repeatedly switching weapons for a few seconds:
```
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017418
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 3900
                },
                "match_stats": {
                        "kills": 5,
                        "assists": 2,
                        "deaths": 11,
                        "mvps": 0,
                        "score": 68
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
                                "name": "weapon_ak47",
                                "paintkit": "cu_ak47_mastery",
                                "type": "Rifle",
                                "ammo_clip": 30,
                                "ammo_clip_max": 30,
                                "ammo_reserve": 90,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017424
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 2700
                },
                "match_stats": {
                        "kills": 5,
                        "assists": 2,
                        "deaths": 11,
                        "mvps": 0,
                        "score": 68
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
                                "name": "weapon_mp7",
                                "paintkit": "sp_spray_army",
                                "type": "Submachine Gun",
                                "ammo_clip": 30,
                                "ammo_clip_max": 30,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "state": {
                                "equip_value": 3900
                        },
                        "weapons": {
                                "weapon_2": {
                                        "name": "weapon_ak47",
                                        "paintkit": "cu_ak47_mastery",
                                        "type": "Rifle",
                                        "ammo_reserve": 90
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017425
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 2700
                },
                "match_stats": {
                        "kills": 5,
                        "assists": 2,
                        "deaths": 11,
                        "mvps": 0,
                        "score": 68
                },
                "weapons": {
                        "weapon_0": {
                                "name": "weapon_knife_t",
                                "paintkit": "default",
                                "type": "Knife",
                                "state": "active"
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
                                "name": "weapon_mp7",
                                "paintkit": "sp_spray_army",
                                "type": "Submachine Gun",
                                "ammo_clip": 30,
                                "ammo_clip_max": 30,
                                "ammo_reserve": 120,
                                "state": "holstered"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_0": {
                                        "state": "holstered"
                                },
                                "weapon_2": {
                                        "state": "active"
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017425
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 2700
                },
                "match_stats": {
                        "kills": 5,
                        "assists": 2,
                        "deaths": 11,
                        "mvps": 0,
                        "score": 68
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
                                "name": "weapon_mp7",
                                "paintkit": "sp_spray_army",
                                "type": "Submachine Gun",
                                "ammo_clip": 30,
                                "ammo_clip_max": 30,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_0": {
                                        "state": "active"
                                },
                                "weapon_2": {
                                        "state": "holstered"
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017426
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 2700
                },
                "match_stats": {
                        "kills": 5,
                        "assists": 2,
                        "deaths": 11,
                        "mvps": 0,
                        "score": 68
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
                                "state": "active"
                        },
                        "weapon_2": {
                                "name": "weapon_mp7",
                                "paintkit": "sp_spray_army",
                                "type": "Submachine Gun",
                                "ammo_clip": 30,
                                "ammo_clip_max": 30,
                                "ammo_reserve": 120,
                                "state": "holstered"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "state": "holstered"
                                },
                                "weapon_2": {
                                        "state": "active"
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017426
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 2700
                },
                "match_stats": {
                        "kills": 5,
                        "assists": 2,
                        "deaths": 11,
                        "mvps": 0,
                        "score": 68
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
                                "name": "weapon_mp7",
                                "paintkit": "sp_spray_army",
                                "type": "Submachine Gun",
                                "ammo_clip": 30,
                                "ammo_clip_max": 30,
                                "ammo_reserve": 120,
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "player": {
                        "weapons": {
                                "weapon_1": {
                                        "state": "active"
                                },
                                "weapon_2": {
                                        "state": "holstered"
                                }
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}

```

When deathmatch round is over. Entering freezetime. Sending two chat messages.
Joining new deathmatch server. Choosing team. Warmup ends. Round starts.

```
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017593
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "observer_slot": 0,
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 27,
                        "armor": 83,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 4,
                        "round_killhs": 1,
                        "equip_value": 2700
                },
                "match_stats": {
                        "kills": 10,
                        "assists": 3,
                        "deaths": 17,
                        "mvps": 0,
                        "score": 135
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
                                "name": "weapon_mp7",
                                "paintkit": "sp_spray_army",
                                "type": "Submachine Gun",
                                "ammo_clip": 30,
                                "ammo_clip_max": 30,
                                "ammo_reserve": 120,
                                "state": "active"
                        },
                        "weapon_3": {
                                "name": "weapon_healthshot",
                                "paintkit": "default",
                                "type": "StackableItem",
                                "ammo_reserve": 1,
                                "state": "holstered"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "gameover",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 1,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "freezetime",
                "win_team": "CT"
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017606
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "observer_slot": 9,
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 27,
                        "armor": 83,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 4,
                        "round_killhs": 1,
                        "equip_value": 2700
                },
                "match_stats": {
                        "kills": 10,
                        "assists": 3,
                        "deaths": 17,
                        "mvps": 0,
                        "score": 135
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
                                "name": "weapon_mp7",
                                "paintkit": "sp_spray_army",
                                "type": "Submachine Gun",
                                "ammo_clip": 30,
                                "ammo_clip_max": 30,
                                "ammo_reserve": 120,
                                "state": "active"
                        },
                        "weapon_3": {
                                "name": "weapon_healthshot",
                                "paintkit": "default",
                                "type": "StackableItem",
                                "ammo_reserve": 1,
                                "state": "holstered"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "gameover",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 1,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "freezetime",
                "win_team": "CT"
        },
        "previously": {
                "player": {
                        "observer_slot": 0
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017607
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "observer_slot": 8,
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 27,
                        "armor": 83,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 4,
                        "round_killhs": 1,
                        "equip_value": 2700
                },
                "match_stats": {
                        "kills": 10,
                        "assists": 3,
                        "deaths": 17,
                        "mvps": 0,
                        "score": 135
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
                                "name": "weapon_mp7",
                                "paintkit": "sp_spray_army",
                                "type": "Submachine Gun",
                                "ammo_clip": 30,
                                "ammo_clip_max": 30,
                                "ammo_reserve": 120,
                                "state": "active"
                        },
                        "weapon_3": {
                                "name": "weapon_healthshot",
                                "paintkit": "default",
                                "type": "StackableItem",
                                "ammo_reserve": 1,
                                "state": "holstered"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "gameover",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 1,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "freezetime",
                "win_team": "CT"
        },
        "previously": {
                "player": {
                        "observer_slot": 9
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017614
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "observer_slot": 7,
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 27,
                        "armor": 83,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 4,
                        "round_killhs": 1,
                        "equip_value": 2700
                },
                "match_stats": {
                        "kills": 10,
                        "assists": 3,
                        "deaths": 17,
                        "mvps": 0,
                        "score": 135
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
                                "name": "weapon_mp7",
                                "paintkit": "sp_spray_army",
                                "type": "Submachine Gun",
                                "ammo_clip": 30,
                                "ammo_clip_max": 30,
                                "ammo_reserve": 120,
                                "state": "active"
                        },
                        "weapon_3": {
                                "name": "weapon_healthshot",
                                "paintkit": "default",
                                "type": "StackableItem",
                                "ammo_reserve": 1,
                                "state": "holstered"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "gameover",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 1,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "freezetime",
                "win_team": "CT"
        },
        "previously": {
                "player": {
                        "observer_slot": 8
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017615
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "observer_slot": 7,
                "team": "T",
                "activity": "textinput",
                "state": {
                        "health": 27,
                        "armor": 83,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 4,
                        "round_killhs": 1,
                        "equip_value": 2700
                },
                "match_stats": {
                        "kills": 10,
                        "assists": 3,
                        "deaths": 17,
                        "mvps": 0,
                        "score": 135
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
                                "name": "weapon_mp7",
                                "paintkit": "sp_spray_army",
                                "type": "Submachine Gun",
                                "ammo_clip": 30,
                                "ammo_clip_max": 30,
                                "ammo_reserve": 120,
                                "state": "active"
                        },
                        "weapon_3": {
                                "name": "weapon_healthshot",
                                "paintkit": "default",
                                "type": "StackableItem",
                                "ammo_reserve": 1,
                                "state": "holstered"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "gameover",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 1,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "freezetime",
                "win_team": "CT"
        },
        "previously": {
                "player": {
                        "activity": "playing"
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017616
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "observer_slot": 7,
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 27,
                        "armor": 83,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 4,
                        "round_killhs": 1,
                        "equip_value": 2700
                },
                "match_stats": {
                        "kills": 10,
                        "assists": 3,
                        "deaths": 17,
                        "mvps": 0,
                        "score": 135
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
                                "name": "weapon_mp7",
                                "paintkit": "sp_spray_army",
                                "type": "Submachine Gun",
                                "ammo_clip": 30,
                                "ammo_clip_max": 30,
                                "ammo_reserve": 120,
                                "state": "active"
                        },
                        "weapon_3": {
                                "name": "weapon_healthshot",
                                "paintkit": "default",
                                "type": "StackableItem",
                                "ammo_reserve": 1,
                                "state": "holstered"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "gameover",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 1,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "freezetime",
                "win_team": "CT"
        },
        "previously": {
                "player": {
                        "activity": "textinput"
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017617
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "observer_slot": 7,
                "team": "T",
                "activity": "textinput",
                "state": {
                        "health": 27,
                        "armor": 83,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 4,
                        "round_killhs": 1,
                        "equip_value": 2700
                },
                "match_stats": {
                        "kills": 10,
                        "assists": 3,
                        "deaths": 17,
                        "mvps": 0,
                        "score": 135
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
                                "name": "weapon_mp7",
                                "paintkit": "sp_spray_army",
                                "type": "Submachine Gun",
                                "ammo_clip": 30,
                                "ammo_clip_max": 30,
                                "ammo_reserve": 120,
                                "state": "active"
                        },
                        "weapon_3": {
                                "name": "weapon_healthshot",
                                "paintkit": "default",
                                "type": "StackableItem",
                                "ammo_reserve": 1,
                                "state": "holstered"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "gameover",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 1,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "freezetime",
                "win_team": "CT"
        },
        "previously": {
                "player": {
                        "activity": "playing"
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017618
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "observer_slot": 7,
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 27,
                        "armor": 83,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 4,
                        "round_killhs": 1,
                        "equip_value": 2700
                },
                "match_stats": {
                        "kills": 10,
                        "assists": 3,
                        "deaths": 17,
                        "mvps": 0,
                        "score": 135
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
                                "name": "weapon_mp7",
                                "paintkit": "sp_spray_army",
                                "type": "Submachine Gun",
                                "ammo_clip": 30,
                                "ammo_clip_max": 30,
                                "ammo_reserve": 120,
                                "state": "active"
                        },
                        "weapon_3": {
                                "name": "weapon_healthshot",
                                "paintkit": "default",
                                "type": "StackableItem",
                                "ammo_reserve": 1,
                                "state": "holstered"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_mirage",
                "phase": "gameover",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 1,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "freezetime",
                "win_team": "CT"
        },
        "previously": {
                "player": {
                        "activity": "textinput"
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017623
        },
        "player": {
                "steamid": "76561197993200126",
                "name": "Andy",
                "activity": "menu"
        },
        "previously": {
                "player": {
                        "clan": "VUKUKARSK",
                        "observer_slot": 7,
                        "team": "T",
                        "activity": "playing",
                        "state": {
                                "health": 27,
                                "armor": 83,
                                "helmet": true,
                                "flashed": 0,
                                "smoked": 0,
                                "burning": 0,
                                "money": 0,
                                "round_kills": 4,
                                "round_killhs": 1,
                                "equip_value": 2700
                        },
                        "match_stats": {
                                "kills": 10,
                                "assists": 3,
                                "deaths": 17,
                                "mvps": 0,
                                "score": 135
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
                                        "name": "weapon_mp7",
                                        "paintkit": "sp_spray_army",
                                        "type": "Submachine Gun",
                                        "ammo_clip": 30,
                                        "ammo_clip_max": 30,
                                        "ammo_reserve": 120,
                                        "state": "active"
                                },
                                "weapon_3": {
                                        "name": "weapon_healthshot",
                                        "paintkit": "default",
                                        "type": "StackableItem",
                                        "ammo_reserve": 1,
                                        "state": "holstered"
                                }
                        }
                },
                "map": true,
                "round": true
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017633
        },
        "player": {
                "steamid": "76561197993200126",
                "name": "unconnected",
                "activity": "playing",
                "state": {
                        "health": 0,
                        "armor": 0,
                        "helmet": false,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 0
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 0,
                        "mvps": 0,
                        "score": 0
                },
                "weapons": {

                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_inferno",
                "phase": "warmup",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "previously": {
                "player": {
                        "name": "Andy",
                        "activity": "menu"
                }
        },
        "added": {
                "player": {
                        "state": true,
                        "match_stats": true,
                        "weapons": true
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017638
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "activity": "playing",
                "state": {
                        "health": 0,
                        "armor": 0,
                        "helmet": false,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 0
                },
                "match_stats": {
                        "kills": 0,
                        "assists": 0,
                        "deaths": 0,
                        "mvps": 0,
                        "score": 0
                },
                "weapons": {

                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_inferno",
                "phase": "warmup",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "previously": {
                "player": {
                        "name": "unconnected"
                }
        },
        "added": {
                "player": {
                        "clan": true
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017643
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "observer_slot": 4,
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
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
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_inferno",
                "phase": "warmup",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "previously": {
                "player": {
                        "state": {
                                "health": 0,
                                "armor": 0,
                                "helmet": false,
                                "equip_value": 0
                        }
                }
        },
        "added": {
                "player": {
                        "weapons": {
                                "weapon_0": true,
                                "weapon_1": true
                        },
                        "observer_slot": true,
                        "team": true
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017643
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "observer_slot": 0,
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
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
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_inferno",
                "phase": "warmup",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "previously": {
                "player": {
                        "observer_slot": 4
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017643
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "observer_slot": 0,
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
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
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_inferno",
                "phase": "warmup",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "name": "Vuku karsk- og skyttarlaug",
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "added": {
                "map": {
                        "team_t": {
                                "name": true
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017645
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "observer_slot": 0,
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
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
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_inferno",
                "phase": "warmup",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "previously": {
                "map": {
                        "team_t": {
                                "name": "Vuku karsk- og skyttarlaug"
                        }
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
{
        "provider": {
                "name": "Counter-Strike: Global Offensive",
                "appid": 730,
                "version": 13779,
                "steamid": "76561197993200126",
                "timestamp": 1612017649
        },
        "player": {
                "steamid": "76561197993200126",
                "clan": "VUKUKARSK",
                "name": "Andy",
                "observer_slot": 0,
                "team": "T",
                "activity": "playing",
                "state": {
                        "health": 100,
                        "armor": 100,
                        "helmet": true,
                        "flashed": 0,
                        "smoked": 0,
                        "burning": 0,
                        "money": 0,
                        "round_kills": 0,
                        "round_killhs": 0,
                        "equip_value": 1200
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
                                "state": "active"
                        }
                }
        },
        "map": {
                "mode": "deathmatch",
                "name": "de_inferno",
                "phase": "live",
                "round": 0,
                "team_ct": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "team_t": {
                        "score": 0,
                        "consecutive_round_losses": 0,
                        "timeouts_remaining": 1,
                        "matches_won_this_series": 0
                },
                "num_matches_to_win_series": 0,
                "current_spectators": 0,
                "souvenirs_total": 0
        },
        "round": {
                "phase": "live"
        },
        "previously": {
                "map": {
                        "phase": "warmup"
                }
        },
        "auth": {
                "token": "Q79v5tcxVQ8u"
        }
}
```


## Steam OAuth2 login

We can authenticate users in the following way:

1. Let them sign in using Steam OAuth2 provider. We get the steam ID.
2. Generate the cfg file with a random secret token in the auth section.
3. User manually install cfg file.
4. User launches CS:GO. CS:GO sends HTTP request to API.
5. Reverse lookup secret token -> steam ID. It's not valid if steamID in message
is different from the stored steam ID.

This way all requests are connected to a SteamID. 
We can simply drop all unauthenticated requests.

Steam OpenID provider:

    $ curl https://steamcommunity.com/openid
    <?xml version="1.0" encoding="UTF-8"?>
    <xrds:XRDS xmlns:xrds="xri://$xrds" xmlns="xri://$xrd*($v*2.0)">
    <XRD>
    <Service priority="0">
    <Type>http://specs.openid.net/auth/2.0/server</Type>
    <URI>https://steamcommunity.com/openid/login</URI>
    </Service>
    </XRD>
    </xrds:XRDS>

Curling the URI gives this. We are redirected to Steamcommunity.com with a new cookie.

    > GET /openid/login HTTP/1.1
    > Host: steamcommunity.com
    > User-Agent: curl/7.47.0
    > Accept: */*
    >
    < HTTP/1.1 302 Moved Temporarily
    < Server: nginx
    < Content-Type: text/html; charset=UTF-8
    < X-Frame-Options: DENY
    < Content-Security-Policy: default-src blob: data: https: 'unsafe-inline' 'unsafe-eval'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://community.cloudflare.steamstatic.com/ https://cdn.cloudflare.steamstatic.com/steamcommunity/public/assets/ https://api.steampowered.com/ *.google-analytics.com https://www.google.com https://www.gstatic.com https://apis.google.com https://recaptcha.net https://www.gstatic.cn/recaptcha/; object-src 'none'; connect-src 'self' https://community.cloudflare.steamstatic.com/ https://store.steampowered.com/ wss://community.steam-api.com/websocket/ https://api.steampowered.com/ *.google-analytics.com https://*.valvesoftware.com https://*.steambeta.net https://*.steamcontent.com https://steambroadcast.akamaized.net https://steambroadcastchat.akamaized.net https://broadcast.st.dl.bscstorage.net https://broadcast.st.dl.eccdnx.com http://127.0.0.1:27060 ws://127.0.0.1:27060; frame-src 'self' steam: https://store.steampowered.com/ https://www.youtube.com https://www.google.com https://sketchfab.com https://player.vimeo.com https://medal.tv https://www.google.com/recaptcha/ https://recaptcha.net/recaptcha/ https://help.steampowered.com/;
    < Expires: Mon, 26 Jul 1997 05:00:00 GMT
    < Cache-Control: no-cache
    < Location: https://steamcommunity.com
    < Content-Encoding: gzip
    < Vary: Accept-Encoding
    < Content-Length: 20
    < Date: Mon, 01 Feb 2021 17:24:04 GMT
    < Connection: keep-alive
    < Set-Cookie: sessionid=d4350b4b59afe4b20820ef49; Path=/; Secure; SameSite=None
    < Set-Cookie: steamCountry=DE%7C60698fcbe40561212d3a1b665a038208; Path=/; Secure; SameSite=None
