package store

import (
	"testing"

	"github.com/bdharris08/scorekeeper/score"
)

func TestMemoryStoreSimple(t *testing.T) {
	ms := MemoryStore{
		S: map[string][]score.Score{},
	}

	s := score.TestScore{
		TName:  "test",
		TValue: 1,
	}

	if err := ms.Store(&s); err != nil {
		t.Error(err)
	}

	scores, err := ms.Retrieve()
	if err != nil {
		t.Error(err)
	}
	if expected, got := s.TName, scores[s.Name][0].Name(); expected != got {
		t.Errorf("Expected %s but got %s", expected, got)
	}
}
