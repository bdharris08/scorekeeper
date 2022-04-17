package scorekeeper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bdharris08/scorekeeper/score"
	"github.com/bdharris08/scorekeeper/stat"
	"github.com/bdharris08/scorekeeper/store"
)

// ScoreKeeper keeps scores using some store.
// It is the top level object for the ScoreKeeper library.
// It fulfills the requested AddAction and GetStats methods.
type ScoreKeeper struct {
	// factory for constructiong registered store types
	f score.ScoreFactory
	// scoreStore for storing scores in
	s store.ScoreStore
	// Scores chan will allow clients (through AddAction) to send scores to the worker.
	// Constrain scores channel to only receive, ensuring only the worker reads.
	// errors can be returned by the included channel, like an addressed envelope in an envelope.
	scores chan<- scoreEnvelope
	// Requests chan will be used by clients (through GetStats) to request stats from the worker.
	// Constrain requests channel to only receive, ensuring only the worker reads.
	requests chan<- requestEnvelope
	// close(exit) to stop the worker.
	quit chan bool
}

// New creates and returns a ScoreKeeper with the provided ScoreStore.
// It starts
func New(st store.ScoreStore, sf score.ScoreFactory) (*ScoreKeeper, error) {
	sk := &ScoreKeeper{}

	// default to memoryStore if none was provided
	if st == nil {
		st = &store.MemoryStore{S: map[string]map[string][]score.Score{}}
	}
	sk.s = st

	if len(sf) == 0 {
		return nil, fmt.Errorf("scoreTypes must be provided")
	}
	sk.f = sf

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

// ValidScoreType checks for the presence of scoreType in the score factory
func ValidScoreType(sk *ScoreKeeper, scoreType string) bool {
	_, ok := sk.f[scoreType]
	return ok
}

// result stats from the scorekeeper, or an error
type result struct {
	result string
	err    error
}

// scoreEnvelope allows the caller to send a Score and receive an error from the worker
type scoreEnvelope struct {
	score score.Score
	err   chan error
}

// requestEnvelope encapsulates a request for a type of score and a channel to receive the result
type requestEnvelope struct {
	scoreType string
	r         chan result
}

// work on new scores sent from AddAction.
func (sk *ScoreKeeper) work() (chan<- scoreEnvelope, chan<- requestEnvelope) {
	scores := make(chan scoreEnvelope)
	requests := make(chan requestEnvelope)
	go func() {
		for {
			select {
			case <-sk.quit:
				return

			case s := <-scores:
				err := sk.s.Store(s.score)
				s.err <- err

			case re := <-requests:
				res, err := sk.get(re.scoreType)
				re.r <- result{
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
func (sk *ScoreKeeper) AddAction(scoreType, action string) error {
	if sk.s == nil {
		return ErrNoKeeper
	}
	if sk.scores == nil || sk.quit == nil {
		return ErrNotRunning
	}
	if valid := ValidScoreType(sk, scoreType); !valid {
		return score.ErrBadScoreType
	}

	s, err := score.Create(sk.f, scoreType)
	if err != nil {
		return fmt.Errorf("failed to AddAction of type %s: %w", scoreType, err)
	}
	if err := s.Read(action); err != nil {
		return err
	}

	errCh := make(chan error)
	sk.scores <- scoreEnvelope{
		score: s,
		err:   errCh,
	}
	err = <-errCh

	return err
}

// GetStats computes some statistics about the actions stored in the ScoreKeeper.
// It returns those statistics as a json-encoded string
func (sk *ScoreKeeper) GetStats(scoreType string) (string, error) {
	if sk.s == nil {
		return "", ErrNoKeeper
	}
	if sk.scores == nil || sk.quit == nil {
		return "", ErrNotRunning
	}
	if valid := ValidScoreType(sk, scoreType); !valid {
		return "", score.ErrBadScoreType
	}

	// pass a channel to the worker and wait for it to return the result
	requestCh := make(chan result)
	sk.requests <- requestEnvelope{
		scoreType: scoreType,
		r:         requestCh,
	}
	res := <-requestCh

	return res.result, res.err
}

// get scores from the store and compute averages, returning a json-encoded string
func (sk *ScoreKeeper) get(scoreType string) (string, error) {
	if sk.s == nil {
		return "", store.ErrNoStore
	}

	scoreMap, err := sk.s.Retrieve(sk.f, scoreType)
	if err != nil {
		return "", err
	}

	if len(scoreMap) == 0 {
		return "", stat.ErrNoData
	}

	avgs := make([]score.AverageTime, 0, len(scoreMap))

	for name, scores := range scoreMap {
		a := stat.Average{}

		res, err := a.Compute(scores)
		if err != nil {
			return "", err
		}

		avg, ok := res.(float64)
		if !ok {
			return "", stat.ErrTypeInvalid
		}

		avgs = append(avgs, score.AverageTime{
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
