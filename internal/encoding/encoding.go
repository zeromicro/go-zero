package encoding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"

	"github.com/pelletier/go-toml/v2"
	"github.com/titanous/json5"
	"github.com/zeromicro/go-zero/core/lang"
	"gopkg.in/yaml.v2"
)

// Json5ToJson converts JSON5 data into its JSON representation.
func Json5ToJson(data []byte) ([]byte, error) {
	var val any
	if err := json5.Unmarshal(data, &val); err != nil {
		return nil, err
	}

	// Validate that there are no unsupported values like Infinity or NaN
	if err := validateJSONCompatible(val); err != nil {
		return nil, err
	}

	return encodeToJSON(val)
}

// validateJSONCompatible checks if the value can be represented in standard JSON.
// JSON5 allows Infinity and NaN, but standard JSON does not support these values.
func validateJSONCompatible(val any) error {
	switch v := val.(type) {
	case float64:
		if math.IsInf(v, 0) {
			return fmt.Errorf("JSON5 value Infinity cannot be represented in standard JSON")
		}
		if math.IsNaN(v) {
			return fmt.Errorf("JSON5 value NaN cannot be represented in standard JSON")
		}
	case []any:
		for _, item := range v {
			if err := validateJSONCompatible(item); err != nil {
				return err
			}
		}
	case map[string]any:
		for _, value := range v {
			if err := validateJSONCompatible(value); err != nil {
				return err
			}
		}
	case map[any]any:
		for _, value := range v {
			if err := validateJSONCompatible(value); err != nil {
				return err
			}
		}
	}
	return nil
}

// TomlToJson converts TOML data into its JSON representation.
func TomlToJson(data []byte) ([]byte, error) {
	var val any
	if err := toml.NewDecoder(bytes.NewReader(data)).Decode(&val); err != nil {
		return nil, err
	}

	return encodeToJSON(val)
}

// YamlToJson converts YAML data into its JSON representation.
func YamlToJson(data []byte) ([]byte, error) {
	var val any
	if err := yaml.Unmarshal(data, &val); err != nil {
		return nil, err
	}

	return encodeToJSON(toStringKeyMap(val))
}

// convertKeyToString ensures all keys of the map are of type string.
func convertKeyToString(in map[any]any) map[string]any {
	res := make(map[string]any)
	for k, v := range in {
		res[lang.Repr(k)] = toStringKeyMap(v)
	}
	return res
}

// convertNumberToJsonNumber converts numbers into json.Number type for compatibility.
func convertNumberToJsonNumber(in any) json.Number {
	return json.Number(lang.Repr(in))
}

// convertSlice processes slice items to ensure key compatibility.
func convertSlice(in []any) []any {
	res := make([]any, len(in))
	for i, v := range in {
		res[i] = toStringKeyMap(v)
	}
	return res
}

// encodeToJSON encodes the given value into its JSON representation.
func encodeToJSON(val any) ([]byte, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(val); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// toStringKeyMap processes the data to ensure that all map keys are of type string.
func toStringKeyMap(v any) any {
	switch v := v.(type) {
	case []any:
		return convertSlice(v)
	case map[any]any:
		return convertKeyToString(v)
	case bool, string:
		return v
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		return convertNumberToJsonNumber(v)
	default:
		return lang.Repr(v)
	}
}
