package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
)

type test struct {
	source      []string
	expected    string
	expectedErr error
}

func Test_GetSourceProto(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		console.Error(err.Error())
		return
	}

	testData := []test{
		{
			source:   []string{"a.proto"},
			expected: filepath.Join(pwd, "a.proto"),
		},
		{
			source:   []string{"/foo/bar/a.proto"},
			expected: "/foo/bar/a.proto",
		},
		{
			source:      []string{"a.proto", "b.proto"},
			expectedErr: errMultiInput,
		},
		{
			source:      []string{"", "--go_out=."},
			expectedErr: errInvalidInput,
		},
	}

	for _, d := range testData {
		ret, err := getSourceProto(d.source, pwd)
		if d.expectedErr != nil {
			assert.Equal(t, d.expectedErr, err)
			continue
		}

		assert.Equal(t, d.expected, ret)
	}
}

func Test_RemoveGoctlFlag(t *testing.T) {
	testData := []test{
		{
			source:   strings.Fields("protoc foo.proto --go_out=. --go_opt=bar --zrpc_out=. --style go-zero --home=foo"),
			expected: "protoc foo.proto --go_out=. --go_opt=bar",
		},
		{
			source:   strings.Fields("foo bar foo.proto"),
			expected: "foo bar foo.proto",
		},
		{
			source:   strings.Fields("protoc foo.proto --go_out . --style=go_zero --home ."),
			expected: "protoc foo.proto --go_out .",
		},
		{
			source:   strings.Fields(`protoc foo.proto --go_out . --style="go_zero" --home="."`),
			expected: "protoc foo.proto --go_out .",
		},
		{
			source:   strings.Fields(`protoc foo.proto --go_opt=. --zrpc_out . --style=goZero  --home=bar`),
			expected: "protoc foo.proto --go_opt=.",
		},
		{
			source:   strings.Fields(`protoc foo.proto --go_opt=. --zrpc_out="bar" --style=goZero  --home=bar`),
			expected: "protoc foo.proto --go_opt=.",
		},
		{
			source:   strings.Fields(`protoc --go_opt=. --go-grpc_out=. --zrpc_out=. foo.proto`),
			expected: "protoc --go_opt=. --go-grpc_out=. foo.proto",
		},
		{
			source:   strings.Fields(`protoc --go_opt=. --go-grpc_out=. --zrpc_out=. --remote=foo --branch=bar foo.proto`),
			expected: "protoc --go_opt=. --go-grpc_out=. foo.proto",
		},
		{
			source:   strings.Fields(`protoc --go_opt=. --go-grpc_out=. --zrpc_out=. --remote foo --branch bar foo.proto`),
			expected: "protoc --go_opt=. --go-grpc_out=. foo.proto",
		},
	}
	for _, e := range testData {
		cmd := strings.Join(removeGoctlFlag(e.source), " ")
		assert.Equal(t, e.expected, cmd)
	}
}

func Test_RemovePluginFlag(t *testing.T) {
	testData := []test{
		{
			source:   strings.Fields("plugins=grpc:."),
			expected: ".",
		},
		{
			source:   strings.Fields("plugins=g1,g2:."),
			expected: ".",
		},
		{
			source:   strings.Fields("g1,g2:."),
			expected: ".",
		},
		{
			source:   strings.Fields("plugins=g1,g2:foo"),
			expected: "foo",
		},
	}

	for _, e := range testData {
		data := removePluginFlag(e.source[0])
		assert.Equal(t, e.expected, data)
	}
}
