package radixs

import (
	"testing"
)

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
