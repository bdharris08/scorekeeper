package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bdharris08/scorekeeper/score"
)

/* Example Schema (postgres):
// why bigint? https://www.cybertec-postgresql.com/en/uuid-serial-or-identity-columns-for-postgresql-auto-generated-primary-keys/
TODO CREATE TABLE cohort
CREATE TABLE trials (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	action text NOT NULL,
	value numeric NOT NULL
);
*/

// SQLStore keeps scores in a postgres database.
// It uses the `database/sql` interfaces to remain driver agnostic.
// It must be a postgres driver due to syntax differences.
type SQLStore struct {
	DB         *sql.DB
	TrialTable string
}

// NewSQLStore returns a new *SQLStore
// Specify scoreTable and trialTable or use "" for defaults
func NewSQLStore(db *sql.DB, scoreTable, trialTable string) *SQLStore {
	s := &SQLStore{
		DB:         db,
		TrialTable: "trials",
	}

	if trialTable != "" {
		s.TrialTable = trialTable
	}

	return s
}

var ErrDBUninitialized = errors.New("db not initialized")
var ErrTrialTable = errors.New("trial table name not specified")

// Store score `s` to the database
func (st *SQLStore) Store(s score.Score) error {
	if err := st.CheckInit(); err != nil {
		return err
	}

	tx, err := st.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := fmt.Sprintf("INSERT INTO %s(action, value) values($1,$2)", st.TrialTable)
	_, err = st.DB.Exec(query, s.Name(), s.Value())
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert trial: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Retrieve Scores from the database by name
func (st *SQLStore) Retrieve() (map[string][]score.Score, error) {
	if err := st.CheckInit(); err != nil {
		return nil, err
	}

	ret := map[string][]score.Score{}

	query := fmt.Sprintf("SELECT action, value FROM %s", st.TrialTable)
	rows, err := st.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query scores: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			action string
			value  float64
		)
		if err := rows.Scan(&action, &value); err != nil {
			return nil, fmt.Errorf("error scanning: %w", err)
		}
		ret[action] = append(ret[action], &score.Trial{Action: action, Time: value})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return ret, nil
}

func (st *SQLStore) CheckInit() error {
	switch {
	case st.DB == nil:
		return ErrDBUninitialized
	case st.TrialTable == "":
		return ErrTrialTable
	default:
		return nil
	}
}
