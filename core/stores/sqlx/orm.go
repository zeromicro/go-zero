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
type cachedFields struct {
	flat   [][]int        // all scannable fields in unwrapFields order (tagged or not)
	byName map[string]int // db tag name -> index into flat (later fields win on duplicates)
	tagged bool           // whether any field carries a db tag
	// ptrIndex lists every settable pointer field (embedded pointer intermediates
	// first, ordered by path length) that must be initialized before scanning,
	// preserving the previous side effect of unwrapFields on every row.
	ptrIndex [][]int
}

var fieldCache sync.Map // reflect.Type -> *cachedFields

// getCachedFields returns the cached mapping for a struct type, building it once.
func getCachedFields(rt reflect.Type) (*cachedFields, error) {
	if cached, ok := fieldCache.Load(rt); ok {
		return cached.(*cachedFields), nil
	}

	cf := &cachedFields{byName: make(map[string]int)}
	if err := appendFieldIndexes(cf, rt, nil); err != nil {
		return nil, err // don't cache failures; behavior matches per-row errors
	}
	cf.tagged = len(cf.byName) > 0

	fieldCache.Store(rt, cf)
	return cf, nil
}

// appendFieldIndexes recursively collects field metadata,
// mirroring the traversal of the previous per-row implementation:
//   - unexported fields are skipped, but fail with ErrNotReadableValue when
//     they carry a db tag (same error the per-row path returned)
//   - db:"-" fields are excluded from the strict count but still matched by name
//   - embedded structs (anonymous, dereferenced) are flattened recursively,
//     with pointer intermediates recorded for per-row initialization
func appendFieldIndexes(cf *cachedFields, rt reflect.Type, prefix []int) error {
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		idx := append(append([]int{}, prefix...), i)
		unexported := field.PkgPath != ""

		if field.Anonymous && mapping.Deref(field.Type).Kind() == reflect.Struct {
			if parseTagName(field) == tagIgnore {
				continue
			}
			// Unexported pointer embedding cannot be initialized (reflect forbids
			// Set through unexported fields); the previous per-row path panicked
			// on it, here the subtree is skipped instead.
			if field.Type.Kind() == reflect.Pointer && unexported {
				continue
			}
			if field.Type.Kind() == reflect.Pointer {
				cf.ptrIndex = append(cf.ptrIndex, idx)
			}
			if err := appendFieldIndexes(cf, mapping.Deref(field.Type), idx); err != nil {
				return err
			}
			continue
		}

		column := parseTagName(field)
		if unexported && len(column) > 0 {
			return ErrNotReadableValue
		}
		if unexported || column == tagIgnore {
			continue
		}

		cf.flat = append(cf.flat, idx)
		if field.Type.Kind() == reflect.Pointer && !unexported {
			cf.ptrIndex = append(cf.ptrIndex, idx)
		}
		if len(column) > 0 {
			cf.byName[column] = len(cf.flat) - 1
		}
	}

	return nil
}

// initPtrFields materializes nil pointer fields, preserving the side effect
// unwrapFields had on every row: every settable nil pointer gets allocated,
// even for columns that are not selected.
func (cf *cachedFields) initPtrFields(v reflect.Value) {
	for _, idx := range cf.ptrIndex {
		f := v.FieldByIndex(idx)
		if f.IsNil() {
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

// matchColumns resolves column names to flat field indexes once per result set;
// unmatched columns get -1 and are discarded via an anonymous sink at scan time.
func (cf *cachedFields) matchColumns(columns []string) []int {
	indexes := make([]int, len(columns))
	if !cf.tagged {
		// untagged structs map columns to fields positionally
		for i := range indexes {
			indexes[i] = i
		}
		return indexes
	}

	for i, column := range columns {
		if fi, ok := cf.byName[column]; ok {
			indexes[i] = fi
		} else {
			indexes[i] = -1
		}
	}
	return indexes
}

// buildValues assembles the scan targets for a single row.
func (cf *cachedFields) buildValues(indirect reflect.Value, colIdx []int) ([]any, error) {
	cf.initPtrFields(indirect)

	values := make([]any, len(colIdx))
	if !cf.tagged {
		for i := range values {
			field := indirect.FieldByIndex(cf.flat[i])
			if field.Kind() == reflect.Pointer {
				// unwrapFields used to hand out the pointed-to value here,
				// keep that behavior for the untagged path
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

	for i, fi := range colIdx {
		if fi < 0 {
			var anonymous any
			values[i] = &anonymous
			continue
		}
		valueData, err := getValueInterface(indirect.FieldByIndex(cf.flat[fi]))
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
