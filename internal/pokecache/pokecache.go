package pokecache

import (
	"time"
	"sync"
	"fmt"
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
	// key is url for request
	// print("starting addition\n")
	c.mux.Lock()
	defer c.mux.Unlock()
	// defer print("mutex unlocked\n\n")
	// print("mutex locked\n")
	if entry, exists := c.entry[key]; exists {
		// reset timer
		// print("entry exists\n")
		entry.createdAt = time.Now()
		c.entry[key] = entry
		return
	}
	// print("entry does not exist\n")
	c.entry[key] = cacheEntry{createdAt: time.Now(), val: valBytes}
	// print("checking val\n")
	// if givenVal, exists := c.Get(key); exists {
	// 	print("value: ", string(givenVal), "\n")
	// }
	// print("ending addition\n")
}

// gets an entry from the cache. 
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	entry, exists := c.entry[key]
	if !exists {
		return []byte{}, false
	}
	return entry.val, true
}

// called when the cache is created
func (c *Cache) reapLoop(interval time.Duration) {
// Each time an interval (the time.Duration passed to NewCache) passes it should remove any entries that are older than the interval. This makes sure that the cache doesn't grow too large over time. For example, if the interval is 5 seconds, and an entry was added 7 seconds ago, that entry should be removed.
	// for entries in cache, if createdAt - now > duration, then remove
	for key, entry := range c.entry {
		if time.Since(entry.createdAt) > interval {
			delete(c.entry, key)
			fmt.Printf("Entry deleted: %s", key)
		}
	}
}