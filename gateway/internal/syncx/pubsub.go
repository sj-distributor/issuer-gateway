package syncx

import (
	"context"
	"fmt"
	"github.com/pygzfei/issuer-gateway/gateway/internal/cache"
	"github.com/pygzfei/issuer-gateway/gateway/internal/config"
	"github.com/pygzfei/issuer-gateway/grpc/pb"
	"github.com/pygzfei/issuer-gateway/pkg/driver"
	"github.com/pygzfei/issuer-gateway/utils"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"time"
)

var GlobalPubSub driver.IProvider

func Init(c *config.Config) {
	podId := os.Getenv("podId")
	if podId == "" {
		id := utils.GetLocalIP()
		podId = id
	}
	switch c.Sync.Target {
	case driver.GRPC:
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("Authorization", "Bearer "+c.Secret))
		GlobalPubSub = driver.NewGrpcClient(c.Sync.GrpcClient.Listen, ctx)
		err := GlobalPubSub.GatewaySubscribe(podId, handlerMessage, func(err error) {
			log.Println(err)
		})
		if err != nil {
			log.Panicln(fmt.Sprintf("Grpc init fail: %s", err.Error()))
		}
		tick := time.Tick(time.Second * 1)
		select {
		case <-tick:
			err = GlobalPubSub.SendCertificateToGateway(podId)
			if err != nil {
				log.Panicln(fmt.Sprintf("Grpc init fail: %s", err.Error()))
			}
		}

	case driver.REDIS:
		redis := c.Sync.Redis

		GlobalPubSub = driver.NewRedisClient(redis.Addrs, redis.User, redis.Pass, redis.MasterName, redis.Db)
		err := GlobalPubSub.GatewaySubscribe(podId, handlerMessage, func(err error) {
			log.Println(err)
		})
		if err != nil {
			log.Panicln(fmt.Sprintf("Redis GatewaySubscribe fail: %s", err.Error()))
		}

		err = GlobalPubSub.SendCertificateToGateway(podId)
		if err != nil {
			log.Panicln(fmt.Sprintf("Redis SendCertificateToGateway fail: %s", err.Error()))
		}

	}

}

func handlerMessage(list []*pb.Cert) {

	var certs []cache.Cert

	for _, cert := range list {
		certs = append(certs, cache.Cert{
			Id:          cert.Id,
			PrivateKey:  cert.PrivateKey,
			Certificate: cert.Certificate,
			Domain:      cert.Domain,
			Target:      cert.Target,
		})
	}

	err := cache.GlobalCache.SetRange(&certs)
	if err != nil {
		log.Println(err)
	}
}
