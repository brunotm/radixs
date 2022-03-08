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

func TestSetGet(t *testing.T) {
	assert := newAssert(t)
	tr, err := FromMap(pairs)
	assert(err == nil, "error creating tree from map", "err:", err)

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
	tr, err := FromMap(kv)
	assert(err == nil, "error creating tree from map", "err:", err)

	var count int
	for k := range kv {
		kv[k] = count
		_ = tr.Set(k, count)
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
	tr, err := FromMap(pairs)
	assert(err == nil, "error creating tree from map", "err:", err)

	k := "smash"
	v := "potato"
	err = tr.Set(k, v)
	assert(err == nil, "error setting key:", k, "error:", err)

	value, err := tr.Get(k)
	assert(err == nil, "key:", k, "not found, err:", err)
	assert(value == v, "value differ for key:", k, "expected:", v, "got:", value)
}

func TestLongestMatch(t *testing.T) {
	assert := newAssert(t)
	tr, err := FromMap(pairs)
	assert(err == nil, "error creating tree from map", "err:", err)

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
	tr, err := FromMap(pairs)
	assert(err == nil, "error creating tree from map", "err:", err)

	key := "toma"
	err = tr.Delete(key)
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
	tr, err := FromMap(pairs)
	assert(err == nil, "error creating tree from map", "err:", err)

	key := "rubbe"
	err = tr.DeletePrefix(key)
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

	err := tr.SetWithParams("/api/v1/projects/:project", "ProjectHandler")
	assert(err == nil, "error setting key:", "/api/v1/projects/:project", "error:", err)

	err = tr.SetWithParams("/api/v1/projects/:project/instances/:instance", "InstanceHandler")
	assert(err == nil, "error setting key:", "/api/v1/projects/:project/instances/:instance", "error:", err)

	err = tr.SetWithParams("/api/v1/projects/:project/instances/:instance/databases/:database", "DatabaseHandler")
	assert(err == nil, "error setting key:", "/api/v1/projects/:project/instances/:instance/databases/:database", "error:", err)

	err = tr.SetWithParams("/api/v1/projects/:project/instances/:instance/applications/:application", "ApplicationHandler")
	assert(err == nil, "error setting key:", "/api/v1/projects/:project/instances/:instance/applications/:application", "error:", err)

	err = tr.SetWithParams("/api/v1/accounts", "AccountsHandler")
	assert(err == nil, "error setting key:", "/api/v1/accounts", "error:", err)

	params := map[string]string{}
	key := "/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W"
	value, err := tr.GetWithParams(key, params)
	assert(err == nil, "failed to set key:", key, "error:", err, "params", params)
	assert(value == "ProjectHandler", "wrong value for key:", key, "got:", value, "expected:", "ProjectHandler", "params", params)
	assert(len(params) == 1 && params["project"] == "01FW1D5RWNR6MEZDJZZYJX8G2W", "invalid parameters for key:", key, "params", params)

	params = map[string]string{}
	key = "/api/v1/accounts"
	value, err = tr.GetWithParams(key, params)
	assert(err == nil, "failed to get key:", key, "error:", err, "params", params)
	assert(value == "AccountsHandler", "wrong value for key:", key, "got:", value, "expected:", "AccountsHandler", "params", params)
	assert(len(params) == 0, "invalid parameters for key:", key, "params", params)

	params = map[string]string{}
	key = "/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459"
	value, err = tr.GetWithParams(key, params)
	assert(err == nil, "failed to get key:", key, "error:", err, "params", params)
	assert(value == "InstanceHandler", "wrong value for key:", key, "got:", value, "expected:", "InstanceHandler", "params", params)
	assert(
		len(params) == 2 && params["project"] == "01FW1D5RWNR6MEZDJZZYJX8G2W" && params["instance"] == "31459",
		"invalid parameters for key:", key, "params", params,
	)

	params = map[string]string{}
	key = "/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459/databases/ordersdb"
	value, err = tr.GetWithParams(key, params)
	assert(err == nil, "failed to set key:", key, "error:", err, "params", params)
	assert(value == "DatabaseHandler", "wrong value for key:", key, "got:", value, "expected:", "DatabaseHandler", "params", params)
	assert(
		len(params) == 3 && params["project"] == "01FW1D5RWNR6MEZDJZZYJX8G2W" &&
			params["instance"] == "31459" && params["database"] == "ordersdb",
		"invalid parameters for key:", key, "params", params,
	)

	params = map[string]string{}
	key = "/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459/applications/application1"
	value, err = tr.GetWithParams(key, params)
	assert(err == nil, "failed to get key:", key, "error:", err, "params", params)
	assert(value == "ApplicationHandler", "wrong value for key:", key, "got:", value, "expected:", "ApplicationHandler", "params", params)
	assert(
		len(params) == 3 && params["project"] == "01FW1D5RWNR6MEZDJZZYJX8G2W" &&
			params["instance"] == "31459" && params["application"] == "application1",
		"invalid parameters for key:", key, "params", params,
	)
}

func TestNeighborMatch(t *testing.T) {
	assert := newAssert(t)

	tr, err := FromMap(pairs)
	assert(err == nil, "error creating tree from map", "err:", err)

	_ = tr.Set("small", 67)
	_ = tr.Set("sma", 677)

	expect := map[string]interface{}{
		"sma": 677, "small": 67, "smaller": 81, "smallish": 82, "smart": 83,
	}

	neighboors := make(map[string]interface{})
	err = tr.NeighborMatch("smalle", neighboors)
	assert(err == nil, "error in neighbor match", "err:", err)

	for k, v := range expect {
		assert(neighboors[k] == v, "invalid result for key:", k, "expected:", v, "got:", neighboors[k])
	}
}

func TestSearchKeyExhaustion(t *testing.T) {
	assert := newAssert(t)

	tr, err := FromMap(pairs)
	assert(err == nil, "error creating tree from map", "err:", err)
	_ = tr.Set("small", 67)

	_, err = tr.Get("smalle")
	assert(err != nil, "get: key should not exist", "err:", err)

	match, _, err := tr.LongestMatch("smalle")
	assert(err == nil && match == "small", "longest match: invalid match", match, "err:", err)

	neighboors := make(map[string]interface{})
	err = tr.NeighborMatch("smalle", neighboors)
	assert(err == nil && len(neighboors) == 4, "neighbor match: invalid matches:", neighboors, "err:", err)
}

func TestFirstParamElement(t *testing.T) {
	assert := newAssert(t)
	tr := New(WithParams(':', '@'))
	err := tr.SetWithParams("@namespace:documents:accounts:@accountId:@subscriptionId:@resourceType:@resourceId", "value")
	assert(err == nil, "error setting with params:", err)

	err = tr.SetWithParams("@namespace:files:accounts:@accountId:@subscriptionId:@resourceType:@resourceId", "value")
	assert(err == nil, "error setting with params:", err)

	params := map[string]string{}
	value, err := tr.GetWithParams("my-company:documents:accounts:E7B4320A06A1:DBCAB1AD:document:46D05077510E", params)
	assert(err == nil, "error getting with params:", err, "params", params)
	assert(value != nil && value.(string) == "value", "wrong value:", value, "params", params)
	assert(
		params["namespace"] == "my-company" &&
			params["accountId"] == "E7B4320A06A1" &&
			params["subscriptionId"] == "DBCAB1AD" &&
			params["resourceType"] == "document" &&
			params["resourceId"] == "46D05077510E",
		"wrong parameters", params)

	params = map[string]string{}
	value, err = tr.GetWithParams("my-company:files:accounts:E7B4320A06A1:DBCAB1AD:file:46D05077510E", params)
	assert(err == nil, "error getting with params:", err, "params", params)
	assert(value != nil && value.(string) == "value", "wrong value:", value, "params", params)
	assert(
		params["namespace"] == "my-company" &&
			params["accountId"] == "E7B4320A06A1" &&
			params["subscriptionId"] == "DBCAB1AD" &&
			params["resourceType"] == "file" &&
			params["resourceId"] == "46D05077510E",
		"wrong parameters", params)
}

func TestMultipleParamsSameKey(t *testing.T) {
	assert := newAssert(t)
	tr := New(WithParams(':', '@'))
	err := tr.SetWithParams("urn:documents:accounts:@accountId:@subscriptionId:@resourceType:@resourceId", "value")
	assert(err == nil, "error setting with params:", err)

	params := map[string]string{}
	value, err := tr.GetWithParams("urn:documents:accounts:E7B4320A06A1:DBCAB1AD:document:46D05077510E", params)
	assert(err == nil, "error getting with params:", err, "params", params)
	assert(value != nil && value.(string) == "value", "wrong value:", value, "params", params)
	assert(
		params["accountId"] == "E7B4320A06A1" &&
			params["subscriptionId"] == "DBCAB1AD" &&
			params["resourceType"] == "document" &&
			params["resourceId"] == "46D05077510E",
		"wrong parameters", params)

	err = tr.SetWithParams("urn:documents:accounts:@accountId:@subscriptionId:@resourceType:@resourceId:admin", "admin")
	assert(err == nil, "error setting with params:", err)

	params = map[string]string{}
	value, err = tr.GetWithParams("urn:documents:accounts:E7B4320A06A1:DBCAB1AD:document:46D05077510E:admin", params)
	assert(err == nil, "error getting with params:", err, "params", params)
	assert(value != nil && value.(string) == "admin", "wrong value:", value, "params", params)
	assert(
		params["accountId"] == "E7B4320A06A1" &&
			params["subscriptionId"] == "DBCAB1AD" &&
			params["resourceType"] == "document" &&
			params["resourceId"] == "46D05077510E",
		"wrong parameters", params)

	err = tr.SetWithParams("urn:documents:accounts:@accountId:@subscriptionId:@resourceType:@resourceId:admin:@adminId", "admin")
	assert(err == nil, "error setting with params:", err)

	params = map[string]string{}
	value, err = tr.GetWithParams("urn:documents:accounts:E7B4320A06A1:DBCAB1AD:document:46D05077510E:admin:XYZ", params)
	assert(err == nil, "error getting with params:", err, "params", params)
	assert(value != nil && value.(string) == "admin", "wrong value:", value, "params", params)
	assert(
		params["accountId"] == "E7B4320A06A1" &&
			params["subscriptionId"] == "DBCAB1AD" &&
			params["resourceType"] == "document" &&
			params["resourceId"] == "46D05077510E",
		params["adminId"] == "XYZ",
		"wrong parameters", params)

	params = map[string]string{}
	_, err = tr.GetWithParams("urn:documents:accounts:E7B4320A06A1", params)
	assert(err != nil, "expected key not found:", err, "params", params)
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
