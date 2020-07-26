package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockedRoute struct {
	route string
	value int
}

func TestSearch(t *testing.T) {
	routes := []mockedRoute{
		{"/", 1},
		{"/api", 2},
		{"/img", 3},
		{"/:layer1", 4},
		{"/api/users", 5},
		{"/img/jpgs", 6},
		{"/img/jpgs", 7},
		{"/api/:layer2", 8},
		{"/:layer1/:layer2", 9},
		{"/:layer1/:layer2/users", 10},
	}

	tests := []struct {
		query    string
		expect   int
		params   map[string]string
		contains bool
	}{
		{
			query:    "",
			contains: false,
		},
		{
			query:    "/",
			expect:   1,
			contains: true,
		},
		{
			query:  "/wildcard",
			expect: 4,
			params: map[string]string{
				"layer1": "wildcard",
			},
			contains: true,
		},
		{
			query:  "/wildcard/",
			expect: 4,
			params: map[string]string{
				"layer1": "wildcard",
			},
			contains: true,
		},
		{
			query:    "/a/b/c",
			contains: false,
		},
		{
			query:  "/a/b",
			expect: 9,
			params: map[string]string{
				"layer1": "a",
				"layer2": "b",
			},
			contains: true,
		},
		{
			query:  "/a/b/",
			expect: 9,
			params: map[string]string{
				"layer1": "a",
				"layer2": "b",
			},
			contains: true,
		},
		{
			query:  "/a/b/users",
			expect: 10,
			params: map[string]string{
				"layer1": "a",
				"layer2": "b",
			},
			contains: true,
		},
	}

	for _, test := range tests {
		t.Run(test.query, func(t *testing.T) {
			tree := NewTree()
			for _, r := range routes {
				tree.Add(r.route, r.value)
			}
			result, ok := tree.Search(test.query)
			assert.Equal(t, test.contains, ok)
			if ok {
				actual := result.Item.(int)
				assert.EqualValues(t, test.params, result.Params)
				assert.Equal(t, test.expect, actual)
			}
		})
	}
}

func TestStrictSearch(t *testing.T) {
	routes := []mockedRoute{
		{"/api/users", 1},
		{"/api/:layer", 2},
	}
	query := "/api/users"

	tree := NewTree()
	for _, r := range routes {
		tree.Add(r.route, r.value)
	}

	for i := 0; i < 1000; i++ {
		result, ok := tree.Search(query)
		assert.True(t, ok)
		assert.Equal(t, 1, result.Item.(int))
	}
}

func TestStrictSearchSibling(t *testing.T) {
	routes := []mockedRoute{
		{"/api/:user/profile/name", 1},
		{"/api/:user/profile", 2},
		{"/api/:user/name", 3},
		{"/api/:layer", 4},
	}
	query := "/api/123/name"

	tree := NewTree()
	for _, r := range routes {
		tree.Add(r.route, r.value)
	}

	for i := 0; i < 1000; i++ {
		result, ok := tree.Search(query)
		assert.True(t, ok)
		assert.Equal(t, 3, result.Item.(int))
	}
}

func TestAddDuplicate(t *testing.T) {
	tree := NewTree()
	err := tree.Add("/a/b", 1)
	assert.Nil(t, err)
	err = tree.Add("/a/b", 2)
	assert.Equal(t, ErrDupItem, err)
	err = tree.Add("/a/b/", 2)
	assert.Equal(t, ErrDupItem, err)
}

func TestPlain(t *testing.T) {
	tree := NewTree()
	err := tree.Add("/a/b", 1)
	assert.Nil(t, err)
	err = tree.Add("/a/c", 2)
	assert.Nil(t, err)
	_, ok := tree.Search("/a/d")
	assert.False(t, ok)
}

func TestSearchWithDoubleSlashes(t *testing.T) {
	tree := NewTree()
	err := tree.Add("//a", 1)
	assert.Error(t, ErrDupSlash, err)
}

func TestSearchInvalidRoute(t *testing.T) {
	tree := NewTree()
	err := tree.Add("", 1)
	assert.Equal(t, ErrNotFromRoot, err)
	err = tree.Add("bad", 1)
	assert.Equal(t, ErrNotFromRoot, err)
}

func TestSearchInvalidItem(t *testing.T) {
	tree := NewTree()
	err := tree.Add("/", nil)
	assert.Equal(t, ErrEmptyItem, err)
}
