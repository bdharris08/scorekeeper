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
	// Scores chan will allow clients (through AddAction) to send scores to the worker.
	// Constrain scores channel to only receive, ensuring only the worker reads.
	// errors can be returned by the included channel, like an addressed envelope in an envelope.
	scores chan<- scoreEnvelope
	// Requests chan will be used by clients (through GetStats) to request stats from the worker.
	// Constrain requests channel to only receive, ensuring only the worker reads.
	requests chan<- chan result
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
	sk.scores, sk.requests = sk.work()
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

// result stats from the scorekeeper, or an error
type result struct {
	result string
	err    error
}

// scoreEnvelope allows the caller to send a Score and receive an error from the worker
type scoreEnvelope struct {
	score Score
	err   chan error
}

// work on new scores sent from AddAction.
func (sk *ScoreKeeper) work() (chan<- scoreEnvelope, chan<- chan result) {
	scores := make(chan scoreEnvelope)
	requests := make(chan chan result)
	go func() {
		for {
			select {
			case <-sk.quit:
				return

			case s := <-scores:
				err := sk.s.Store(s.score)
				s.err <- err

			case r := <-requests:
				res, err := sk.get()
				r <- result{
					result: res,
					err:    err,
				}

			default:
				// loop until quit
			}
		}
	}()
	return scores, requests
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

	errCh := make(chan error)
	sk.scores <- scoreEnvelope{
		score: &s,
		err:   errCh,
	}
	err := <-errCh

	return err
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

	// pass a channel to the worker and wait for it to return the result
	request := make(chan result)
	sk.requests <- request
	res := <-request

	return res.result, res.err
}

// get scores from the store and compute averages, returning a json-encoded string
func (sk *ScoreKeeper) get() (string, error) {
	if sk.s == nil {
		return "", ErrNoData
	}

	scoreMap, err := sk.s.Retrieve()
	if err != nil {
		return "", err
	}

	if len(scoreMap) == 0 {
		return "", ErrNoData
	}

	avgs := make([]AverageTime, 0, len(scoreMap))

	for name, scores := range scoreMap {
		a := Average{}

		res, err := a.Compute(scores)
		if err != nil {
			return "", err
		}

		avg, ok := res.(float64)
		if !ok {
			return "", ErrTypeInvalid
		}

		avgs = append(avgs, AverageTime{
			Action:  name,
			Average: avg,
		})
	}

	b, err := json.Marshal(avgs)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
