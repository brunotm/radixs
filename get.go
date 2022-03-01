package radixs

import (
	"sort"
	"strings"
)

// GetWithParams is like Get but extracts path parameters and stores them
// into the provided params argument which will be accessible after GetWithParams returns.
func (t *Tree) GetWithParams(key string, params map[string]string) (value interface{}, err error) {
	if key == "" {
		return nil, ErrEmptyKey
	}

	n := t.root
	for {
		// obtain the longest common prefix for the current
		// search key and node key
		pi := longestPrefix(n.key, key)
		key = key[pi:]

		// binary search for prefix
		i := sort.Search(len(n.children), func(x int) bool {
			pIdx := longestPrefix(n.children[x].key, key)
			if pIdx > -1 {
				if len(n.children[x].key) > pIdx && n.children[x].key[pIdx:][0] == t.parameter {
					return true
				}
			}

			return n.children[x].key[0] >= key[0]
		})

		// end of search no node with prefix found
		if i >= len(n.children) {
			return nil, ErrKeyNotFound
		}

		pIdx := strings.IndexByte(n.children[i].key, t.parameter)
		if pIdx > 0 {
			key = key[pIdx:]
			paramName := n.children[i].key[pIdx+1:]

			dIdx := strings.IndexByte(key, t.delimiter)
			if dIdx >= 0 {
				params[paramName] = key[:dIdx]
				key = key[dIdx:]
			} else {
				params[paramName] = key
			}

			if dIdx < 0 {
				return n.children[i].value, nil
			}

			key = n.children[i].key + key
			n = n.children[i]
			continue

		}

		// exact match found
		if n.children[i].key == key {
			if n.children[i].value == nil {
				return nil, ErrKeyNotFound
			}

			return n.children[i].value, nil
		}

		// child is a prefix of the search key, continue
		n = n.children[i]
	}
}

// Get retrieves the value for the given key.
// It returns false if the key was not found.
func (t *Tree) Get(key string) (value interface{}, err error) {
	if key == "" {
		return nil, ErrEmptyKey
	}

	n := t.root
	for {
		// obtain the longest common prefix for the current
		// search key and node key
		pi := longestPrefix(n.key, key)
		key = key[pi:]

		if key == "" {
			return nil, ErrKeyNotFound
		}

		// binary search for prefix
		i := sort.Search(len(n.children), func(x int) bool {
			return n.children[x].key[0] >= key[0]
		})

		// end of search no node with prefix found
		if i >= len(n.children) {
			return nil, ErrKeyNotFound
		}

		// exact match found
		if n.children[i].key == key {
			if n.children[i].value == nil {
				return nil, ErrKeyNotFound
			}

			return n.children[i].value, nil
		}

		// child is a prefix of the search key, continue
		n = n.children[i]
	}
}

// LongestMatch is like Get, but instead of an
// exact match, it will return the longest prefix match.
func (t *Tree) LongestMatch(key string) (match string, value interface{}, err error) {
	match, n, err := t.longestMatch(key)
	if err != nil {
		return "", nil, err
	}

	return match, n.value, nil
}

// NeighborMatch is like LongestMatch, but returns the longest match and surrounding keys:
// parent, match, siblings, children and stores them into the provided matches map.
func (t *Tree) NeighborMatch(key string, matches map[string]interface{}) (err error) {
	match, n, err := t.longestMatch(key)
	if err != nil {
		return err
	}

	// add current node
	matches[match] = n.value

	// add current node children
	for x := 0; x < len(n.children); x++ {
		if n.children[x].value != nil {
			matches[match+n.children[x].key] = n.children[x].value
		}
	}

	// add current node parent
	pKey := match[:len(n.parent.key)]
	if n.parent.key != "" && n.parent.value != nil {
		matches[pKey] = n.parent.value
	}

	// add current node siblings
	if len(n.parent.children) > 1 {
		for x := 0; x < len(n.parent.children); x++ {
			if n.parent.children[x].key != n.key && n.parent.children[x].value != nil {
				matches[pKey+n.parent.children[x].key] = n.parent.children[x].value
			}
		}
	}

	return nil
}

func (t *Tree) longestMatch(key string) (match string, v *node, err error) {
	if key == "" {
		return "", nil, ErrEmptyKey
	}

	n := t.root
	for {
		// obtain the longest common prefix for the current search key and node key
		// keep track of the accumulated prefix for returning the longest match
		pi := longestPrefix(n.key, key)
		match += key[:pi]
		key = key[pi:]

		var i int
		if key != "" {
			// binary search for prefix
			i = sort.Search(len(n.children), func(x int) bool {
				return n.children[x].key[0] >= key[0]
			})
		} else {
			i = len(n.children) + 1
		}

		// end of search, reverse walk the tree until the longest match
		if i >= len(n.children) {
			for {
				n = n.parent
				match = match[:len(match)-pi]

				if n.value != nil {
					return match, n, nil
				}

				if n.parent == nil {
					return "", nil, ErrKeyNotFound
				}
			}
		}

		// exact match found
		if n.children[i].key == key {
			if n.children[i].value == nil {
				return "", nil, ErrKeyNotFound
			}
			return match + n.children[i].key, n.children[i], nil
		}

		// child is a prefix of the search key, continue
		n = n.children[i]
	}
}
