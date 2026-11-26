package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

const inlineTagAPI = `
syntax = "v1"

type (
	Auth {
		Token string ` + "`header:\"Authorization\"`" + `
	}
	Middle {
		Auth
	}
	PointerRequest {
		*Auth
	}
	NestedRequest {
		Middle
	}
	RecursiveRequest {
		Token string ` + "`header:\"X-Token\"`" + `
		*RecursiveRequest
	}
)

service test-api {
	@handler Pointer
	get /pointer (PointerRequest)

	@handler Nested
	get /nested (NestedRequest)

	@handler Recursive
	get /recursive (RecursiveRequest)
}
`

func TestParseContentResolvesInlineTypesForTagLookup(t *testing.T) {
	apiSpec, err := ParseContent(inlineTagAPI)
	require.NoError(t, err)

	for _, name := range []string{"PointerRequest", "NestedRequest", "RecursiveRequest"} {
		t.Run(name, func(t *testing.T) {
			tp := findStructByName(t, apiSpec.Types, name)
			require.NotEmpty(t, tp.GetTagMembers("header"))
			require.Empty(t, tp.GetTagMembers("path"))
		})
	}
}

func findStructByName(t *testing.T, types []spec.Type, name string) spec.DefineStruct {
	t.Helper()
	for _, tp := range types {
		if tp.Name() == name {
			defined, ok := tp.(spec.DefineStruct)
			require.True(t, ok)
			return defined
		}
	}

	t.Fatalf("type %s not found", name)
	return spec.DefineStruct{}
}
