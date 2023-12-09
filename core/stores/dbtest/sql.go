package dbtest

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// RunTest runs a test function with a mock database.
func RunTest(t *testing.T, fn func(db *sql.DB, mock sqlmock.Sqlmock)) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		_ = db.Close()
	}()

	fn(db, mock)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// RunTxTest runs a test function with a mock database in a transaction.
func RunTxTest(t *testing.T, f func(tx *sql.Tx, mock sqlmock.Sqlmock)) {
	RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		tx, err := db.Begin()
		if assert.NoError(t, err) {
			f(tx, mock)
		}
	})
}
