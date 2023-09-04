package cache

import (
	"cert-gateway/gateway/configs"
	"github.com/dgraph-io/ristretto"
	"log"
)

var GlobalCache *MemoryCache

func Init(config *configs.Config) {
	GlobalCache = newMemoryCache(config)
}

type MemoryCache struct {
	DB             *ristretto.Cache
	LastCertNumber uint64
	Config         *configs.Config
}

func newMemoryCache(config *configs.Config) *MemoryCache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		log.Fatalln(err)
	}

	return &MemoryCache{
		DB:             cache,
		LastCertNumber: 0,
		Config:         config,
	}
}

// Get returns the value for the given key.
func (c *MemoryCache) Get(key interface{}) (interface{}, bool) {
	return c.DB.Get(key)
}

// Set sets the value for the given key.
func (c *MemoryCache) Set(key interface{}, value interface{}) bool {
	return c.DB.Set(key, value, 1)
}

// Delete deletes the value for the given key.
func (c *MemoryCache) Delete(key string) {
	c.DB.Del(key)
}
