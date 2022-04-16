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
CREATE TABLE <scoreType> (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	name text NOT NULL,
	value numeric NOT NULL
);
*/

// SQLStore keeps scores in a postgres database.
// It uses the `database/sql` interfaces to remain driver agnostic.
// It must be a postgres driver due to syntax differences.
type SQLStore struct {
	DB *sql.DB
}

// NewSQLStore returns a new *SQLStore
// Specify scoreTable and trialTable or use "" for defaults
func NewSQLStore(db *sql.DB) (*SQLStore, error) {
	if db == nil {
		return nil, ErrDBUninitialized
	}

	return &SQLStore{DB: db}, nil
}

var ErrDBUninitialized = errors.New("db not initialized")
var ErrTrialTable = errors.New("trial table name not specified")

// Store score `s` to the database
func (st *SQLStore) Store(s score.Score) error {
	tx, err := st.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := fmt.Sprintf("INSERT INTO %s(name, value) values($1,$2)", s.Type())
	_, err = st.DB.Exec(query, s.Name(), s.Value())
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert score: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Retrieve Scores from the database by name
// Use `database/sql` pattern rather than talking directly to driver.
// This should allow for swapping out drivers.
func (st *SQLStore) Retrieve(f score.ScoreFactory, scoreType string) (map[string][]score.Score, error) {
	ret := map[string][]score.Score{}

	/* typical database/sql pattern:
	- query rows
	- defer rows.Close() to avoid memory leaks if we exit before rows.Next() == false
	- iterate with rows.Next()
		- scan values, load them into return object
		- sql package only natively supports largest datatypes,
			lucky for us we are already working with float64 and strings
	- after looping, check for errors with rows.Err()
	*/

	query := fmt.Sprintf("SELECT name, value FROM %s", scoreType)
	rows, err := st.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query scores: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name  string
			value float64
		)

		if err := rows.Scan(&name, &value); err != nil {
			return nil, fmt.Errorf("error scanning: %w", err)
		}

		s, err := score.Create(f, scoreType)
		if err != nil {
			return nil, fmt.Errorf("failed to create score: %v", err)
		}

		s.Set(name, value)
		ret[name] = append(ret[name], s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return ret, nil
}
