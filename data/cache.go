package data

import (
	"sync"
	"time"
)

type Item struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	items map[string]Item
	mu    sync.RWMutex
}

func (item Item) isExpired() bool {
	return time.Now().UnixNano() > item.Expiration
}

func (c *Cache) Set(k string, x interface{}, d time.Duration) {
	c.mu.Lock()
	c.items[k] = Item{
		Value:      x,
		Expiration: time.Now().Add(d).UnixNano(),
	}
	c.mu.Unlock()
}

func (c *Cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()

	item, found := c.items[k]
	c.mu.RUnlock()
	if !found {
		return "", false
	}
	if item.isExpired() {
		return "", false
	}
	return item.Value, true
}

func CreateCache() *Cache {
	c := &Cache{
		items: make(map[string]Item),
	}
	return c
}
