// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.14.0
// source: pubsub.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CertificateServiceClient is the client API for CertificateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CertificateServiceClient interface {
	// 发送证书同步给某个 Gateway
	SendCertificateToGateway(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (*Empty, error)
	// Issuer发送证书给Provider
	SyncCertificateToProvider(ctx context.Context, in *CertificateList, opts ...grpc.CallOption) (*Empty, error)
	// Gateway订阅
	GatewaySubscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (CertificateService_GatewaySubscribeClient, error)
	Check(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type certificateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCertificateServiceClient(cc grpc.ClientConnInterface) CertificateServiceClient {
	return &certificateServiceClient{cc}
}

func (c *certificateServiceClient) SendCertificateToGateway(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.CertificateService/SendCertificateToGateway", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *certificateServiceClient) SyncCertificateToProvider(ctx context.Context, in *CertificateList, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.CertificateService/SyncCertificateToProvider", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *certificateServiceClient) GatewaySubscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (CertificateService_GatewaySubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &CertificateService_ServiceDesc.Streams[0], "/pb.CertificateService/GatewaySubscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &certificateServiceGatewaySubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CertificateService_GatewaySubscribeClient interface {
	Recv() (*CertificateList, error)
	grpc.ClientStream
}

type certificateServiceGatewaySubscribeClient struct {
	grpc.ClientStream
}

func (x *certificateServiceGatewaySubscribeClient) Recv() (*CertificateList, error) {
	m := new(CertificateList)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *certificateServiceClient) Check(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.CertificateService/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CertificateServiceServer is the server API for CertificateService service.
// All implementations must embed UnimplementedCertificateServiceServer
// for forward compatibility
type CertificateServiceServer interface {
	// 发送证书同步给某个 Gateway
	SendCertificateToGateway(context.Context, *SubscribeRequest) (*Empty, error)
	// Issuer发送证书给Provider
	SyncCertificateToProvider(context.Context, *CertificateList) (*Empty, error)
	// Gateway订阅
	GatewaySubscribe(*SubscribeRequest, CertificateService_GatewaySubscribeServer) error
	Check(context.Context, *Empty) (*Empty, error)
	mustEmbedUnimplementedCertificateServiceServer()
}

// UnimplementedCertificateServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCertificateServiceServer struct {
}

func (UnimplementedCertificateServiceServer) SendCertificateToGateway(context.Context, *SubscribeRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCertificateToGateway not implemented")
}
func (UnimplementedCertificateServiceServer) SyncCertificateToProvider(context.Context, *CertificateList) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncCertificateToProvider not implemented")
}
func (UnimplementedCertificateServiceServer) GatewaySubscribe(*SubscribeRequest, CertificateService_GatewaySubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method GatewaySubscribe not implemented")
}
func (UnimplementedCertificateServiceServer) Check(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedCertificateServiceServer) mustEmbedUnimplementedCertificateServiceServer() {}

// UnsafeCertificateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CertificateServiceServer will
// result in compilation errors.
type UnsafeCertificateServiceServer interface {
	mustEmbedUnimplementedCertificateServiceServer()
}

func RegisterCertificateServiceServer(s grpc.ServiceRegistrar, srv CertificateServiceServer) {
	s.RegisterService(&CertificateService_ServiceDesc, srv)
}

func _CertificateService_SendCertificateToGateway_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscribeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CertificateServiceServer).SendCertificateToGateway(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.CertificateService/SendCertificateToGateway",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CertificateServiceServer).SendCertificateToGateway(ctx, req.(*SubscribeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CertificateService_SyncCertificateToProvider_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CertificateList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CertificateServiceServer).SyncCertificateToProvider(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.CertificateService/SyncCertificateToProvider",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CertificateServiceServer).SyncCertificateToProvider(ctx, req.(*CertificateList))
	}
	return interceptor(ctx, in, info, handler)
}

func _CertificateService_GatewaySubscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CertificateServiceServer).GatewaySubscribe(m, &certificateServiceGatewaySubscribeServer{stream})
}

type CertificateService_GatewaySubscribeServer interface {
	Send(*CertificateList) error
	grpc.ServerStream
}

type certificateServiceGatewaySubscribeServer struct {
	grpc.ServerStream
}

func (x *certificateServiceGatewaySubscribeServer) Send(m *CertificateList) error {
	return x.ServerStream.SendMsg(m)
}

func _CertificateService_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CertificateServiceServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.CertificateService/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CertificateServiceServer).Check(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// CertificateService_ServiceDesc is the grpc.ServiceDesc for CertificateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CertificateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.CertificateService",
	HandlerType: (*CertificateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendCertificateToGateway",
			Handler:    _CertificateService_SendCertificateToGateway_Handler,
		},
		{
			MethodName: "SyncCertificateToProvider",
			Handler:    _CertificateService_SyncCertificateToProvider_Handler,
		},
		{
			MethodName: "Check",
			Handler:    _CertificateService_Check_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GatewaySubscribe",
			Handler:       _CertificateService_GatewaySubscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pubsub.proto",
}
