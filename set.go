package radixs

import (
	"sort"
)

// Set or update the value for the given key
func (t *Tree) Set(key string, value interface{}) (ok bool) {
	if len(key) == 0 || value == nil {
		return false
	}

	n := t.root
	for {
		// existing key, update its value
		if n.key == key {
			// updating an existing prefix increase tree size
			if n.value == nil {
				t.size++
			}
			n.value = value
			return true
		}

		// obtain the longest common prefix for the current search key
		pi := longestPrefix(n.key, key)

		// key segment is less than current node key
		// insert and add existing node as child
		if pi > 0 && len(n.key) > pi {
			// pnode for updating children parent after split
			var pnode *node

			// common prefix is full search key segment
			// split and add current node as a child
			if pi == len(key) {
				pnode = &node{
					key:      n.key[pi:],
					value:    n.value,
					parent:   n,
					children: n.children,
				}

				n.children = []*node{pnode}
				n.key = key
				n.value = value
			} else {
				// key segment shares a common prefix with the current node
				// split at the common prefix and add children nodes

				pnode = &node{
					key:      n.key[pi:],
					value:    n.value,
					parent:   n,
					children: n.children,
				}

				n.children = []*node{
					pnode,
					{
						key:    key[pi:],
						value:  value,
						parent: n,
					}}

				n.key = n.key[:pi]
				n.value = nil

				// ensure nodes are sorted
				if n.children[0].key[0] > n.children[1].key[0] {
					n.children[0], n.children[1] = n.children[1], n.children[0]
				}
			}

			// update split node children parent
			for x := 0; x < len(pnode.children); x++ {
				pnode.children[x].parent = pnode
			}

			t.size++
			return true
		}

		key = key[pi:]
		// node key is a prefix of the current key search for insertion index
		i := sort.Search(len(n.children), func(x int) bool {
			return n.children[x].key[0] >= key[0]
		})

		// if insertion index is less than edges size call insertion at
		// the node at index position
		if i < len(n.children) {
			// child at index position shares a prefix with key,
			// continue iteration
			if n.children[i].key[0] == key[0] {
				n = n.children[i]
				continue
			}

			// insert node at index position
			n.children = append(n.children[:i+1], n.children[i:]...)
			n.children[i] = &node{
				key:    key,
				value:  value,
				parent: n,
			}

			t.size++
			return true
		}

		// insertion index is bigger than children size, append to it
		n.children = append(
			n.children,
			&node{
				key:    key,
				value:  value,
				parent: n,
			})

		t.size++
		return true
	}
}
