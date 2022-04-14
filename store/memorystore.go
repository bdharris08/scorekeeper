package store

import (
	"errors"

	"github.com/bdharris08/scorekeeper/score"
)

// MemoryStore keeps scores in memory.
// It will be used if no other store is provided.
// Organize scores in labeled lists.
type MemoryStore struct {
	S map[string][]score.Score
}

// Store a Score in memory.
func (ms *MemoryStore) Store(s score.Score) error {
	if ms.S == nil {
		ms.S = map[string][]score.Score{}
	}

	ms.S[s.Name()] = append(ms.S[s.Name()], s)

	return nil
}

var ErrNoScores = errors.New("no scores found")

// Retrieve Scores from memory by name.
func (ms *MemoryStore) Retrieve() (map[string][]score.Score, error) {
	if ms.S == nil {
		return nil, ErrNoScores
	}

	return ms.S, nil
}

// Names returns the score category names currently stored
func (ms *MemoryStore) Names() []string {
	names := make([]string, 0, len(ms.S))

	for name := range ms.S {
		names = append(names, name)
	}

	return names
}
