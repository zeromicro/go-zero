package sqlx

import (
	"context"
	"errors"
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
		ID  int64 `db:"id"`
		sub `db:"-"`
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

// Regress #4: a non-nil unexported embedded pointer keeps scanning by name on
// the single-row path, which scans the caller's struct in place. (Slice scans
// always build fresh rows, where such a pointer stays nil — that is #4b.)
// The old per-row map recursed through it; flagEmbedRO does not propagate
// into the exported inner fields.
func TestRegressUnexportedEmbeddedPtrNamed(t *testing.T) {
	type privateSub struct {
		Ptr *int64 `db:"ptr"`
	}
	type row struct {
		ID int64 `db:"id"`
		*privateSub
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"id", "ptr"}).AddRow(int64(2), int64(42)))

	conn := NewSqlConnFromDB(db)
	out := row{privateSub: &privateSub{}}
	if err := conn.QueryRowPartialCtx(context.Background(), &out, "SELECT id, ptr FROM t"); err != nil {
		t.Fatalf("named scan through non-nil unexported embedded pointer failed: %v", err)
	}
	if out.Ptr == nil || *out.Ptr != 42 {
		t.Fatalf("field behind unexported embedded pointer not scanned: %+v", out)
	}
}

// Regress #4b: an embedded pointer that ptrIndex cannot pre-allocate
// (unexported here, db:"-" in #4c) and is nil must fail the scan with
// ErrNotReadableValue, not the reflect panic the old per-row code hit.
// Covered on the slice path, where fresh rows keep such pointers nil.
func TestRegressUnexportedEmbeddedPtrNilErrors(t *testing.T) {
	type privateSub struct {
		Ptr *int64 `db:"ptr"`
	}
	type row struct {
		ID int64 `db:"id"`
		*privateSub
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"id", "ptr"}).AddRow(int64(1), int64(42)))

	conn := NewSqlConnFromDB(db)
	var out []row
	err := conn.QueryRowsPartialCtx(context.Background(), &out, "SELECT id, ptr FROM t")
	if !errors.Is(err, ErrNotReadableValue) {
		t.Fatalf("want ErrNotReadableValue on nil unexported embedded pointer, got: %v", err)
	}
}

// Regress #4c: exported embedded pointer tagged db:"-" — collectFlat skips the
// whole subtree, so it is never pre-allocated. Set, its inner tagged field
// scans on the single-row path (as in the old code); on a fresh row it errors.
func TestRegressIgnoredEmbeddedPtrSelectedColumn(t *testing.T) {
	type sub struct {
		V *int64 `db:"v"`
	}
	type row struct {
		ID   int64 `db:"id"`
		*sub `db:"-"`
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"id", "v"}).AddRow(int64(1), int64(7)))

	conn := NewSqlConnFromDB(db)
	out := row{sub: &sub{}}
	if err := conn.QueryRowPartialCtx(context.Background(), &out, "SELECT id, v FROM t"); err != nil {
		t.Fatalf("scan through non-nil db:\"-\" embedded pointer failed: %v", err)
	}
	if out.V == nil || *out.V != 7 {
		t.Fatalf("field behind db:\"-\" embedded pointer not scanned: %+v", out)
	}

	db2, mock2, _ := sqlmock.New()
	defer db2.Close()
	mock2.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"id", "v"}).AddRow(int64(1), int64(7)))

	conn2 := NewSqlConnFromDB(db2)
	var out2 []row
	err := conn2.QueryRowsPartialCtx(context.Background(), &out2, "SELECT id, v FROM t")
	if !errors.Is(err, ErrNotReadableValue) {
		t.Fatalf("want ErrNotReadableValue on nil db:\"-\" embedded pointer, got: %v", err)
	}
}

// Regress #5: the old per-row getTaggedFieldValueMap allocated the nil pointer
// of every visited tagged field as a side effect — including db:"-" leaves —
// even when the column was not selected. The cached path must do the same.
func TestRegressTagIgnoreLeafPtrAllocated(t *testing.T) {
	type row struct {
		A    *int64  `db:"a"`
		Skip *string `db:"-"`
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"a"}).AddRow(int64(1)))

	conn := NewSqlConnFromDB(db)
	var out []row
	if err := conn.QueryRowsPartialCtx(context.Background(), &out, "SELECT a FROM t"); err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(out) != 1 || out[0].A == nil || *out[0].A != 1 {
		t.Fatalf("values: %+v", out)
	}
	if out[0].Skip == nil {
		t.Fatal("db:\"-\" tagged pointer must stay allocated (old map-build side effect)")
	}
}

// Regress #5b: same side effect for tagged fields inside a db:"-" embedded
// struct, while an untagged pointer in that subtree stays nil (the old map
// build never touched untagged fields, and unwrapFields skipped the subtree).
func TestRegressTagIgnoreEmbedPtrSideEffects(t *testing.T) {
	type sub struct {
		Foo    *string `db:"foo"`
		Untag  *string
		LeafIg *string `db:"-"`
	}
	type row struct {
		sub `db:"-"`
		Top *int64 `db:"top"`
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"foo", "top"}).AddRow("v1", int64(9)))

	conn := NewSqlConnFromDB(db)
	var out []row
	if err := conn.QueryRowsPartialCtx(context.Background(), &out, "SELECT foo, top FROM t"); err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	r := out[0]
	if r.Foo == nil || *r.Foo != "v1" {
		t.Fatalf("tagged field in db:\"-\" embed not scanned: %+v", r)
	}
	if r.Untag != nil {
		t.Fatalf("untagged pointer in db:\"-\" embed must stay nil: %+v", r)
	}
	if r.LeafIg == nil {
		t.Fatalf("db:\"-\" leaf in db:\"-\" embed must stay allocated: %+v", r)
	}
}

