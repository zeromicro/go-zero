package mapping

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// To make .json & .yaml consistent, we just use json as the tag key.
const yamlTagKey = "json"

var (
	ErrUnsupportedType = errors.New("only map-like configs are suported")

	yamlUnmarshaler = NewUnmarshaler(yamlTagKey)
)

func UnmarshalYamlBytes(content []byte, v interface{}) error {
	return unmarshalYamlBytes(content, v, yamlUnmarshaler)
}

func UnmarshalYamlReader(reader io.Reader, v interface{}) error {
	return unmarshalYamlReader(reader, v, yamlUnmarshaler)
}

func unmarshalYamlBytes(content []byte, v interface{}, unmarshaler *Unmarshaler) error {
	var o interface{}
	if err := yamlUnmarshal(content, &o); err != nil {
		return err
	}

	if m, ok := o.(map[string]interface{}); ok {
		return unmarshaler.Unmarshal(m, v)
	} else {
		return ErrUnsupportedType
	}
}

func unmarshalYamlReader(reader io.Reader, v interface{}, unmarshaler *Unmarshaler) error {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	return unmarshalYamlBytes(content, v, unmarshaler)
}

// yamlUnmarshal YAML to map[string]interface{} instead of map[interface{}]interface{}.
func yamlUnmarshal(in []byte, out interface{}) error {
	var res interface{}
	if err := yaml.Unmarshal(in, &res); err != nil {
		return err
	}

	*out.(*interface{}) = cleanupMapValue(res)
	return nil
}

func cleanupInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[Repr(k)] = cleanupMapValue(v)
	}
	return res
}

func cleanupInterfaceNumber(in interface{}) json.Number {
	return json.Number(Repr(in))
}

func cleanupInterfaceSlice(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = cleanupMapValue(v)
	}
	return res
}

func cleanupMapValue(v interface{}) interface{} {
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
		return Repr(v)
	}
}
