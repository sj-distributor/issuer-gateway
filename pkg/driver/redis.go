package driver

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

// RedisClient redis client
type RedisClient struct {
	redis          *redis.Client
	gatewayChannel string
	certChannel    string
}

// NewRedisClient new a redis client
func NewRedisClient(addr, pass string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass, // no password set
		DB:       db,   // use default DB
	})
	return &RedisClient{
		redis:          rdb,
		gatewayChannel: "gateway-sync-pub-sub",
		certChannel:    "cert-sync-pub-sub",
	}
}

func (r *RedisClient) Subscribe(_ string, received OnMessageReceived, receiving ...OnErrReceiving) error {
	ctx := context.Background()
	// 创建一个订阅者
	pubsub := r.redis.Subscribe(ctx, r.gatewayChannel)

	// 创建一个通道来接收订阅的消息
	subChannel := pubsub.Channel()

	// 接收订阅的消息
	for msg := range subChannel {
		log.Printf("收到消息: %s\n", msg.Payload)
		received(msg.Payload)
	}

	return nil
}

func (r *RedisClient) Publish(msg string) error {
	ctx := context.Background()
	err := r.redis.Publish(ctx, r.gatewayChannel, msg).Err()
	if err != nil {
		log.Printf("发布消息失败:%s\n", err)
		return err
	}
	log.Printf("已发布消息: %s\n", msg)
	return nil
}
