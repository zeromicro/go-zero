package conf

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/hash"
)

func TestConfigJson(t *testing.T) {
	tests := []string{
		".json",
		".yaml",
		".yml",
	}
	text := `{
	"a": "foo",
	"b": 1
}`
	for _, test := range tests {
		test := test
		t.Run(test, func(t *testing.T) {
			t.Parallel()

			tmpfile, err := createTempFile(test, text)
			assert.Nil(t, err)
			defer os.Remove(tmpfile)

			var val struct {
				A string `json:"a"`
				B int    `json:"b"`
			}
			MustLoad(tmpfile, &val)
			assert.Equal(t, "foo", val.A)
			assert.Equal(t, 1, val.B)
		})
	}
}

func createTempFile(ext, text string) (string, error) {
	tmpfile, err := ioutil.TempFile(os.TempDir(), hash.Md5Hex([]byte(text))+"*"+ext)
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(tmpfile.Name(), []byte(text), os.ModeTemporary); err != nil {
		return "", err
	}

	filename := tmpfile.Name()
	if err = tmpfile.Close(); err != nil {
		return "", err
	}

	return filename, nil
}
