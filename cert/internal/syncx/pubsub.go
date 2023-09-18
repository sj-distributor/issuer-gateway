package syncx

import (
	"cert-gateway/cert/internal/config"
	"cert-gateway/pkg/driver"
	"context"
	"google.golang.org/grpc/metadata"
	"log"
)

var GlobalPubSub driver.IPubSubDriver

func Init(c *config.Config) {

	switch c.Sync.Type {
	case driver.GRPC:
		md := metadata.Pairs("Authorization", "Bearer "+c.Secret)
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		GlobalPubSub = driver.NewGrpcClient(c.Sync.Address, ctx)
		break
	case driver.REDIS:
		GlobalPubSub = driver.NewRedisClient(c.Sync.Address, c.Sync.Pass, 0)
		err := GlobalPubSub.Publish("redis channel init...")
		if err != nil {
			log.Panicln("redis init fail: ", err)
		}
		break
	case driver.AMQP:
		break
	case driver.ETCD:
		break
	}

}
