package score

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Trial is a kind of Score. It is a timed action.
type Trial struct {
	Action string `json:"action"`
	// Since the "time" units weren't specified, let's assume milliseconds as a reasonably precise
	// human-scale time measurement. Max float64 is 1.7976931348623157e+308,
	// or rougly 5.700447535712569×10^297 years, which seems like plenty of time to jump.
	// We only need an int to store this data, but I don't know the edginess of the edge cases
	// that will be used in testing.
	Time float64 `json:"time"`
}

// AverageTime will be used to report an average time
type AverageTime struct {
	Action  string  `json:"action"`
	Average float64 `json:"avg"`
}

var (
	ErrNoInput   = errors.New("no input provided")
	ErrBadTime   = errors.New("invalid time")
	ErrNoTime    = errors.New("missing time")
	ErrBadAction = errors.New("invalid action")
	ErrBadInput  = errors.New("bad input")
)

func NewTrial() Score {
	return &Trial{}
}

// Type returns the type of Score
func (t *Trial) Type() string {
	return "trial"
}

// Name returns the trial's action
func (t *Trial) Name() string {
	return t.Action
}

// Value returns the trial's time
func (t *Trial) Value() interface{} {
	return t.Time
}

// Set the value and name of the Trial
func (t *Trial) Set(name string, value interface{}) error {
	f, ok := value.(float64)
	if !ok {
		return fmt.Errorf("failed to assert value type")
	}

	t.Action = name
	t.Time = f
	return nil
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
