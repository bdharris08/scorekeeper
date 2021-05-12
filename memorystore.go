package scorekeeper

import "errors"

// MemoryStore keeps scores in memory.
// It will be used if no other store is provided.
// Organize scores in labeled lists.
type MemoryStore struct {
	s map[string][]Score
}

// Store a Score in memory.
func (ms *MemoryStore) Store(s Score) error {
	if ms.s == nil {
		ms.s = map[string][]Score{}
	}

	ms.s[s.Name()] = append(ms.s[s.Name()], s)

	return nil
}

var ErrNoScores = errors.New("no scores found")

// Retrieve Scores from memory by name.
func (ms *MemoryStore) Retrieve(name string) ([]Score, error) {
	if ms.s == nil {
		return nil, ErrNoScores
	}

	scores, ok := ms.s[name]
	if !ok {
		return nil, ErrNoScores
	}

	return scores, nil
}
