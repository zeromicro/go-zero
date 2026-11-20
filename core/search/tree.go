package search

import (
	"errors"
	"fmt"
)

const (
	colon = ':'
	slash = '/'
)

var (
	// errDupItem means adding duplicated item.
	errDupItem = errors.New("duplicated item")
	// errDupSlash means item is started with more than one slash.
	errDupSlash = errors.New("duplicated slash")
	// errEmptyItem means adding empty item.
	errEmptyItem = errors.New("empty item")
	// errInvalidState means search tree is in an invalid state.
	errInvalidState = errors.New("search tree is in an invalid state")
	// errNotFromRoot means path is not starting with slash.
	errNotFromRoot = errors.New("path should start with /")

	// NotFound is used to hold the not found result.
	NotFound Result
)

type (
	node struct {
		item     any
		children [2]map[string]*node
	}

	// A Tree is a search tree.
	Tree struct {
		root *node
	}

	// A Result is a search result from tree.
	Result struct {
		Item   any
		Params map[string]string
	}
)

// NewTree returns a Tree.
func NewTree() *Tree {
	return &Tree{
		root: newNode(nil),
	}
}

// Add adds item to associate with route.
func (t *Tree) Add(route string, item any) error {
	if len(route) == 0 || route[0] != slash {
		return errNotFromRoot
	}

	if item == nil {
		return errEmptyItem
	}

	err := add(t.root, route[1:], item)
	switch {
	case errors.Is(err, errDupItem):
		return duplicatedItem(route)
	case errors.Is(err, errDupSlash):
		return duplicatedSlash(route)
	default:
		return err
	}
}

// Search searches item that associates with given route.
func (t *Tree) Search(route string) (Result, bool) {
	if len(route) == 0 || route[0] != slash {
		return NotFound, false
	}

	var result Result
	ok := t.next(t.root, route[1:], &result)
	return result, ok
}

func (t *Tree) next(n *node, route string, result *Result) bool {
	if len(route) == 0 && n.item != nil {
		result.Item = n.item
		return true
	}

	for i := range route {
		if route[i] != slash {
			continue
		}

		token := route[:i]
		rest := route[i+1:]
		if child, ok := n.children[0][token]; ok && t.next(child, rest, result) {
			return true
		}

		for k, v := range n.children[1] {
			if t.next(v, rest, result) {
				addParam(result, k[1:], token)
				return true
			}
		}

		return false
	}

	if child, ok := n.children[0][route]; ok && child.item != nil {
		result.Item = child.item
		return true
	}

	for k, v := range n.children[1] {
		if v.item != nil {
			result.Item = v.item
			addParam(result, k[1:], route)
			return true
		}
	}

	return false
}

func (nd *node) getChildren(route string) map[string]*node {
	if len(route) > 0 && route[0] == colon {
		return nd.children[1]
	}

	return nd.children[0]
}

func add(nd *node, route string, item any) error {
	if len(route) == 0 {
		if nd.item != nil {
			return errDupItem
		}

		nd.item = item
		return nil
	}

	if route[0] == slash {
		return errDupSlash
	}

	for i := range route {
		if route[i] != slash {
			continue
		}

		token := route[:i]
		children := nd.getChildren(token)
		if child, ok := children[token]; ok {
			if child == nil {
				return errInvalidState
			}

			return add(child, route[i+1:], item)
		}

		child := newNode(nil)
		children[token] = child
		return add(child, route[i+1:], item)
	}

	children := nd.getChildren(route)
	if child, ok := children[route]; ok {
		if child.item != nil {
			return errDupItem
		}

		child.item = item
	} else {
		children[route] = newNode(item)
	}

	return nil
}

func addParam(result *Result, k, v string) {
	if result.Params == nil {
		result.Params = make(map[string]string)
	}

	result.Params[k] = v
}

func duplicatedItem(item string) error {
	return fmt.Errorf("duplicated item for %s", item)
}

func duplicatedSlash(item string) error {
	return fmt.Errorf("duplicated slash for %s", item)
}

func newNode(item any) *node {
	return &node{
		item: item,
		children: [2]map[string]*node{
			make(map[string]*node),
			make(map[string]*node),
		},
	}
}
