package scorekeeper

import "errors"

type ScoreKeeper struct{}

var (
	ErrNoInput     = errors.New("no input provided")
	ErrOutOfBounds = errors.New("time out of bounds")
	ErrBadInput    = errors.New("bad input")
)

func (s *ScoreKeeper) AddAction(json string) error {
	return errors.New("not implemented")
}

func (s *ScoreKeeper) GetStats() string {
	return ""
}
