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
	n.fail = n
	for _, child := range n.children {
		child.fail = n
		n.buildNode(child)
	}
}

func (n *node) buildNode(nd *node) {
	if nd.children == nil {
		return
	}

	var fifo []*node
	for key, child := range nd.children {
		fifo = append(fifo, child)

		if fail, ok := nd.fail.children[key]; ok {
			child.fail = fail
		} else {
			child.fail = n
		}
	}

	for _, val := range fifo {
		n.buildNode(val)
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
		} else if cur == n {
			continue
		} else {
			cur = cur.fail
			if child, ok = cur.children[chars[i]]; !ok {
				continue
			}
			cur = child
		}

		if child.end {
			scopes = append(scopes, scope{
				start: i + 1 - child.depth,
				stop:  i + 1,
			})
		}
		if child.fail != n && child.fail.end {
			scopes = append(scopes, scope{
				start: i + 1 - child.fail.depth,
				stop:  i + 1,
			})
		}
	}

	return scopes
}
