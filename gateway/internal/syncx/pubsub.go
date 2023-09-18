package syncx

import (
	"cert-gateway/gateway/internal/cache"
	"cert-gateway/gateway/internal/config"
	"cert-gateway/pkg/driver"
	"cert-gateway/utils"
	"context"
	"fmt"
	"github.com/go-jose/go-jose/v3/json"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
)

var GlobalPubSub driver.IPubSubDriver

func Init(c *config.Config) {
	podId := os.Getenv("podId")
	if podId == "" {
		id := utils.GetLocalId()
		podId = id
	}
	switch c.Sync.Type {
	case driver.GRPC:
		md := metadata.Pairs("Authorization", "Bearer "+c.Secret)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		GlobalPubSub = driver.NewGrpcClient(c.Sync.Address, ctx)

		err := GlobalPubSub.Subscribe(podId, handlerMessage, func(err error) {
			log.Println(err)
		})
		if err != nil {
			log.Panicln(fmt.Sprintf("Grpc init fail: %s", err.Error()))
		}
		break
	case driver.REDIS:
		GlobalPubSub = driver.NewRedisClient(c.Sync.Address, c.Sync.Pass, 0)
		err := GlobalPubSub.Subscribe(podId, handlerMessage, func(err error) {
			log.Println(err)
		})

		if err != nil {
			log.Panicln(fmt.Sprintf("redis init fail: %s", err.Error()))
		}
		break
	case driver.AMQP:
		break
	case driver.ETCD:
		break
	}
}

func handlerMessage(msg string) {

	var certs []cache.Cert

	err := json.Unmarshal([]byte(msg), &certs)
	if err != nil {
		log.Println(err)
	}

	err = cache.GlobalCache.SetRange(&certs)
	if err != nil {
		log.Println(err)
	}
}
