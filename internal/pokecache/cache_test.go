package pokecache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second

	cache := NewCache(interval)

	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for _, c := range cases {
		cache.Add(c.key, c.val)

		val, ok := cache.Get(c.key)
		if !ok {
			t.Errorf("expected to find key %q", c.key)
			continue
		}

		if string(val) != string(c.val) {
			t.Errorf(
				"expected %q, got %q",
				string(c.val),
				string(val),
			)
		}
	}
}

func TestReapLoop(t *testing.T) {
	const interval = 20 * time.Millisecond

	cache := NewCache(interval)

	key := "https://example.com"
	cache.Add(key, []byte("testdata"))

	_, ok := cache.Get(key)
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(interval * 3)

	_, ok = cache.Get(key)
	if ok {
		t.Errorf("expected cache entry to be removed")
	}
}
