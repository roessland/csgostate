package teamevents

type PlayerJoinedPayload struct {
	TeamID     int
	PlayerID   string
	PlayerNick string // Only provided when spectated player joins a team.
}

type playerJoined struct {
	handlers []func(payload PlayerJoinedPayload)
}

func (e *playerJoined) String() string {
	return "team_playerjoined"
}

func (e *playerJoined) Register(handler func(payload PlayerJoinedPayload)) {
	e.handlers = append(e.handlers, handler)
}

func (e *playerJoined) Trigger(payload PlayerJoinedPayload) {
	for _, handler := range e.handlers {
		go handler(payload)
	}
}
