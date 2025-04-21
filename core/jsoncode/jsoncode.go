package jsoncode

import "encoding/json"

type (
	MarshalFn   func(v any) ([]byte, error)
	UnmarshalFn func(by []byte, v any) error
)

var (
	Marshal MarshalFn = func(v any) ([]byte, error) {
		return json.Marshal(v)
	}

	Unmarshal UnmarshalFn = func(by []byte, v any) error {
		return json.Unmarshal(by, v)
	}
)

func SetMarshalFn(fn MarshalFn) {
	Marshal = fn
}

func SetUnmarshalFn(fn UnmarshalFn) {
	Unmarshal = fn
}
