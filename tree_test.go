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

func newAssert(t testing.TB) func(cond bool, kvs ...interface{}) {
	return func(cond bool, args ...interface{}) {
		if !cond {
			t.Helper()
			t.Error(args...)
		}
	}
}

func TestTreeSize(t *testing.T) {
	assert := newAssert(t)

	tr := FromMap(pairs)
	assert(tr.Size() == uint64(len(pairs)), "expected size:", len(pairs), "got:", tr.Size())

	err := tr.Delete("smart")
	assert(err == nil && tr.Size() == uint64(len(pairs))-1, "expected size:", len(pairs)-1, "got:", tr.Size(), "err:", err)

	err = tr.DeletePrefix("rubber")
	assert(err == nil && tr.Size() == uint64(len(pairs))-4, "expected size:", len(pairs)-4, "got:", tr.Size(), "err:", err)
}

func TestStringRep(t *testing.T) {
	assert := newAssert(t)
	tr := FromMap(pairs)
	assert(tr.String() == stringRep, "expected:", stringRep, "got:", tr.String())
}

func TestSetGet(t *testing.T) {
	assert := newAssert(t)
	tr := FromMap(pairs)

	for k, v := range pairs {
		value, err := tr.Get(k)
		assert(err == nil, "get key:", k, "error:", err)
		assert(value == v, "key:", k, "expected value:", v, "got:", value)
	}

	value, err := tr.Get("smalerishy")
	assert(err == ErrKeyNotFound, "key: smalerishy, should not exist, value:", value)

	value, err = tr.Get("romanei")
	assert(err == ErrKeyNotFound, "key: romanei, should not exist, value:", value)

	value, err = tr.Get("")
	assert(err == ErrEmptyKey, "empty key, should not exist, value:", value)

	err = tr.Set("", "abc")
	assert(err == ErrEmptyKey, "empty key was set, err:", err)

	err = tr.Set("abc", nil)
	assert(err == ErrNilValue, "nil value was set, err:", err)
}

func TestSetUpdate(t *testing.T) {
	assert := newAssert(t)
	kv := copyMap(pairs)
	tr := FromMap(kv)

	var count int
	for k := range kv {
		kv[k] = count
		tr.Set(k, count)
		count++
	}

	for k, v := range kv {
		value, err := tr.Get(k)
		assert(err == nil, "error fetching key:", k)
		assert(value == v, "value differ for key:", k, "expected:", v, "got:", value)
	}
}

func TestSetSplit(t *testing.T) {
	assert := newAssert(t)
	tr := FromMap(pairs)
	k := "smash"
	v := "potato"
	tr.Set(k, v)

	value, err := tr.Get(k)
	assert(err == nil, "key:", k, "not found, err:", err)
	assert(value == v, "value differ for key:", k, "expected:", v, "got:", value)
}

func TestLongestMatch(t *testing.T) {
	assert := newAssert(t)
	tr := FromMap(pairs)

	key := "smarties"
	expected := "smart"
	prefix, value, err := tr.LongestMatch(key)
	assert(err == nil, "longest match for key:", key, "not found, err:", err)
	assert(prefix == expected, "unmatched longest prefix, expected:", "smart", "got:", prefix)
	assert(value == pairs[expected], "unmatched longest prefix value, expected:", pairs[expected], "got:", value)

	key = "rubberized"
	expected = "rubberized"
	prefix, value, err = tr.LongestMatch(key)
	assert(err == nil, "longest match for key:", key, "not found, err:", err)
	assert(prefix == expected, "unmatched longest prefix, expected:", expected, "got:", prefix)
	assert(value == pairs[expected], "unmatched longest prefix value, expected:", pairs[expected], "got:", value)

	_, _, err = tr.LongestMatch("smallest")
	assert(err != nil, "longest match for key:", key, "should not exist")
}

