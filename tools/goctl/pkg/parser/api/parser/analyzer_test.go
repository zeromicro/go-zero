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

func TestParse_importedGroupOrderIsStable(t *testing.T) {
	dir := t.TempDir()
	for i, prefix := range []string{"/s1", "/s2", "/s3", "/s4"} {
		n := i + 1
		body := fmt.Sprintf("syntax = \"v1\"\n\n@server(\n    group: demo\n    prefix: %s\n)\nservice demo-api {\n	@handler H%d\n	get /h%d returns (Resp)\n}\n", prefix, n, n)
		err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("svc%d.api", n)), []byte(body), 0o666)
		assert.NoError(t, err)
	}
	mainAPI := "syntax = \"v1\"\n\nimport (\n	\"svc1.api\"\n	\"svc2.api\"\n	\"svc3.api\"\n	\"svc4.api\"\n)\n\ntype Resp {\n	Msg string `json:\"msg\"`\n}\n"
	mainPath := filepath.Join(dir, "main.api")
	assert.NoError(t, os.WriteFile(mainPath, []byte(mainAPI), 0o666))

	var first []string
	for i := 0; i < 20; i++ {
		apiSpec, err := Parse(mainPath, nil)
		if !assert.NoError(t, err) {
			return
		}
		var prefixes []string
		for _, g := range apiSpec.Service.Groups {
			prefixes = append(prefixes, g.Annotation.Properties["prefix"])
		}
		if first == nil {
			first = prefixes
			continue
		}
		assert.Equal(t, first, prefixes)
	}
	assert.Equal(t, []string{"/s1", "/s2", "/s3", "/s4"}, first)
}
