package stringx

type node struct {
	children map[rune]*node
	fail     *node
	depth    int
	end      bool
	preEnd   *preEnd // current node path, a end word node exist
}

type preEnd struct {
	pre *node
	gap int
}

func (n *node) add(word string) {
	chars := []rune(word)
	if len(chars) == 0 {
		return
	}

	nd := n
	for i, char := range chars {
		if nd.children == nil {
			child := new(node)
			child.depth = i + 1
			nd.children = map[rune]*node{char: child}
			nd = child
		} else if child, ok := nd.children[char]; ok {
			nd = child
		} else {
			child := new(node)
			child.depth = i + 1
			nd.children[char] = child
			nd = child
		}
	}

	nd.end = true
}

func (n *node) linkFail() {
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

func (n *node) build() {
	// bfs
	n.linkFail()
	// dfs
	n.linkPreEnd(nil)
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

func (n *node) linkPreEnd(preGap *preEnd) {
	var curGap preEnd
	if preGap != nil {
		curGap = *preGap
		curGap.gap++
		n.preEnd = &curGap
	}

	if n.end {
		preGap = &preEnd{pre: n, gap: 0}
	} else if preGap != nil {
		newGap := *preGap
		newGap.gap++
		preGap = &newGap
	}
	n.preEnd = preGap

	for _, node := range n.children {
		node.linkPreEnd(preGap)
	}
}

func (n *node) longestMatch(chars []rune, start int) (uselessLen, matchLen int, jump *node) {
	cur := n
	var longestMatched *node
	findLongestMatch := func(n *node) *preEnd {
		var match *preEnd
		icur := n
		for icur.fail != nil {
			icur = icur.fail
			if icur.preEnd != nil {
				match = icur.preEnd
				break
			}
		}
		return match
	}
	for i := start; i < len(chars); i++ {
		char := chars[i]
		child, ok := cur.children[char]
		if ok {
			cur = child
			if cur.end {
				longestMatched = cur
			}
		} else {
			if longestMatched != nil {
				return 0, longestMatched.depth, nil
			}
			if n.end {
				return 0, start, nil
			}
			// old path pre longest match
			longestMatch := findLongestMatch(cur)
			if longestMatch != nil {
				return i - longestMatch.gap - longestMatch.pre.depth, longestMatch.pre.depth, nil
			}
			// new path match
			var jump *node
			icur := cur
			for icur.fail != nil {
				jump, ok = icur.fail.children[char]
				if ok {
					break
				}
				icur = icur.fail
			}
			if jump != nil {
				return i + 1 - jump.depth, 0, jump
			}
			return i + 1, 0, nil
		}
	}
	// this longest matched node
	if longestMatched != nil {
		return 0, longestMatched.depth, nil
	}
	// jumped node matched node
	if n.end {
		return 0, start, nil
	}
	// old path pre longest match
	longestMatch := findLongestMatch(cur)
	if longestMatch != nil {
		return len(chars) - longestMatch.gap - longestMatch.pre.depth, longestMatch.pre.depth, nil
	}
	return len(chars), 0, nil
}
