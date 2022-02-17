package radixs

import (
	"crypto/rand"
	"fmt"
	"testing"
)

var pairs = map[string]interface{}{
	"roma":       0,
	"romane":     1,
	"romanus":    2,
	"romulus":    3,
	"rubens":     4,
	"rube":       5,
	"rubber":     51,
	"rubberized": 511,
	"rubberize":  512,
	"rubicon":    6,
	"rubicundus": 7,
	"smaller":    81,
	"smallerish": 811,
	"smallish":   82,
	"smart":      83,
	"smarter":    84,
	"smarting":   85,
}

var stringRep = `D, W
0, 24    root
1, 16        key: r -> <nil>
2, 6            key: om -> <nil>
3, 4                key: a -> 0
4, 3                    key: n -> <nil>
5, 1                        key: e -> 1
5, 1                        key: us -> 2
3, 1                key: ulus -> 3
2, 9            key: ub -> <nil>
3, 3                key: ber -> 51
4, 2                    key: ize -> 512
5, 1                        key: d -> 511
3, 2                key: e -> 5
4, 1                    key: ns -> 4
3, 3                key: ic -> <nil>
4, 1                    key: on -> 6
4, 1                    key: undus -> 7
1, 8        key: sma -> <nil>
2, 4            key: ll -> <nil>
3, 2                key: er -> 81
4, 1                    key: ish -> 811
3, 1                key: ish -> 82
2, 3            key: rt -> 83
3, 1                key: er -> 84
3, 1                key: ing -> 85
`

func TestTreeSize(t *testing.T) {
	tr := FromMap(pairs)

	if tr.Size() != uint64(len(pairs)) {
		t.Errorf("expected size: %d, got %d\n%s\n", len(pairs), tr.size, tr.String())
	}

	tr.Delete("smart")
	if tr.Size() != uint64(len(pairs))-1 {
		t.Errorf("expected size: %d, got %d\n%s\n", len(pairs), tr.size, tr.String())
	}

	tr.DeletePrefix("rubber")
	if tr.Size() != uint64(len(pairs))-4 {
		t.Errorf("expected size: %d, got %d\n%s\n", len(pairs), tr.size, tr.String())
	}
}

func TestStringRep(t *testing.T) {
	tr := FromMap(pairs)
	if tr.String() != stringRep {
		t.Errorf("expected:\n%s \ngot:\n%s\n", stringRep, tr.String())
	}
}

func TestSetGet(t *testing.T) {
	tr := FromMap(pairs)

	for k, v := range pairs {
		value, ok := tr.Get(k)
		if !ok {
			t.Errorf("key %s should be present in tree\n%s", k, tr.String())
		}

		if value != v {
			t.Errorf("key: %s, value: %#v, expected: %#v\n%s", k, value, v, tr.String())
		}
	}

	value, ok := tr.Get("smalerishy")
	if ok {
		t.Errorf("key: sma, value: %#v, should not exist\n%s", value, tr.String())
	}

	value, ok = tr.Get("romanei")
	if ok {
		t.Errorf("key: sma, value: %#v, should not exist\n%s", value, tr.String())
	}

	value, ok = tr.Get("")
	if ok {
		t.Errorf("get empty key value: %#v, should not exist\n%s", value, tr.String())
	}

	if tr.Set("", "abc") {
		t.Errorf("empty key was set\n%s", tr.String())
	}

	if tr.Set("abc", nil) {
		t.Errorf("nil value was set\n%s", tr.String())
	}
}

func TestSetUpdate(t *testing.T) {
	kv := copyMap(pairs)
	tr := FromMap(kv)

	var count int
	for k := range kv {
		kv[k] = count
		tr.Set(k, count)
		count++
	}

	for k, v := range kv {
		value, ok := tr.Get(k)
		if !ok {
			t.Errorf("key %s should be present in tree\n%s", k, tr.String())
		}

		if value != v {
			t.Errorf("key: %s, value: %#v, expected: %#v\n%s", k, value, v, tr.String())
		}
	}
}

func TestSetSplit(t *testing.T) {
	tr := FromMap(pairs)
	k := "smash"
	v := "potato"
	tr.Set(k, v)

	value, ok := tr.Get(k)
	if !ok {
		t.Errorf("key %s should be present in tree\n%s", k, tr.String())
	}

	if value != v {
		t.Errorf("key: %s, value: %#v, expected: %#v\n%s", k, value, v, tr.String())
	}
}

