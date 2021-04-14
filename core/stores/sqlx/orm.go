package sqlx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/tal-tech/go-zero/core/mapping"
)

const tagName = "db"

var (
	ErrNotMatchDestination  = errors.New("not matching destination to scan")
	ErrNotReadableValue     = errors.New("value not addressable or interfaceable")
	ErrNotSettable          = errors.New("passed in variable is not settable")
	ErrUnsupportedValueType = errors.New("unsupported unmarshal type")
)

type rowsScanner interface {
	Columns() ([]string, error)
	Err() error
	Next() bool
	Scan(v ...interface{}) error
}

func getTaggedFieldValueMap(v reflect.Value) (map[string]interface{}, error) {
	rt := mapping.Deref(v.Type())
	size := rt.NumField()
	result := make(map[string]interface{}, size)

	for i := 0; i < size; i++ {
		key := parseTagName(rt.Field(i))
		if len(key) == 0 {
			return nil, nil
		}

		valueField := reflect.Indirect(v).Field(i)
		switch valueField.Kind() {
		case reflect.Ptr:
			if !valueField.CanInterface() {
				return nil, ErrNotReadableValue
			}
			if valueField.IsNil() {
				baseValueType := mapping.Deref(valueField.Type())
				valueField.Set(reflect.New(baseValueType))
			}
			result[key] = valueField.Interface()
		default:
			if !valueField.CanAddr() || !valueField.Addr().CanInterface() {
				return nil, ErrNotReadableValue
			}
			result[key] = valueField.Addr().Interface()
		}
	}

	return result, nil
}

func mapStructFieldsIntoSlice(v reflect.Value, columns []string, strict bool) ([]interface{}, error) {
	fields := unwrapFields(v)
	if strict && len(columns) < len(fields) {
		return nil, ErrNotMatchDestination
	}

	taggedMap, err := getTaggedFieldValueMap(v)
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	if len(taggedMap) == 0 {
		for i := 0; i < len(values); i++ {
			valueField := fields[i]
			switch valueField.Kind() {
			case reflect.Ptr:
				if !valueField.CanInterface() {
					return nil, ErrNotReadableValue
				}
				if valueField.IsNil() {
					baseValueType := mapping.Deref(valueField.Type())
					valueField.Set(reflect.New(baseValueType))
				}
				values[i] = valueField.Interface()
			default:
				if !valueField.CanAddr() || !valueField.Addr().CanInterface() {
					return nil, ErrNotReadableValue
				}
				values[i] = valueField.Addr().Interface()
			}
		}
	} else {
		for i, column := range columns {
			if tagged, ok := taggedMap[column]; ok {
				values[i] = tagged
			} else {
				var anonymous interface{}
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
	} else {
		options := strings.Split(key, ",")
		return options[0]
	}
}

func unmarshalRow(v interface{}, scanner rowsScanner, strict bool) error {
	if !scanner.Next() {
		if err := scanner.Err(); err != nil {
			return err
		}
		return ErrNotFound
	}

	rv := reflect.ValueOf(v)
	if err := mapping.ValidatePtr(&rv); err != nil {
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
		if rve.CanSet() {
			return scanner.Scan(v)
		} else {
			return ErrNotSettable
		}
	case reflect.Struct:
		columns, err := scanner.Columns()
		if err != nil {
			return err
		}
		if values, err := mapStructFieldsIntoSlice(rve, columns, strict); err != nil {
			return err
		} else {
			return scanner.Scan(values...)
		}
	case reflect.Map:
		if item, err := sqlRowToMap(scanner); err != nil {
			return err
		} else {
			*v.(*map[string]interface{}) = *item
			return nil
		}
	default:
		return ErrUnsupportedValueType
	}
}

func unmarshalRows(v interface{}, scanner rowsScanner, strict bool) error {
	rv := reflect.ValueOf(v)
	if err := mapping.ValidatePtr(&rv); err != nil {
		return err
	}

	rt := reflect.TypeOf(v)
	rte := rt.Elem()
	rve := rv.Elem()
	switch rte.Kind() {
	case reflect.Slice:
		if rve.CanSet() {
			ptr := rte.Elem().Kind() == reflect.Ptr
			appendFn := func(item reflect.Value) {
				if ptr {
					rve.Set(reflect.Append(rve, item))
				} else {
					rve.Set(reflect.Append(rve, reflect.Indirect(item)))
				}
			}
			fillFn := func(value interface{}) error {
				if rve.CanSet() {
					if err := scanner.Scan(value); err != nil {
						return err
					} else {
						appendFn(reflect.ValueOf(value))
						return nil
					}
				}
				return ErrNotSettable
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
					if values, err := mapStructFieldsIntoSlice(value, columns, strict); err != nil {
						return err
					} else {
						if err := scanner.Scan(values...); err != nil {
							return err
						} else {
							appendFn(value)
						}
					}
				}
			case reflect.Map:
				list := make([]map[string]interface{}, 0)
				for scanner.Next() {
					if item, err := sqlRowToMap(scanner); err != nil {
						return err
					} else {
						list = append(list, *item)
					}
				}
				*v.(*[]map[string]interface{}) = list
			default:
				fmt.Println("==============ErrUnsupportedValueType==============")
				return ErrUnsupportedValueType
			}

			return nil
		} else {
			return ErrNotSettable
		}
	default:
		return ErrUnsupportedValueType
	}
}

func unwrapFields(v reflect.Value) []reflect.Value {
	var fields []reflect.Value
	indirect := reflect.Indirect(v)

	for i := 0; i < indirect.NumField(); i++ {
		child := indirect.Field(i)
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

func sqlRowToMap(rows rowsScanner) (*map[string]interface{}, error) {
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength)
	for index, _ := range cache {
		var a interface{}
		cache[index] = &a
	}
	_ = rows.Scan(cache...)
	item := make(map[string]interface{})
	for i, data := range cache {
		item[columns[i]] = Readval(*data.(*interface{})) //取实际类型
	}
	return &item, nil
}

func Readval(value interface{}) interface{} {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		return value.(float64)
	case float32:
		return value.(float32)
	case int:
		return value.(int)
	case uint:
		return value.(uint)
	case int8:
		return value.(int8)
	case uint8:
		return value.(uint8)
	case int16:
		return value.(int16)
	case uint16:
		return value.(uint16)
	case int32:
		return value.(int32)
	case uint32:
		return value.(uint32)
	case int64:
		return value.(int64)
	case uint64:
		return value.(uint64)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		return value
	}
	return key
}
