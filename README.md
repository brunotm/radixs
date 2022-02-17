# radixs

A Go implementation of a radix tree for building fast and compact in-memory indexes, that uses binary searches to speed up insert, retrieve and delete operations on dense trees.

This implementation in addition of using binary searches for:

- eliminate edge specific types and pointers in order to save memory and allocations on dense trees.
- insert, retrieve and delete operations are non recursive in order to avoid the lack of tail call optimization in the Go compiler.
- tree nodes are memory aligned for optimal space utilization.
- supports longest prefix partial matches
- supports key parameters and delimiters

___

## Usage Without Parameters

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

## Usage With Parameters
```go

tr := radixs.New(radixs.WithParams('/', ':'))

	fmt.Println(tr.SetWithParams("/api/v1/projects", "ProjectsHandler"))
	fmt.Println(tr.SetWithParams("/api/v1/projects/:project", "ProjectHandler"))
	fmt.Println(tr.SetWithParams("/api/v1/projects/:project/instances/:instance", "InstanceHandler"))
	fmt.Println(tr.SetWithParams("/api/v1/projects/:project/instances/:instance/operations/:operation", "OperationHandler"))

	fmt.Println(tr.SetWithParams("/api/v1/projects/:project/instances/:instance/databases/:database", "DatabaseHandler"))
	fmt.Println(tr.SetWithParams("/api/v1/projects/:project/instances/:instance/databases/:database/resources/:rsrc", "ResourceHandler"))
	fmt.Println(tr.SetWithParams("/api/v1/projects/:project/instances/:instance/databases/:database/sessions/:session", "SessionHandler"))
	fmt.Println(tr.SetWithParams("/api/v1/users", "UsersHandler"))
	fmt.Println(tr.SetWithParams("/api/v1/users/:users", "UserHandler"))
	fmt.Println(tr.SetWithParams("/api/v1/users/:users/messages", "MessagesHandler"))

	fmt.Println(tr.String())

	params := map[string]string{}
	value, ok := tr.GetWithParams("/api/v1/projects", params)
	fmt.Printf("value: %#v, params: %#v, ok: %t\n", value, params, ok)

	params = map[string]string{}
	value, ok = tr.GetWithParams("/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W", params)
	fmt.Printf("value: %#v, params: %#v, ok: %t\n", value, params, ok)

	params = map[string]string{}
	value, ok = tr.GetWithParams(
		"/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459",
		params)
	fmt.Printf("value: %#v, params: %#v, ok: %t\n", value, params, ok)

	params = map[string]string{}
	value, ok = tr.GetWithParams(
		"/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459/operations/upgrade",
		params)
	fmt.Printf("value: %#v, params: %#v, ok: %t\n", value, params, ok)

	params = map[string]string{}
	value, ok = tr.GetWithParams(
		"/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459/databases/ordersdb",
		params)
	fmt.Printf("value: %#v, params: %#v, ok: %t\n", value, params, ok)

	params = map[string]string{}
	value, ok = tr.GetWithParams(
		"/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459/databases/ordersdb/resources/order_items",
		params)
	fmt.Printf("value: %#v, params: %#v, ok: %t\n", value, params, ok)

	params = map[string]string{}
	value, ok = tr.GetWithParams(
		"/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459/databases/ordersdb/sessions/281474976710655",
		params)
	fmt.Printf("value: %#v, params: %#v, ok: %t\n", value, params, ok)

// D, W
// 0, 13    root
// 1, 13        key: /api/v1/ -> <nil>
// 2, 9            key: projects -> "ProjectsHandler"
// 3, 8                key: /:project -> "ProjectHandler"
// 4, 7                    key: /instances/:instance -> "InstanceHandler"
// 5, 6                        key: / -> <nil>
// 6, 4                            key: databases/:database -> "DatabaseHandler"
// 7, 3                                key: / -> <nil>
// 8, 1                                    key: resources/:rsrc -> "ResourceHandler"
// 8, 1                                    key: sessions/:session -> "SessionHandler"
// 6, 1                            key: operations/:operation -> "OperationHandler"
// 2, 3            key: users -> "UsersHandler"
// 3, 2                key: /:users -> "UserHandler"
// 4, 1                    key: /messages -> "MessagesHandler"

// value: "ProjectsHandler", params: map[string]string{}, ok: true
// value: "ProjectHandler", params: map[string]string{"project":"01FW1D5RWNR6MEZDJZZYJX8G2W"}, ok: true
// value: "InstanceHandler", params: map[string]string{"instance":"31459", "project":"01FW1D5RWNR6MEZDJZZYJX8G2W"}, ok: true
// value: "OperationHandler", params: map[string]string{"instance":"31459", "operation":"upgrade", "project":"01FW1D5RWNR6MEZDJZZYJX8G2W"}, ok: true
// value: "DatabaseHandler", params: map[string]string{"database":"ordersdb", "instance":"31459", "project":"01FW1D5RWNR6MEZDJZZYJX8G2W"}, ok: true
// value: "ResourceHandler", params: map[string]string{"database":"ordersdb", "instance":"31459", "project":"01FW1D5RWNR6MEZDJZZYJX8G2W", "rsrc":"order_items"}, ok: true
// value: "SessionHandler", params: map[string]string{"database":"ordersdb", "instance":"31459", "project":"01FW1D5RWNR6MEZDJZZYJX8G2W", "session":"281474976710655"}, ok: true
```

## Benchmarks
```
goos: darwin
goarch: amd64
pkg: github.com/brunotm/radixs
cpu: Intel(R) Core(TM) i5-8279U CPU @ 2.40GHz
BenchmarkSingleRead-8                      41779908       26.79 ns/op      0 B/op      0 allocs/op
BenchmarkSingleReadWithParameters-8        22239025       49.24 ns/op      0 B/op      0 allocs/op
BenchmarkSingleInsert-8                    19442536       59.50 ns/op      8 B/op      0 allocs/op
BenchmarkSingleInsertWithParameters-8      11472224      102.2 ns/op       8 B/op      0 allocs/op
BenchmarkLongestMatch-8                    14452323       73.01 ns/op      5 B/op      1 allocs/op
```