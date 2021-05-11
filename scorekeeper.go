package scorekeeper

import (
	"encoding/json"
	"errors"
	"strings"
)

type ScoreKeeper struct{}

type Score struct {
	Action string `json:"action"`
	// Since the time units weren't specified, let's assume ms for reasonable precise human-scale time measurements.
	// Max uint64 is 18446744073709551615, or rougly 500 million years, which seems like plenty of time to jump.
	// We only need an int to store this data, but I don't know the edginess the edge cases that will be used in testing.
	Time uint64 `json:"time"`
}

var (
	ErrNoInput   = errors.New("no input provided")
	ErrBadTime   = errors.New("invalid time")
	ErrNoTime    = errors.New("missing time")
	ErrBadAction = errors.New("invalid action")
	ErrBadInput  = errors.New("bad input")
)

func (s *Score) Read(action string) error {
	if action == "" {
		return ErrNoInput
	}

	if !strings.Contains(action, "time") {
		return ErrNoTime
	}

	err := json.Unmarshal([]byte(action), s)
	if err != nil {
		if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
			switch jsonErr.Field {
			case "action":
				return ErrBadAction
			case "time":
				return ErrBadTime
			}
		}

		return ErrBadInput
	}

	if s.Action == "" {
		return ErrBadAction
	}

	return nil
}

func (sk *ScoreKeeper) AddAction(action string) error {
	var s Score
	return s.Read(action)
}

func (s *ScoreKeeper) GetStats() string {
	return ""
}
