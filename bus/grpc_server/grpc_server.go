package grpc_server

import (
	"cert-gateway/bus/pb"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"sync"
)

// PublisherServer 实现了发布者服务
type PublisherServer struct {
	mu      sync.Mutex
	clients map[string]pb.PubSubService_SubscribeServer
	pb.PubSubServiceServer
}

// NewPublisherServer 创建一个新的发布者服务
func NewPublisherServer() *PublisherServer {
	return &PublisherServer{
		clients: make(map[string]pb.PubSubService_SubscribeServer),
	}
}

// Subscribe 实现了订阅者订阅操作
func (s *PublisherServer) Subscribe(req *pb.SubscribeRequest, stream pb.PubSubService_SubscribeServer) error {

	// 获取订阅者的唯一标识符
	subscriberID := req.SubscriberId

	// 将订阅者的流添加到订阅者列表
	s.mu.Lock()
	s.clients[subscriberID] = stream
	s.mu.Unlock()

	defer func() {
		// 在订阅者断开连接时，从订阅者列表中移除
		s.mu.Lock()
		delete(s.clients, subscriberID)
		s.mu.Unlock()
	}()

	// 此处可以添加额外的订阅逻辑，例如发送欢迎消息
	welcomeMessage := &pb.Message{
		Message: "Welcome to the Pub/Sub service!",
	}

	err := stream.Send(welcomeMessage)
	if err != nil {
		log.Printf("Error sending welcome message to subscriber %s: %v", subscriberID, err)
	}

	// 此处可以根据业务需求继续监听客户端的请求，例如接收消息或执行其他操作
	for {
		// 处理订阅者的请求，例如接收消息或执行其他操作
		select {
		case <-stream.Context().Done():
			// 订阅者断开连接时的处理
			return nil
		}
	}
}

// Publish  实现了发布操作
func (s *PublisherServer) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.Empty, error) {
	message := req.Message

	// 向所有订阅者发送消息
	for _, client := range s.clients {
		err := client.Send(&pb.Message{Message: message})
		if err != nil {
			log.Printf("Error sending message to subscriber: %v", err)
		}
	}

	return &pb.Empty{}, nil
}

// StreamInterceptor 自定义流拦截器
func StreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	md, b := metadata.FromIncomingContext(ss.Context())
	if !b {
		logx.Errorf("Metadata not found in context.")
	}
	val := md.Get("authorization")
	if len(val) == 0 {
		return errors.New(401, "UnAuthorization")
	}
	token := val[0]

	logx.Infof("bearer token is %s", token)

	err := handler(srv, ss)
	// grpc 客户端 断连提示
	log.Printf("Stream Interceptor: After method %s is called", info.FullMethod)
	return err
}

// UnaryInterceptor 自定义一元拦截器
func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Unary Interceptor: Before method %s is called", info.FullMethod)
	resp, err := handler(ctx, req)
	log.Printf("Unary Interceptor: After method %s is called", info.FullMethod)
	return resp, err
}
