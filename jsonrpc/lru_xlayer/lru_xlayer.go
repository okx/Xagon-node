package lru_xlayer

import (
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
	lru "github.com/hashicorp/golang-lru"
)

var (
	once     sync.Once
	instance *LRUCache
)

// LRUCache is a simple interface for a LRU cache.
type LRUCache struct {
	cache *lru.Cache
}

func newLRUCache() *LRUCache {
	c, err := lru.New(GetConfig().Size)
	if err != nil {
		log.Fatal("Failed to create LRU cache", "err", err)
		return nil
	}
	return &LRUCache{cache: c}
}

// Get retrieves a value from the cache.
func (l *LRUCache) Get(key string) (interface{}, bool) {
	if l == nil || !GetConfig().Enable {
		return nil, false
	}

	return l.cache.Get(key)
}

// Set adds a value to the cache.
func (l *LRUCache) Set(key string, value interface{}) error {
	if l == nil || !GetConfig().Enable {
		return nil
	}
	l.cache.Add(key, value)
	return nil
}

// GetLRU returns the LRU cache instance.
func GetLRU() *LRUCache {
	once.Do(func() {
		instance = newLRUCache()
	})
	return instance
}

// Init initializes the LRU cache.
func Init(config Config) {
	SetConfig(config)
	GetLRU()
}
