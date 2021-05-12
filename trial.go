package scorekeeper

import (
	"encoding/json"
	"errors"
	"strings"
)

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

func (t *Trial) Name() string {
	return t.Action
}

// Read a json-encoded string into the Trial struct.
func (t *Trial) Read(action string) error {
	if action == "" {
		return ErrNoInput
	}

	if !strings.Contains(action, "time") {
		return ErrNoTime
	}

	err := json.Unmarshal([]byte(action), t)
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

	if t.Action == "" {
		return ErrBadAction
	}

	return nil
}
