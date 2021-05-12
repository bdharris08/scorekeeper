package scorekeeper

import "testing"

// TestScore is a simple score for testing the Store.
type TestScore struct {
	name  string
	value int
}

// Name the test score.
func (t *TestScore) Name() string {
	return t.name
}

// Read nothing, nowhere.
func (t *TestScore) Read(action string) error {
	return nil
}

func TestSimple(t *testing.T) {
	ms := MemoryStore{
		s: map[string][]Score{},
	}

	s := TestScore{
		name:  "test",
		value: 1,
	}

	if err := ms.Store(&s); err != nil {
		t.Error(err)
	}

	scores, err := ms.Retrieve(s.name)
	if err != nil {
		t.Error(err)
	}
	if expected, got := s.name, scores[0].Name(); expected != got {
		t.Errorf("Expected %s but got %s", expected, got)
	}
}

func TestRetrieve(t *testing.T) {
	// TODO
}
