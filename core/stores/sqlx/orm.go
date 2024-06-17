package sqlx

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/zeromicro/go-zero/core/mapping"
)

const tagName = "db"

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

func getTaggedFieldValueMap(v reflect.Value) (map[string]any, error) {
	rt := mapping.Deref(v.Type())
	size := rt.NumField()
	result := make(map[string]any, size)

	for i := 0; i < size; i++ {
		field := rt.Field(i)
		if field.Anonymous && mapping.Deref(field.Type).Kind() == reflect.Struct {
			inner, err := getTaggedFieldValueMap(reflect.Indirect(v).Field(i))
			if err != nil {
				return nil, err
			}

			for key, val := range inner {
				result[key] = val
			}

			continue
		}

		key := parseTagName(field)
		if len(key) == 0 {
			continue
		}

		valueField := reflect.Indirect(v).Field(i)
		valueData, err := getValueInterface(valueField)
		if err != nil {
			return nil, err
		}

		result[key] = valueData
	}

	return result, nil
}

func getValueInterface(value reflect.Value) (any, error) {
	switch value.Kind() {
	case reflect.Ptr:
		if !value.CanInterface() {
			return nil, ErrNotReadableValue
		}

		if value.IsNil() {
			baseValueType := mapping.Deref(value.Type())
			value.Set(reflect.New(baseValueType))
		}

		return value.Interface(), nil
	default:
		if !value.CanAddr() || !value.Addr().CanInterface() {
			return nil, ErrNotReadableValue
		}

		return value.Addr().Interface(), nil
	}
}

func isScanFailed(err error) bool {
	return err != nil && !errors.Is(err, context.DeadlineExceeded)
}

func mapStructFieldsIntoSlice(v reflect.Value, columns []string, strict bool) ([]any, error) {
	fields := unwrapFields(v)
	if strict && len(columns) < len(fields) {
		return nil, ErrNotMatchDestination
	}

	taggedMap, err := getTaggedFieldValueMap(v)
	if err != nil {
		return nil, err
	}

	values := make([]any, len(columns))
	if len(taggedMap) == 0 {
		if len(fields) < len(values) {
			return nil, ErrNotMatchDestination
		}

		for i := 0; i < len(values); i++ {
			valueField := fields[i]
			valueData, err := getValueInterface(valueField)
			if err != nil {
				return nil, err
			}

			values[i] = valueData
		}
	} else {
		for i, column := range columns {
			if tagged, ok := taggedMap[column]; ok {
				values[i] = tagged
			} else {
				var anonymous any
				values[i] = &anonymous
			}
		}
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

			for scanner.Next() {
				value := reflect.New(base)
				values, err := mapStructFieldsIntoSlice(value, columns, strict)
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

func unwrapFields(v reflect.Value) []reflect.Value {
	var fields []reflect.Value
	indirect := reflect.Indirect(v)

	for i := 0; i < indirect.NumField(); i++ {
		child := indirect.Field(i)
		if !child.CanSet() {
			continue
		}

		if child.Kind() == reflect.Ptr && child.IsNil() {
			baseValueType := mapping.Deref(child.Type())
			child.Set(reflect.New(baseValueType))
		}

		child = reflect.Indirect(child)
		childType := indirect.Type().Field(i)
		if child.Kind() == reflect.Struct && childType.Anonymous {
			fields = append(fields, unwrapFields(child)...)
		} else {
			fields = append(fields, child)
		}
	}

	return fields
}
