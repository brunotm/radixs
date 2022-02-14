package radixs

import (
	"sort"
)

// Get retrieves the value for the given key.
// It returns false if the key was not found.
func (t *Tree) Get(key string) (value interface{}, ok bool) {
	if len(key) == 0 {
		return nil, false
	}

	n := t.root
	for {
		// obtain the longest common prefix for the current
		// search key and node key
		pi := longestPrefix(n.key, key)
		key = key[pi:]

		// binary search for prefix
		i := sort.Search(len(n.children), func(x int) bool {
			return n.children[x].key[0] >= key[0]
		})

		// end of search no node with prefix found
		if i >= len(n.children) {
			return nil, false
		}

		// exact match found
		if n.children[i].key == key {
			return n.children[i].value, n.children[i].value != nil
		}

		// child is a prefix of the search key, continue
		n = n.children[i]
	}
}

// LongestMatch is like Get, but instead of an
// exact match, it will return the longest prefix match.
func (t *Tree) LongestMatch(key string) (match string, value interface{}, ok bool) {
	n := t.root
	for {
		// obtain the longest common prefix for the current search key and node key
		// keep track of the accumulated prefix for returning the longest match
		pi := longestPrefix(n.key, key)
		match += key[:pi]
		key = key[pi:]

		// binary search for prefix
		i := sort.Search(len(n.children), func(x int) bool {
			return n.children[x].key[0] >= key[0]
		})

		// end of search, reverse walk the tree until the longest match
		if i >= len(n.children) {
			for {
				n = n.parent
				match = match[:len(match)-pi]

				if n.value != nil {
					return match, n.value, true
				}

				if n.parent == nil {
					return "", nil, false
				}
			}
		}

		// exact match found
		if n.children[i].key == key {
			return match + n.children[i].key, n.children[i].value, n.children[i].value != nil
		}

		// child is a prefix of the search key, continue
		n = n.children[i]
	}
}