func TestLongestMatch(t *testing.T) {
	tr := FromMap(pairs)

	prefix, value, ok := tr.LongestMatch("smarties")
	if !ok {
		t.Errorf("longest match for smarties not found\n%s", tr.String())
	}

	if prefix != "smart" {
		t.Errorf(
			"longest match for smarties expected: smart, got: %s\n%s",
			prefix, tr.String(),
		)
	}

	if value != pairs["smart"] {
		t.Errorf(
			"longest match for smarties expected value: %#v, got: %s\n%s",
			pairs["smart"], value, tr.String())
	}

	prefix, value, ok = tr.LongestMatch("rubberized")
	if !ok {
		t.Errorf("longest match for rubberized not found\n%s", tr.String())
	}

	if prefix != "rubberized" {
		t.Errorf(
			"longest match for rubberized expected: rubberized, got: %s\n%s",
			prefix, tr.String(),
		)
	}

	if value != pairs["rubberized"] {
		t.Errorf(
			"longest match for rubberized expected value: %#v, got: %s\n%s",
			pairs["rubberized"], value, tr.String())
	}

	_, _, ok = tr.LongestMatch("smallest")
	if ok {
		t.Errorf("longest match for smallest should not exist\n%s", tr.String())
	}

}

func TestSetDelete(t *testing.T) {
	tr := FromMap(pairs)

	if tr.Delete("toma") {
		t.Errorf("deleted non existing key: romarish\n%s", tr.String())
	}

	if tr.Delete("romarish") {
		t.Errorf("deleted non existing key: romarish\n%s", tr.String())
	}

	if !tr.Delete("roma") {
		t.Errorf("failed to delete key: roma\n%s", tr.String())
	}

	if !tr.Delete("smart") {
		t.Errorf("failed to delete key: smart\n%s", tr.String())
	}

	if !tr.Delete("rubberized") {
		t.Errorf("failed to delete key: rubberized\n%s", tr.String())
	}

	if !tr.Delete("smallish") {
		t.Errorf("failed to delete key: rubberized\n%s", tr.String())
	}

	k := "romanus"
	value, ok := tr.Get(k)
	if !ok {
		t.Errorf("key %s should be present in tree\n%s", k, tr.String())
	}

	if value != pairs[k] {
		t.Errorf("key: %s, value: %#v, expected: %#v\n%s", k, value, pairs[k], tr.String())
	}

	k = "smarter"
	value, ok = tr.Get(k)
	if !ok {
		t.Errorf("key %s should be present in tree\n%s", k, tr.String())
	}

	if value != pairs[k] {
		t.Errorf("key: %s, value: %#v, expected: %#v\n%s", k, value, pairs[k], tr.String())
	}

	k = "rubberize"
	value, ok = tr.Get(k)
	if !ok {
		t.Errorf("key %s should be present in tree\n%s", k, tr.String())
	}

	if value != pairs[k] {
		t.Errorf("key: %s, value: %#v, expected: %#v\n%s", k, value, pairs[k], tr.String())
	}
}

func TestDeletePrefix(t *testing.T) {
	tr := FromMap(pairs)

	if !tr.DeletePrefix("rubbe") {
		t.Errorf("failed to delete existing prefix: rubbe\n%s", tr.String())
	}

	k := "rube"
	value, ok := tr.Get(k)
	if !ok {
		t.Errorf("key %s should be present in tree\n%s", k, tr.String())
	}

	if value != pairs[k] {
		t.Errorf("key: %s, value: %#v, expected: %#v\n%s", k, value, pairs[k], tr.String())
	}

	k = "rubber"
	value, ok = tr.Get(k)
	if ok {
		t.Errorf("key %s, value: %#v should not be present in tree\n%s", k, value, tr.String())
	}

	k = "rubberized"
	value, ok = tr.Get(k)
	if ok {
		t.Errorf("key %s, value: %#v should not be present in tree\n%s", k, value, tr.String())
	}

	k = "rubberize"
	value, ok = tr.Get(k)
	if ok {
		t.Errorf("key %s, value: %#v should not be present in tree\n%s", k, value, tr.String())
	}

	if !tr.DeletePrefix("small") {
		t.Errorf("failed to delete existing prefix: rubbe\n%s", tr.String())
	}

	k = "smaller"
	value, ok = tr.Get(k)
	if ok {
		t.Errorf("key %s, value: %#v should not be present in tree\n%s", k, value, tr.String())
	}
}

func TestIter(t *testing.T) {
	tr := FromMap(pairs)
	kvs := copyMap(pairs)

	tr.Iter(func(key string, value interface{}) bool {
		v, ok := kvs[key]
		if !ok {
			t.Errorf("key %s, not present in source\n%s", key, tr.String())
		}

		if value != v {
			t.Errorf("key %s, value: %#v incorrect source: %#v\n%s", key, value, v, tr.String())
		}

		delete(kvs, key)
		return true
	})

	if len(kvs) > 0 {
		t.Errorf("all keys should be deleted in source, remaining: %#v\n%s", kvs, tr.String())
	}

	var stop bool
	tr.Iter(func(key string, value interface{}) bool {
		if stop {
			t.Errorf("iter should have stopped the iteration")
		}

		stop = true
		return false
	})
}

