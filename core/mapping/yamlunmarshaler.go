package mapping

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// To make .json & .yaml consistent, we just use json as the tag key.
const yamlTagKey = "json"

var (
	// ErrUnsupportedType is an error that indicates the config format is not supported.
	ErrUnsupportedType = errors.New("only map-like configs are suported")

	yamlUnmarshaler = NewUnmarshaler(yamlTagKey)
)

// UnmarshalYamlBytes unmarshals content into v.
func UnmarshalYamlBytes(content []byte, v interface{}) error {
	return unmarshalYamlBytes(content, v, yamlUnmarshaler)
}

// UnmarshalYamlReader unmarshals content from reader into v.
func UnmarshalYamlReader(reader io.Reader, v interface{}) error {
	return unmarshalYamlReader(reader, v, yamlUnmarshaler)
}

func unmarshalYamlBytes(content []byte, v interface{}, unmarshaler *Unmarshaler) error {
	var o interface{}
	if err := yamlUnmarshal(content, &o); err != nil {
		return err
	}

	if m, ok := o.(map[string]interface{}); ok {
		env := make(map[string]string)
		for kEnv, vEnv := range m {
			if vEnvStr, yes := vEnv.(string); yes {
				if vParse, yes := parseIfEnv(vEnvStr, env); yes {
					m[kEnv] = vParse
				}
			}

		}

		return unmarshaler.Unmarshal(m, v)
	}

	return ErrUnsupportedType
}

func parseIfEnv(value string, env map[string]string) (string, bool) {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		envKey := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")
		if v, ok := env[envKey]; !ok {
			v = os.Getenv(envKey)
			env[envKey] = v
			return v, true
		} else {
			return v, true
		}
	}
	return "", false
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
