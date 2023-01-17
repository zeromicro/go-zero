package lang

import (
	"fmt"
	"reflect"
	"strconv"
)

// Placeholder is a placeholder object that can be used globally.
var Placeholder PlaceholderType

type (
	// AnyType can be used to hold any type.
	AnyType = any
	// PlaceholderType represents a placeholder type.
	PlaceholderType = struct{}
)

// Repr returns the string representation of v.
func Repr(v any) string {
	if v == nil {
		return ""
	}

	// if func (v *Type) String() string, we can't use Elem()
	switch vt := v.(type) {
	case fmt.Stringer:
		return vt.String()
	}

	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	return reprOfValue(val)
}

func reprOfValue(val reflect.Value) string {
	switch vt := val.Interface().(type) {
	case bool:
		return strconv.FormatBool(vt)
	case error:
		return vt.Error()
	case float32:
		return strconv.FormatFloat(float64(vt), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(vt, 'f', -1, 64)
	case fmt.Stringer:
		return vt.String()
	case int:
		return strconv.Itoa(vt)
	case int8:
		return strconv.Itoa(int(vt))
	case int16:
		return strconv.Itoa(int(vt))
	case int32:
		return strconv.Itoa(int(vt))
	case int64:
		return strconv.FormatInt(vt, 10)
	case string:
		return vt
	case uint:
		return strconv.FormatUint(uint64(vt), 10)
	case uint8:
		return strconv.FormatUint(uint64(vt), 10)
	case uint16:
		return strconv.FormatUint(uint64(vt), 10)
	case uint32:
		return strconv.FormatUint(uint64(vt), 10)
	case uint64:
		return strconv.FormatUint(vt, 10)
	case []byte:
		return string(vt)
	default:
		return fmt.Sprint(val.Interface())
	}
}
