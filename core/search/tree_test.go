package search

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stringx"
)

type mockedRoute struct {
	route string
	value any
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

	result, ok := tree.Search(query)
	assert.True(t, ok)
	assert.Equal(t, 3, result.Item.(int))
}

func TestAddDuplicate(t *testing.T) {
	tree := NewTree()
	err := tree.Add("/a/b", 1)
	assert.Nil(t, err)
	err = tree.Add("/a/b", 2)
	assert.Error(t, errDupItem, err)
	err = tree.Add("/a/b/", 2)
	assert.Error(t, errDupItem, err)
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
	assert.Error(t, errDupSlash, err)
}

func TestSearchInvalidRoute(t *testing.T) {
	tree := NewTree()
	err := tree.Add("", 1)
	assert.Equal(t, errNotFromRoot, err)
	err = tree.Add("bad", 1)
	assert.Equal(t, errNotFromRoot, err)
}

func TestSearchInvalidItem(t *testing.T) {
	tree := NewTree()
	err := tree.Add("/", nil)
	assert.Equal(t, errEmptyItem, err)
}

func TestSearchInvalidState(t *testing.T) {
	nd := newNode("0")
	nd.children[0]["1"] = nil
	assert.Error(t, add(nd, "1/2", "2"))
}

func BenchmarkSearchTree(b *testing.B) {
	const (
		avgLen  = 1000
		entries = 10000
	)

	tree := NewTree()
	generate := func() string {
		var buf strings.Builder
		size := rand.Intn(avgLen) + avgLen/2
		val := stringx.Randn(size)
		prev := 0
		for j := rand.Intn(9) + 1; j < size; j += rand.Intn(9) + 1 {
			buf.WriteRune('/')
			buf.WriteString(val[prev:j])
			prev = j
		}
		if prev < size {
			buf.WriteRune('/')
			buf.WriteString(val[prev:])
		}
		return buf.String()
	}
	index := rand.Intn(entries)
	var query string
	for i := 0; i < entries; i++ {
		val := generate()
		if i == index {
			query = val
		}
		tree.Add(val, i)
	}

	for i := 0; i < b.N; i++ {
		tree.Search(query)
	}
}
