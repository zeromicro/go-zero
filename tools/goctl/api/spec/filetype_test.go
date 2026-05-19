package spec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestFileType_Name(t *testing.T) {
	ft := spec.FileType{RawName: "File"}
	assert.Equal(t, "File", ft.Name())
}

func TestFileType_Comments(t *testing.T) {
	ft := spec.FileType{RawName: "File"}
	assert.Nil(t, ft.Comments())
}

func TestFileType_Documents(t *testing.T) {
	ft := spec.FileType{RawName: "File"}
	assert.Nil(t, ft.Documents())
}

func TestFileType_InStruct(t *testing.T) {
	ds := spec.DefineStruct{
		RawName: "UploadRequest",
		Members: []spec.Member{
			{
				Name: "Id",
				Type: spec.PrimitiveType{RawName: "string"},
				Tag:  `form:"id"`,
			},
			{
				Name: "File",
				Type: spec.FileType{RawName: "File"},
				Tag:  `form:"file"`,
			},
		},
	}
	assert.Equal(t, "UploadRequest", ds.Name())
	assert.Len(t, ds.Members, 2)
	assert.IsType(t, spec.FileType{}, ds.Members[1].Type)
}

func TestFileType_InArray(t *testing.T) {
	arr := spec.ArrayType{
		RawName: "[]File",
		Value:   spec.FileType{RawName: "File"},
	}
	assert.Equal(t, "[]File", arr.Name())
	assert.IsType(t, spec.FileType{}, arr.Value)
}
