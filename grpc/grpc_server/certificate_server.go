package grpc_server

import (
	"context"
	"fmt"
	conf "github.com/pygzfei/issuer-gateway/grpc/config"
	"github.com/pygzfei/issuer-gateway/grpc/pb"
	"github.com/pygzfei/issuer-gateway/pkg/errs"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
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
		mu:       sync.Mutex{},
	}
}

func (s *CertificatePubSubServer) Check(context.Context, *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// SyncCertificateToProvider 发送证书给provider
func (s *CertificatePubSubServer) SyncCertificateToProvider(_ context.Context, certs *pb.CertificateList) (*pb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache.SetRange(certs)

	for gatewayIp, stream := range s.gateways {
		err := stream.Send(certs)

		if err != nil {
			logx.Errorf("Error sending message to GatewayIP(%s): %v", gatewayIp, err)
		}
	}

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
		logx.Errorf("Error sending message to Gateway(%s): %v", gatewayIp, err)
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

	logx.Debugf("Subscription from %s ...", localIp)

	// 此处可以根据业务需求继续监听客户端的请求，例如接收消息或执行其他操作
	for {
		// 处理订阅者的请求，例如接收消息或执行其他操作
		select {
		case <-stream.Context().Done():
			// 订阅者断开连接时的处理
			logx.Errorw("Subscriber disconnected", logx.Field("localIp", localIp))
			break
		}
	}
}

// StreamInterceptor streamInterceptor
func StreamInterceptor(conf *conf.Config) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		token, err := getTokenFromCtx(ss.Context())
		if err != nil {
			return err
		}
		logx.Debugf("bearer token is %s", *token)

		if *token != conf.Secret {
			return errs.UnAuthorizationException
		}

		err = handler(srv, ss)
		// grpc 客户端 断连提示
		logx.Error("Stream Interceptor: Executed error", err)
		return err
	}
}

// UnaryInterceptor unaryInterceptor
func UnaryInterceptor(conf *conf.Config) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logx.Debugf("Unary Interceptor: Before method %s is called", info.FullMethod)

		token, err := getTokenFromCtx(ctx)
		if err != nil {
			return nil, err
		}
		logx.Debugf("bearer token is %s", *token)

		if *token != conf.Secret {
			return nil, errs.UnAuthorizationException
		}

		resp, err := handler(ctx, req)

		logx.Debugf("Unary Interceptor: After method %s is called", info.FullMethod)
		return resp, err
	}
}

func getTokenFromCtx(ctx context.Context) (token *string, err error) {
	md, b := metadata.FromIncomingContext(ctx)
	if !b {
		logx.Errorf("Metadata not found in context.")
	}
	val := md.Get("authorization")
	if len(val) == 0 {
		return nil, errs.UnAuthorizationException
	}

	replace := strings.TrimPrefix(val[0], "Bearer ")

	return &replace, nil
}
