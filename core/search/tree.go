package search

import (
	"errors"
	"strings"
)

const (
	colon = ":"
	slash = '/'
)

var (
	ErrDupItem      = errors.New("duplicated item")
	ErrDupSlash     = errors.New("duplicated slash")
	ErrEmptyItem    = errors.New("empty item")
	ErrInvalidState = errors.New("search tree is in an invalid state")
	ErrNotFromRoot  = errors.New("path should start with /")

	NotFound Result
)

type (
	node struct {
		item interface{}

		exactNode map[string]*node
		colonNode map[string]*node
	}

	Tree struct {
		root *node
	}

	Result struct {
		Item   interface{}
		Params map[string]string
	}
)

func NewTree() *Tree {
	return &Tree{
		root: newNode(),
	}
}

func (t *Tree) Add(route string, item interface{}) error {
	if len(route) == 0 || route[0] != slash {
		return ErrNotFromRoot
	}

	if item == nil {
		return ErrEmptyItem
	}
	route = strings.Trim(route, "/")
	return add(t.root, route, item)
}

func (t *Tree) Search(route string) (Result, bool) {
	if len(route) == 0 || route[0] != slash {
		return NotFound, false
	}

	var result Result
	route = strings.Trim(route, "/")
	ok := next(t.root, strings.Split(route, "/"), &result)
	return result, ok
}

func next(n *node, paths []string, result *Result) bool {
	if len(paths) == 0 {
		if n.item == nil {
			return false
		}
		result.Item = n.item
		return true
	}
	p := paths[0]
	if p == "" {
		return next(n, paths[1:], result)
	}
	if child, ok := n.exactNode[p]; ok {
		if ok := next(child, paths[1:], result); ok {
			return true
		}
	}

	for k, n := range n.colonNode {
		if ok := next(n, paths[1:], result); ok {
			addParam(result, k[1:], p)
			return true
		}
	}
	return false
}

func addParam(result *Result, k, v string) {
	if result.Params == nil {
		result.Params = make(map[string]string)
	}

	result.Params[k] = v
}

func newNode() *node {
	return &node{
		item:      nil,
		exactNode: make(map[string]*node),
		colonNode: make(map[string]*node),
	}
}

func add(root *node, route string, item interface{}) error {
	if len(route) == 0 {
		if root.item != nil {
			return ErrDupItem
		}
		root.item = item
		return nil
	}
	paths := strings.Split(route, "/")

	nd := root
	for _, char := range paths {
		if char == "" {
			return ErrDupItem
		}
		var nextM map[string]*node
		if strings.HasPrefix(char, colon) {
			nextM = nd.colonNode
		} else {
			nextM = nd.exactNode
		}
		if child, ok := nextM[char]; ok {
			nd = child
		} else {
			child = newNode()
			nextM[char] = child
			nd = child
		}
	}
	if nd.item != nil {
		return ErrDupItem
	}

	nd.item = item
	return nil
}
