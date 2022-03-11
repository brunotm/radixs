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

	tr, err := FromMap(pairs)
	assert(err == nil, "error creating tree from map", "err:", err)
	assert(tr.Size() == uint64(len(pairs)), "expected size:", len(pairs), "got:", tr.Size())

	err = tr.Delete("smart")
	assert(err == nil && tr.Size() == uint64(len(pairs))-1, "expected size:", len(pairs)-1, "got:", tr.Size(), "err:", err)

	err = tr.DeletePrefix("rubber")
	assert(err == nil && tr.Size() == uint64(len(pairs))-4, "expected size:", len(pairs)-4, "got:", tr.Size(), "err:", err)
}

func TestStringRep(t *testing.T) {
	assert := newAssert(t)
	tr, err := FromMap(pairs)
	assert(err == nil, "error creating tree from map", "err:", err)
	assert(tr.String() == stringRep, "expected:", stringRep, "got:", tr.String())
}

func TestIter(t *testing.T) {
	assert := newAssert(t)
	tr, err := FromMap(pairs)
	assert(err == nil, "error creating tree from map", "err:", err)

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
	assert := newAssert(b)
	tr, err := FromMap(pairs)
	assert(err == nil, "error creating tree from map", "err:", err)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_, _ = tr.Get("smart")
	}
}

func BenchmarkSingleReadWithParameters(b *testing.B) {
	tr, _ := FromMap(pairs, WithParams('/', ':'))
	_ = tr.SetWithParams("/api/v1/projects/:project", "CITY")
	_ = tr.SetWithParams("/api/v1/projects/:project/instances/:instance", "MONUMENT")
	_ = tr.SetWithParams("/api/v1/projects/:project/instances/:instance/databases/:database", "DatabaseHandler")
	params := map[string]string{}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_, _ = tr.GetWithParams("/api/v1/projects/Lisbon/instances/:instancer", params)
	}
}

func BenchmarkSingleInsert(b *testing.B) {
	tr, _ := FromMap(pairs)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = tr.Set("smarty", n)
	}
}

func BenchmarkSingleInsertWithParameters(b *testing.B) {
	tr, _ := FromMap(pairs, WithParams('/', ':'))
	_ = tr.SetWithParams("/api/v1/projects/:project", "CITY")
	_ = tr.SetWithParams("/api/v1/projects/:project/instances/:instance", "MONUMENT")

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = tr.SetWithParams("/api/v1/projects/:project/instances/:instance", n)
	}
}

func BenchmarkLongestMatch(b *testing.B) {
	tr, _ := FromMap(pairs)

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
