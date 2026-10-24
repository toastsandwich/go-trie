package trie

import (
	"cmp"
	"sort"
)

type Node struct {
	Value    rune
	IsEnd    bool
	len      int
	Children []*Node
}

func createNode(r rune) *Node {
	return &Node{
		Value:    r,
		Children: make([]*Node, 0),
	}
}

func (n *Node) addChild(p *Node) {
	n.len++
	n.Children = append(n.Children, p)
	sort.Slice(n.Children, func(i, j int) bool {
		return cmp.Less(n.Children[i].Value, n.Children[j].Value)
	})
}

func (n *Node) findChild(k rune) (int, bool) {
	l := 0
	r := n.len - 1
	for l <= r {
		m := (r-l)/2 + l
		switch cmp.Compare(n.Children[m].Value, k) {
		case 0:
			return m, true
		case -1:
			l = m + 1
		case 1:
			r = m - 1
		}
	}
	return -1, false
}

func (n *Node) getChild(idx int) *Node {
	if n.len <= idx {
		return nil
	}
	return n.Children[idx]
}

type Trie struct {
	Root *Node
}

func New() *Trie {
	return &Trie{
		Root: createNode(-1),
	}
}

func (t *Trie) Reset() {
	t.Root = createNode(-1)
}

func (t *Trie) Insert(s string) {
	curr := t.Root
	runes := []rune(s)
	rlen := len(runes)
	for i, r := range runes {
		idx, ok := curr.findChild(r)
		if !ok {
			newNode := createNode(r)
			newNode.IsEnd = ((rlen - 1) == i)

			curr.addChild(newNode)
			curr = newNode
			continue
		}
		if f := curr.getChild(idx); f != nil {
			if rlen-1 == i {
				f.IsEnd = true
				break
			}
			curr = f
		}
	}
}

func (t *Trie) Search(s string) bool {
	curr := t.Root
	for _, r := range s {
		idx, ok := curr.findChild(r)
		if !ok {
			return false
		}
		if f := curr.getChild(idx); f == nil {
			return false
		} else {
			curr = f
		}
	}
	return true
}
