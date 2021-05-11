package scorekeeper

import (
	"testing"
)

func TestAddAction(t *testing.T) {
	type testCase struct {
		name   string
		action string
		err    error
	}
	testCases := []testCase{
		{
			name:   "simple",
			action: `{"action":"jump", "time":100}`,
		},
		{
			name:   "zero",
			action: `{"action":"jump", "time":0}`,
		},
		{
			name:   "empty",
			action: ``,
			err:    ErrNoInput,
		},
		{
			name:   "negative",
			action: `{"action":"levitate", "time":-1}`,
			err:    ErrOutOfBounds,
		},
		{
			name:   "huge",
			action: `{"action":"jump", "time":18446744073709551615}`,
		},
		{
			name:   "too huge",
			action: `{"action":"jump", "time":18446744073709551616}`,
			err:    ErrOutOfBounds,
		},
		{
			name:   "NaN",
			action: `{"action":"jump", "time":"1s"}`,
			err:    ErrBadInput,
		},
		{
			name:   "empty action",
			action: `{"action":"", "time":1}`,
			err:    ErrBadInput,
		},
		{
			name:   "missing time",
			action: `{"action":"exist"}`,
			err:    ErrBadInput,
		},
		{
			name:   "missing action",
			action: `{"time":1}`,
			err:    ErrBadInput,
		},
		{
			name:   "missing both",
			action: `{}`,
			err:    ErrBadInput,
		},
	}

	for _, tc := range testCases {
		s := ScoreKeeper{}

		err := s.AddAction(tc.action)
		if expected, got := tc.err, err; expected != got {
			t.Errorf("[%s] Expected error to be '%v' but got '%v'", tc.name, expected, got)
		}
	}
}

func TestGetStats(t *testing.T) {
	type testCase struct {
		name    string
		actions []string
		stats   string
		errs    []error
	}
	testCases := []testCase{
		{
			name: "provided",
			actions: []string{
				`{"action":"jump", "time":100}`,
				`{"action":"run", "time":75}`,
				`{"action":"jump", "time":200}`,
			},
			stats: `[
				{"action":"jump", "avg":150},
				{"action":"run", "avg":75}
			]`,
			errs: []error{nil, nil, nil},
		},
		{
			name:    "empty",
			actions: []string{},
			stats:   `[]`,
			errs:    []error{},
		},
		{
			name: "zero",
			actions: []string{
				`{"action":"stand", "time":0}`,
			},
			stats: `[
				{"action":"stand", "avg":0}
			]`,
			errs: []error{nil},
		},
		{
			name: "unique",
			actions: []string{
				`{"action":"hop", "time":1}`,
				`{"action":"skip", "time":2}`,
				`{"action":"jump", "time":3}`,
			},
			stats: `[
				{"action":"hop", "time":1},
            	{"action":"skip", "time":2},
            	{"action":"jump", "time":3}
			]`,
			errs: []error{nil, nil, nil},
		},
		{
			name: "negative",
			actions: []string{
				`{"action":"sink", "time":-100}`,
			},
			stats: `[]`,
			errs:  []error{ErrOutOfBounds},
		},
		{
			name: "robust",
			actions: []string{
				`{"action":"sink", "time":-100}`,
				`{"action":"exist"}`,
				`{"action":"jump", "time":100}`,
				`{"action":"run", "time":75}`,
				`{"action":"jump", "time":200}`,
			},
			stats: `[
				{"action":"jump", "avg":150},
				{"action":"run", "avg":75}
			]`,
			errs: []error{ErrOutOfBounds, ErrBadInput, nil, nil, nil},
		},
	}

	for _, tc := range testCases {
		s := ScoreKeeper{}

		for i, a := range tc.actions {
			errs := make([]error, len(tc.actions))
			errs[i] = s.AddAction(a)
			if expected, got := tc.errs[i], errs[i]; expected != got {
				t.Errorf("[%s] Expected error for action %d to be '%v' but got '%v'", tc.name, i, expected, got)
			}
		}

		if expected, got := tc.stats, s.GetStats(); expected != got {
			t.Errorf("[%s] Expected stats to be '%s' but got '%s'", tc.name, expected, got)
		}
	}
}
