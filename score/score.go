package score

import (
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
	Value() interface{}
}

// TestScore is a simple score for testing.
type TestScore struct {
	TName  string
	TValue float64
}

// Name the test score.
func (t *TestScore) Name() string {
	return t.TName
}

// Read nothing, nowhere.
func (t *TestScore) Read(action string) error {
	return nil
}

// Value returns the value of the test score.
func (t *TestScore) Value() interface{} {
	return t.TValue
}

	}

	return nil
}
