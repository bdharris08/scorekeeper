package scorekeeper

import (
	"encoding/json"
	"errors"
	"strings"
)

// Score is kept by ScoreKeeper and tracks something.
// We will only have one kind of Score for this project (Trial),
// but through this interface we could extend to other kinds easily
type Score interface {
	// Generate a unique identifier for later organization.
	Name() string
	// Read a json-encoded string into the Score struct.
	Read(s string) error
	// Value returns the value of the Score.
	// By returning an empty interface, we can
	Value() interface{}
}

// TestScore is a simple score for testing.
type TestScore struct {
	name  string
	value int64
}

// Name the test score.
func (t *TestScore) Name() string {
	return t.name
}

// Read nothing, nowhere.
func (t *TestScore) Read(action string) error {
	return nil
}

// Value returns the value of the test score.
func (t *TestScore) Value() interface{} {
	return t.value
}

// Trial is a kind of Score. It is a timed action.
type Trial struct {
	Action string `json:"action"`
	// Since the "time" units weren't specified, let's assume ms as a reasonably precise
	// human-scale time measurement. Max int64 is 9223372036854775807,
	// or rougly 300 million years, which seems like plenty of time to jump.
	// We only need an int to store this data, but I don't know the edginess of the edge cases
	// that will be used in testing.
	// Even though time can't be negative, we'll use int64 to avoid a loss of precision
	Time int64 `json:"time"`
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

// Name returns the trial's action
func (t *Trial) Name() string {
	return t.Action
}

// Value returns the trial's time
func (t *Trial) Value() interface{} {
	return t.Time
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
