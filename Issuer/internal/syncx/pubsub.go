package syncx

import (
	"cert-gateway/grpc/pb"
	"cert-gateway/issuer/internal/config"
	"cert-gateway/issuer/internal/database"
	"cert-gateway/issuer/internal/database/entity"
	"cert-gateway/pkg/driver"
	"context"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
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

	err := database.DB().Where("certificate != ''").Where("expire > ?", time.Now()).Find(&entityCerts).Order("id").Error

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
