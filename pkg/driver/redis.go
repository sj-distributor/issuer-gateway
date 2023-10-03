package driver

import (
	"context"
	"fmt"
	"github.com/go-jose/go-jose/v3/json"
	"github.com/pygzfei/issuer-gateway/grpc/pb"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"sync"
)

// RedisClient redis client
type RedisClient struct {
	mu                    sync.Mutex
	redis                 redis.UniversalClient
	globalGatewayChannel  string
	selfChannel           string
	certificateStorageKey string
}

// NewRedisClient new a redis client
func NewRedisClient(addr []string, user, pass, masterName string, db int) *RedisClient {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:      addr,
		Username:   user,
		Password:   pass,
		DB:         db,
		MasterName: masterName,
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("redis init failed: %s", err)
	}
	return &RedisClient{
		mu:                    sync.Mutex{},
		redis:                 rdb,
		globalGatewayChannel:  "global-gateway-sync-channel",
		selfChannel:           "self-gateway-sync-",
		certificateStorageKey: "issuer-gateway-certs-key",
	}
}

func (r *RedisClient) GatewaySubscribe(localIP string, received OnMessageReceived, receiving ...OnErrReceiving) error {
	ctx := context.Background()
	// 创建一个订阅者
	global := r.redis.Subscribe(ctx, r.globalGatewayChannel)

	self := r.redis.Subscribe(ctx, fmt.Sprintf("%s%s", r.selfChannel, localIP))

	// 创建一个通道来接收订阅的消息
	globalChannel := global.Channel()
	selfChannel := self.Channel()

	go func() {
		// 接收订阅的消息
		for msg := range selfChannel {
			logx.Infof("selfChannel received: %s", msg.Payload)
			var certs []*pb.Cert
			err := json.Unmarshal([]byte(msg.Payload), &certs)
			if err != nil {
				logx.Infof("selfChannel sync failed: %s", err)
				if len(receiving) > 0 {
					receiving[0](err)
				}
				continue
			}
			received(certs)
		}
	}()

	go func() {
		// 接收订阅的消息
		for msg := range globalChannel {
			logx.Infof("globalChannel received: %s", msg.Payload)
			var certs []*pb.Cert
			err := json.Unmarshal([]byte(msg.Payload), &certs)
			if err != nil {
				logx.Infof("globalChannel sync failed: %s", err)
				if len(receiving) > 0 {
					receiving[0](err)
				}
				continue
			}
			received(certs)
		}
	}()

	return nil
}

func (r *RedisClient) SyncCertificateToProvider(certificateList *pb.CertificateList) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ctx := context.Background()

	stringCmd := r.redis.Get(ctx, r.certificateStorageKey)

	s := sync.Map{}
	var certs []*pb.Cert

	// 处理已存在的证书
	if stringCmd.Val() != "" {
		err := json.Unmarshal([]byte(stringCmd.Val()), &certs)
		if err != nil {
			return err
		}

		for _, cert := range certs {
			temp := *cert
			s.Store(cert.Id, &temp)
		}

		certs = []*pb.Cert{}
	}

	// 处理新增证书
	for _, cert := range certificateList.Certs {
		temp := *cert
		s.Store(cert.Id, &temp)
	}

	s.Range(func(_, value any) bool {
		temp := value.(*pb.Cert)
		certs = append(certs, temp)
		return true
	})

	marshal, err := json.Marshal(certs)
	if err != nil {
		return err
	}
	err = r.redis.Set(ctx, r.certificateStorageKey, marshal, 0).Err()

	if err != nil {
		return err
	}

	// 同步给所有节点
	err = r.redis.Publish(ctx, r.globalGatewayChannel, marshal).Err()
	if err != nil {
		logx.Infof("send globalChannel message failed: %s", err)
		return err
	}

	return nil
}

func (r *RedisClient) SendCertificateToGateway(localIP string) error {
	ctx := context.Background()
	stringCmd := r.redis.Get(ctx, r.certificateStorageKey)
	if stringCmd.Err() != nil {
		return stringCmd.Err()
	}

	if stringCmd.Val() != "" || stringCmd.Err() == nil {
		err := r.redis.Publish(ctx, fmt.Sprintf("%s%s", r.selfChannel, localIP), stringCmd.Val()).Err()
		if err != nil {
			return err
		}
	}

	return stringCmd.Err()
}
