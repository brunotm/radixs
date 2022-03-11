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
		// binary search for prefix
		i := sort.Search(len(n.children), func(x int) bool {
			if n.children[x].key[0] == t.parameter {
				return true
			}

			return n.children[x].key[0] >= key[0]
		})

		// end of search no node with prefix found
		if i >= len(n.children) {
			return nil, ErrKeyNotFound
		}

		pi := longestPrefix(n.children[i].key, key)
		nodeKey := n.children[i].key[pi:]
		key = key[pi:]

		// parameter found, start consuming until last parameter or end of key/nodeKey
		for len(key) > 0 && len(nodeKey) > 0 && nodeKey[0] == t.parameter {
			name := nodeKey[1:]
			if pdIdx := strings.IndexByte(name, t.delimiter); pdIdx > -1 {
				name = name[:pdIdx]
			}

			value := key
			if vdIdx := strings.IndexByte(value, t.delimiter); vdIdx > -1 {
				value = value[:vdIdx]
			}

			params[name] = value
			key = key[len(value):]
			nodeKey = nodeKey[len(name)+1:] // include the parameter placeholder

			if pi := longestPrefix(nodeKey, key); pi > 0 {
				key = key[pi:]
				nodeKey = nodeKey[pi:]
			}
		}

		if key == "" {
			if nodeKey == "" {
				return n.children[i].value, nil
			}

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

// Get retrieves the value for the given key.
// It returns false if the key was not found.
func (t *Tree) Get(key string) (value interface{}, err error) {
	if key == "" {
		return nil, ErrEmptyKey
	}

	n := t.root
	for {
		switch {
		case key == n.key:
			if n.value == nil {
				return nil, ErrKeyNotFound
			}
			return n.value, nil

		case strings.HasPrefix(key, n.key):
			key = key[len(n.key):]
			i := sort.Search(len(n.children), func(x int) bool {
				return n.children[x].key[0] >= key[0]
			})

			if i >= len(n.children) {
				return nil, ErrKeyNotFound
			}
			n = n.children[i]

		default:
			return nil, ErrKeyNotFound
		}
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

func (t *Tree) longestMatch(key string) (prefix string, n *node, err error) {
	if key == "" {
		return "", nil, ErrEmptyKey
	}

	n = t.root
	for {
		prefix += n.key

		if key == n.key {
			if n.value != nil {
				return prefix, n, nil
			}
			break
		}

		if strings.HasPrefix(key, n.key) {
			key = key[len(n.key):]

			i := sort.Search(len(n.children), func(x int) bool {
				return n.children[x].key[0] >= key[0]
			})

			if i >= len(n.children) {
				break
			}
			n = n.children[i]
			continue
		}

		break
	}

	for n.parent != nil {
		prefix = prefix[:len(prefix)-len(n.key)]
		n = n.parent
		if n.value != nil {
			return prefix, n, nil
		}
	}

	return "", nil, ErrKeyNotFound
}

func (t *Tree) get(key string) (n *node, match bool) {
	n = t.root
	for {
		switch {
		case key == n.key:
			return n, n.value != nil

		case strings.HasPrefix(key, n.key):
			key = key[len(n.key):]

			i := sort.Search(len(n.children), func(x int) bool {
				return n.children[x].key[0] >= key[0]
			})

			if i >= len(n.children) {
				return n, false
			}
			n = n.children[i]

		default:
			return n, false
		}
	}
}
