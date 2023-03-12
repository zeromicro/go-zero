//go:build debug

package search

import "fmt"

func (t *Tree) Print() {
	if t.root.item == nil {
		fmt.Println("/")
	} else {
		fmt.Printf("/:%#v\n", t.root.item)
	}
	printNode(t.root, 1)
}

func printNode(n *node, depth int) {
	indent := make([]byte, depth)
	for i := 0; i < len(indent); i++ {
		indent[i] = '\t'
	}

	for _, children := range n.children {
		for k, v := range children {
			if v.item == nil {
				fmt.Printf("%s%s\n", string(indent), k)
			} else {
				fmt.Printf("%s%s:%#v\n", string(indent), k, v.item)
			}
			printNode(v, depth+1)
		}
	}
}
