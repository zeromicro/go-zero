package sqlx

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/mapping"
)

const (
	tagIgnore = "-"
	tagName   = "db"
)

var (
	// ErrNotMatchDestination is an error that indicates not matching destination to scan.
	ErrNotMatchDestination = errors.New("not matching destination to scan")
	// ErrNotReadableValue is an error that indicates value is not addressable or interfaceable.
	ErrNotReadableValue = errors.New("value not addressable or interfaceable")
	// ErrNotSettable is an error that indicates the passed in variable is not settable.
	ErrNotSettable = errors.New("passed in variable is not settable")
	// ErrUnsupportedValueType is an error that indicates unsupported unmarshal type.
	ErrUnsupportedValueType = errors.New("unsupported unmarshal type")
)

type rowsScanner interface {
	Columns() ([]string, error)
	Err() error
	Next() bool
	Scan(v ...any) error
}

// cachedFields caches the type-level mapping from db tag names to field indexes,
// so that per-row scanning no longer re-parses struct tags, re-walks the struct
// and rebuilds the tagged value map on every row. Same approach as jmoiron/sqlx
// StructScan, which also caches the reflect work per type.
//
// The two collections intentionally mirror the two traversals the previous
// per-row code performed, because they had different semantics:
//   - flat follows unwrapFields: only settable (exported) fields, db:"-"
//     excluded (whole subtree when on an embedding). It drives the strict
//     count, the positional (untagged) path and the pointer init list.
//   - byName follows getTaggedFieldValueMap: embedded structs are flattened
//     regardless of their own tag or export status, and every non-empty db tag
//     (including "-") is kept for name matching. A tagged unexported field is
//     kept as well; scanning it later fails in getValueInterface with
//     ErrNotReadableValue, matching the error the old map build returned.
//     Paths may cross embedded pointers that ptrIndex cannot pre-allocate
//     (unexported, or db:"-" so collectFlat skipped the subtree): crossing a
//     nil one fails with ErrNotReadableValue in fieldByIndex instead of the
//     reflect panic the old per-row code hit, non-nil ones scan as before.
//
// The cache has no eviction: entries live for the process lifetime, keyed by
// reflect.Type, so the number of entries is bounded by the set of scanned
// model types. Concurrent first builds may duplicate work but are idempotent.
type cachedFields struct {
	flat   [][]int          // unwrapFields order: strict count + positional path
	byName map[string][]int // db tag name -> field index path (later fields win)
	tagged bool             // whether any field carries a db tag
	// ptrIndex lists every settable pointer field (embedded pointer intermediates
	// first, ordered by path length) that must be initialized before scanning,
	// preserving the previous side effect of unwrapFields on every row.
	ptrIndex [][]int
	// taggedPtrIndex lists the pointer fields visited with a non-empty db tag
	// during the collectByName walk. The old per-row getTaggedFieldValueMap
	// allocated the nil pointer of every such field while building the map —
	// including db:"-" leaves, tagged fields inside db:"-" or unexported
	// embeddings, and fields whose tag a later duplicate overwrote — even when
	// the column was not selected. This index preserves that side effect.
	taggedPtrIndex [][]int
}

var fieldCache sync.Map // reflect.Type -> *cachedFields

// getCachedFields returns the cached mapping for a struct type, building it once.
func getCachedFields(rt reflect.Type) (*cachedFields, error) {
	if cached, ok := fieldCache.Load(rt); ok {
		return cached.(*cachedFields), nil
	}

	cf := &cachedFields{byName: make(map[string][]int)}
	collectFlat(cf, rt, nil)
	collectByName(cf, rt, nil)
	cf.tagged = len(cf.byName) > 0

	fieldCache.Store(rt, cf)
	return cf, nil
}

// collectFlat mirrors unwrapFields: unexported fields (and whole unexported or
// db:"-" embeddings) are skipped; exported pointer fields are recorded for
// per-row initialization.
func collectFlat(cf *cachedFields, rt reflect.Type, prefix []int) {
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" { // not settable through reflect
			continue
		}
		if parseTagName(field) == tagIgnore {
			continue
		}

		idx := append(append([]int{}, prefix...), i)
		if field.Anonymous && mapping.Deref(field.Type).Kind() == reflect.Struct {
			if field.Type.Kind() == reflect.Pointer {
				cf.ptrIndex = append(cf.ptrIndex, idx)
			}
			collectFlat(cf, mapping.Deref(field.Type), idx)
			continue
		}

		if field.Type.Kind() == reflect.Pointer {
			cf.ptrIndex = append(cf.ptrIndex, idx)
		}
		cf.flat = append(cf.flat, idx)
	}
}

