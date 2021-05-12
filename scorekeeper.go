package scorekeeper

// ScoreKeeper keeps scores using some store.
// It is the top level object for the ScoreKeeper library.
type ScoreKeeper struct {
	store ScoreStore
}

// ScoreStore stores scores for ScoreKeeper.
// It could be in memory or backed by a database.
type ScoreStore interface {
	Store() error
}

// Score is kept by ScoreKeeper and tracks something.
// We will only have one kind of Score for this project,
// but through this interface we could extend to other kinds easily
type Score interface {
	Read(s string) error
}

func (sk *ScoreKeeper) AddAction(action string) error {
	var s Trial
	return s.Read(action)
}

// GetStats computes some statistics about the actions stored in the ScoreKeeper
func (s *ScoreKeeper) GetStats() string {
	return ""
}
