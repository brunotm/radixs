package radixs

// Tree is a compact radix (compact prefix) tree which is guaranteed
// to be sorted. Key/Value pairs are always inserted, retrieved and updated
// using binary searches making the tree operations very efficient
// for large trees.
type Tree struct {
	size      uint64
	root      *node
	delimiter byte
	parameter byte
}

// New creates a new radix tree
func New(opts ...OptFunc) (t *Tree) {
	t = &Tree{
		root: &node{},
	}

	for x := 0; x < len(opts); x++ {
		opts[x](t)
	}

	return t
}

// OptFunc functional options for tree creation
type OptFunc func(t *Tree)

// WithParams sets the tree key delimiters and parameter placeholder
// when working with path parameter in keys
func WithParams(delimiter, parameter byte) (opt OptFunc) {
	return func(t *Tree) {
		t.delimiter = delimiter
		t.parameter = parameter
	}
}

func FromMap(m map[string]interface{}, opts ...OptFunc) (t *Tree) {
	t = New(opts...)

	for k, v := range m {
		t.Set(k, v)
	}

	return t
}

// Set or update the value for the given key
func (t *Tree) Set(key string, value interface{}) (ok bool) {
	return t.set(key, value, false)
}

// SetWithParams is like Set, but provides additional validation
// to prevent invalid keys and conflicts when work with key parameters
func (t *Tree) SetWithParams(key string, value interface{}) (ok bool) {
	return t.set(key, value, true)
}

// Delete removes the provided key from the tree.
// It returns false if the key was not found.
func (t *Tree) Delete(key string) (ok bool) {
	return t.delete(key, false)
}

// DeletePrefix deletes all keys under the given prefix
func (t *Tree) DeletePrefix(key string) (ok bool) {
	return t.delete(key, true)
}

// Iter calls f sequentially for each key and value present in the tree.
// If f returns false it stops the iteration.
// Iter is guaranteed to iterate the tree in ascending lexicographic order
func (t *Tree) Iter(f func(key string, value interface{}) bool) {
	t.root.iter("", f)
}

// Size returns the number of leaf nodes in the tree
func (t *Tree) Size() (sz uint64) {
	return t.size
}

// String returns a string representation of the tree
func (t *Tree) String() (s string) {
	b := &stringBuilder{}
	b.WriteString("D, W\n")
	t.root.string(b, 0)
	return b.String()
}

func longestPrefix(s1, s2 string) int {
	max := len(s1)
	if l := len(s2); l < max {
		max = l
	}

	var i int
	for i = 0; i < max; i++ {
		if s1[i] != s2[i] {
			break
		}
	}

	return i
}

type stringBuilder struct {
	buf []byte
}

func (b *stringBuilder) WriteString(s string) {
	b.buf = append(b.buf, s...)
}

func (b *stringBuilder) String() string {
	return string(b.buf)
}
