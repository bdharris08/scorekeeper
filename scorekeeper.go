package scorekeeper

import (
	"encoding/json"
	"errors"
)

// ScoreKeeper keeps scores using some store.
// It is the top level object for the ScoreKeeper library.
// It fulfills the requested AddAction and GetStats methods.
type ScoreKeeper struct {
	s ScoreStore
	// constrain scores channel to only receive, ensuring only the worker reads
	scores chan<- Score
	// close(exit) to stop the worker.
	quit chan bool
}

// New creates and returns a ScoreKeeper with the provided ScoreStore.
// It starts
func New(store ScoreStore) (*ScoreKeeper, error) {
	sk := &ScoreKeeper{}

	// default to memoryStore if none was provided
	if store == nil {
		store = &MemoryStore{s: map[string][]Score{}}
	}
	sk.s = store

	return sk, nil
}

// Start a worker routine to listen on the scores channel.
func (sk *ScoreKeeper) Start() {
	sk.quit = make(chan bool)
	sk.scores = sk.work()
}

var ErrNotRunning = errors.New("scorekeeper not running. Use Start()")

// Stop the worker goroutine.
func (sk *ScoreKeeper) Stop() error {
	if sk.quit != nil {
		close(sk.quit)
		return nil
	}
	return ErrNotRunning
}

// work on new scores sent from AddAction.
func (sk *ScoreKeeper) work() chan<- Score {
	scores := make(chan Score)
	go func() {
		for {
			select {
			case <-sk.quit:
				return
			case s := <-scores:
				if err := sk.s.Store(s); err != nil {
					panic(err) // TODO
				}
			default:
				// loop until quit
			}
		}
	}()
	return scores
}

var ErrNoKeeper = errors.New("scorekeeper uninitialized. Use New()")

// AddAction takes a json-encoded string action and keeps it for later.
func (sk *ScoreKeeper) AddAction(action string) error {
	if sk.s == nil {
		return ErrNoKeeper
	}
	if sk.scores == nil || sk.quit == nil {
		return ErrNotRunning
	}

	var s Trial
	if err := s.Read(action); err != nil {
		return err
	}

	sk.scores <- &s
	return nil
}

// GetStats computes some statistics about the actions stored in the ScoreKeeper.
// It returns those statistics as a json-encoded string
func (sk *ScoreKeeper) GetStats() (string, error) {
	if sk.s == nil {
		return "", ErrNoKeeper
	}
	if sk.scores == nil || sk.quit == nil {
		return "", ErrNotRunning
	}

	avgs, err := sk.get()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(avgs)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// get scores from the store and compute averages
func (sk *ScoreKeeper) get() ([]AverageTime, error) {
	if sk.s == nil {
		return nil, ErrNoData
	}

	scoreMap, err := sk.s.Retrieve()
	if err != nil {
		return nil, err
	}

	if len(scoreMap) == 0 {
		return nil, ErrNoData
	}

	avgs := make([]AverageTime, 0, len(scoreMap))

	for name, scores := range scoreMap {
		a := Average{}

		res, err := a.Compute(scores)
		if err != nil {
			return nil, err
		}

		avg, ok := res.(float64)
		if !ok {
			return nil, ErrTypeInvalid
		}

		avgs = append(avgs, AverageTime{
			Action:  name,
			Average: avg,
		})
	}

	return avgs, nil
}
