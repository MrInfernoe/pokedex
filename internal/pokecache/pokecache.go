package pokecache

import (
	"time"
	"sync"
	// "fmt"
)

type cacheEntry struct {
	createdAt 	time.Time
	val 		[]byte
}

type Cache struct {
	mux 	*sync.Mutex
	entry 	map[string]cacheEntry
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{mux: &sync.Mutex{}, entry: make(map[string]cacheEntry)}
	go func(){
		ticker := time.NewTicker(interval)
		for range ticker.C {
			cache.reapLoop(interval)
		}
	}()
	return &cache
}

// adds a new entry to the cache.
func (c *Cache) Add(key string, valBytes []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if entry, exists := c.entry[key]; exists {
		entry.createdAt = time.Now()
		c.entry[key] = entry
		return
	}
	c.entry[key] = cacheEntry{createdAt: time.Now(), val: valBytes}
}

// gets an entry from the cache. 
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	entry, exists := c.entry[key]
	if !exists {
		return []byte{}, false
	}
	entry.createdAt = time.Now()
	c.entry[key] = entry
	return entry.val, true
}

// goroutine spawned when the cache is created
// for entries in cache, if createdAt - now > duration, then remove
func (c *Cache) reapLoop(interval time.Duration) {
	for key, entry := range c.entry {
		if time.Since(entry.createdAt) > interval {
			delete(c.entry, key)
			// fmt.Printf("Entry deleted: %s", key)
		}
	}
}