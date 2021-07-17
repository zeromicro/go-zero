package mapping

import (
	"net/textproto"
	"strings"
)

type (
	// A Valuer interface defines the way to get values from the underlying object with keys.
	Valuer interface {
		// Value gets the value associated with the given key.
		Value(key string) (interface{}, bool)
	}

	// A MapValuer is a map that can use Value method to get values with given keys.
	MapValuer map[string]interface{}
	// A HeaderMapValuer is a map that can use Value method to get values with given keys from headers.
	HeaderMapValuer map[string]interface{}
	// A FormMapValuer is a map that can use Value method to get values with given keys from form.
	FormMapValuer map[string]interface{}
)

// Value gets the value associated with the given key from mv.
func (mv MapValuer) Value(key string) (interface{}, bool) {
	v, ok := mv[key]
	return v, ok
}

// Value gets the value associated with the given key from HeaderMapValuer.
func (mv HeaderMapValuer) Value(key string) (interface{}, bool) {
	return handleValue(mv, key, textproto.CanonicalMIMEHeaderKey)
}

// Value gets the value associated with the given key from FormMapValuer.
func (mv FormMapValuer) Value(key string) (interface{}, bool) {
	return handleValue(mv, key, func(key string) string {
		return key
	})
}

func handleValue(mv map[string]interface{}, key string, funcKey func(key string) string) (interface{}, bool) {
	var isSlice bool
	if strings.HasSuffix(key, "[]") {
		length := len(key)
		if length > 2 {
			isSlice = true
			key = key[:length-2]
		} else {
			return nil, false
		}
	}
	v, ok := mv[funcKey(key)]
	if !isSlice && ok {
		if newV, yes := v.([]string); yes && len(newV) > 0 {
			v = newV[0]
		} else {
			return nil, false
		}
	}

	return v, ok
}
