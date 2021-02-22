package search

import "errors"

const (
	colon = ':'
	slash = '/'
)

var (
	// ErrDupItem means adding duplicated item.
	ErrDupItem = errors.New("duplicated item")
	// ErrDupSlash means item is started with more than one slash.
	ErrDupSlash = errors.New("duplicated slash")
	// ErrEmptyItem means adding empty item.
	ErrEmptyItem = errors.New("empty item")
	// ErrInvalidState means search tree is in an invalid state.
	ErrInvalidState = errors.New("search tree is in an invalid state")
	// ErrNotFromRoot means path is not starting with slash.
	ErrNotFromRoot = errors.New("path should start with /")

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
		children [2]map[string]*node
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
		return ErrNotFromRoot
	}

	if item == nil {
		return ErrEmptyItem
	}

	return add(t.root, route[1:], item)
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
		if route[i] == slash {
			token := route[:i]
			return n.forEach(func(k string, v *node) bool {
				if r := match(k, token); r.found {
					if t.next(v, route[i+1:], result) {
						if r.named {
							addParam(result, r.key, r.value)
						}

						return true
					}
				}

				return false
			})
		}
	}

	return n.forEach(func(k string, v *node) bool {
		if r := match(k, route); r.found && v.item != nil {
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

func (nd *node) getChildren(route string) map[string]*node {
	if len(route) > 0 && route[0] == colon {
		return nd.children[1]
	}

	return nd.children[0]
}

func add(nd *node, route string, item interface{}) error {
	if len(route) == 0 {
		if nd.item != nil {
			return ErrDupItem
		}

		nd.item = item
		return nil
	}

	if route[0] == slash {
		return ErrDupSlash
	}

	for i := range route {
		if route[i] == slash {
			token := route[:i]
			children := nd.getChildren(token)
			if child, ok := children[token]; ok {
				if child != nil {
					return add(child, route[i+1:], item)
				}

				return ErrInvalidState
			}

			child := newNode(nil)
			children[token] = child
			return add(child, route[i+1:], item)
		}
	}

	children := nd.getChildren(route)
	if child, ok := children[route]; ok {
		if child.item != nil {
			return ErrDupItem
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

func match(pat, token string) innerResult {
	if pat[0] == colon {
		return innerResult{
			key:   pat[1:],
			value: token,
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
		children: [2]map[string]*node{
			make(map[string]*node),
			make(map[string]*node),
		},
	}
}
