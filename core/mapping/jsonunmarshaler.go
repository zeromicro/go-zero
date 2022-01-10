package mapping

import (
	"io"

	"github.com/zeromicro/go-zero/core/jsonx"
)

const jsonTagKey = "json"

var jsonUnmarshaler = NewUnmarshaler(jsonTagKey)

// UnmarshalJsonBytes unmarshals content into v.
func UnmarshalJsonBytes(content []byte, v interface{}) error {
	return unmarshalJsonBytes(content, v, jsonUnmarshaler)
}

// UnmarshalJsonMap unmarshals content from m into v.
func UnmarshalJsonMap(m map[string]interface{}, v interface{}) error {
	return jsonUnmarshaler.Unmarshal(m, v)
}

// UnmarshalJsonReader unmarshals content from reader into v.
func UnmarshalJsonReader(reader io.Reader, v interface{}) error {
	return unmarshalJsonReader(reader, v, jsonUnmarshaler)
}

func unmarshalJsonBytes(content []byte, v interface{}, unmarshaler *Unmarshaler) error {
	var m map[string]interface{}
	if err := jsonx.Unmarshal(content, &m); err != nil {
		return err
	}

	return unmarshaler.Unmarshal(m, v)
}

func unmarshalJsonReader(reader io.Reader, v interface{}, unmarshaler *Unmarshaler) error {
	var m map[string]interface{}
	if err := jsonx.UnmarshalFromReader(reader, &m); err != nil {
		return err
	}

	return unmarshaler.Unmarshal(m, v)
}
