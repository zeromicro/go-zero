package ast

import (
	"fmt"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/placeholder"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

const (
	_ CollectionFlag = iota << 1
	NotIn
	LeftIn
	RightIn
	AllIn
)

var initNode = NewTokenNode(token.Token{
	Position: token.Position{
		Line: 1,
	},
})

type CollectionFlag int

type NodeSet struct {
	m    map[Node]int
	list []Node
}

func NewNodeSet() *NodeSet {
	return &NodeSet{
		m: map[Node]int{
			initNode: 0,
		},
		list: []Node{initNode},
	}
}

func (s *NodeSet) Append(node Node) {
	if _, ok := s.m[node]; ok {
		return
	}
	s.list = append(s.list, node)
	s.m[node] = len(s.list) - 1
}

func (s *NodeSet) InsertAfter(node, after Node) {
	if _, ok := s.m[node]; ok {
		return
	}
	idx, ok := s.m[after]
	if !ok {
		panic(fmt.Sprintf("node <%T> not exists", after))
	}

	s.m[node] = idx
	pos := idx + 1
	if pos > len(s.list)-1 {
		s.list = append(s.list, node)
		return
	}

	var list []Node
	list = append(list, s.list[:pos]...)
	list = append(list, node)
	remainList := s.list[pos:]
	list = append(list, remainList...)
	for _, e := range remainList {
		s.m[e] += 1
	}
	s.list = list
}

func (s *NodeSet) between(left, right Node, flag CollectionFlag, onlyComment bool) []Node {
	if len(s.list) == 0 {
		return nil
	}

	leftIdx, ok := s.m[left]
	if !ok {
		return nil
	}
	rightIdx, ok := s.m[right]
	if !ok {
		return nil
	}

	if leftIdx > rightIdx {
		return nil
	}

	var results []Node
	var bg, end int
	switch flag {
	case NotIn:
		bg = leftIdx + 1
		end = rightIdx
	case LeftIn:
		bg = leftIdx
		end = rightIdx
	case RightIn:
		bg = leftIdx + 1
		end = rightIdx + 1
	case AllIn:
		bg = leftIdx
		end = rightIdx + 1
	}
	if bg > len(s.list) {
		return nil
	}
	if end > len(s.list) {
		end = len(s.list)
	}

	if bg > end {
		return results
	}

	for i := bg; i < end; i++ {
		if onlyComment {
			if _, ok := s.list[i].(*CommentStmt); ok {
				results = append(results, s.list[i])
			}
		} else {
			results = append(results, s.list[i])
		}
	}

	return results
}

func (s *NodeSet) Between(left, right Node, flag CollectionFlag) []Node {
	return s.between(left, right, flag, false)
}

func (s *NodeSet) CommentBetween(left, right Node, flag CollectionFlag) []Node {
	return s.between(left, right, flag, true)
}

func (s *NodeSet) LineCommentAfter(node Node,skip map[Node]placeholder.Type) []Node {
	if len(s.list) == 0 {
		return nil
	}

	var bgIdx = s.m[node] + 1
	if bgIdx > len(s.list)-1 {
		return nil
	}

	var results []Node
	for i := bgIdx; i < len(s.list); i++ {
		if _,ok:=skip[s.list[i]];ok{
			continue
		}
		line := s.list[i].Pos().Line
		_, ok := s.list[i].(*CommentStmt)
		if ok {
			if line == node.End().Line {
				results = append(results, s.list[i])
			}
		}else{
			break
		}
	}
	return results
}

func (s *NodeSet) FirstToken() Node {
	if len(s.list) == 0 {
		return nil
	}
	return s.list[0]
}

func (s *NodeSet) LastToken() Node {
	if len(s.list) == 0 {
		return nil
	}
	return s.list[len(s.list)-1]
}
