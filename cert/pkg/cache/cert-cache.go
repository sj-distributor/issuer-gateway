package cache

import (
	"cert-gateway/cert/internal/config"
	"cert-gateway/cert/internal/database/entity"
	"cert-gateway/cert/internal/types"
	"gorm.io/gorm"
	"sync"
)

// CertCache types.Cert 缓存
var CertCache *MemoryCache

func Init(db *gorm.DB, c *config.Config) {
	CertCache = newMemoryCache(db, c)
}

type MemoryCache struct {
	cache         *sync.Map
	maximum       uint64
	Database      *gorm.DB
	config        *config.Config
	isAllUpdating bool
}

func newMemoryCache(db *gorm.DB, c *config.Config) *MemoryCache {
	return &MemoryCache{
		cache:         &sync.Map{},
		maximum:       0,
		Database:      db,
		config:        c,
		isAllUpdating: false,
	}
}

// Get returns the value for the given key.
func (c *MemoryCache) Get(key interface{}) (interface{}, bool) {
	return c.cache.Load(key)
}

// Set sets the value for the given key.
func (c *MemoryCache) Set(key interface{}, value interface{}) {
	c.cache.Store(key, value)
}

// Delete deletes the value for the given key.
func (c *MemoryCache) Delete(key string) {
	c.cache.Delete(key)
}

// CaptureChanged 获取增量数据
func (c *MemoryCache) CaptureChanged(maximum uint64) *[]types.Cert {
	certs := []types.Cert{}

	if c.isAllUpdating == true {
		return &certs
	}

	c.cache.Range(func(key, value any) bool {
		u := key.(uint64)
		if maximum > u {
			cert := value.(types.Cert)
			certs = append(certs, cert)
		}
		if u > c.maximum {
			c.maximum = u
		}
		return true
	})
	return &certs
}

func (c *MemoryCache) SetupCacheFormDatabase() error {
	c.isAllUpdating = true
	defer func() {
		c.isAllUpdating = false
	}()

	var certs []entity.Cert

	_ = c.Database.Where("id > ?", c.maximum).Find(&certs)

	for _, cert := range certs {
		c.Set(cert.Id, types.Cert{
			Domain:      cert.Domain,
			Certificate: cert.Certificate,
			PrivateKey:  cert.PrivateKey,
		})
	}

	return nil
}

func (c *MemoryCache) name() {

}
