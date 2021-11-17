package teamevents

type EventRepo struct {
	Created           created
	PlayerJoined      playerJoined
	Spawned           spawned
	RoundPhaseChanged roundPhaseChanged
}

func NewRepo() *EventRepo {
	repo := &EventRepo{}
	repo.Created = created{}
	repo.PlayerJoined = playerJoined{}
	repo.Spawned = spawned{}
	repo.RoundPhaseChanged = roundPhaseChanged{}
	return repo
}
