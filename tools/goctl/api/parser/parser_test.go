package parser

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

//go:embed testdata/test.api
var testApi string

func TestParseContent(t *testing.T) {
	sp, err := ParseContent(testApi)
	assert.Nil(t, err)
	assert.Equal(t, spec.Doc{`// syntax doc`}, sp.Syntax.Doc)
	assert.Equal(t, spec.Doc{`// syntax comment`}, sp.Syntax.Comment)
	for _, tp := range sp.Types {
		if tp.Name() == "Request" {
			assert.Equal(t, []string{`// type doc`}, tp.Documents())
		}
	}
	for _, e := range sp.Service.Routes() {
		if e.Handler == "GreetHandler" {
			assert.Equal(t, spec.Doc{"// handler doc"}, e.HandlerDoc)
			assert.Equal(t, spec.Doc{"// handler comment"}, e.HandlerComment)
		}
	}
}

func TestParseContentWithFileType(t *testing.T) {
	t.Run("pointer to File should be rejected", func(t *testing.T) {
		content := `
syntax = "v1"
type UploadRequest {
	File *File ` + "`form:\"file\"`" + `
}
service upload-api {
	@handler upload
	post /upload (UploadRequest)
}`
		// Legacy parser panics for *File
		defer func() {
			r := recover()
			assert.NotNil(t, r)
			assert.Contains(t, r.(string), "File type cannot be used as pointer")
		}()
		_, _ = ParseContent(content, "test.api")
	})

	t.Run("single File field", func(t *testing.T) {
		content := `
syntax = "v1"
type UploadRequest {
	Id   string ` + "`form:\"id\"`" + `
	File File   ` + "`form:\"file\"`" + `
}
service upload-api {
	@handler upload
	post /upload (UploadRequest)
}`
		sp, err := ParseContent(content, "test.api")
		assert.Nil(t, err)

		var uploadReq spec.DefineStruct
		for _, tp := range sp.Types {
			if tp.Name() == "UploadRequest" {
				uploadReq = tp.(spec.DefineStruct)
			}
		}
		assert.Equal(t, "UploadRequest", uploadReq.Name())

		// Find File member
		var fileMember spec.Member
		for _, m := range uploadReq.Members {
			if m.Name == "File" {
				fileMember = m
			}
		}
		assert.Equal(t, "File", fileMember.Name)
		_, ok := fileMember.Type.(spec.FileType)
		assert.True(t, ok, "expected FileType, got %T", fileMember.Type)
	})

	t.Run("slice of File field", func(t *testing.T) {
		content := `
syntax = "v1"
type MultiUploadRequest {
	Id    string ` + "`form:\"id\"`" + `
	Files []File ` + "`form:\"files\"`" + `
}
service upload-api {
	@handler multiUpload
	post /multi-upload (MultiUploadRequest)
}`
		sp, err := ParseContent(content, "test.api")
		assert.Nil(t, err)

		var multiReq spec.DefineStruct
		for _, tp := range sp.Types {
			if tp.Name() == "MultiUploadRequest" {
				multiReq = tp.(spec.DefineStruct)
			}
		}

		var filesMember spec.Member
		for _, m := range multiReq.Members {
			if m.Name == "Files" {
				filesMember = m
			}
		}
		arrType, ok := filesMember.Type.(spec.ArrayType)
		assert.True(t, ok, "expected ArrayType, got %T", filesMember.Type)
		_, ok = arrType.Value.(spec.FileType)
		assert.True(t, ok, "expected FileType as array value, got %T", arrType.Value)
	})
}

func TestMissingService(t *testing.T) {
	sp, err := ParseContent("")
	assert.Nil(t, err)
	err = sp.Validate()
	assert.Equal(t, spec.ErrMissingService, err)
}