// collectByName mirrors getTaggedFieldValueMap: anonymous structs are flattened
// regardless of their own db tag or export status, and every non-empty tag is
// kept. Unlike collectFlat this recurses into unexported and db:"-" embedded
// pointers as well: their inner exported fields stay scannable by name when the
// pointer is set, and fieldByIndex reports a nil crossing instead of the panic
// the old per-row code hit there.
func collectByName(cf *cachedFields, rt reflect.Type, prefix []int) {
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		idx := append(append([]int{}, prefix...), i)

		if field.Anonymous && mapping.Deref(field.Type).Kind() == reflect.Struct {
			collectByName(cf, mapping.Deref(field.Type), idx)
			continue
		}

		column := parseTagName(field)
		if len(column) == 0 {
			continue
		}
		cf.byName[column] = idx
		// The old map build ran getValueInterface on every visited tagged
		// field, allocating its nil pointer as a side effect even when the
		// column was not selected or a later duplicate tag overwrote the map
		// entry. Unexported tagged leaves errored there and still do at scan
		// time, so they are not collected for initialization.
		if field.PkgPath == "" && field.Type.Kind() == reflect.Pointer {
			cf.taggedPtrIndex = append(cf.taggedPtrIndex, idx)
		}
	}
}

// initPtrFields materializes nil pointer fields, preserving both side effects
// the previous per-row code had on every row: unwrapFields allocated every
// settable nil pointer (ptrIndex, embedded pointers first), and the tagged map
// build additionally allocated the nil pointer of every visited tagged field
// (taggedPtrIndex), even for columns that were not selected.
func (cf *cachedFields) initPtrFields(v reflect.Value) {
	// ptrIndex paths only cross exported non-ignored embedded pointers, and
	// ancestors precede descendants, so plain FieldByIndex cannot panic here.
	for _, idx := range cf.ptrIndex {
		f := v.FieldByIndex(idx)
		if f.IsNil() {
			f.Set(reflect.New(f.Type().Elem()))
		}
	}
	// taggedPtrIndex paths may cross embedded pointers outside ptrIndex
	// (unexported or db:"-" tagged); skip those when nil — the old code
	// panicked there, fieldByIndex keeps scan-time errors consistent for
	// columns that actually match. Run after ptrIndex so paths crossing an
	// initialized embedded pointer are walkable.
	for _, idx := range cf.taggedPtrIndex {
		f, err := fieldByIndex(v, idx)
		if err != nil {
			continue
		}
		if f.IsNil() && f.CanSet() {
			f.Set(reflect.New(f.Type().Elem()))
		}
	}
}

func getValueInterface(value reflect.Value) (any, error) {
	if !value.CanAddr() || !value.Addr().CanInterface() {
		return nil, ErrNotReadableValue
	}

	if value.Kind() == reflect.Pointer && value.IsNil() {
		baseValueType := mapping.Deref(value.Type())
		value.Set(reflect.New(baseValueType))
	}

	return value.Addr().Interface(), nil
}

func isScanFailed(err error) bool {
	return err != nil && !errors.Is(err, context.DeadlineExceeded)
}

func mapStructFieldsIntoSlice(v reflect.Value, columns []string, strict bool) ([]any, error) {
	cf, err := getCachedFields(mapping.Deref(v.Type()))
	if err != nil {
		return nil, err
	}
	if strict && len(columns) < len(cf.flat) {
		return nil, ErrNotMatchDestination
	}
	if !cf.tagged && len(cf.flat) < len(columns) {
		return nil, ErrNotMatchDestination
	}

	return cf.buildValues(reflect.Indirect(v), cf.matchColumns(columns))
}

// matchColumns resolves column names to field index paths once per result set;
// unmatched columns get nil and are discarded via an anonymous sink at scan time.
func (cf *cachedFields) matchColumns(columns []string) [][]int {
	indexes := make([][]int, len(columns))
	if !cf.tagged {
		// untagged structs map columns to fields positionally
		for i := range indexes {
			indexes[i] = cf.flat[i]
		}
		return indexes
	}

	for i, column := range columns {
		indexes[i] = cf.byName[column]
	}
	return indexes
}

