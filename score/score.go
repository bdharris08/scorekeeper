package score

import (
	"fmt"
)

// Score is kept by ScoreKeeper and tracks something.
// We will only have one kind of Score for this project (Trial),
// but through this interface we could extend to other kinds easily
type Score interface {
	// Type of score, for example "trial"
	Type() string
	// Generate a unique identifier for later organization.
	Name() string
	// Read a json-encoded string into the Score struct.
	Read(s string) error
	// Value returns the value of the Score.
	Value() interface{}
	// Set the name and value of the score
	Set(name string, value interface{}) error
}

type ScoreConstructor func() Score

type ScoreFactory map[string]ScoreConstructor

func Create(f ScoreFactory, scoreType string) (Score, error) {
	constructor, ok := f[scoreType]
	if !ok {
		return nil, fmt.Errorf("unregistered scoreType %s", scoreType)
	}

	return constructor(), nil
}

// TestScore is a simple score for testing.
type TestScore struct {
	TName  string
	TValue float64
}

func NewTestScore() Score {
	return &TestScore{}
}

// Type returns the type of Score
func (t *TestScore) Type() string {
	return "test"
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

// Set the value and name of the TestScore
func (t *TestScore) Set(name string, value interface{}) error {
	f, ok := value.(float64)
	if !ok {
		return fmt.Errorf("failed to assert value type")
	}

	t.TName = name
	t.TValue = f
	return nil
}
