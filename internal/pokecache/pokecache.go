package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	m       *sync.Mutex
}

type cacheEntry struct {
	craetedAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		entries: map[string]cacheEntry{},
		m:       &sync.Mutex{},
	}

	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			c.reapLoop()
		}
	}()

	return c
}

func (c Cache) Add(key string, val []byte) {
	c.m.Lock()
	defer c.m.Unlock()

	newEntry := cacheEntry{
		craetedAt: time.Now(),
		val:       val,
	}
	// fmt.Println("\n[POKECACHE]: adding to", key)
	// fmt.Println()
	// fmt.Print("◓ > ")
	c.entries[key] = newEntry
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	entry, ok := c.entries[key]
	return entry.val, ok
}

func (c Cache) reapLoop() {
	c.m.Lock()
	defer c.m.Unlock()

	for key, entry := range c.entries {
		diff := entry.craetedAt.Sub(time.Now())
		if diff < -5*time.Second {
			// fmt.Println("\n[POKECACHE]: removing", key)
			// fmt.Println()
			// fmt.Print("◓ > ")
			delete(c.entries, key)
		}
	}
}
