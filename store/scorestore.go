package store

import (
	"errors"

	"github.com/bdharris08/scorekeeper/score"
)

// ScoreStore stores scores for ScoreKeeper.
// It could be in memory or backed by a database.
type ScoreStore interface {
	Store(s score.Score) error
	Retrieve(f score.ScoreFactory, scoreType string) (map[string][]score.Score, error)
}

var ErrNoStore = errors.New("scoreStore uninitialized")
