package mcp

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/mapping"
)

// ParseArguments parses the arguments and populates the request object
func ParseArguments(args any, req any) error {
	switch arguments := args.(type) {
	case map[string]string:
		m := make(map[string]any, len(arguments))
		for k, v := range arguments {
			m[k] = v
		}
		return mapping.UnmarshalJsonMap(m, req, mapping.WithStringValues())
	case map[string]any:
		return mapping.UnmarshalJsonMap(arguments, req)
	default:
		return fmt.Errorf("unsupported argument type: %T", arguments)
	}
}
