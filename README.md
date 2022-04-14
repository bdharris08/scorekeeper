# ScoreKeeper
A library for keeping score, a practice project

[![test](https://github.com/bdharris08/scorekeeper/actions/workflows/test.yml/badge.svg)](https://github.com/bdharris08/scorekeeper/actions/workflows/test.yml)

### Features
`scoreKeeper.AddAction(action string) error`
AddAction accepts a json-encoded action string like `"{"action":"hop", "time":100}"`.
It will store the action as a Score in a ScoreStore for later stats calculations.
It can return the errors:
- `ErrNoKeeper` = "scorekeeper uninitialized. Use New()"
- `ErrNotRunning` = "scorekeeper not running. Use Start()"

`scoreKeeper.GetStats() (string, error)`
GetStats will return a json-encoded list of average scores like `"[{"action":"hop", "avg":100}]"`.
It can return the errors:
- `ErrNoKeeper` = "scorekeeper uninitialized. Use New()"
- `ErrNotRunning` = "scorekeeper not running. Use Start()"
- `ErrNoData` = "no data to report"
- `ErrTypeInvalid` = "invalid type" of score sent to the Average calculator (it expects float64). This can't happen with the provided MemoryStore, but could happen with custom ScoreStore implementations.

#### The ScoreStore
ScoreKeeper comes with an in-memory implementation of the `ScoreStore interface`.
You could implement your own using that interface. Just pass an initialized store to `New`.

### Example Usage

#### Setup
```go

import (
    "github.com/bdharris08/scorekeeper"
    "github.com/bdharris08/scorekeeper/store"
)

func main() {
    scoreKeeper, err := scorekeeper.New(&store.MemoryStore{})
	if err != nil {
		panic(fmt.Errorf("error creating scoreKeeper: %w", err))
	}

	scoreKeeper.Start()
	defer scoreKeeper.Stop()
}
```

#### A simple example
```go
actions := []string{
    `{"action":"hop", "time":100}`,
    `{"action":"skip", "time":100}`,
    `{"action":"hop", "time":100}`,
}

// A simple example
for _, a := range actions {
    if err := scoreKeeper.AddAction(a); err != nil {
        fmt.Println(err)
    }
}
result, err := scoreKeeper.GetStats()
if err != nil {
    fmt.Println(err)
}
// Do something with the result
fmt.Println(result)
```

#### Concurrent example
```go
// You can even access the scoreKeeper concurrently!
var wg sync.WaitGroup
for i := 0; i < numWorkers; i++ {
    wg.Add(1)

    go func() {
        defer wg.Done()
        for _, a := range actions {
            if err := scoreKeeper.AddAction(a); err != nil {
                fmt.Println(err)
            }
        }
        result, err := scoreKeeper.GetStats()
        if err != nil {
            fmt.Println(err)
        }
        // Do something with the result, 
        // this will likely be an intermediate result 
        // since not all workers will have finished
        fmt.Println(result)
    }()
}

wg.Wait()

result, err = scoreKeeper.GetStats()
if err != nil {
    fmt.Println(err)
}
// Do something with the final result
fmt.Println(result)
```
