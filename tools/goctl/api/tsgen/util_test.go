package tsgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestGenTsType(t *testing.T) {
	member := spec.Member{
		Name:     "foo",
		Type:     spec.PrimitiveType{RawName: "string"},
		Tag:      `json:"foo,options=foo|bar|options|123"`,
		Comment:  "",
		Docs:     nil,
		IsInline: false,
	}
	ty, err := genTsType(member, 1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `'foo' | 'bar' | 'options' | '123'`, ty)

	member.IsInline = true
	ty, err = genTsType(member, 1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `'foo' | 'bar' | 'options' | '123'`, ty)

	member.Type = spec.PrimitiveType{RawName: "int"}
	member.Tag = `json:"foo,options=1|3|4|123"`
	ty, err = genTsType(member, 1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `1 | 3 | 4 | 123`, ty)
}

func TestHasActualTagMembers(t *testing.T) {
	// Test with no members
	emptyStruct := spec.DefineStruct{
		RawName: "Empty",
		Members: []spec.Member{},
	}
	assert.False(t, hasActualTagMembers(emptyStruct, "form"))
	assert.False(t, hasActualTagMembers(emptyStruct, "header"))

	// Test with direct form members
	directFormStruct := spec.DefineStruct{
		RawName: "DirectForm",
		Members: []spec.Member{
			{
				Name: "Field1",
				Type: spec.PrimitiveType{RawName: "string"},
				Tag:  `form:"field1"`,
			},
		},
	}
	assert.True(t, hasActualTagMembers(directFormStruct, "form"))
	assert.False(t, hasActualTagMembers(directFormStruct, "header"))

	// Test with inline struct containing form members
	inlineFormStruct := spec.DefineStruct{
		RawName: "PaginationReq",
		Members: []spec.Member{
			{
				Name: "PageNum",
				Type: spec.PrimitiveType{RawName: "int"},
				Tag:  `form:"pageNum"`,
			},
			{
				Name: "PageSize",
				Type: spec.PrimitiveType{RawName: "int"},
				Tag:  `form:"pageSize"`,
			},
		},
	}
	parentStruct := spec.DefineStruct{
		RawName: "ParentReq",
		Members: []spec.Member{
			{
				Name:     "",
				Type:     inlineFormStruct,
				IsInline: true,
			},
		},
	}
	assert.True(t, hasActualTagMembers(parentStruct, "form"))
	assert.False(t, hasActualTagMembers(parentStruct, "header"))

	// Test with both direct and inline members
	mixedStruct := spec.DefineStruct{
		RawName: "MixedReq",
		Members: []spec.Member{
			{
				Name: "Sth",
				Type: spec.PrimitiveType{RawName: "string"},
				Tag:  `form:"sth"`,
			},
			{
				Name:     "",
				Type:     inlineFormStruct,
				IsInline: true,
			},
		},
	}
	assert.True(t, hasActualTagMembers(mixedStruct, "form"))
	assert.False(t, hasActualTagMembers(mixedStruct, "header"))

	// Test with inline struct containing only json members (body members)
	inlineJsonStruct := spec.DefineStruct{
		RawName: "JsonStruct",
		Members: []spec.Member{
			{
				Name: "Code",
				Type: spec.PrimitiveType{RawName: "int64"},
				Tag:  `json:"code"`,
			},
			{
				Name: "Msg",
				Type: spec.PrimitiveType{RawName: "string"},
				Tag:  `json:"msg"`,
			},
		},
	}
	parentJsonStruct := spec.DefineStruct{
		RawName: "ParentResp",
		Members: []spec.Member{
			{
				Name:     "",
				Type:     inlineJsonStruct,
				IsInline: true,
			},
		},
	}
	assert.False(t, hasActualTagMembers(parentJsonStruct, "form"))
	assert.False(t, hasActualTagMembers(parentJsonStruct, "header"))
}

