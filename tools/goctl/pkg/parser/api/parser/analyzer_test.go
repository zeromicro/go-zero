package parser

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/assertx"
)

func Test_Parse(t *testing.T) {
	t.Run(
		"valid", func(t *testing.T) {
			apiSpec, err := Parse("./testdata/example.api", nil)
			assert.Nil(t, err)
			ast := assert.New(t)
			ast.Equal(
				spec.Info{
					Title:   "type title here",
					Desc:    "type desc here",
					Version: "type version here",
					Author:  "type author here",
					Email:   "type email here",
					Properties: map[string]string{
						"title":   "type title here",
						"desc":    "type desc here",
						"version": "type version here",
						"author":  "type author here",
						"email":   "type email here",
					},
				}, apiSpec.Info,
			)
			ast.True(
				func() bool {
					for _, group := range apiSpec.Service.Groups {
						value, ok := group.Annotation.Properties["summary"]
						if ok {
							return value == "test"
						}
					}
					return false
				}(),
			)
		},
	)

	t.Run(
		"invalid", func(t *testing.T) {
			data, err := os.ReadFile("./testdata/invalid.api")
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
				err := os.WriteFile(filename, split, 0666)
				assert.NoError(t, err)
				testFile = append(testFile, filename)
			}
			for _, v := range testFile {
				_, err := Parse(v, nil)
				assertx.Error(t, err)
			}
		},
	)

	t.Run(
		"circleImport", func(t *testing.T) {
			_, err := Parse("./testdata/base.api", nil)
			assertx.Error(t, err)
		},
	)

	t.Run(
		"link_import", func(t *testing.T) {
			_, err := Parse("./testdata/link_import.api", nil)
			assert.Nil(t, err)
		},
	)

	t.Run(
		"duplicate_types", func(t *testing.T) {
			_, err := Parse("./testdata/duplicate_type.api", nil)
			assertx.Error(t, err)
		},
	)

	t.Run(
		"duplicate_path_expression", func(t *testing.T) {
			_, err := Parse("./testdata/duplicate_path_expression.api", nil)
			assertx.Error(t, err)
		},
	)
	t.Run(
		"duplicate_path_expression_different_prefix", func(t *testing.T) {
			_, err := Parse("./testdata/duplicate_path_expression_different_prefix.api", nil)

			assert.Nil(t, err)
		},
	)
}
