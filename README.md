# ScoreKeeper
A library for keeping score, a practice project

[![test](https://github.com/bdharris08/scorekeeper/actions/workflows/test.yml/badge.svg)](https://github.com/bdharris08/scorekeeper/actions/workflows/test.yml)

### Features
`scoreKeeper.AddAction(scoreType, action string) error`
AddAction accepts a json-encoded action string like `"{"action":"hop", "time":100}"`.
It will store the action as a Score in a ScoreStore (categorized by scoreType) for later stats calculations.
It can return the errors:
- `ErrNoKeeper` = "scorekeeper uninitialized. Use New()"
- `ErrNotRunning` = "scorekeeper not running. Use Start()"
- Errors defined by the score type used, for example `score.ErrNoInput`

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

Memory Store: see [example/memory/README.md](./example/memory/README.md)

Postgres SQL Store: see [example/pgsql/README.md](./example/pgsql/README.md)
