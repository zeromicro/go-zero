package search

import (
	"errors"
	"fmt"
	"strings"
)

const (
	colon    = ':'
	slash    = '/'
	asterisk = '*'
)

var (
	// errDupItem means adding duplicated item.
	errDupItem = errors.New("duplicated item")
	// errConflictItem means adding conflicted item.
	errConflictItem = errors.New("conflicted item")
	// errDupSlash means item is started with more than one slash.
	errDupSlash = errors.New("duplicated slash")
	// errEmptyItem means adding empty item.
	errEmptyItem = errors.New("empty item")
	// errInvalidState means search tree is in an invalid state.
	errInvalidState = errors.New("search tree is in an invalid state")
	// errNotFromRoot means path is not starting with slash.
	errNotFromRoot = errors.New("path should start with /")
	// errSlashAfterAsterisk means path contains / after *.
	errSlashAfterAsterisk = errors.New("path contains / after *")

	// NotFound is used to hold the not found result.
	NotFound Result
)

type (
	innerResult struct {
		key   string
		value string
		named bool
		found bool
	}

	node struct {
		item     interface{}
		children [3]map[string]*node
	}

	// A Tree is a search tree.
	Tree struct {
		root *node
	}

	// A Result is a search result from tree.
	Result struct {
		Item   interface{}
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
func (t *Tree) Add(route string, item interface{}) error {
	if len(route) == 0 || route[0] != slash {
		return errNotFromRoot
	}

	if item == nil {
		return errEmptyItem
	}

	lastAsteriskIndex := strings.LastIndex(route, string(asterisk))
	posLastSlashIndex := strings.LastIndex(route, string(slash))
	if lastAsteriskIndex >= 0 && posLastSlashIndex >= 0 && lastAsteriskIndex < posLastSlashIndex {
		return errSlashAfterAsterisk
	}

	err := add(t.root, route[1:], item)
	switch err {
	case errDupItem:
		return duplicatedItem(route)
	case errDupSlash:
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
		return n.forEach(func(k string, v *node) bool {
			r := match(k, token, route)
			pos := i + 1
			if k[0] == asterisk {
				pos = len(route)
			}
			if !r.found || !t.next(v, route[pos:], result) {
				return false
			}

			if r.named {
				addParam(result, r.key, r.value)
			}

			return true
		})
	}

	return n.forEach(func(k string, v *node) bool {
		if r := match(k, route, route); r.found && v.item != nil {
			result.Item = v.item
			if r.named {
				addParam(result, r.key, r.value)
			}

			return true
		}

		return false
	})
}

func (nd *node) forEach(fn func(string, *node) bool) bool {
	for _, children := range nd.children {
		for k, v := range children {
			if fn(k, v) {
				return true
			}
		}
	}

	return false
}

func (nd *node) existChildren() bool {

	if len(nd.children[0]) > 0 || len(nd.children[1]) > 0 || len(nd.children[2]) > 0 {
		return true
	}

	return false
}

func (nd *node) existCatchAll() bool {
	return len(nd.children[2]) > 0
}

func (nd *node) getChildren(route string) map[string]*node {
	if len(route) > 0 {

		if route[0] == colon {
			return nd.children[1]
		}
		if route[0] == asterisk {
			return nd.children[2]
		}
	}

	return nd.children[0]
}

func add(nd *node, route string, item interface{}) error {
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
		if nd.existCatchAll() {
			return errConflictItem
		}

		token := route[:i]
		children := nd.getChildren(token)
		if child, ok := children[token]; ok {
			if child != nil {
				return add(child, route[i+1:], item)
			}

			return errInvalidState
		}

		child := newNode(nil)
		children[token] = child
		return add(child, route[i+1:], item)
	}

	if route[0] == asterisk && nd.existChildren() {
		return errConflictItem
	}
	if nd.existCatchAll() {
		return errConflictItem
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

func match(pat, token, route string) innerResult {
	if pat[0] == colon {
		return innerResult{
			key:   pat[1:],
			value: token,
			named: true,
			found: true,
		}
	}

	if pat[0] == asterisk {
		return innerResult{
			key:   pat[1:],
			value: route,
			named: true,
			found: true,
		}
	}

	return innerResult{
		found: pat == token,
	}
}

func newNode(item interface{}) *node {
	return &node{
		item: item,
		children: [3]map[string]*node{
			make(map[string]*node),
			make(map[string]*node),
			make(map[string]*node),
		},
	}
}
