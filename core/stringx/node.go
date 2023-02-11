package stringx

type node struct {
	children map[rune]*node
	fail     *node
	depth    int
	end      bool
}

func (n *node) add(word string) {
	chars := []rune(word)
	if len(chars) == 0 {
		return
	}

	nd := n
	var depth int
	for i, char := range chars {
		if nd.children == nil {
			child := new(node)
			child.depth = i + 1
			nd.children = map[rune]*node{char: child}
			nd = child
		} else if child, ok := nd.children[char]; ok {
			nd = child
			depth++
		} else {
			child := new(node)
			child.depth = i + 1
			nd.children[char] = child
			nd = child
		}
	}

	nd.end = true
}

func (n *node) build() {
	var nodes []*node
	for _, child := range n.children {
		child.fail = n
		nodes = append(nodes, child)
	}
	for len(nodes) > 0 {
		nd := nodes[0]
		nodes = nodes[1:]
		for key, child := range nd.children {
			nodes = append(nodes, child)
			cur := nd
			for cur != nil {
				if cur.fail == nil {
					child.fail = n
					break
				}
				if fail, ok := cur.fail.children[key]; ok {
					child.fail = fail
					break
				}
				cur = cur.fail
			}
		}
	}
}

func (n *node) find(chars []rune) []scope {
	var scopes []scope
	size := len(chars)
	cur := n

	for i := 0; i < size; i++ {
		child, ok := cur.children[chars[i]]
		if ok {
			cur = child
		} else {
			for cur != n {
				cur = cur.fail
				if child, ok = cur.children[chars[i]]; ok {
					cur = child
					break
				}
			}

			if child == nil {
				continue
			}
		}

		for child != n {
			if child.end {
				scopes = append(scopes, scope{
					start: i + 1 - child.depth,
					stop:  i + 1,
				})
			}
			child = child.fail
		}
	}

	return scopes
}

func (n *node) longestMatch(chars []rune, start int) (used int, jump *node, matched bool) {
	cur := n
	var matchedNode *node

	for i := start; i < len(chars); i++ {
		child, ok := cur.children[chars[i]]
		if ok {
			cur = child
			if cur.end {
				matchedNode = cur
			}
		} else {
			if matchedNode != nil {
				return matchedNode.depth, nil, true
			}

			if n.end {
				return start, nil, true
			}

			var jump *node
			for cur.fail != nil {
				jump, ok = cur.fail.children[chars[i]]
				if ok {
					break
				}
				cur = cur.fail
			}
			if jump != nil {
				return i + 1 - jump.depth, jump, false
			}

			return i + 1, nil, false
		}
	}

	// longest matched node
	if matchedNode != nil {
		return matchedNode.depth, nil, true
	}

	// last matched node
	if n.end {
		return start, nil, true
	}

	return len(chars), nil, false
}
