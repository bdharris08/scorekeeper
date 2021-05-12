package scorekeeper

import (
	"testing"
)

func TestAverage(t *testing.T) {
	type testCase struct {
		name string
		ss   []Score
		errs []error
		res  string
		err  error
	}

	testCases := []testCase{
		{
			name: "provided",
			ss: []Score{
				&TestScore{value: int64(100)},
				&TestScore{value: int64(200)},
			},
			errs: []error{nil, nil},
			res:  "150",
		},
		{
			name: "empty",
			ss:   []Score{},
			errs: []error{},
			res:  "",
			err:  ErrNoData,
		},
		{
			name: "zero",
			ss: []Score{
				&TestScore{name: "jump", value: int64(0)},
				&TestScore{name: "jump", value: int64(200)},
			},
			errs: []error{nil, nil},
			res:  "100",
		},
		{
			name: "one",
			ss: []Score{
				&TestScore{name: "jump", value: int64(1)},
			},
			errs: []error{nil},
			res:  "1",
		},
		{
			// Presumably not possible for Trial but just in case
			name: "negative",
			ss: []Score{
				&TestScore{name: "jump", value: int64(-200)},
				&TestScore{name: "jump", value: int64(200)},
			},
			errs: []error{nil, nil},
			res:  "0",
		},
		{
			name: "duplicate",
			ss: []Score{
				&TestScore{name: "jump", value: int64(100)},
				&TestScore{name: "jump", value: int64(100)},
			},
			errs: []error{nil, nil},
			res:  "100",
		},
		{
			name: "floating",
			ss: []Score{
				&TestScore{name: "jump", value: int64(101)},
				&TestScore{name: "jump", value: int64(100)},
			},
			errs: []error{nil, nil},
			res:  "100.5",
		},
		{
			name: "repeating of course",
			ss: []Score{
				&TestScore{name: "jump", value: int64(2)},
				&TestScore{name: "jump", value: int64(3)},
				&TestScore{name: "jump", value: int64(5)},
			},
			errs: []error{nil, nil, nil},
			res:  "3.3333333333333335",
		},
	}

	for _, tc := range testCases {
		a := Average{}

		for i, s := range tc.ss {
			if expected, got := tc.errs[i], a.Step(s); expected != got {
				t.Errorf("[%s] Expected error to be '%v' but got '%v'", tc.name, expected, got)
			}
		}

		res, err := a.Report()
		if expected, got := tc.err, err; expected != got {
			t.Errorf("[%s] Expected error to be '%v' but got '%v'", tc.name, expected, got)
		}
		if expected, got := tc.res, res; expected != got {
			t.Errorf("[%s] Expected %s but got %s", tc.name, expected, got)
		}

		res2, err := a.Compute(tc.ss)
		if expected, got := tc.err, err; expected != got {
			t.Errorf("[%s] Compute: Expected error to be '%v' but got '%v'", tc.name, expected, got)
		}
		if expected, got := tc.res, res2; expected != got {
			t.Errorf("[%s] Compute: Expected %s but got %s", tc.name, expected, got)
		}

		// sanity check
		if res != res2 {
			t.Errorf("[%s] results didn't match: running: %s, computed: %s", tc.name, res, res2)
		}
	}
}
