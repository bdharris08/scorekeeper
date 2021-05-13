package scorekeeper

import (
	"testing"
)

func TestAverage(t *testing.T) {
	type testCase struct {
		name string
		ss   []Score
		errs []error
		res  float64
		err  error
	}

	testCases := []testCase{
		{
			name: "provided",
			ss: []Score{
				&TestScore{value: float64(100)},
				&TestScore{value: float64(200)},
			},
			errs: []error{nil, nil},
			res:  float64(150),
		},
		{
			name: "empty",
			ss:   []Score{},
			errs: []error{},
			err:  ErrNoData,
		},
		{
			name: "zero",
			ss: []Score{
				&TestScore{name: "jump", value: float64(0)},
				&TestScore{name: "jump", value: float64(200)},
			},
			errs: []error{nil, nil},
			res:  float64(100),
		},
		{
			name: "one",
			ss: []Score{
				&TestScore{name: "jump", value: float64(1)},
			},
			errs: []error{nil},
			res:  float64(1),
		},
		{
			// Presumably not possible for Trial but just in case
			name: "negative",
			ss: []Score{
				&TestScore{name: "jump", value: float64(-200)},
				&TestScore{name: "jump", value: float64(200)},
			},
			errs: []error{nil, nil},
			res:  float64(0),
		},
		{
			name: "duplicate",
			ss: []Score{
				&TestScore{name: "jump", value: float64(100)},
				&TestScore{name: "jump", value: float64(100)},
			},
			errs: []error{nil, nil},
			res:  float64(100),
		},
		{
			name: "floating",
			ss: []Score{
				&TestScore{name: "jump", value: float64(101)},
				&TestScore{name: "jump", value: float64(100)},
			},
			errs: []error{nil, nil},
			res:  float64(100.5),
		},
		{
			name: "repeating of course",
			ss: []Score{
				&TestScore{name: "jump", value: float64(2)},
				&TestScore{name: "jump", value: float64(3)},
				&TestScore{name: "jump", value: float64(5)},
			},
			errs: []error{nil, nil, nil},
			res:  float64(3.3333333333333335),
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
			t.Errorf("[%s] Expected %f but got %f", tc.name, expected, got)
		}

		res2, err := a.Compute(tc.ss)
		if expected, got := tc.err, err; expected != got {
			t.Errorf("[%s] Compute: Expected error to be '%v' but got '%v'", tc.name, expected, got)
		}
		if expected, got := tc.res, res2; expected != got {
			t.Errorf("[%s] Compute: Expected %f but got %f", tc.name, expected, got)
		}

		// sanity check
		if res != res2 {
			t.Errorf("[%s] results didn't match: running: %f, computed: %f", tc.name, res, res2)
		}
	}
}
