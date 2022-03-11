package radixs

import "testing"

func TestDelete(t *testing.T) {
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
