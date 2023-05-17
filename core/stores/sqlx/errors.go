package sqlx

import (
	"database/sql"
	"errors"
)

var (
	// ErrNotFound is an alias of sql.ErrNoRows
	ErrNotFound = sql.ErrNoRows

	errCantNestTx    = errors.New("cannot nest transactions")
	errNoRawDBFromTx = errors.New("cannot get raw db from transaction")
)
