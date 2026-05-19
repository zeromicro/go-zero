package gogen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestGolangExpr_FileType(t *testing.T) {
	t.Run("single File type", func(t *testing.T) {
		result := golangExpr(spec.FileType{RawName: "File"})
		assert.Equal(t, "*multipart.FileHeader", result)
	})

	t.Run("slice of File type", func(t *testing.T) {
		result := golangExpr(spec.ArrayType{
			RawName: "[]File",
			Value:   spec.FileType{RawName: "File"},
		})
		assert.Equal(t, "[]*multipart.FileHeader", result)
	})

	t.Run("primitive type unchanged", func(t *testing.T) {
		result := golangExpr(spec.PrimitiveType{RawName: "string"})
		assert.Equal(t, "string", result)
	})

	t.Run("array of primitive unchanged", func(t *testing.T) {
		result := golangExpr(spec.ArrayType{
			RawName: "[]string",
			Value:   spec.PrimitiveType{RawName: "string"},
		})
		assert.Equal(t, "[]string", result)
	})
}

func TestContainsFile(t *testing.T) {
	t.Run("struct with File field", func(t *testing.T) {
		types := []spec.Type{
			spec.DefineStruct{
				RawName: "UploadRequest",
				Members: []spec.Member{
					{Name: "Id", Type: spec.PrimitiveType{RawName: "string"}},
					{Name: "File", Type: spec.FileType{RawName: "File"}},
				},
			},
		}
		assert.True(t, containsFile(types))
	})

	t.Run("struct with []File field", func(t *testing.T) {
		types := []spec.Type{
			spec.DefineStruct{
				RawName: "MultiUploadRequest",
				Members: []spec.Member{
					{Name: "Id", Type: spec.PrimitiveType{RawName: "string"}},
					{Name: "Files", Type: spec.ArrayType{
						RawName: "[]File",
						Value:   spec.FileType{RawName: "File"},
					}},
				},
			},
		}
		assert.True(t, containsFile(types))
	})

	t.Run("struct without File field", func(t *testing.T) {
		types := []spec.Type{
			spec.DefineStruct{
				RawName: "Response",
				Members: []spec.Member{
					{Name: "Url", Type: spec.PrimitiveType{RawName: "string"}},
				},
			},
		}
		assert.False(t, containsFile(types))
	})

	t.Run("empty types", func(t *testing.T) {
		assert.False(t, containsFile(nil))
		assert.False(t, containsFile([]spec.Type{}))
	})
}

func TestHasFileField(t *testing.T) {
	t.Run("File field", func(t *testing.T) {
		members := []spec.Member{
			{Name: "Id", Type: spec.PrimitiveType{RawName: "string"}},
			{Name: "File", Type: spec.FileType{RawName: "File"}},
		}
		assert.True(t, hasFileField(members))
	})

	t.Run("[]File field", func(t *testing.T) {
		members := []spec.Member{
			{Name: "Files", Type: spec.ArrayType{
				RawName: "[]File",
				Value:   spec.FileType{RawName: "File"},
			}},
		}
		assert.True(t, hasFileField(members))
	})

	t.Run("no File field", func(t *testing.T) {
		members := []spec.Member{
			{Name: "Id", Type: spec.PrimitiveType{RawName: "string"}},
			{Name: "Name", Type: spec.PrimitiveType{RawName: "string"}},
		}
		assert.False(t, hasFileField(members))
	})
}
