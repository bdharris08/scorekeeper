package scorekeeper

import (
	"errors"
)

// Stat does some math on a set of scores and returns the result.
// We will only have one kind of Statistic for this project (Average),
// but through this interface we could extend to other kinds easily.
type Stat interface {
	// Compute generates the statistic on a set of scores immediately.
	Compute(ss []Score) (interface{}, error)
	// Step takes a score and, if possible, includes it in the running computation
	// For example, Average.Step() will add a score to the running total
	// We won't need this for Average really, but it could be useful in other types of computations
	Step(s Score) error
	// Report returns the result of the running computation.
	Report() (interface{}, error)
}

// Average is a Stat that computes an average of scores with int64 values.
type Average struct {
	n float64
	s float64
}

var ErrTypeInvalid = errors.New("score value is invalid")
var ErrNoData = errors.New("no data to report")

// Compute a floating point average from a list of scores with int64 values.
func (a *Average) Compute(ss []Score) (interface{}, error) {
	var (
		sum float64
		n   float64
	)

	for _, s := range ss {
		v, ok := s.Value().(int64)
		if !ok {
			return float64(0), ErrTypeInvalid
		}
		sum += float64(v)
		n++
	}

	if n < float64(1) {
		return float64(0), ErrNoData
	}

	avg := sum / n

	return avg, nil
}

// Step adds a score to the running sum.
func (a *Average) Step(s Score) error {
	v, ok := s.Value().(int64)
	if !ok {
		return ErrTypeInvalid
	}

	a.s += float64(v)
	a.n++
	return nil
}

// Report the average from the running sum as a string.
func (a *Average) Report() (interface{}, error) {
	if a.n < float64(1) {
		return float64(0), ErrNoData
	}

	avg := a.s / a.n
	return avg, nil
}
