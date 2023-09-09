// main.go
package main

import (
	"cert-gateway/bus/pb"
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
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

func (s *PublisherServer) mustEmbedUnimplementedPubSubServiceServer() {}

func main() {
	// 监听 gRPC 服务端口
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 创建发布者服务
	publisher := NewPublisherServer()

	// 注册发布者服务到 gRPC 服务器
	pb.RegisterPubSubServiceServer(server, publisher)

	fmt.Println("Server is listening on :50051...")
	// 启动 gRPC 服务器
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
