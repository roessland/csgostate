package teamevents

type CreatedPayload struct {
	TeamID int
}

type created struct {
	handlers []func(payload CreatedPayload)
}

func (e *created) String() string {
	return "team_created"
}

func (e *created) Register(handler func(payload CreatedPayload)) {
	e.handlers = append(e.handlers, handler)
}

func (e *created) Trigger(payload CreatedPayload) {
	for _, handler := range e.handlers {
		go handler(payload)
	}
}
