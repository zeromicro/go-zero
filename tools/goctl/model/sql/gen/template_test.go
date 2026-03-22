package gen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/parser"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/template"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

func TestGenTemplates(t *testing.T) {
	err := pathx.InitTemplates(category, templates)
	assert.Nil(t, err)
	dir, err := pathx.GetTemplateDir(category)
	assert.Nil(t, err)
	file := filepath.Join(dir, "model-new.tpl")
	data, err := os.ReadFile(file)
	assert.Nil(t, err)
	assert.Equal(t, string(data), template.New)
}

func TestRevertTemplate(t *testing.T) {
	name := "model-new.tpl"
	err := pathx.InitTemplates(category, templates)
	assert.Nil(t, err)

	dir, err := pathx.GetTemplateDir(category)
	assert.Nil(t, err)

	file := filepath.Join(dir, name)
	data, err := os.ReadFile(file)
	assert.Nil(t, err)

	modifyData := string(data) + "modify"
	err = pathx.CreateTemplate(category, name, modifyData)
	assert.Nil(t, err)

	data, err = os.ReadFile(file)
	assert.Nil(t, err)

	assert.Equal(t, string(data), modifyData)

	assert.Nil(t, RevertTemplate(name))

	data, err = os.ReadFile(file)
	assert.Nil(t, err)
	assert.Equal(t, template.New, string(data))
}

func TestClean(t *testing.T) {
	name := "model-new.tpl"
	err := pathx.InitTemplates(category, templates)
	assert.Nil(t, err)

	assert.Nil(t, Clean())

	dir, err := pathx.GetTemplateDir(category)
	assert.Nil(t, err)

	file := filepath.Join(dir, name)
	_, err = os.ReadFile(file)
	assert.NotNil(t, err)
}

func TestUpdate(t *testing.T) {
	name := "model-new.tpl"
	err := pathx.InitTemplates(category, templates)
	assert.Nil(t, err)

	dir, err := pathx.GetTemplateDir(category)
	assert.Nil(t, err)

	file := filepath.Join(dir, name)
	data, err := os.ReadFile(file)
	assert.Nil(t, err)

	modifyData := string(data) + "modify"
	err = pathx.CreateTemplate(category, name, modifyData)
	assert.Nil(t, err)

	data, err = os.ReadFile(file)
	assert.Nil(t, err)

	assert.Equal(t, string(data), modifyData)

	assert.Nil(t, Update())

	data, err = os.ReadFile(file)
	assert.Nil(t, err)
	assert.Equal(t, template.New, string(data))
}

func TestHasField(t *testing.T) {
	tests := []struct {
		name       string
		table      Table
		fieldName  string
		wantResult bool
	}{
		{
			name: "field exists",
			table: Table{
				Table: parser.Table{
					Fields: []*parser.Field{
						{NameOriginal: "id"},
						{NameOriginal: "name"},
						{NameOriginal: "created_at"},
					},
				},
			},
			fieldName:  "name",
			wantResult: true,
		},
		{
			name: "field does not exist",
			table: Table{
				Table: parser.Table{
					Fields: []*parser.Field{
						{NameOriginal: "id"},
						{NameOriginal: "name"},
					},
				},
			},
			fieldName:  "email",
			wantResult: false,
		},
		{
			name: "empty table",
			table: Table{
				Table: parser.Table{
					Fields: []*parser.Field{},
				},
			},
			fieldName:  "id",
			wantResult: false,
		},
		{
			name: "case sensitive match",
			table: Table{
				Table: parser.Table{
					Fields: []*parser.Field{
						{NameOriginal: "ID"},
						{NameOriginal: "Name"},
					},
				},
			},
			fieldName:  "id",
			wantResult: false,
		},
		{
			name: "exact match required",
			table: Table{
				Table: parser.Table{
					Fields: []*parser.Field{
						{NameOriginal: "user_name"},
					},
				},
			},
			fieldName:  "user_name",
			wantResult: true,
		},
		{
			name: "partial match should fail",
			table: Table{
				Table: parser.Table{
					Fields: []*parser.Field{
						{NameOriginal: "user_name"},
					},
				},
			},
			fieldName:  "user",
			wantResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := hasField(tt.table)
			result := fn(tt.fieldName)
			assert.Equal(t, tt.wantResult, result)
		})
	}
}

func TestHasFieldWithRealTable(t *testing.T) {
	// Create a realistic table structure
	table := Table{
		Table: parser.Table{
			Name: stringx.From("users"),
			Fields: []*parser.Field{
				{NameOriginal: "id", DataType: "int64"},
				{NameOriginal: "username", DataType: "string"},
				{NameOriginal: "email", DataType: "string"},
				{NameOriginal: "password", DataType: "string"},
				{NameOriginal: "created_at", DataType: "time.Time"},
				{NameOriginal: "updated_at", DataType: "time.Time"},
			},
		},
	}

	fn := hasField(table)

	// Test all existing fields
	assert.True(t, fn("id"))
	assert.True(t, fn("username"))
	assert.True(t, fn("email"))
	assert.True(t, fn("password"))
	assert.True(t, fn("created_at"))
	assert.True(t, fn("updated_at"))

	// Test non-existing fields
	assert.False(t, fn("deleted_at"))
	assert.False(t, fn("ID"))
	assert.False(t, fn("Username"))
	assert.False(t, fn(""))
}

func TestHasFieldPerformance(t *testing.T) {
	// Create a table with many fields to test performance optimization
	var fields []*parser.Field
	for i := 0; i < 1000; i++ {
		fields = append(fields, &parser.Field{
			NameOriginal: "field_" + string(rune('0'+i%10)) + string(rune('a'+i%26)),
		})
	}

	table := Table{
		Table: parser.Table{
			Fields: fields,
		},
	}

	fn := hasField(table)

	// Verify the function works correctly
	assert.True(t, fn(fields[0].NameOriginal))
	assert.True(t, fn(fields[999].NameOriginal))
	assert.False(t, fn("non_existent_field"))
}
