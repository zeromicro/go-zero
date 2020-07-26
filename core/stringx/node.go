package stringx

type node struct {
	children map[rune]*node
	end      bool
}

func (n *node) add(word string) {
	chars := []rune(word)
	if len(chars) == 0 {
		return
	}

	nd := n
	for _, char := range chars {
		if nd.children == nil {
			child := new(node)
			nd.children = map[rune]*node{
				char: child,
			}
			nd = child
		} else if child, ok := nd.children[char]; ok {
			nd = child
		} else {
			child := new(node)
			nd.children[char] = child
			nd = child
		}
	}

	nd.end = true
}
