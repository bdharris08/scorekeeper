package scorekeeper

// ScoreKeeper keeps scores using some store.
// It is the top level object for the ScoreKeeper library.
// It fulfills the requested AddAction and GetStats methods.
type ScoreKeeper struct {
	s ScoreStore
}

// ScoreStore stores scores for ScoreKeeper.
// It could be in memory or backed by a database.
type ScoreStore interface {
	Store(s Score) error
	Retrieve(name string) ([]Score, error)
}

// Score is kept by ScoreKeeper and tracks something.
// We will only have one kind of Score for this project,
// but through this interface we could extend to other kinds easily
type Score interface {
	Read(s string) error
}

// Keep a Score in the ScoreStore, using a memory store if none was specified.
func (sk *ScoreKeeper) Keep(score Score) error {
	if sk.s == nil {
		sk.s = &MemoryStore{
			s: map[string][]Score{},
		}
	}

	return sk.s.Store(score)
}

// AddAction takes a json-encoded string and keeps it for later.
func (sk *ScoreKeeper) AddAction(action string) error {
	var s Trial
	if err := s.Read(action); err != nil {
		return err
	}

	return sk.Keep(&s)
}

// GetStats computes some statistics about the actions stored in the ScoreKeeper.
func (s *ScoreKeeper) GetStats() string {
	return ""
}
