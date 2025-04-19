package jsoncode

import "encoding/json"

var defaultJsonCode JsonCode = &stdJson{}

func SetJsonCode(jsonCode JsonCode) {
	defaultJsonCode = jsonCode
}

func Marshal(v any) ([]byte, error) {
	return defaultJsonCode.Marshal(v)
}

func Unmarshal(by []byte, v any) error {
	return defaultJsonCode.Unmarshal(by, v)
}

type stdJson struct{}

func (s *stdJson) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (s *stdJson) Unmarshal(by []byte, v any) error {
	return json.Unmarshal(by, v)
}

type JsonCode interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}
