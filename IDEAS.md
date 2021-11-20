# Ideas for CS:GO State

* Make a proper event based system

* Make graceful update/restart possible without dropping requests.
	* Systemd stuff?
	* Save whatever state to file, load it again upon restart
	* Run migrations automatically

* Move states to separate DB, only keep user state and important stuff in main
  the DB.

* Make CI/CD with autorelease

* Make systemd config file

* CS:GO State client binary that runs on user's computer -- This way we can
  turn up the sampling rate (since state updates are sent only to localhost)
  and deliver more accurate events with less bandwidth, and also dynamically
  update the CFG file based on online state. Of course this would have to be
  only for the most engaged users since running a random binary on your
  computer feels unsafe.
	* But the architecture should be setup such that this "client" runs
	  server-side for every player, and can later be swapped with running the
	  binary on each computer.
	* Implications: Incoming events need auth, and a serialization format.
	  Event generation should run as a goroutine for each user, and accept
	  gamestate updates directly from a http handler.

* Make it possible to create teams
	* Invitation system
	* Suggest friends based on people with same
		* Round/score/Map

* Make it possible to add message service to a user profile
	* If you are the team leader you can choose the service being used
	* Possible services:
		* Discord voice bot
		* Discord channel message via webhook

* Add some kind of "start" button where any person in the team can force-add
  friends (if they are playing same side/map) to the group.

* Strat roulette for a group.

* "Diss bot" -- calling out players reloading or throwing nades while getting
  killed.

* When rolling session secret, log users out instead of cryptic error message
  "unable to get session"
  
* Make sure DB is u+rw,g-rwx,o-rwx
