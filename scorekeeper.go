package scorekeeper

import "errors"

type ScoreKeeper struct{}

func (s *ScoreKeeper) AddAction(json string) error {
	return errors.New("not implemented")
}

func (s *ScoreKeeper) GetStats() string {
	return ""
}
