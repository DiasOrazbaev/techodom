package inmemory_test

import (
	"techodom/pkg/cache/inmemory"
	"testing"
	"time"
)

func TestInMemoryCache_Add(t *testing.T) {
	c := inmemory.NewCache(100 * time.Millisecond)
	c.Add("foo", "bar")

	if _, ok := c.Get("foo"); !ok {
		t.Errorf("Expected key 'foo' to exist in cache, but it did not")
	}
}

func TestInMemoryCache_Get(t *testing.T) {
	c := inmemory.NewCache(100 * time.Millisecond)
	c.Add("foo", "bar")

	if value, ok := c.Get("foo"); !ok || value != "bar" {
		t.Errorf("Expected value 'bar' for key 'foo', but got '%s'", value)
	}
}

func TestInMemoryCache_Get_NonExistentKey(t *testing.T) {
	c := inmemory.NewCache(100 * time.Millisecond)

	if _, ok := c.Get("foo"); ok {
		t.Errorf("Expected key 'foo' to not exist in cache, but it did")
	}
}

func TestInMemoryCache_Len(t *testing.T) {
	c := inmemory.NewCache(100 * time.Millisecond)
	c.Add("foo", "bar")

	if len := c.Len(); len != 1 {
		t.Errorf("Expected cache to have length 1, but got %d", len)
	}
}

func TestInMemoryCache_Len_EmptyCache(t *testing.T) {
	c := inmemory.NewCache(100 * time.Millisecond)

	if len := c.Len(); len != 0 {
		t.Errorf("Expected empty cache to have length 0, but got %d", len)
	}
}

func TestInMemoryCache_Expiration(t *testing.T) {
	c := inmemory.NewCache(100 * time.Millisecond)
	c.Add("foo", "bar")

	time.Sleep(200 * time.Millisecond)

	if _, ok := c.Get("foo"); ok {
		t.Errorf("Expected key 'foo' to have expired and be evicted from cache, but it still existed")
	}
}
