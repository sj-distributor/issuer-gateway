package cache

import (
	"cert-gateway/gateway/internal/config"
	"cert-gateway/pkg/acme"
	"crypto/tls"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/go-jose/go-jose/v3/json"
	"io"
	"log"
	"net/http"
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

func (c *MemoryCache) Sync() error {
	// 发起 GET 请求
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%d", c.Config.Server.Url, c.LastCertNumber), nil)
	if err != nil {
		log.Println("GET 请求错误:", err)
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Config.Server.Secret))
	// 创建 HTTP 客户端
	client := &http.Client{}
	// 发送请求
	response, err := client.Do(req)

	defer response.Body.Close()

	// 读取响应内容
	var resp Resp
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("读取响应内容错误:", err)
		return err
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Println("解析响应内容错误:", err)
		return err
	}
	log.Println("GET 响应内容:")
	log.Println(resp)
	for _, cert := range resp.Data.Certs {
		if c.LastCertNumber < cert.Id {
			c.LastCertNumber = cert.Id
		}

		certificateDecrypt, privateKeyDecrypt, _, err := acme.DecryptCertificate(cert.Certificate, cert.PrivateKey, "", c.Config.Server.Secret)
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
