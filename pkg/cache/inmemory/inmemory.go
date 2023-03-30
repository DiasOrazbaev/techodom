package inmemory

import (
	"sync"
	"time"
)

type Cache interface {
	Add(key, value string)
	Get(key string) (value string, ok bool)
	Len() int
}

type cacheEntry struct {
	value string
	timer *time.Timer
}

type inMemoryCache struct {
	entries map[string]*cacheEntry
	mu      sync.RWMutex
	timeout time.Duration
}

func NewCache(timeout time.Duration) Cache {
	return &inMemoryCache{
		entries: make(map[string]*cacheEntry, 1000),
		timeout: timeout,
	}
}

func (c *inMemoryCache) addWithLock(key, value string) {
	entry := &cacheEntry{
		value: value,
		timer: time.AfterFunc(c.timeout, func() {
			c.mu.Lock()
			defer c.mu.Unlock()
			delete(c.entries, key)
		}),
	}
	c.entries[key] = entry
}

func (c *inMemoryCache) Add(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.addWithLock(key, value)
}

func (c *inMemoryCache) Get(key string) (value string, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	if ok {
		value = entry.value
	}
	return
}

func (c *inMemoryCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
