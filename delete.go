package radixs

import (
	"sort"
)

// delete removes both keys and prefixes if specified
func (t *Tree) delete(key string, prefix bool) (err error) {
	if key == "" {
		return ErrEmptyKey
	}

	n := t.root
	for {
		// find and remove the longest common prefix for
		// the current key segment and node key
		pi := longestPrefix(n.key, key)
		key = key[pi:]

		// do a binary search for the key prefix in the current node children
		i := sort.Search(len(n.children), func(x int) bool {
			if key == "" {
				return true
			}

			return n.children[x].key[0] >= key[0]
		})

		// no results for prefix found in search
		if i >= len(n.children) {
			return ErrKeyNotFound
		}

		// delete prefix remaining key segment is a prefix of next node
		if prefix && len(key) == longestPrefix(key, n.children[i].key) {
			var subSize int
			n.children[i].dfsI(func(n *node) bool {
				if n.value != nil {
					subSize++
				}
				return true
			})

			n.children = append(n.children[:i], n.children[i+1:]...)

			t.size -= uint64(subSize)
			return nil
		}

		// key found
		if n.children[i].key == key {
			// checks if the current node is also
			// a prefix for underlying child nodes and:
			// - set the node value to nil if node is a prefix
			// - remove the node if it has no underlying children
			switch len(n.children[i].children) > 0 {
			case true:
				n.children[i].value = nil
			case false:
				n.children = append(n.children[:i], n.children[i+1:]...)
			}

			// if the node has single child left, merge with its parent
			if len(n.children) == 1 && n != t.root {
				n.key += n.children[0].key
				n.value = n.children[0].value
				n.children = n.children[0].children

				// update all children parent node
				for x := 0; x < len(n.children); x++ {
					n.children[x].parent = n
				}
			}

			t.size--
			return nil
		}

		// if child at index i shares a common prefix with the current
		// key segment descend into it
		n = n.children[i]
	}
}
