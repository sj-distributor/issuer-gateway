package driver

import (
	"fmt"
	"github.com/pygzfei/issuer-gateway/grpc/grpc_server"
	"github.com/pygzfei/issuer-gateway/grpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewGrpcServiceAndListen(addr string) {
	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer(
		// 注册一元拦截器
		grpc.UnaryInterceptor(grpc_server.UnaryInterceptor),
		// 注册流拦截器
		grpc.StreamInterceptor(grpc_server.StreamInterceptor),
	)

	certificatePubSubServer := grpc_server.NewCertificatePubSubServer()

	// 监听 gRPC 服务端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", addr))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 注册发布者服务到 gRPC 服务器
	pb.RegisterCertificateServiceServer(grpcServer, certificatePubSubServer)

	logx.Infof("Grpc Server is listening on : %s", addr)
	// 启动 gRPC 服务器
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
