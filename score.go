package scorekeeper

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
