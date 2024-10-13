package ethtxmanager

import (
	"sync"
)

const (
	capacity = 10
)

var instance *lruCache
var once sync.Once

func getAuthInstance() *lruCache {
	once.Do(func() {
		instance = newLRUCache(capacity)
	})

	return instance
}
