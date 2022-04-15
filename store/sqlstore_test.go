package store

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bdharris08/scorekeeper/score"
)

func TestStore(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	score := &score.TestScore{TName: "test", TValue: float64(0)}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO trials").WithArgs(score.Name(), score.Value()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	st := NewSQLStore(db, "trials")
	if err := st.Store(score); err != nil {
		t.Fatalf("failed to store test score: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestRetrieve(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"action", "value"}).
		AddRow("test", float64(0)).
		AddRow("test", float64(1))

	mock.ExpectQuery("SELECT action, value FROM trials").WillReturnRows(rows)

	st := NewSQLStore(db, "trials")
	got, err := st.Retrieve()
	if err != nil {
		t.Fatalf("failed to retrieve rows: %v", err)
	}

	if len(got) != 1 {
		t.Errorf("expected result to have one score")
	}

	values, ok := got["test"]
	if !ok {
		t.Error("expected to find score 'test' in result")
	}

	expectedValues := map[float64]bool{
		0: true,
		1: true,
	}

	if g, e := len(values), len(expectedValues); g != e {
		t.Errorf("expected %d values but got %d", e, g)
	}

	gotMap := map[float64]bool{}
	for _, s := range values {
		e, ok := s.Value().(float64)
		gotMap[e] = true
		if !ok {
			t.Error("failed to type assert value")
		}
		if _, ok := expectedValues[e]; !ok {
			t.Errorf("found unexpected value %f", s.Value())
		}
	}

	for k, _ := range expectedValues {
		if _, ok := gotMap[k]; !ok {
			t.Errorf("expected but did not get value %f", k)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
