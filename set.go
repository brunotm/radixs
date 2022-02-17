package radixs

import (
	"sort"
	"strings"
)

// Set or update the value for the given key
func (t *Tree) set(key string, value interface{}, params bool) (err error) {
	if key == "" {
		return ErrEmptyKey
	}

	if value == nil {
		return ErrNilValue
	}

	if params {
		// scan key and check for invalid constructs with delimiters and parameters
		for x := 0; x < len(key); x++ {
			// delim followed by delim
			if key[x] == t.delimiter && key[x+1] == t.delimiter {
				return ErrInvalidKey
			}

			// param followed by delim or delim
			if key[x] == t.parameter && (key[x+1] == t.delimiter || key[x+1] == t.parameter) {
				return ErrInvalidKey
			}
		}
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
			return nil
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

				parentK := n.key[:pi]
				childK1 := n.key[pi:]
				childK2 := key[pi:]

				// if working with parameters and the current node key at split is not
				// an exact prefix of the current key segment we have an invalid key
				if params {
					nIdx := strings.IndexByte(childK1, t.parameter)
					sIdx := strings.IndexByte(childK2, t.parameter)

					// check if we have a conflicting parameter in either the search key or node key
					// by verifying that n.key[pi:] is a prefix of key[pi:] if:
					// we have a parameter placeholder in the search or node key after the common prefix at index 0 or 1.
					// we have parameter placeholder at the last index of the common prefix
					// TODO: this is plain ugly, fix it
					if (nIdx == 0 || nIdx == 1) || (sIdx == 0 || sIdx == 1) ||
						(key[pi-1] == t.parameter || n.key[pi-1] == t.parameter) {

						if !strings.HasPrefix(childK2, childK1) {
							return ErrConflictKey
						}

					}
				}

				pnode = &node{
					key:      childK1,
					value:    n.value,
					parent:   n,
					children: n.children,
				}

				n.children = []*node{
					pnode,
					{
						key:    childK2,
						value:  value,
						parent: n,
					}}

				n.key = parentK
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
			return nil
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
			return nil
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
		return nil
	}
}
