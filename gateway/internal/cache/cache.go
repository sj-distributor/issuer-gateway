package cache

import (
	"cert-gateway/gateway/internal/config"
	"cert-gateway/pkg/acme"
	"crypto/tls"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"log"
)

var GlobalCache *MemoryCache

func Init(config *config.Config) {
	GlobalCache = newMemoryCache(config)
}

type MemoryCache struct {
	DB             *ristretto.Cache
	LastCertNumber uint64
	Config         *config.Config
}

func newMemoryCache(config *config.Config) *MemoryCache {
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
func (c *MemoryCache) Get(domain string) (*Cert, bool) {
	if val, b := c.DB.Get(domain); b {
		cert := val.(Cert)
		return &cert, true
	}
	return nil, false
}

// Set sets the value for the given key.
func (c *MemoryCache) Set(domain string, value Cert) bool {
	return c.DB.Set(domain, value, 1)
}

// Delete deletes the value for the given key.
func (c *MemoryCache) Delete(key string) {
	c.DB.Del(key)
}

// SetRange setRange
func (c *MemoryCache) SetRange(certs *[]Cert) error {

	for _, cert := range *certs {
		if c.LastCertNumber < cert.Id {
			c.LastCertNumber = cert.Id
		}

		certificateDecrypt, privateKeyDecrypt, _, err := acme.DecryptCertificate(cert.Certificate, cert.PrivateKey, "", c.Config.Secret)
		if err != nil {
			return err
		}
		cert.Certificate = certificateDecrypt
		cert.PrivateKey = privateKeyDecrypt

		pair, err := tls.X509KeyPair([]byte(certificateDecrypt), []byte(privateKeyDecrypt))
		if err != nil {
			log.Println(fmt.Sprintf("tls.X509KeyPair error: %s", err.Error()))
			continue
		}
		cert.TlS = pair

		c.Set(cert.Domain, cert)
	}

	return nil
}
