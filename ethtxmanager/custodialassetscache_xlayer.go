package ethtxmanager

import (
	"container/list"
)

// lruCache represents a cache with a fixed size
type lruCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
}

// entry is used to store key-value pairs in the cache
type entry struct {
	key   string
	value string
}

// newLRUCache creates a new LRU cache with the given capacity
func newLRUCache(capacity int) *lruCache {
	return &lruCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

// get retrieves the value for the given key if present, otherwise returns -1
func (l *lruCache) get(key string) string {
	if element, ok := l.cache[key]; ok {
		// Move the accessed element to the front of the list
		l.list.MoveToFront(element)
		// Return the value
		return element.Value.(*entry).value
	}
	// Key not found
	return ""
}

// put adds a key-value pair to the cache or updates the value if the key already exists
func (l *lruCache) put(key string, value string) {
	// Check if the key already exists
	if element, ok := l.cache[key]; ok {
		// Update the value
		element.Value.(*entry).value = value
		// Move the element to the front of the list
		l.list.MoveToFront(element)
		return
	}

	// If the cache has reached capacity, remove the least recently used item
	if l.list.Len() == l.capacity {
		// Get the least recently used element (tail of the list)
		lru := l.list.Back()
		if lru != nil {
			// Remove it from the list and the cache
			l.list.Remove(lru)
			delete(l.cache, lru.Value.(*entry).key)
		}
	}

	// Add the new item to the front of the list
	element := l.list.PushFront(&entry{key: key, value: value})
	// Add it to the cache
	l.cache[key] = element
}
