package scorekeeper

import (
	"encoding/json"
)

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
// It returns those statistics as a json-encoded string
func (sk *ScoreKeeper) GetStats() (string, error) {
	avgs, err := sk.get()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(avgs)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// keep a Score in the ScoreStore, using a memory store if none was specified.
func (sk *ScoreKeeper) keep(score Score) error {
	if sk.s == nil {
		sk.s = &MemoryStore{
			s: map[string][]Score{},
		}
	}

	return sk.s.Store(score)
}

// get scores from the store and compute averages
func (sk *ScoreKeeper) get() ([]AverageTime, error) {
	if sk.s == nil {
		return nil, ErrNoData
	}

	scoreMap, err := sk.s.Retrieve()
	if err != nil {
		return nil, err
	}

	avgs := make([]AverageTime, 0, len(scoreMap))

	for name, scores := range scoreMap {
		a := Average{}

		res, err := a.Compute(scores)
		if err != nil {
			return nil, err
		}

		avg, ok := res.(float64)
		if !ok {
			return nil, ErrTypeInvalid
		}

		avgs = append(avgs, AverageTime{
			Action:  name,
			Average: avg,
		})
	}

	return avgs, nil
}
