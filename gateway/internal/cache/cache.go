package cache

import (
	"crypto/tls"
	"fmt"
	"github.com/pygzfei/issuer-gateway/gateway/internal/config"
	"github.com/pygzfei/issuer-gateway/pkg/acme"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

var GlobalCache *MemoryCache

func Init(config *config.Config) {
	GlobalCache = newMemoryCache(config)
}

type MemoryCache struct {
	DB             sync.Map
	LastCertNumber uint64
	Config         *config.Config
}

func newMemoryCache(config *config.Config) *MemoryCache {
	return &MemoryCache{
		DB:             sync.Map{},
		LastCertNumber: 0,
		Config:         config,
	}
}

// Get returns the value for the given key.
func (c *MemoryCache) Get(domain string) (*Cert, bool) {
	if val, b := c.DB.Load(domain); b {
		cert := val.(Cert)
		return &cert, true
	}
	return nil, false
}

// Set sets the value for the given key.
func (c *MemoryCache) Set(domain string, value Cert) bool {
	c.DB.Store(domain, value)
	return true
}

// Delete deletes the value for the given key.
func (c *MemoryCache) Delete(key string) {
	c.DB.Delete(key)
}

// SetRange setRange
func (c *MemoryCache) SetRange(certs *[]Cert) error {

	logx.Debugw("accept certs", logx.Field("certs number", len(*certs)))

	for _, cert := range *certs {
		if cert.PrivateKey == "" || cert.Certificate == "" && cert.Id != 0 {
			c.Set(cert.Domain, cert)
			continue
		}
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
			logx.Errorf("tls.X509KeyPair error: %s", err)
			continue
		}
		cert.TlS = pair

		fmt.Println("set cert domain :", cert.Domain)

		c.Set(cert.Domain, cert)
	}

	return nil
}