// fieldByIndex is FieldByIndex that reports an error instead of panicking when
// the path crosses a nil embedded pointer. Such an intermediate is not in
// ptrIndex when it is unexported (reflect cannot set it) or tagged db:"-"
// (collectFlat skips the subtree); the previous per-row code panicked there.
func fieldByIndex(v reflect.Value, index []int) (reflect.Value, error) {
	for i, x := range index {
		if i > 0 && v.Kind() == reflect.Pointer {
			if v.IsNil() {
				return reflect.Value{}, ErrNotReadableValue
			}
			v = v.Elem()
		}
		v = v.Field(x)
	}
	return v, nil
}

// buildValues assembles the scan targets for a single row.
func (cf *cachedFields) buildValues(indirect reflect.Value, colIdx [][]int) ([]any, error) {
	cf.initPtrFields(indirect)

	values := make([]any, len(colIdx))
	for i, idx := range colIdx {
		if idx == nil {
			var anonymous any
			values[i] = &anonymous
			continue
		}

		field, err := fieldByIndex(indirect, idx)
		if err != nil {
			return nil, err
		}
		if !cf.tagged && field.Kind() == reflect.Pointer {
			// unwrapFields used to hand out the pointed-to value on the
			// positional (untagged) path, keep that behavior
			field = field.Elem()
		}
		valueData, err := getValueInterface(field)
		if err != nil {
			return nil, err
		}
		values[i] = valueData
	}

	return values, nil
}

func parseTagName(field reflect.StructField) string {
	key := field.Tag.Get(tagName)
	if len(key) == 0 {
		return ""
	}

	options := strings.Split(key, ",")
	return strings.TrimSpace(options[0])
}

func unmarshalRow(v any, scanner rowsScanner, strict bool) error {
	if !scanner.Next() {
		if err := scanner.Err(); err != nil {
			return err
		}
		return ErrNotFound
	}

	rv := reflect.ValueOf(v)
	if err := mapping.ValidatePtr(rv); err != nil {
		return err
	}

	rte := reflect.TypeOf(v).Elem()
	rve := rv.Elem()
	switch rte.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		if !rve.CanSet() {
			return ErrNotSettable
		}

		return scanner.Scan(v)
	case reflect.Struct:
		columns, err := scanner.Columns()
		if err != nil {
			return err
		}

		values, err := mapStructFieldsIntoSlice(rve, columns, strict)
		if err != nil {
			return err
		}

		return scanner.Scan(values...)
	default:
		return ErrUnsupportedValueType
	}
}

func unmarshalRows(v any, scanner rowsScanner, strict bool) error {
	rv := reflect.ValueOf(v)
	if err := mapping.ValidatePtr(rv); err != nil {
		return err
	}

	rt := reflect.TypeOf(v)
	rte := rt.Elem()
	rve := rv.Elem()
	if !rve.CanSet() {
		return ErrNotSettable
	}

	switch rte.Kind() {
	case reflect.Slice:
		ptr := rte.Elem().Kind() == reflect.Ptr
		appendFn := func(item reflect.Value) {
			if ptr {
				rve.Set(reflect.Append(rve, item))
			} else {
				rve.Set(reflect.Append(rve, reflect.Indirect(item)))
			}
		}
		fillFn := func(value any) error {
			if err := scanner.Scan(value); err != nil {
				return err
			}

			appendFn(reflect.ValueOf(value))
			return nil
		}

		base := mapping.Deref(rte.Elem())
		switch base.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.String:
			for scanner.Next() {
				value := reflect.New(base)
				if err := fillFn(value.Interface()); err != nil {
					return err
				}
			}
		case reflect.Struct:
			columns, err := scanner.Columns()
			if err != nil {
				return err
			}

			cf, err := getCachedFields(base)
			if err != nil {
				return err
			}
			if strict && len(columns) < len(cf.flat) {
				return ErrNotMatchDestination
			}
			if !cf.tagged && len(cf.flat) < len(columns) {
				return ErrNotMatchDestination
			}

			// columns are fixed for the whole result set, match once
			colIdx := cf.matchColumns(columns)

			for scanner.Next() {
				value := reflect.New(base)
				values, err := cf.buildValues(value.Elem(), colIdx)
				if err != nil {
					return err
				}

				if err := scanner.Scan(values...); err != nil {
					return err
				}

				appendFn(value)
			}
		default:
			return ErrUnsupportedValueType
		}

		return scanner.Err()
	default:
		return ErrUnsupportedValueType
	}
}