func TestSeWithParams(t *testing.T) {
	tr := New(WithParams('/', ':'))

	if !tr.SetWithParams("/api/v1/projects/:project", "ProjectsHandler") {
		t.Errorf("key should be set")
	}

	if !tr.SetWithParams("/api/v1/projects/:project/instances/:instance", "InstanceHandler") {
		t.Errorf("key should be set")
	}

	if !tr.SetWithParams("/api/v1/projects/:project/instances/:instance/databases/:database", "DatabaseHandler") {
		t.Errorf("key should be set")
	}

	if tr.SetWithParams("/api/v1/projects//:project/instances/:instance/databases/:database", "DatabaseHandler") {
		t.Errorf("key should not be set")
	}

	if tr.SetWithParams("/api/v1:/projects/:project/instances/:instance/databases/:database", "DatabaseHandler") {
		t.Errorf("key should not be set")
	}

	if tr.SetWithParams("/api/v1/projects/::project/instances/:instance/databases/:database", "DatabaseHandler") {
		t.Errorf("key should not be set")
	}

	if tr.SetWithParams("/api/v1/projects/project/instances/:instance/databases/:database", "DatabaseHandler") {
		fmt.Println(tr.String())
		t.Errorf("key should not be set")
	}

	if tr.SetWithParams("/api/v1/projects/:state/instances/:instance/databases/:database", "DatabaseHandler") {
		t.Errorf("key should not be set")
	}
}

func TestGetWithParams(t *testing.T) {
	tr := New(WithParams('/', ':'))

	tr.SetWithParams("/api/v1/projects/:project", "ProjectHandler")
	tr.SetWithParams("/api/v1/projects/:project/instances/:instance", "InstanceHandler")
	tr.SetWithParams("/api/v1/projects/:project/instances/:instance/databases/:database", "DatabaseHandler")

	params := map[string]string{}
	v, ok := tr.GetWithParams("/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W", params)
	if !ok {
		t.Errorf("key should be set: %#v, %#v, %t", v, params, ok)
	}
	if v != "ProjectHandler" {
		t.Errorf("invalid value")
	}
	if len(params) != 1 || params["project"] != "01FW1D5RWNR6MEZDJZZYJX8G2W" {
		t.Errorf("invalid parameters")
	}

	params = map[string]string{}
	v, ok = tr.GetWithParams("/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459", params)
	if !ok {
		t.Errorf("key should be set")
	}
	if v != "InstanceHandler" {
		t.Errorf("invalid value")
	}
	if len(params) != 2 || params["project"] != "01FW1D5RWNR6MEZDJZZYJX8G2W" || params["instance"] != "31459" {
		t.Errorf("invalid parameters")
	}

	params = map[string]string{}
	v, ok = tr.GetWithParams("/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459/databases/ordersdb", params)
	if !ok {
		t.Errorf("key should be set")
	}
	if v != "DatabaseHandler" {
		t.Errorf("invalid value")
	}
	if len(params) != 3 || params["project"] != "01FW1D5RWNR6MEZDJZZYJX8G2W" || params["instance"] != "31459" || params["database"] != "ordersdb" {
		t.Errorf("invalid parameters")
	}
}

func TestRandomLoad(t *testing.T) {
	keyCount := 100000
	tr := New()

	for x := 0; x < keyCount; x++ {
		key := generateUUID()
		if !tr.Set(key, x) {
			t.Errorf("failure to set: %s", key)
		}
	}

	if tr.Size() != uint64(keyCount) {
		t.Errorf("leafs expected: %d, got: %d", keyCount, tr.Size())
	}

	var remove []string
	tr.Iter(func(key string, value interface{}) bool {
		remove = append(remove, key)
		return true
	})

	for x := 0; x < len(remove); x++ {
		if !tr.Delete(remove[x]) {
			fmt.Println(tr.String())
			t.Errorf("failure to delete: %s", remove[x])
		}
	}
}

func BenchmarkSingleRead(b *testing.B) {
	tr := FromMap(pairs)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tr.Get("smart")
	}
}

func BenchmarkSingleReadWithParameters(b *testing.B) {
	tr := FromMap(pairs, WithParams('/', ':'))
	tr.SetWithParams("/api/v1/projects/:project", "CITY")
	tr.SetWithParams("/api/v1/projects/:project/instances/:instance", "MONUMENT")
	tr.SetWithParams("/api/v1/projects/:project/instances/:instance/databases/:database", "DatabaseHandler")
	params := map[string]string{}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tr.GetWithParams("/api/v1/projects/Lisbon/instances/:instancer", params)
	}
}

func BenchmarkSingleInsert(b *testing.B) {
	tr := FromMap(pairs)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tr.Set("smarty", n)
	}
}

func BenchmarkSingleInsertWithParameters(b *testing.B) {
	tr := FromMap(pairs, WithParams('/', ':'))
	tr.SetWithParams("/api/v1/projects/:project", "CITY")
	tr.SetWithParams("/api/v1/projects/:project/instances/:instance", "MONUMENT")

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tr.SetWithParams("/api/v1/projects/:project/instances/:instance", n)
	}
}

func BenchmarkLongestMatch(b *testing.B) {
	tr := FromMap(pairs)
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_, _, _ = tr.LongestMatch("smart")
	}
}

// generateUUID is used to generate a random UUID
func generateUUID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Errorf("generateUUID: failed to read random: %w", err))
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16])
}

func copyMap(src map[string]interface{}) (dst map[string]interface{}) {
	dst = make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
