# radixs


<a href="https://en.wikipedia.org/wiki/Radix_tree">
<img style="float: right;" src="https://upload.wikimedia.org/wikipedia/commons/a/ae/Patricia_trie.svg"></a>
A Go implementation of a radix tree, that uses binary searches to speeding up insert, retrieve and delete operations specially on dense trees. This implementation additionally eliminate edge specific types and pointers in order to save memory and allocations on dense trees.

___
## Usage

```go
package main

import (
	"fmt"

	"github.com/brunotm/radixs"
)

func main() {
	tr := radixs.New() // radix.FromMap() alternatively build from an existing map
	tr.Set("romane", 0)
	tr.Set("romanus", 1)
	tr.Set("romulus", "remus brother")
	tr.Set("rubens", "51.2170° N, 4.4093° E")
	tr.Set("ruber", 4)
	tr.Set("rubicon", 108)
	tr.Set("rubicundus", func() bool { return true })

	value, ok := tr.Get("rubens")
	fmt.Printf("value for rubens: %#v, ok: %t\n", value, ok)

	size := tr.Size()
	fmt.Printf("tree size: %d\n", size)

	key, value, ok := tr.LongestMatch("romanesco")
	fmt.Printf("longest prefix for romanesco: %s, value: %#v, ok: %t\n",
		key, value, ok)

	fmt.Printf("\ntree string representation:\n%s\n", tr.String())

	ok = tr.Delete("romulus")
	fmt.Printf("romulus deleted: %t\n", ok)

	ok = tr.DeletePrefix("rube")
	fmt.Printf("prefix rube deleted: %t\n", ok)

	fmt.Printf("\ntree string representation:\n%s\n", tr.String())

	fmt.Printf("ordered iteration of the tree key/value pairs:\n")
	tr.Iter(func(key string, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
}


// value for rubens: "51.2170° N, 4.4093° E", ok: true
// tree size: 7
// longest prefix for romanesco: , value: <nil>, ok: false

// tree string representation:
// D, W
// 0, 13    root
// 1, 13        key: r -> <nil>
// 2, 5            key: om -> <nil>
// 3, 3                key: an -> <nil>
// 4, 1                    key: e -> 0
// 4, 1                    key: us -> 1
// 3, 1                key: ulus -> "remus brother"
// 2, 7            key: ub -> <nil>
// 3, 3                key: e -> <nil>
// 4, 1                    key: ns -> "51.2170° N, 4.4093° E"
// 4, 1                    key: r -> 4
// 3, 3                key: ic -> <nil>
// 4, 1                    key: on -> 108
// 4, 1                    key: undus -> (func() bool)(0x108e200)

// romulus deleted: true
// prefix rube deleted: true

// tree string representation:
// D, W
// 0, 8    root
// 1, 8        key: r -> <nil>
// 2, 3            key: oman -> <nil>
// 3, 1                key: e -> 0
// 3, 1                key: us -> 1
// 2, 4            key: ub -> <nil>
// 3, 3                key: ic -> <nil>
// 4, 1                    key: on -> 108
// 4, 1                    key: undus -> (func() bool)(0x108e200)

// ordered iteration of the tree key/value pairs:
// romane 0
// romanus 1
// rubicon 108
// rubicundus 0x108e200
```