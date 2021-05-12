package scorekeeper

// ScoreKeeper keeps scores using some store.
// It is the top level object for the ScoreKeeper library.
// It fulfills the requested AddAction and GetStats methods.
type ScoreKeeper struct {
	s ScoreStore
}

// AddAction takes a json-encoded string and keeps it for later.
func (sk *ScoreKeeper) AddAction(action string) error {
	var s Trial
	if err := s.Read(action); err != nil {
		return err
	}

	return sk.keep(&s)
}

// GetStats computes some statistics about the actions stored in the ScoreKeeper.
func (s *ScoreKeeper) GetStats() string {
	return ""
}

// Keep a Score in the ScoreStore, using a memory store if none was specified.
func (sk *ScoreKeeper) keep(score Score) error {
	if sk.s == nil {
		sk.s = &MemoryStore{
			s: map[string][]Score{},
		}
	}

	return sk.s.Store(score)
}
