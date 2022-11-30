package mapping

import (
	"errors"
	"io"

	"github.com/zeromicro/go-zero/core/internal/types"
	"gopkg.in/yaml.v2"
)

// To make .json & .yaml consistent, we just use json as the tag key.
const yamlTagKey = "json"

var (
	// ErrUnsupportedType is an error that indicates the config format is not supported.
	ErrUnsupportedType = errors.New("only map-like configs are supported")

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

func unmarshal(unmarshaler *Unmarshaler, o, v interface{}) error {
	if m, ok := o.(map[string]interface{}); ok {
		return unmarshaler.Unmarshal(m, v)
	}

	return ErrUnsupportedType
}

func unmarshalYamlBytes(content []byte, v interface{}, unmarshaler *Unmarshaler) error {
	var o interface{}
	if err := yamlUnmarshal(content, &o); err != nil {
		return err
	}

	return unmarshal(unmarshaler, o, v)
}

func unmarshalYamlReader(reader io.Reader, v interface{}, unmarshaler *Unmarshaler) error {
	var res interface{}
	if err := yaml.NewDecoder(reader).Decode(&res); err != nil {
		return err
	}

	return unmarshal(unmarshaler, types.ToStringKeyMap(res), v)
}

// yamlUnmarshal YAML to map[string]interface{} instead of map[interface{}]interface{}.
func yamlUnmarshal(in []byte, out interface{}) error {
	var res interface{}
	if err := yaml.Unmarshal(in, &res); err != nil {
		return err
	}

	*out.(*interface{}) = types.ToStringKeyMap(res)
	return nil
}
