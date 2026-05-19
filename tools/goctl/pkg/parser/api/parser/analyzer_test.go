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

func Test_Parse_FileType(t *testing.T) {
	t.Run("pointer to File should be rejected", func(t *testing.T) {
		content := `syntax = "v1"
type UploadRequest {
	File *File ` + "`form:\"file\"`" + `
}
service upload-api {
	@handler upload
	post /upload (UploadRequest)
}`
		_, err := Parse("test.api", content)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "File type cannot be used as pointer")
	})

	t.Run("single File field", func(t *testing.T) {
		content := `syntax = "v1"
type UploadRequest {
	Id   string ` + "`form:\"id\"`" + `
	File File   ` + "`form:\"file\"`" + `
}
service upload-api {
	@handler upload
	post /upload (UploadRequest)
}`
		apiSpec, err := Parse("test.api", content)
		assert.Nil(t, err)

		var uploadReq spec.DefineStruct
		for _, tp := range apiSpec.Types {
			if tp.Name() == "UploadRequest" {
				uploadReq = tp.(spec.DefineStruct)
			}
		}
		assert.Equal(t, "UploadRequest", uploadReq.Name())
		assert.Len(t, uploadReq.Members, 2)

		// Check File field is FileType
		fileMember := uploadReq.Members[1]
		assert.Equal(t, "File", fileMember.Name)
		_, ok := fileMember.Type.(spec.FileType)
		assert.True(t, ok, "expected FileType, got %T", fileMember.Type)
	})

	t.Run("slice of File field", func(t *testing.T) {
		content := `syntax = "v1"
type MultiUploadRequest {
	Id    string ` + "`form:\"id\"`" + `
	Files []File ` + "`form:\"files\"`" + `
}
service upload-api {
	@handler multiUpload
	post /multi-upload (MultiUploadRequest)
}`
		apiSpec, err := Parse("test.api", content)
		assert.Nil(t, err)

		var multiReq spec.DefineStruct
		for _, tp := range apiSpec.Types {
			if tp.Name() == "MultiUploadRequest" {
				multiReq = tp.(spec.DefineStruct)
			}
		}
		assert.Equal(t, "MultiUploadRequest", multiReq.Name())

		// Check []File field is ArrayType with FileType value
		filesMember := multiReq.Members[1]
		assert.Equal(t, "Files", filesMember.Name)
		arrType, ok := filesMember.Type.(spec.ArrayType)
		assert.True(t, ok, "expected ArrayType, got %T", filesMember.Type)
		_, ok = arrType.Value.(spec.FileType)
		assert.True(t, ok, "expected FileType as array value, got %T", arrType.Value)
	})

	t.Run("mixed struct with File and non-File", func(t *testing.T) {
		content := `syntax = "v1"
type MixedRequest {
	Name   string ` + "`form:\"name\"`" + `
	File   File   ` + "`form:\"file\"`" + `
	Files  []File ` + "`form:\"files\"`" + `
	Count  int    ` + "`form:\"count\"`" + `
}
service upload-api {
	@handler mixed
	post /mixed (MixedRequest)
}`
		apiSpec, err := Parse("test.api", content)
		assert.Nil(t, err)

		var mixedReq spec.DefineStruct
		for _, tp := range apiSpec.Types {
			if tp.Name() == "MixedRequest" {
				mixedReq = tp.(spec.DefineStruct)
			}
		}
		assert.Len(t, mixedReq.Members, 4)

		// Name -> PrimitiveType
		_, ok := mixedReq.Members[0].Type.(spec.PrimitiveType)
		assert.True(t, ok)

		// File -> FileType
		_, ok = mixedReq.Members[1].Type.(spec.FileType)
		assert.True(t, ok)

		// Files -> ArrayType(FileType)
		arrType, ok := mixedReq.Members[2].Type.(spec.ArrayType)
		assert.True(t, ok)
		_, ok = arrType.Value.(spec.FileType)
		assert.True(t, ok)

		// Count -> PrimitiveType
		_, ok = mixedReq.Members[3].Type.(spec.PrimitiveType)
		assert.True(t, ok)
	})
}

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
