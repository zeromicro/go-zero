package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMember_IsTagMember(t *testing.T) {
	t.Run("non-inline member with matching tag returns true", func(t *testing.T) {
		m := Member{Tag: `header:"Authorization"`}
		assert.True(t, m.IsTagMember("header"))
	})

	t.Run("non-inline member without matching tag returns false", func(t *testing.T) {
		m := Member{Tag: `json:"username"`}
		assert.False(t, m.IsTagMember("header"))
		assert.False(t, m.IsTagMember("path"))
		assert.False(t, m.IsTagMember("form"))
	})

	t.Run("non-inline member without any tag returns false", func(t *testing.T) {
		m := Member{}
		assert.False(t, m.IsTagMember("header"))
	})

	t.Run("inline struct without matching child tag returns false (#4800)", func(t *testing.T) {
		m := Member{
			Name:     "Pagination",
			IsInline: true,
			Type: DefineStruct{
				RawName: "Pagination",
				Members: []Member{
					{Name: "Page", Tag: `json:"page"`},
					{Name: "PageSize", Tag: `json:"pageSize"`},
				},
			},
		}
		assert.False(t, m.IsTagMember("header"))
		assert.False(t, m.IsTagMember("path"))
		assert.False(t, m.IsTagMember("form"))
	})

	t.Run("inline struct whose child has matching tag returns true", func(t *testing.T) {
		m := Member{
			Name:     "Auth",
			IsInline: true,
			Type: DefineStruct{
				RawName: "Auth",
				Members: []Member{
					{Name: "Token", Tag: `header:"Authorization"`},
				},
			},
		}
		assert.True(t, m.IsTagMember("header"))
	})

	t.Run("nested inline structs are recursed", func(t *testing.T) {
		m := Member{
			Name:     "Outer",
			IsInline: true,
			Type: DefineStruct{
				RawName: "Outer",
				Members: []Member{
					{
						Name:     "Inner",
						IsInline: true,
						Type: DefineStruct{
							RawName: "Inner",
							Members: []Member{
								{Name: "Token", Tag: `header:"X-Token"`},
							},
						},
					},
				},
			},
		}
		assert.True(t, m.IsTagMember("header"))
	})

	t.Run("nested inline structs without matching child return false", func(t *testing.T) {
		m := Member{
			Name:     "Outer",
			IsInline: true,
			Type: DefineStruct{
				RawName: "Outer",
				Members: []Member{
					{
						Name:     "Inner",
						IsInline: true,
						Type: DefineStruct{
							RawName: "Inner",
							Members: []Member{
								{Name: "Page", Tag: `json:"page"`},
							},
						},
					},
				},
			},
		}
		assert.False(t, m.IsTagMember("header"))
	})

	t.Run("inline NestedStruct whose child has matching tag returns true", func(t *testing.T) {
		m := Member{
			Name:     "Auth",
			IsInline: true,
			Type: NestedStruct{
				RawName: "Auth",
				Members: []Member{
					{Name: "Token", Tag: `header:"Authorization"`},
				},
			},
		}
		assert.True(t, m.IsTagMember("header"))
	})

	t.Run("empty inline struct returns false", func(t *testing.T) {
		m := Member{
			Name:     "Empty",
			IsInline: true,
			Type:     DefineStruct{RawName: "Empty"},
		}
		assert.False(t, m.IsTagMember("header"))
	})
}

func TestDefineStruct_GetTagMembers_InlineRegression(t *testing.T) {
	s := DefineStruct{
		RawName: "QueryUserListReq",
		Members: []Member{
			{Name: "Username", Tag: `json:"username,optional"`},
			{
				Name:     "Pagination",
				IsInline: true,
				Type: DefineStruct{
					RawName: "Pagination",
					Members: []Member{
						{Name: "Page", Tag: `json:"page"`},
						{Name: "PageSize", Tag: `json:"pageSize"`},
					},
				},
			},
		},
	}
	assert.Empty(t, s.GetTagMembers("header"),
		"inline struct without header-tagged children must not match header (#4800)")
	assert.Empty(t, s.GetTagMembers("path"))
}
