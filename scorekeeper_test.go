package scorekeeper

import (
	"strings"
	"testing"
)

func TestAddActionErrors(t *testing.T) {
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
		},
		{
			name:   "huge",
			action: `{"action":"jump", "time":9223372036854775807}`,
		},
		{
			name:   "NaN",
			action: `{"action":"jump", "time":"1s"}`,
			err:    ErrBadTime,
		},
		{
			name:   "empty action",
			action: `{"action":"", "time":1}`,
			err:    ErrBadAction,
		},
		{
			name:   "missing time",
			action: `{"action":"exist"}`,
			err:    ErrNoTime,
		},
		{
			name:   "missing action",
			action: `{"time":1}`,
			err:    ErrBadAction,
		},
		{
			name:   "missing both",
			action: `{}`,
			err:    ErrNoTime,
		},
	}

	for _, tc := range testCases {
		s, err := New(&MemoryStore{map[string][]Score{}})
		if err != nil {
			t.Fatal(err)
		}

		s.Start()
		defer s.Stop()

		err = s.AddAction(tc.action)
		if expected, got := tc.err, err; expected != got {
			t.Errorf("[%s] Expected error to be '%v' but got '%v'", tc.name, expected, got)
		}
	}
}

func TestAddActionConcurrent(t *testing.T) {
	s, err := New(&MemoryStore{map[string][]Score{}})
	if err != nil {
		t.Fatal(err)
	}

	s.Start()
	defer s.Stop()

	actions := []string{
		`{"action":"hop", "time":100}`,
		`{"action":"skip", "time":100}`,
		`{"action":"jump", "time":100}`,
		`{"action":"trip", "time":100}`,
		`{"action":"fall", "time":100}`,
		`{"action":"lay", "time":1000}`,
		`{"action":"curse", "time":100}`,
		`{"action":"stand", "time":100}`,
		`{"action":"walk", "time":100}`,
		`{"action":"run", "time":100}`,
		`{"action":"hop", "time":100}`,
	}

	chaos := func(actions []string) {
		for _, a := range actions {
			if err := s.AddAction(a); err != nil {
				t.Error(err)
			}
		}
	}

	for i := 0; i < 10; i++ {
		go chaos(actions)
	}
}

func TestGetStats(t *testing.T) {
	type testCase struct {
		name    string
		actions []string
		stats   string
		errs    []error
		err     error
	}
	testCases := []testCase{
		{
			name: "provided",
			actions: []string{
				`{"action":"jump", "time":100}`,
				`{"action":"run", "time":75}`,
				`{"action":"jump", "time":200}`,
			},
			stats: `[{"action":"jump","avg":150},{"action":"run","avg":75}]`,
			errs:  []error{nil, nil, nil},
		},
		{
			name:    "empty",
			actions: []string{},
			errs:    []error{},
			err:     ErrNoData,
		},
		{
			name: "zero",
			actions: []string{
				`{"action":"stand", "time":0}`,
			},
			stats: `[{"action":"stand","avg":0}]`,
			errs:  []error{nil},
		},
		{
			name: "unique",
			actions: []string{
				`{"action":"hop", "time":1}`,
				`{"action":"skip", "time":2}`,
				`{"action":"jump", "time":3}`,
			},
			stats: `[{"action":"hop","avg":1},{"action":"skip","avg":2},{"action":"jump","avg":3}]`,
			errs:  []error{nil, nil, nil},
		},
		{
			name: "negative",
			actions: []string{
				`{"action":"sink","time":-100}`,
			},
			stats: `[{"action":"sink","avg":-100}]`,
			errs:  []error{nil},
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
			stats: `[{"action":"sink","avg":-100},{"action":"jump","avg":150},{"action":"run","avg":75}]`,
			errs:  []error{nil, ErrNoTime, nil, nil, nil},
		},
	}

	for _, tc := range testCases {
		s, err := New(&MemoryStore{map[string][]Score{}})
		if err != nil {
			t.Fatal(err)
		}

		s.Start()
		defer s.Stop()

		for i, a := range tc.actions {
			errs := make([]error, len(tc.actions))
			errs[i] = s.AddAction(a)
			if expected, got := tc.errs[i], errs[i]; expected != got {
				t.Errorf("[%s] Expected error for action %d to be '%v' but got '%v'", tc.name, i, expected, got)
			}
		}

		stats, err := s.GetStats()
		if expected, got := tc.err, err; expected != got {
			t.Errorf("[%s] Expected GetStats err to be '%v' but got '%v'", tc.name, expected, got)
		}

		if expected, got := tc.stats, stats; expected != got {
			t.Errorf("[%s] Expected stats to be '%s' but got '%s'", tc.name, expected, got)
		}
	}
}

func statsEquivalent(a, b string) bool {
	// trim the outside brackets [ string ] => string
	a, b = strings.Trim(a, "[]"), strings.Trim(b, "[]")

	// split the string into a list on commas
	al, bl := strings.Split(a, ","), strings.Split(b, ",")

	if len(al) != len(bl) {
		return false
	}

	// use a map to confirm both lists contain the same members
	am := make(map[string]bool, len(al))
	for _, s := range al {
		am[s] = true
	}

	for _, s := range bl {
		if present := am[s]; !present {
			return false
		}
	}

	return true
}

func TestStatsEquivalent(t *testing.T) {
	type testCase struct {
		name string
		a, b string
		e    bool
	}
	testCases := []testCase{
		{
			name: "one",
			a:    `[{"action":"hop","avg":1}]`,
			b:    `[{"action":"hop","avg":1}]`,
			e:    true,
		},
		{
			name: "def not",
			a:    `[]`,
			b:    `[{"action":"hop","avg":1}]`,
			e:    false,
		},
		{
			name: "changed values",
			a:    `[{"action":"skip","avg":1}]`,
			b:    `[{"action":"hop","avg":1}]`,
			e:    false,
		},
		{
			name: "different lengths",
			a:    `[{"action":"skip","avg":1}, {"action":"hop","avg":1}]`,
			b:    `[{"action":"hop","avg":1}]`,
			e:    false,
		},
		{
			name: "three",
			a:    `[{"action":"hop","avg":1},{"action":"skip","avg":2},{"action":"jump","avg":3}]`,
			b:    `[{"action":"skip","avg":2},{"action":"jump","avg":3},{"action":"hop","avg":1}]`,
			e:    true,
		},
	}

	for _, tc := range testCases {
		if expected, got := tc.e, statsEquivalent(tc.a, tc.b); expected != got {
			t.Errorf("[%s] expected a===b to be %t but got %t", tc.name, expected, got)
		}
	}
}
