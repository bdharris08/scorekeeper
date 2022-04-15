package store

import (
	"github.com/bdharris08/scorekeeper/score"
)

// ScoreStore stores scores for ScoreKeeper.
// It could be in memory or backed by a database.
type ScoreStore interface {
	Store(s score.Score) error
	Retrieve() (map[string][]score.Score, error)
}
