package encoding

import (
	"bytes"
	"encoding/json"

	"github.com/pelletier/go-toml/v2"
	"github.com/zeromicro/go-zero/core/lang"
	"gopkg.in/yaml.v2"
)

func TomlToJson(data []byte) ([]byte, error) {
	var val interface{}
	if err := toml.NewDecoder(bytes.NewReader(data)).Decode(&val); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(val); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func YamlToJson(data []byte) ([]byte, error) {
	var val interface{}
	if err := yaml.Unmarshal(data, &val); err != nil {
		return nil, err
	}

	val = toStringKeyMap(val)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(val); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func convertKeyToString(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[lang.Repr(k)] = toStringKeyMap(v)
	}
	return res
}

func convertNumberToJsonNumber(in interface{}) json.Number {
	return json.Number(lang.Repr(in))
}

func convertSlice(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = toStringKeyMap(v)
	}
	return res
}

func toStringKeyMap(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return convertSlice(v)
	case map[interface{}]interface{}:
		return convertKeyToString(v)
	case bool, string:
		return v
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		return convertNumberToJsonNumber(v)
	default:
		return lang.Repr(v)
	}
}
