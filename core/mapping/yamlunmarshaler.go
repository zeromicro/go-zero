package mapping

import (
	"errors"
	"io"

	"github.com/zeromicro/go-zero/core/internal/encoding"
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
	b, err := encoding.YamlToJson(content)
	if err != nil {
		return err
	}

	return unmarshalJsonBytes(b, v, yamlUnmarshaler)
}

// UnmarshalYamlReader unmarshals content from reader into v.
func UnmarshalYamlReader(reader io.Reader, v interface{}) error {
	b, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	return UnmarshalYamlBytes(b, v)
}
