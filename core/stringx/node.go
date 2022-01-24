package stringx

import "fmt"

var idx = 1

type node struct {
	children map[rune]*node
	fail     *node
	depth    int
	end      bool
	word     string
	id       int
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
			child.id = idx
			idx++
			nd.children = map[rune]*node{char: child}
			nd = child
		} else if child, ok := nd.children[char]; ok {
			nd = child
			depth++
		} else {
			child := new(node)
			child.depth = i + 1
			child.id = idx
			idx++
			nd.children[char] = child
			nd = child
		}
	}

	nd.word = word
	nd.end = true
}

func (n *node) build() {
	n.fail = n
	n.id = 0
	idx++
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
			for {
				if fail, ok := cur.fail.children[key]; ok {
					child.fail = fail
					break
				}
				if cur == n {
					break
				}
				cur = cur.fail
			}
			if child.fail == nil {
				child.fail = n
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

func (n *node) print() {
	fmt.Println("/")
	printNode(n)
}

func printNode(nd *node) {
	for k, v := range nd.children {
		for i := 0; i < v.depth; i++ {
			fmt.Print(" ")
		}
		fmt.Printf("%c,%d,%d", k, v.id, v.fail.id)
		if v.end {
			fmt.Printf(",%t\n", v.end)
		} else {
			fmt.Println()
		}
		printNode(v)
	}
}
