package cache

import (
	"cert-gateway/gateway/configs"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/swxctx/ghttp"
	"log"
)

var GlobalCache *memoryCache

func Init(config *configs.Config) {
	GlobalCache = newMemoryCache(config)
}

type memoryCache struct {
	DB             *ristretto.Cache
	LastCertNumber uint64
	Config         *configs.Config
}

func newMemoryCache(config *configs.Config) *memoryCache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		log.Fatalln(err)
	}

	return &memoryCache{
		DB:             cache,
		LastCertNumber: 0,
		Config:         config,
	}
}

// Get returns the value for the given key.
func (c *memoryCache) Get(key interface{}) (interface{}, bool) {
	return c.DB.Get(key)
}

// Set sets the value for the given key.
func (c *memoryCache) Set(key interface{}, value interface{}) bool {
	return c.DB.Set(key, value, 1)
}

// Delete deletes the value for the given key.
func (c *memoryCache) Delete(key string) {
	c.DB.Del(key)
}

// Fetch gets the certs form the server
func (c *memoryCache) Fetch() {
	req := &ghttp.Request{
		Url: fmt.Sprintf("%s/ping?last=%d", c.Config.Server, c.LastCertNumber),
	}

	req.AddHeader("Content-Type", "application/json")
	req.AddHeader("X-Client-Secret", c.Config.Server.Secret)

	resp, err := req.Do()
	if err != nil {
		log.Fatalln(err)
	}

	certs := &[]Cert{}
	err = resp.Body.FromToJson(certs)
	if err != nil {
		log.Fatalln(err)
	}
}
