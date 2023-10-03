package syncx

import (
	"context"
	"github.com/pygzfei/issuer-gateway/grpc/pb"
	"github.com/pygzfei/issuer-gateway/issuer/internal/config"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database/entity"
	"github.com/pygzfei/issuer-gateway/pkg/driver"
	"google.golang.org/grpc/metadata"
	"log"
)

func Init(c *config.Config) driver.IProvider {

	var provider driver.IProvider

	switch c.Sync.Target {
	case driver.GRPC:
		ctx := metadata.NewOutgoingContext(
			context.Background(),
			metadata.Pairs("Authorization", "Bearer "+c.Secret))

		provider = driver.NewGrpcClient(c.Sync.GrpcClient.Listen, ctx)
		break
	case driver.REDIS:
		redis := c.Sync.Redis
		provider = driver.NewRedisClient(redis.Addrs, redis.User, redis.Pass, redis.MasterName, redis.Db)
		break
	}

	go setUpCertificate(provider)

	return provider
}

// setUpCertificate set up certificate
func setUpCertificate(provider driver.IProvider) {

	var entityCerts []entity.Cert

	err := database.DB().Find(&entityCerts).Order("id").Error

	if err != nil {
		log.Panicf("set up certificate error: %s", err)
	}

	var certs []*pb.Cert

	for _, cert := range entityCerts {
		temp := &pb.Cert{
			Id:                cert.Id,
			PrivateKey:        cert.PrivateKey,
			Certificate:       cert.Certificate,
			Domain:            cert.Domain,
			Target:            cert.Target,
			IssuerCertificate: cert.IssuerCertificate,
		}
		certs = append(certs, temp)
	}

	err = provider.SyncCertificateToProvider(&pb.CertificateList{
		Certs: certs,
	})

	if err != nil {
		log.Panicf("set up certificate error: %s", err)
	}

}
