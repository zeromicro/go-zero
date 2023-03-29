package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/assertx"
)

func Test_Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		_, err := Parse("./testdata/example.api", nil)
		assert.Nil(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		data, err := ioutil.ReadFile("./testdata/invalid.api")
		assert.NoError(t, err)
		splits := bytes.Split(data, []byte("-----"))
		var testFile []string
		for idx, split := range splits {
			replacer := strings.NewReplacer(" ", "", "\t", "", "\n", "", "\r", "", "\f", "")
			r := replacer.Replace(string(split))
			if len(r) == 0 {
				continue
			}
			filename := filepath.Join(t.TempDir(), fmt.Sprintf("invalid%d.api", idx))
			err := ioutil.WriteFile(filename, split, 0666)
			assert.NoError(t, err)
			testFile = append(testFile, filename)
		}
		for _, v := range testFile {
			_, err := Parse(v, nil)
			assertx.Error(t, err)
		}
	})
	t.Run("circleImport", func(t *testing.T) {
		_, err := Parse("./testdata/base.api", nil)
		assertx.Error(t, err)
	})
}
