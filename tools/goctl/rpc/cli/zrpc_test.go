package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
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

	var testData = []test{
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
			expectedErr: errMutilInput,
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
