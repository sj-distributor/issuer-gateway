package grpc_server

import (
	"cert-gateway/grpc/pb"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"sync"
)

type CertificatePubSubServer struct {
	cache    *MemoryCache
	mu       sync.Mutex
	gateways map[string]pb.CertificateService_GatewaySubscribeServer
	pb.CertificateServiceServer
}

func NewCertificatePubSubServer() *CertificatePubSubServer {
	return &CertificatePubSubServer{
		cache:    MewMemoryCache(),
		gateways: make(map[string]pb.CertificateService_GatewaySubscribeServer),
	}
}

// SyncCertificateToProvider 发送证书给provider
func (s *CertificatePubSubServer) SyncCertificateToProvider(_ context.Context, certs *pb.CertificateList) (*pb.Empty, error) {
	s.cache.SetRange(certs)
	return &pb.Empty{}, nil
}

// SendCertificateToGateway 发送证书同步给某个 Gateway
func (s *CertificatePubSubServer) SendCertificateToGateway(_ context.Context, req *pb.SubscribeRequest) (*pb.Empty, error) {
	gatewayIp := req.LocalIp
	
	s.mu.Lock()
	stream := s.gateways[gatewayIp]
	s.mu.Unlock()

	if stream == nil {
		return &pb.Empty{}, errors.New(404, fmt.Sprintf("gateway ip: %s,  not found", gatewayIp))
	}
	err := stream.Send(s.cache.GetAll())

	if err != nil {
		log.Printf(fmt.Sprintf("Error sending message to Gateway(%s): %v", gatewayIp, err))
	}

	return &pb.Empty{}, nil
}

// GatewaySubscribe Gateway订阅
func (s *CertificatePubSubServer) GatewaySubscribe(req *pb.SubscribeRequest, stream pb.CertificateService_GatewaySubscribeServer) error {
	// 获取gateway的ip作为唯一标识符
	localIp := req.LocalIp

	// 将订阅者的流添加到订阅者列表
	s.mu.Lock()
	s.gateways[localIp] = stream
	s.mu.Unlock()

	defer func() {
		// 在订阅者断开连接时，从订阅者列表中移除
		s.mu.Lock()
		delete(s.gateways, localIp)
		s.mu.Unlock()
	}()

	log.Printf("Subscription from %s ...\n", localIp)

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