// Regress #5c: when a duplicate tag overwrites an earlier map entry, the old
// map build had already allocated the earlier field's pointer while walking.
// The cached path must keep that allocation even though only the later field
// stays scannable. The plain case (both fields exported, no db:"-") is already
// covered by ptrIndex; the diverging case is the overwritten field living
// inside a db:"-" embedded struct.
func TestRegressDuplicateTagOverwrittenPtrAllocated(t *testing.T) {
	type sub struct {
		A *string `db:"x"`
	}
	type row struct {
		sub `db:"-"`
		B   *string `db:"x"`
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"x"}).AddRow("later"))

	conn := NewSqlConnFromDB(db)
	var out []row
	if err := conn.QueryRowsPartialCtx(context.Background(), &out, "SELECT x FROM t"); err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if out[0].B == nil || *out[0].B != "later" {
		t.Fatalf("later duplicate must win the column: %+v", out[0])
	}
	if out[0].A == nil {
		t.Fatal("overwritten duplicate's pointer must stay allocated (old map-build side effect)")
	}
}

// Regress #6: empty result sets must not fail field-count validation. The old
// implementation validated counts per row inside mapStructFieldsIntoSlice, so
// a zero-row query returned success no matter how columns and fields mismatched.
func TestRegressZeroRowsSkipCountValidation(t *testing.T) {
	// strict mode, tagged struct with more fields than selected columns
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}))

	conn := NewSqlConnFromDB(db)
	var out []struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}
	if err := conn.QueryRowsCtx(context.Background(), &out, "SELECT id FROM t"); err != nil {
		t.Fatalf("strict zero-row scan must succeed, got: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %+v", out)
	}

	// untagged struct with more columns than fields, non-strict
	db2, mock2, _ := sqlmock.New()
	defer db2.Close()
	mock2.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"a", "b"}))

	conn2 := NewSqlConnFromDB(db2)
	var out2 []struct {
		A int64
	}
	if err := conn2.QueryRowsPartialCtx(context.Background(), &out2, "SELECT a, b FROM t"); err != nil {
		t.Fatalf("untagged zero-row scan must succeed, got: %v", err)
	}
}

// Regress #7: an unexported field with a non-empty db tag failed every scanned
// row with ErrNotReadableValue while the old code built the tagged map per
// row, even when no column selected the field. Zero rows stayed a success,
// because the map was never built.
func TestRegressUnexportedTaggedFieldErrors(t *testing.T) {
	type row struct {
		A      int64 `db:"a"`
		hidden int64 `db:"hidden"`
	}

	// column not selected, one row: still ErrNotReadableValue
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"a"}).AddRow(int64(1)))

	conn := NewSqlConnFromDB(db)
	var out []row
	if err := conn.QueryRowsPartialCtx(context.Background(), &out, "SELECT a FROM t"); !errors.Is(err, ErrNotReadableValue) {
		t.Fatalf("want ErrNotReadableValue on unexported tagged field, got: %v", err)
	}

	// same shape on the single-row path
	db2, mock2, _ := sqlmock.New()
	defer db2.Close()
	mock2.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"a"}).AddRow(int64(1)))

	conn2 := NewSqlConnFromDB(db2)
	one := row{}
	if err := conn2.QueryRowPartialCtx(context.Background(), &one, "SELECT a FROM t"); !errors.Is(err, ErrNotReadableValue) {
		t.Fatalf("want ErrNotReadableValue on single-row path, got: %v", err)
	}

	// zero rows: no error, the old map build never ran
	db3, mock3, _ := sqlmock.New()
	defer db3.Close()
	mock3.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"a"}))

	conn3 := NewSqlConnFromDB(db3)
	var out3 []row
	if err := conn3.QueryRowsPartialCtx(context.Background(), &out3, "SELECT a FROM t"); err != nil {
		t.Fatalf("zero rows must stay success with unexported tagged field, got: %v", err)
	}

	// strict count mismatch wins over the unexported error, matching the old
	// order (unwrapFields + strict check ran before the map build)
	db4, mock4, _ := sqlmock.New()
	defer db4.Close()
	mock4.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"a"}).AddRow(int64(1)))

	conn4 := NewSqlConnFromDB(db4)
	var out4 []struct {
		A      int64 `db:"a"`
		B      int64 `db:"b"`
		hidden int64 `db:"hidden"`
	}
	if err := conn4.QueryRowsCtx(context.Background(), &out4, "SELECT a FROM t"); !errors.Is(err, ErrNotMatchDestination) {
		t.Fatalf("strict count mismatch must win over unexported error, got: %v", err)
	}
}

// Regress #7b: an unexported field tagged db:"-" errored as well in the old
// map build — any non-empty tag counts, "-" is only filtered for name
// matching, not for the getValueInterface call the old build made.
func TestRegressUnexportedTagIgnoreFieldErrors(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"a"}).AddRow(int64(1)))

	conn := NewSqlConnFromDB(db)
	var out []struct {
		A    int64 `db:"a"`
		skip int64 `db:"-"`
	}
	if err := conn.QueryRowsPartialCtx(context.Background(), &out, "SELECT a FROM t"); !errors.Is(err, ErrNotReadableValue) {
		t.Fatalf("want ErrNotReadableValue on unexported db:\"-\" field, got: %v", err)
	}
}
