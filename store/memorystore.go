package store

import (
	"errors"

	"github.com/bdharris08/scorekeeper/score"
)

// MemoryStore keeps scores in memory.
// It will be used if no other store is provided.
// Organize scores in labeled lists.
type MemoryStore struct {
	S map[string]map[string][]score.Score
}

// Store a Score in memory.
func (ms *MemoryStore) Store(s score.Score) error {
	if ms.S == nil {
		ms.S = map[string]map[string][]score.Score{}
	}

	t := s.Type()
	n := s.Name()

	if ms.S[t] == nil {
		ms.S[t] = map[string][]score.Score{}
	}

	ms.S[t][n] = append(ms.S[t][n], s)
	return nil
}

var ErrNoScores = errors.New("no scores found")

// Retrieve Scores from memory by name.
func (ms *MemoryStore) Retrieve(f score.ScoreFactory, scoreType string) (map[string][]score.Score, error) {
	if ms.S == nil {
		return nil, ErrNoScores
	}

	return ms.S[scoreType], nil
}
