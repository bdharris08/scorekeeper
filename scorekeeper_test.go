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
		if err != tc.err {
			t.Errorf("[%s] Expected error to be %v but got %v", tc.name, tc.err, err)
		}
	}
}