func TestSetDelete(t *testing.T) {
	assert := newAssert(t)
	tr := FromMap(pairs)

	key := "toma"
	err := tr.Delete(key)
	assert(err == ErrKeyNotFound, "deleted non existing key:", key, "err:", err)

	key = "romarish"
	err = tr.Delete(key)
	assert(err == ErrKeyNotFound, "deleted non existing key:", key, "err:", err)

	key = "roma"
	err = tr.Delete(key)
	assert(err == nil, "failed to delete existing key:", key, "err:", err)

	key = "smart"
	err = tr.Delete(key)
	assert(err == nil, "failed to delete existing key:", key, "err:", err)

	key = "rubberized"
	err = tr.Delete(key)
	assert(err == nil, "failed to delete existing key:", key, "err:", err)

	key = "smallish"
	err = tr.Delete(key)
	assert(err == nil, "failed to delete existing key:", key, "err:", err)

	key = "romanus"
	value, err := tr.Get(key)
	assert(err == nil, "failed to get existing key:", key, "err:", err)
	assert(value == pairs[key], "wrong value:", value, "for key:", key)

	key = "smarter"
	value, err = tr.Get(key)
	assert(err == nil, "failed to get existing key:", key, "err:", err)
	assert(value == pairs[key], "wrong value:", value, "for key:", key)

	key = "rubberize"
	value, err = tr.Get(key)
	assert(err == nil, "failed to get existing key:", key, "err:", err)
	assert(value == pairs[key], "wrong value:", value, "for key:", key)
}

func TestDeletePrefix(t *testing.T) {
	assert := newAssert(t)
	tr := FromMap(pairs)

	key := "rubbe"
	err := tr.DeletePrefix(key)
	assert(err == nil, "failed to delete prefix:", key, "err:", err)

	key = "rube"
	value, err := tr.Get(key)
	assert(err == nil, "failed to get existing key:", key, "err:", err)
	assert(value == pairs[key], "wrong value:", value, "for key:", key)

	key = "rubber"
	value, err = tr.Get(key)
	assert(err == ErrKeyNotFound, "get non existing key:", key, "err:", err)
	assert(value == nil, "wrong value:", value, "for key:", key)

	key = "rubberized"
	value, err = tr.Get(key)
	assert(err == ErrKeyNotFound, "get non existing key:", key, "err:", err)
	assert(value == nil, "wrong value:", value, "for key:", key)

	key = "rubberize"
	value, err = tr.Get(key)
	assert(err == ErrKeyNotFound, "get non existing key:", key, "err:", err)
	assert(value == nil, "wrong value:", value, "for key:", key)

	key = "small"
	err = tr.DeletePrefix(key)
	assert(err == nil, "failed to delete prefix:", key, "err:", err)

	key = "smaller"
	value, err = tr.Get(key)
	assert(err == ErrKeyNotFound, "failed to get existing key:", key, "err:", err)
	assert(value != pairs[key], "wrong value:", value, "for key:", key)
}

