package sqlx

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// Regression tests from PR review: these shapes must behave exactly like the
// previous per-row implementation.

// Regress #1: unexported embedded value struct must not enter the strict count
// (old unwrapFields skipped it via CanSet), while its inner tagged fields stay
// scannable by name (old getTaggedFieldValueMap recursed regardless).
func TestRegressUnexportedEmbeddedValueStrict(t *testing.T) {
	type privateSub struct {
		Ptr *int64 `db:"ptr"`
	}
	type row struct {
		ID int64 `db:"id"`
		privateSub
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	// only the id column selected, strict mode: must pass (old behavior)
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	conn := NewSqlConnFromDB(db)
	var out []row
	if err := conn.QueryRowsCtx(context.Background(), &out, "SELECT id FROM t"); err != nil {
		t.Fatalf("strict scan with unexported embedded value failed: %v", err)
	}
	if len(out) != 1 || out[0].ID != 1 {
		t.Fatalf("bad result: %+v", out)
	}
}

// Regress #1b: inner tagged field of an unexported embedded value struct is
// still matched by column name (old tagged map contained it).
func TestRegressUnexportedEmbeddedValueNamed(t *testing.T) {
	type privateSub struct {
		Ptr *int64 `db:"ptr"`
	}
	type row struct {
		ID int64 `db:"id"`
		privateSub
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"id", "ptr"}).AddRow(int64(2), int64(42)))

	conn := NewSqlConnFromDB(db)
	var out []row
	if err := conn.QueryRowsPartialCtx(context.Background(), &out, "SELECT id, ptr FROM t"); err != nil {
		t.Fatalf("named scan failed: %v", err)
	}
	if out[0].Ptr == nil || *out[0].Ptr != 42 {
		t.Fatalf("embedded named field not scanned: %+v", out)
	}
}

// Regress #2: struct whose every field is db:"-" must silently discard columns
// (old taggedMap contained "-" keys -> tagged path -> anonymous sink).
func TestRegressAllIgnoredFields(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"a", "b"}).AddRow(int64(1), int64(2)))

	conn := NewSqlConnFromDB(db)
	var out []struct {
		A int64 `db:"-"`
		B int64 `db:"-"`
	}
	if err := conn.QueryRowsPartialCtx(context.Background(), &out, "SELECT a, b FROM t"); err != nil {
		t.Fatalf("all-ignored struct must discard columns, got err: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 row: %+v", out)
	}
}

// Regress #3: embedded struct carrying db:"-" must still expose its inner
// tagged fields for name matching (old getTaggedFieldValueMap ignored the
// tag on anonymous fields).
func TestRegressEmbeddedWithTagIgnoreInnerTagged(t *testing.T) {
	type sub struct {
		V int64 `db:"v"`
	}
	type row struct {
		ID int64 `db:"id"`
		sub  `db:"-"`
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"id", "v"}).AddRow(int64(1), int64(2)))

	conn := NewSqlConnFromDB(db)
	var out []row
	if err := conn.QueryRowsPartialCtx(context.Background(), &out, "SELECT id, v FROM t"); err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if out[0].V != 2 {
		t.Fatalf("inner field of db:\"-\" embedded struct lost data: %+v", out)
	}
}

// Duplicate tags: later field wins (map assignment order in the old code).
func TestDuplicateTagNameLaterWins(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"x"}).AddRow(int64(7)))

	conn := NewSqlConnFromDB(db)
	out := []struct {
		First  int64 `db:"x"`
		Second int64 `db:"x"`
	}{{}}
	if err := conn.QueryRowPartialCtx(context.Background(), &out[0], "SELECT x FROM t"); err != nil {
		t.Fatal(err)
	}
	if out[0].Second != 7 {
		t.Fatalf("duplicate tag should let the later field win, got %+v", out)
	}
}
