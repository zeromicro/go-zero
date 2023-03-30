package jsonx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Marshal marshals v into json bytes.
func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalToString marshals v into a string.
func MarshalToString(v any) (string, error) {
	data, err := Marshal(v)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Unmarshal unmarshals data bytes into v.
func Unmarshal[T any](data []byte) (v T, err error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	v, err = unmarshalUseNumber[T](decoder)
	if err != nil {
		err = formatError(string(data), err)
	}
	return
}

// UnmarshalFromString unmarshals v from str.
func UnmarshalFromString[T any](str string) (v T, err error) {
	decoder := json.NewDecoder(strings.NewReader(str))
	v, err = unmarshalUseNumber[T](decoder)
	if err != nil {
		err = formatError(str, err)
	}
	return
}

// UnmarshalFromReader unmarshals v from reader.
func UnmarshalFromReader[T any](reader io.Reader) (v T, err error) {
	var buf strings.Builder
	teeReader := io.TeeReader(reader, &buf)
	decoder := json.NewDecoder(teeReader)
	v, err = unmarshalUseNumber[T](decoder)
	if err != nil {
		err = formatError(buf.String(), err)
	}
	return
}

func unmarshalUseNumber[T any](decoder *json.Decoder) (v T, err error) {
	decoder.UseNumber()
	err = decoder.Decode(&v)
	return
}

func formatError(v string, err error) error {
	return fmt.Errorf("string: `%s`, error: `%w`", v, err)
}
