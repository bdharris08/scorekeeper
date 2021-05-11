package scorekeeper

import "errors"

type ScoreKeeper struct{}

var (
	ErrNoInput   = errors.New("no input provided")
	ErrBadTime   = errors.New("invalid time")
	ErrNoTime    = errors.New("missing time")
	ErrBadAction = errors.New("invalid action")
	ErrBadInput  = errors.New("bad input")
)

func (s *ScoreKeeper) AddAction(json string) error {
	return errors.New("not implemented")
}

func (s *ScoreKeeper) GetStats() string {
	return ""
}
