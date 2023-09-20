package grpc_server

import (
	"cert-gateway/bus/pb"
	"sync"
)

type MemoryCache struct {
	cache *sync.Map
}

func MewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cache: &sync.Map{},
	}
}

// Get returns the value for the given key.
func (c *MemoryCache) Get(id uint64) (*pb.Cert, bool) {
	if value, ok := c.cache.Load(id); ok {
		cert := value.(*pb.Cert)
		return cert, true
	}
	return nil, false
}

// Set sets the value for the given key.
func (c *MemoryCache) Set(id uint64, value *pb.Cert) {
	c.cache.Store(id, value)
}

func (c *MemoryCache) SetRange(certificateList *pb.CertificateList) {
	for _, cert := range certificateList.Certs {
		temp := &pb.Cert{
			Id:          cert.Id,
			Domain:      cert.Domain,
			Target:      cert.Target,
			Certificate: cert.Certificate,
			PrivateKey:  cert.PrivateKey,
		}
		c.Set(temp.Id, temp)
	}
}

func (c *MemoryCache) GetAll() *pb.CertificateList {
	result := &pb.CertificateList{
		Certs: []*pb.Cert{},
	}
	c.cache.Range(func(_, value any) bool {
		cert := value.(*pb.Cert)
		result.Certs = append(result.Certs, cert)
		return true
	})
	return result
}
