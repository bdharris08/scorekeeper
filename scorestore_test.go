package scorekeeper

import "testing"

func TestMemoryStoreSimple(t *testing.T) {
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

func TestMemoryStoreRetrieve(t *testing.T) {
	// TODO
}