func TestIter(t *testing.T) {
	assert := newAssert(t)
	tr := FromMap(pairs)
	kvs := copyMap(pairs)

	tr.Iter(func(key string, value interface{}) bool {
		v, ok := kvs[key]
		assert(ok == true, "key:", key, "not present in source")
		assert(value == v, "key:", key, "incorrect value:", value, "expected:", v)

		delete(kvs, key)
		return true
	})

	assert(len(kvs) == 0, "kvs should be empty:", kvs)

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
	assert := newAssert(t)
	tr := New(WithParams('/', ':'))

	key := "/api/v1/projects/:project"
	value := "ProjectsHandler"
	err := tr.SetWithParams(key, value)
	assert(err == nil, "error setting key:", key, "error:", err)

	key = "/api/v1/projects/:project/instances/:instance"
	value = "InstanceHandler"
	err = tr.SetWithParams(key, value)
	assert(err == nil, "error setting key:", key, "error:", err)

	key = "/api/v1/projects/:project/instances/:instance/databases/:database"
	value = "DatabaseHandler"
	err = tr.SetWithParams(key, value)
	assert(err == nil, "error setting key:", key, "error:", err)

	key = "/api/v1/projects//:project/instances/:instance/databases/:database"
	value = "DatabaseHandler"
	err = tr.SetWithParams(key, value)
	assert(err == ErrInvalidKey, "set invalid key:", key, "error:", err)

	key = "/api/v1:/projects/:project/instances/:instance/databases/:database"
	value = "DatabaseHandler"
	err = tr.SetWithParams(key, value)
	assert(err == ErrInvalidKey, "set invalid key:", key, "error:", err)

	key = "/api/v1/projects/::project/instances/:instance/databases/:database"
	value = "DatabaseHandler"
	err = tr.SetWithParams(key, value)
	assert(err == ErrInvalidKey, "set invalid key:", key, "error:", err)

	key = "/api/v1/projects/project/instances/:instance/databases/:database"
	value = "DatabaseHandler"
	err = tr.SetWithParams(key, value)
	assert(err == ErrConflictKey, "set conflicting key:", key, "error:", err)

	key = "/api/v1/projects/:state/instances/:instance/databases/:database"
	value = "DatabaseHandler"
	err = tr.SetWithParams(key, value)
	assert(err == ErrConflictKey, "set conflicting key:", key, "error:", err)
}

func TestGetWithParams(t *testing.T) {
	assert := newAssert(t)
	tr := New(WithParams('/', ':'))

	tr.SetWithParams("/api/v1/projects/:project", "ProjectHandler")
	tr.SetWithParams("/api/v1/projects/:project/instances/:instance", "InstanceHandler")
	tr.SetWithParams("/api/v1/projects/:project/instances/:instance/databases/:database", "DatabaseHandler")

	params := map[string]string{}
	key := "/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W"
	value, err := tr.GetWithParams(key, params)
	assert(err == nil, "failed to set key:", key, "error:", err)
	assert(value == "ProjectHandler", "wrong value for key:", key, "got:", value, "expected:", "ProjectHandler")
	assert(len(params) == 1 && params["project"] == "01FW1D5RWNR6MEZDJZZYJX8G2W", "invalid parameters for key:", key, params)

	params = map[string]string{}
	key = "/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459"
	value, err = tr.GetWithParams(key, params)
	assert(err == nil, "failed to set key:", key, "error:", err)
	assert(value == "InstanceHandler", "wrong value for key:", key, "got:", value, "expected:", "InstanceHandler")
	assert(
		len(params) == 2 && params["project"] == "01FW1D5RWNR6MEZDJZZYJX8G2W" && params["instance"] == "31459",
		"invalid parameters for key:", key, params,
	)

	params = map[string]string{}
	key = "/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459/databases/ordersdb"
	value, err = tr.GetWithParams(key, params)
	assert(err == nil, "failed to set key:", key, "error:", err)
	assert(value == "DatabaseHandler", "wrong value for key:", key, "got:", value, "expected:", "DatabaseHandler")
	assert(
		len(params) == 3 && params["project"] == "01FW1D5RWNR6MEZDJZZYJX8G2W" &&
			params["instance"] == "31459" && params["database"] == "ordersdb",
		"invalid parameters for key:", key, params,
	)
}

func TestRandomLoad(t *testing.T) {
	assert := newAssert(t)
	keyCount := 100000
	tr := New()

	for x := 0; x < keyCount; x++ {
		key := generateUUID()
		err := tr.Set(key, x)
		assert(err == nil, "failed to set key:", key, "err", err)
	}

	assert(tr.Size() == uint64(keyCount), "wrong leaf count, got:", tr.Size(), "expected:", keyCount)

	var remove []string
	tr.Iter(func(key string, value interface{}) bool {
		remove = append(remove, key)
		return true
	})

	for x := 0; x < len(remove); x++ {
		err := tr.Delete(remove[x])
		assert(err == nil, "failed to delete key:", remove[x], "error:", err)
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
