package types

import (
	"encoding/json"

	"github.com/zeromicro/go-zero/core/lang"
)

func ToStringKeyMap(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanupInterfaceSlice(v)
	case map[interface{}]interface{}:
		return cleanupInterfaceMap(v)
	case bool, string:
		return v
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		return cleanupInterfaceNumber(v)
	default:
		return lang.Repr(v)
	}
}

func cleanupInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[lang.Repr(k)] = ToStringKeyMap(v)
	}
	return res
}

func cleanupInterfaceNumber(in interface{}) json.Number {
	return json.Number(lang.Repr(in))
}

func cleanupInterfaceSlice(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = ToStringKeyMap(v)
	}
	return res
}
