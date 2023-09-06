package cache

import (
	"cert-gateway/cert/internal/config"
	"cert-gateway/cert/internal/database/entity"
	"cert-gateway/cert/internal/types"
	"gorm.io/gorm"
	"sync"
	"time"
)

// CertCache types.Cert 缓存
var CertCache *MemoryCache

func Init(db *gorm.DB, c *config.Config) error {
	CertCache = newMemoryCache(db, c)
	err := CertCache.ReLoad(db)
	if err != nil {
		return err
	}
	return nil
}

type MemoryCache struct {
	cache    *sync.Map
	Database *gorm.DB
	config   *config.Config
	lastId   uint64
}

func newMemoryCache(db *gorm.DB, c *config.Config) *MemoryCache {
	return &MemoryCache{
		cache:    &sync.Map{},
		Database: db,
		config:   c,
		lastId:   0,
	}
}

// Get returns the value for the given key.
func (c *MemoryCache) Get(key uint64) (*types.Cert, bool) {
	if value, ok := c.cache.Load(key); ok {
		cert := value.(types.Cert)
		return &cert, true
	}
	return nil, false
}

// Set sets the value for the given key.
func (c *MemoryCache) Set(key uint64, value types.Cert) {
	if value.Id > c.lastId {
		c.lastId = value.Id
	}
	c.cache.Store(key, value)
}

// Delete deletes the value for the given key.
func (c *MemoryCache) Delete(key uint64) {
	c.cache.Delete(key)
}

// CaptureChanged 获取增量数据
func (c *MemoryCache) CaptureChanged(maximum uint64) *[]types.Cert {
	var certs []types.Cert

	c.cache.Range(func(key, value any) bool {
		u := key.(uint64)
		if u > maximum {
			cert := value.(types.Cert)
			certs = append(certs, cert)
		}
		return true
	})
	return &certs
}

func (c *MemoryCache) SetRange(certs *[]entity.Cert) {

	for _, cert := range *certs {

		if cert.Id > c.lastId {
			c.lastId = cert.Id
		}

		c.Set(cert.Id, types.Cert{
			Id:          cert.Id,
			Domain:      cert.Domain,
			Target:      cert.Target,
			Certificate: cert.Certificate,
			PrivateKey:  cert.PrivateKey,
		})
	}

}

func (c *MemoryCache) ReLoad(db *gorm.DB) error {
	var entityCerts []entity.Cert

	err := db.Where("expire > ?", time.Now()).Find(&entityCerts).Order("id").Error

	if err != nil {
		return err
	}

	c.SetRange(&entityCerts)

	return nil
}

func (c *MemoryCache) PartialUpdate(db *gorm.DB) error {
	var entityCerts []entity.Cert

	err := db.Where("id > ?", c.lastId).Find(&entityCerts).Order("id").Error

	if err != nil {
		return err
	}

	c.SetRange(&entityCerts)

	return nil
}
