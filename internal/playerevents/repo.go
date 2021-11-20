package playerevents

type EventRepo struct {
	Appeared   appeared
	Died       died
	Spawned    spawned
	Spectating spectating
}

func NewRepo() *EventRepo {
	repo := &EventRepo{}
	repo.Appeared = appeared{}
	repo.Died = died{}
	repo.Spawned = spawned{}
	repo.Spectating = spectating{}

	return repo
}
