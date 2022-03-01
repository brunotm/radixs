# radixs

[![Build Status](https://github.com/brunotm/radixs/actions/workflows/test.yml/badge.svg)](https://github.com/brunotm/radixs/actions)
[![Go Report Card](https://goreportcard.com/badge/brunotm/radixs?cache=0)](https://goreportcard.com/report/brunotm/radixs)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/brunotm/radixs)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/brunotm/radixs/master/LICENSE)

---

A Go implementation of a radix tree for building fast and compact in-memory indexes, that uses binary searches to speed up insert, retrieve and delete operations on dense trees.

This implementation in addition of using binary searches for:

- eliminate edge specific types and pointers in order to save memory and allocations on dense trees.
- insert, retrieve and delete operations are non recursive in order to avoid the lack of tail call optimization in the Go compiler.
- tree nodes are memory aligned for optimal space utilization.
- supports longest prefix partial matches
- supports longest prefix neighbor matches
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
	_ = tr.Set("romane", 0)
	_ = tr.Set("romanus", 1)
	_ = tr.Set("romulus", "remus brother")
	_ = tr.Set("rubens", "51.2170° N, 4.4093° E")
	_ = tr.Set("ruber", 4)
	_ = tr.Set("rubicon", 108)
	_ = tr.Set("rubicundus", func() bool { return true })

	value, err = tr.Get("rubens")
	fmt.Printf("value for rubens: %#v, err: %s\n", value, err)

	size := tr.Size()
	fmt.Printf("tree size: %d\n", size)

	key, value, err := tr.LongestMatch("romanesco")
	fmt.Printf("longest prefix for romanesco: %s, value: %#v, err: %s\n",
		key, value, err)

	fmt.Printf("\ntree string representation:\n%s\n", tr.String())

	err = tr.Delete("romulus")
	fmt.Printf("romulus deleted: %s\n", err)

	err = tr.DeletePrefix("rube")
	fmt.Printf("prefix rube deleted: %s\n", err)

	fmt.Printf("\ntree string representation:\n%s\n", tr.String())

	fmt.Printf("ordered iteration of the tree key/value pairs:\n")
	tr.Iter(func(key string, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})


// value for rubens: "51.2170° N, 4.4093° E", err: %!s(<nil>)
// tree size: 7
// longest prefix for romanesco: , value: <nil>, err: radixs: key not found

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
// 4, 1                    key: undus -> (func() bool)(0x108f900)

// romulus deleted: %!s(<nil>)
// prefix rube deleted: %!s(<nil>)

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
// 4, 1                    key: undus -> (func() bool)(0x108f900)

// ordered iteration of the tree key/value pairs:
// romane 0
// romanus 1
// rubicon 108
// rubicundus 0x108f900
```

## Usage With Parameters
```go

tr := radixs.New(radixs.WithParams('/', ':')) // delimiter: '/', parameter placeholder: ':'

	_ = tr.SetWithParams("/api/v1/projects", "ProjectsHandler")
	_ = tr.SetWithParams("/api/v1/projects/:project", "ProjectHandler")
	_ = tr.SetWithParams("/api/v1/projects/:project/instances/:instance", "InstanceHandler")
	_ = tr.SetWithParams("/api/v1/projects/:project/instances/:instance/operations/:operation", "OperationHandler")

	_ = tr.SetWithParams("/api/v1/projects/:project/instances/:instance/databases/:database", "DatabaseHandler")
	_ = tr.SetWithParams("/api/v1/projects/:project/instances/:instance/databases/:database/resources/:rsrc", "ResourceHandler")
	_ = tr.SetWithParams("/api/v1/projects/:project/instances/:instance/databases/:database/sessions/:session", "SessionHandler")
	_ = tr.SetWithParams("/api/v1/users", "UsersHandler")
	_ = tr.SetWithParams("/api/v1/users/:users", "UserHandler")
	_ = tr.SetWithParams("/api/v1/users/:users/messages", "MessagesHandler")

	fmt.Println(tr.String())

	params := map[string]string{}
	value, err := tr.GetWithParams("/api/v1/projects", params)
	fmt.Printf("value: %#v, params: %#v, err: %s\n", value, params, err)

	params = map[string]string{}
	value, err = tr.GetWithParams("/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W", params)
	fmt.Printf("value: %#v, params: %#v, err: %s\n", value, params, err)

	params = map[string]string{}
	value, err = tr.GetWithParams(
		"/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459",
		params)
	fmt.Printf("value: %#v, params: %#v, err: %s\n", value, params, err)

	params = map[string]string{}
	value, err = tr.GetWithParams(
		"/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459/operations/upgrade",
		params)
	fmt.Printf("value: %#v, params: %#v, err: %s\n", value, params, err)

	params = map[string]string{}
	value, err = tr.GetWithParams(
		"/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459/databases/ordersdb",
		params)
	fmt.Printf("value: %#v, params: %#v, err: %s\n", value, params, err)

	params = map[string]string{}
	value, err = tr.GetWithParams(
		"/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459/databases/ordersdb/resources/order_items",
		params)
	fmt.Printf("value: %#v, params: %#v, err: %s\n", value, params, err)

	params = map[string]string{}
	value, err = tr.GetWithParams(
		"/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459/databases/ordersdb/sessions/281474976710655",
		params)
	fmt.Printf("value: %#v, params: %#v, err: %s\n", value, params, err)

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

// value: "ProjectsHandler", params: map[string]string{}, err: %!s(<nil>)
// value: "ProjectHandler", params: map[string]string{"project":"01FW1D5RWNR6MEZDJZZYJX8G2W"}, err: %!s(<nil>)
// value: "InstanceHandler", params: map[string]string{"instance":"31459", "project":"01FW1D5RWNR6MEZDJZZYJX8G2W"}, err: %!s(<nil>)
// value: "OperationHandler", params: map[string]string{"instance":"31459", "operation":"upgrade", "project":"01FW1D5RWNR6MEZDJZZYJX8G2W"}, err: %!s(<nil>)
// value: "DatabaseHandler", params: map[string]string{"database":"ordersdb", "instance":"31459", "project":"01FW1D5RWNR6MEZDJZZYJX8G2W"}, err: %!s(<nil>)
// value: "ResourceHandler", params: map[string]string{"database":"ordersdb", "instance":"31459", "project":"01FW1D5RWNR6MEZDJZZYJX8G2W", "rsrc":"order_items"}, err: %!s(<nil>)
// value: "SessionHandler", params: map[string]string{"database":"ordersdb", "instance":"31459", "project":"01FW1D5RWNR6MEZDJZZYJX8G2W", "session":"281474976710655"}, err: %!s(<nil>)
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