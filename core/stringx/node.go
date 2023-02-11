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

func (n *node) longestMatch(chars []rune, paths []*node) (uselessLen, matchLen int, nextPaths []*node) {
	cur := n
	var longestMatched *node
	findMatch := func(path []*node) (*node, int) {
		var (
			result *node
			start  int
		)
		for i := len(path) - 1; i >= 0; i-- {
			icur := path[i]
			var cur *node
			for icur.fail != nil {
				if icur.fail.end {
					cur = icur.fail
					break
				}
				icur = icur.fail
			}
			if cur != nil {
				if result == nil {
					result = cur
					start = i - result.depth + 1
				} else {
					if curStart := i - cur.depth + 1; curStart < start {
						result = cur
						start = curStart
					} else if curStart == start && cur.depth > result.depth {
						result = cur
						start = curStart
					}
				}
			}
		}
		return result, start
	}

	for i := len(paths); i < len(chars); i++ {
		char := chars[i]
		child, ok := cur.children[char]
		if ok {
			cur = child
			if cur.end {
				longestMatched = cur
			}
			paths = append(paths, cur)
		} else {
			if longestMatched != nil {
				return 0, longestMatched.depth, nil
			}
			if n.end {
				return 0, n.depth, nil
			}
			// old path pre longest preMatch
			preMatch, preStart := findMatch(paths)
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
			switch {
			case preMatch != nil && jump != nil:
				if jumpStart := i - jump.depth + 1; jumpStart < preStart {
					return jumpStart, 0, append(paths[jumpStart:], jump)
				} else if jumpStart == preStart {
					if jump.depth > preMatch.depth {
						return jumpStart, 0, append(paths[jumpStart:], jump)
					}
					return preStart, preMatch.depth, nil
				}
			case preMatch != nil && jump == nil:
				return preStart, preMatch.depth, nil
			case preMatch == nil && jump != nil:
				return i - jump.depth + 1, 0, append(paths[i-jump.depth+1:], jump)
			case preMatch == nil && jump == nil:
				return i + 1, 0, nil
			}
		}
	}
	// this longest matched node
	if longestMatched != nil {
		return 0, longestMatched.depth, nil
	}
	if n.end {
		return 0, n.depth, nil
	}
	match, start := findMatch(paths)
	if match != nil {
		return start, match.depth, nil
	}
	return len(chars), 0, nil
}