func TestHasActualBodyMembers(t *testing.T) {
	// Test with no members
	emptyStruct := spec.DefineStruct{
		RawName: "Empty",
		Members: []spec.Member{},
	}
	assert.False(t, hasActualBodyMembers(emptyStruct))

	// Test with direct json members
	directJsonStruct := spec.DefineStruct{
		RawName: "DirectJson",
		Members: []spec.Member{
			{
				Name: "Code",
				Type: spec.PrimitiveType{RawName: "int64"},
				Tag:  `json:"code"`,
			},
		},
	}
	assert.True(t, hasActualBodyMembers(directJsonStruct))

	// Test with inline struct containing json members
	inlineJsonStruct := spec.DefineStruct{
		RawName: "BaseResp",
		Members: []spec.Member{
			{
				Name: "Code",
				Type: spec.PrimitiveType{RawName: "int64"},
				Tag:  `json:"code"`,
			},
			{
				Name: "Msg",
				Type: spec.PrimitiveType{RawName: "string"},
				Tag:  `json:"msg"`,
			},
		},
	}
	parentStruct := spec.DefineStruct{
		RawName: "ParentResp",
		Members: []spec.Member{
			{
				Name:     "",
				Type:     inlineJsonStruct,
				IsInline: true,
			},
		},
	}
	assert.True(t, hasActualBodyMembers(parentStruct))

	// Test with inline struct containing only form members (not body members)
	inlineFormStruct := spec.DefineStruct{
		RawName: "PaginationReq",
		Members: []spec.Member{
			{
				Name: "PageNum",
				Type: spec.PrimitiveType{RawName: "int"},
				Tag:  `form:"pageNum"`,
			},
		},
	}
	parentFormStruct := spec.DefineStruct{
		RawName: "ParentReq",
		Members: []spec.Member{
			{
				Name:     "",
				Type:     inlineFormStruct,
				IsInline: true,
			},
		},
	}
	assert.False(t, hasActualBodyMembers(parentFormStruct))
}

func TestHasActualNonBodyMembers(t *testing.T) {
	// Test with no members
	emptyStruct := spec.DefineStruct{
		RawName: "Empty",
		Members: []spec.Member{},
	}
	assert.False(t, hasActualNonBodyMembers(emptyStruct))

	// Test with direct form members
	directFormStruct := spec.DefineStruct{
		RawName: "DirectForm",
		Members: []spec.Member{
			{
				Name: "Field1",
				Type: spec.PrimitiveType{RawName: "string"},
				Tag:  `form:"field1"`,
			},
		},
	}
	assert.True(t, hasActualNonBodyMembers(directFormStruct))

	// Test with inline struct containing form members
	inlineFormStruct := spec.DefineStruct{
		RawName: "PaginationReq",
		Members: []spec.Member{
			{
				Name: "PageNum",
				Type: spec.PrimitiveType{RawName: "int"},
				Tag:  `form:"pageNum"`,
			},
			{
				Name: "PageSize",
				Type: spec.PrimitiveType{RawName: "int"},
				Tag:  `form:"pageSize"`,
			},
		},
	}
	parentStruct := spec.DefineStruct{
		RawName: "ParentReq",
		Members: []spec.Member{
			{
				Name:     "",
				Type:     inlineFormStruct,
				IsInline: true,
			},
		},
	}
	assert.True(t, hasActualNonBodyMembers(parentStruct))

	// Test with inline struct containing only json members (body members)
	inlineJsonStruct := spec.DefineStruct{
		RawName: "BaseResp",
		Members: []spec.Member{
			{
				Name: "Code",
				Type: spec.PrimitiveType{RawName: "int64"},
				Tag:  `json:"code"`,
			},
		},
	}
	parentJsonStruct := spec.DefineStruct{
		RawName: "ParentResp",
		Members: []spec.Member{
			{
				Name:     "",
				Type:     inlineJsonStruct,
				IsInline: true,
			},
		},
	}
	assert.False(t, hasActualNonBodyMembers(parentJsonStruct))

	// Test with both direct and inline non-body members
	mixedStruct := spec.DefineStruct{
		RawName: "MixedReq",
		Members: []spec.Member{
			{
				Name: "Sth",
				Type: spec.PrimitiveType{RawName: "string"},
				Tag:  `form:"sth"`,
			},
			{
				Name:     "",
				Type:     inlineFormStruct,
				IsInline: true,
			},
		},
	}
	assert.True(t, hasActualNonBodyMembers(mixedStruct))
}
