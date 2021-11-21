package service

import (
	"sync"
	"time"
)

type Challenge string

type CacheEntry struct {
	Created   time.Time
	Challenge Challenge
}

type Cache struct {
	entries map[string]CacheEntry
	mu      sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]CacheEntry),
	}
}

func (c *Cache) Set(key string, ch Challenge) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = CacheEntry{
		Created:   time.Now(),
		Challenge: ch,
	}
}

func (c *Cache) Get(key string) (e CacheEntry, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok = c.entries[key]
	return e, ok
}

func (c *Cache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}
