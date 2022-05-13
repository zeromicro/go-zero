package mapping

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/pelletier/go-toml/v2"
)

// UnmarshalTomlBytes unmarshals TOML bytes into the given v.
func UnmarshalTomlBytes(content []byte, v interface{}) error {
	return UnmarshalTomlReader(bytes.NewReader(content), v)
}

// UnmarshalTomlReader unmarshals TOML from the given io.Reader into the given v.
func UnmarshalTomlReader(r io.Reader, v interface{}) error {
	var val interface{}
	if err := toml.NewDecoder(r).Decode(&val); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(val); err != nil {
		return err
	}

	return UnmarshalJsonReader(&buf, v)
}
