package radixs

import (
	"fmt"
)

type node struct {
	children []*node     // []*radixscratch.node: 0-24 (size 24, align 8)
	key      string      // string: 24-40 (size 16, align 8)
	value    interface{} // interface{}: 40-56 (size 16, align 8)
	parent   *node       // *radixscratch.node: 72-80 (size 8, align 8)
}

func (n *node) iter(prefix string, f func(key string, value interface{}) bool) (ok bool) {
	if n.value != nil {
		if !f(prefix+n.key, n.value) {
			return false
		}
	}

	prefix += n.key
	for x := 0; x < len(n.children); x++ {
		if !n.children[x].iter(prefix, f) {
			return false
		}
	}

	return true
}

// dfsI is like dfs but includes the current node
func (n *node) dfsI(f func(*node) bool) {
	if !f(n) {
		return
	}
	n.dfs(f)
}

// dfs walks all the subtree under the given node calling f for each node.
// If f returns false, dfs stops the iteration.
func (n *node) dfs(f func(*node) bool) {
	for x := 0; x < len(n.children); x++ {
		if !f(n.children[x]) {
			return
		}
		n.children[x].dfs(f)
	}
}

// reverseWalk walks the tree from the current node up to the tree root,
// calling f for each node. If f returns false, reverseWalk stops the iteration.
func (n *node) reverseWalk(f func(*node) bool) {
	p := n.parent

	for p != nil {
		if !f(p) {
			return
		}

		p = p.parent
	}
}

func bfs(n *node, f func(*node) bool) {
	stack := make([]*node, 1, len(n.children)*2+1)
	stack[0] = n

	for len(stack) > 0 {
		p := stack[0]
		if !f(p) {
			return
		}

		stack = stack[1:]
		if len(p.children) > 0 {
			stack = append(stack, p.children...)
		}
	}
}

func (n *node) weight() (p uint64) {
	// exclude tree root
	if n.key != "" {
		p++
	}

	n.dfs(func(n *node) bool {
		p++
		return true
	})

	return p
}

func (n *node) depth() (p uint64) {
	n.reverseWalk(func(n *node) bool {
		p++
		return true
	})

	return p
}

func (n *node) string(b *stringBuilder, tab int) {
	b.WriteString(fmt.Sprintf("%d, %d    ", tab, n.weight()))

	for x := 0; x < tab; x++ {
		b.WriteString("    ")
	}

	if n.key == "" {
		b.WriteString("root\n")
	} else {
		b.WriteString(fmt.Sprintf("key: %s -> %#v\n", n.key, n.value))
	}

	tab++
	for x := 0; x < len(n.children); x++ {
		n.children[x].string(b, tab)
	}
}
