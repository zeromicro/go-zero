package types

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func IsInteger(v any) bool {
	switch v.(type) {
	case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int, uint:
		return true
	case float32, float64:
		s := fmt.Sprintf("%0.1f", v)
		return strings.HasSuffix(s, ".0")
	default:
		return false
	}
}

func IsBool(v any) bool {
	_, ok := v.(bool)
	return ok
}

func IsFloat(v any) bool {
	switch v.(type) {
	case float32, float64:
		return strings.Contains(fmt.Sprint(v), ".")
	default:
		return false
	}
}

func IsTime(v any) bool {
	_, ok := v.(time.Time)
	return ok
}

func IsString(v any) bool {
	_, ok := v.(string)
	return ok
}

func IsNil(v any) bool {
	tp := reflect.TypeOf(v)
	return tp == nil
}
