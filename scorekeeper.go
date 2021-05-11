package scorekeeper

import (
	"encoding/json"
	"errors"
	"strings"
)

// ScoreKeeper keeps scores using some store.
// It is the top level object for the ScoreKeeper library.
type ScoreKeeper struct {
	store ScoreStore
}

// ScoreStore stores scores for ScoreKeeper.
// It could be in memory or backed by a database.
type ScoreStore interface {
	Store() error
}

// Score is kept by ScoreKeeper and tracks something.
// We will only have one kind of Score for this project,
// but through this interface we could extend to other kinds easily
type Score interface {
	Read(s string) error
}

// Trial is a kind of Score. It is a timed action.
type Trial struct {
	Action string `json:"action"`
	// Since the "time" units weren't specified, let's assume ms as a reasonably precise
	// human-scale time measurement. Max uint64 is 18446744073709551615,
	// or rougly 500 million years, which seems like plenty of time to jump.
	// We only need an int to store this data, but I don't know the edginess of the edge cases
	// that will be used in testing.
	Time uint64 `json:"time"`
}

var (
	ErrNoInput   = errors.New("no input provided")
	ErrBadTime   = errors.New("invalid time")
	ErrNoTime    = errors.New("missing time")
	ErrBadAction = errors.New("invalid action")
	ErrBadInput  = errors.New("bad input")
)

// Read a json-encoded string into a Trial struct
func (s *Trial) Read(action string) error {
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

// AddAction takes a json-encoded string and keeps it for later
func (sk *ScoreKeeper) AddAction(action string) error {
	var s Trial
	return s.Read(action)
}

// GetStats computes some statistics about the actions stored in the ScoreKeeper
func (s *ScoreKeeper) GetStats() string {
	return ""
}
