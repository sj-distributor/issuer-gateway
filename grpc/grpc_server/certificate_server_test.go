package grpc_server

import (
	"context"
	"fmt"
	"github.com/go-playground/assert/v2"
	conf "github.com/pygzfei/issuer-gateway/grpc/config"
	"github.com/pygzfei/issuer-gateway/grpc/pb"
	"github.com/pygzfei/issuer-gateway/pkg/errs"
	"github.com/pygzfei/issuer-gateway/utils"
	"github.com/zeromicro/x/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"reflect"
	"sync"
	"testing"
	"time"
)

func wrapServerStream(ctx context.Context) grpc.ServerStream {
	return &mockServerStream{
		ctx: ctx,
	}
}

type mockServerStream struct {
	ctx context.Context
	grpc.ServerStream
}

func (m *mockServerStream) SetHeader(md metadata.MD) error {
	return nil
}

func (m *mockServerStream) SendHeader(metadata.MD) error {
	return nil
}

func (m *mockServerStream) SetTrailer(md metadata.MD) {
	return
}

func (m *mockServerStream) Context() context.Context {
	return m.ctx
}

func (m *mockServerStream) SendMsg(i any) error {
	return nil
}

func (m *mockServerStream) RecvMsg(i any) error {
	return nil
}

func (m *mockServerStream) Send(*pb.CertificateList) error {
	return nil
}

func TestCertificatePubSubServer_Check(t *testing.T) {
	type fields struct {
		cache                    *MemoryCache
		mu                       sync.Mutex
		gateways                 map[string]pb.CertificateService_GatewaySubscribeServer
		CertificateServiceServer pb.CertificateServiceServer
	}
	type args struct {
		in0 context.Context
		in1 *pb.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "run health check success",
			fields: struct {
				cache                    *MemoryCache
				mu                       sync.Mutex
				gateways                 map[string]pb.CertificateService_GatewaySubscribeServer
				CertificateServiceServer pb.CertificateServiceServer
			}{cache: MewMemoryCache(), mu: sync.Mutex{}, gateways: make(map[string]pb.CertificateService_GatewaySubscribeServer), CertificateServiceServer: NewCertificatePubSubServer()},
			args: struct {
				in0 context.Context
				in1 *pb.Empty
			}{in0: context.Background(), in1: &pb.Empty{}},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CertificatePubSubServer{
				cache:                    tt.fields.cache,
				mu:                       tt.fields.mu,
				gateways:                 tt.fields.gateways,
				CertificateServiceServer: tt.fields.CertificateServiceServer,
			}
			_, err := s.Check(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("Check() err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestCertificatePubSubServer_GatewaySubscribe(t *testing.T) {

	ctx, cancelFunc := context.WithCancel(context.Background())

	type fields struct {
		cache                    *MemoryCache
		mu                       sync.Mutex
		gateways                 map[string]pb.CertificateService_GatewaySubscribeServer
		CertificateServiceServer pb.CertificateServiceServer
	}
	type args struct {
		req    *pb.SubscribeRequest
		stream pb.CertificateService_GatewaySubscribeServer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Run gateway subscribe success",
			fields: struct {
				cache                    *MemoryCache
				mu                       sync.Mutex
				gateways                 map[string]pb.CertificateService_GatewaySubscribeServer
				CertificateServiceServer pb.CertificateServiceServer
			}{
				cache:                    MewMemoryCache(),
				mu:                       sync.Mutex{},
				gateways:                 map[string]pb.CertificateService_GatewaySubscribeServer{},
				CertificateServiceServer: NewCertificatePubSubServer(),
			},
			args: struct {
				req    *pb.SubscribeRequest
				stream pb.CertificateService_GatewaySubscribeServer
			}{req: &pb.SubscribeRequest{LocalIp: "123"}, stream: &mockServerStream{ctx: ctx}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CertificatePubSubServer{
				cache:                    tt.fields.cache,
				mu:                       tt.fields.mu,
				gateways:                 tt.fields.gateways,
				CertificateServiceServer: tt.fields.CertificateServiceServer,
			}
			s.mu.Lock()
			defer s.mu.Unlock()

			go func() {
				if err := s.GatewaySubscribe(tt.args.req, tt.args.stream); (err != nil) != tt.wantErr {
					t.Errorf("GatewaySubscribe() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			time.Sleep(time.Second * time.Duration(2))

			total := 0
			for _, _ = range tt.fields.gateways {
				total += 1
			}

			assert.Equal(t, total, 1)
			cancelFunc()

		})
	}
}

func TestCertificatePubSubServer_SendCertificateToGateway(t *testing.T) {

	type fields struct {
		cache                    *MemoryCache
		mu                       sync.Mutex
		gateways                 map[string]pb.CertificateService_GatewaySubscribeServer
		CertificateServiceServer pb.CertificateServiceServer
	}
	type args struct {
		in0 context.Context
		req *pb.SubscribeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "run sendCertificateToGateway not found gateway",
			fields: struct {
				cache                    *MemoryCache
				mu                       sync.Mutex
				gateways                 map[string]pb.CertificateService_GatewaySubscribeServer
				CertificateServiceServer pb.CertificateServiceServer
			}{
				cache:                    MewMemoryCache(),
				mu:                       sync.Mutex{},
				gateways:                 map[string]pb.CertificateService_GatewaySubscribeServer{},
				CertificateServiceServer: NewCertificatePubSubServer(),
			},
			args: struct {
				in0 context.Context
				req *pb.SubscribeRequest
			}{in0: context.Background(), req: &pb.SubscribeRequest{LocalIp: "123"}},
			want:    errors.New(404, fmt.Sprintf("gateway ip: %s,  not found", "123")),
			wantErr: true,
		},
		{
			name: "run sendCertificateToGateway success",
			fields: struct {
				cache                    *MemoryCache
				mu                       sync.Mutex
				gateways                 map[string]pb.CertificateService_GatewaySubscribeServer
				CertificateServiceServer pb.CertificateServiceServer
			}{
				cache: MewMemoryCache(),
				mu:    sync.Mutex{},
				gateways: func() map[string]pb.CertificateService_GatewaySubscribeServer {
					gateways := map[string]pb.CertificateService_GatewaySubscribeServer{}
					gateways["123"] = &mockServerStream{}
					return gateways
				}(),
				CertificateServiceServer: NewCertificatePubSubServer(),
			},
			args: struct {
				in0 context.Context
				req *pb.SubscribeRequest
			}{in0: context.Background(), req: &pb.SubscribeRequest{LocalIp: "123"}},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CertificatePubSubServer{
				cache:                    tt.fields.cache,
				mu:                       tt.fields.mu,
				gateways:                 tt.fields.gateways,
				CertificateServiceServer: tt.fields.CertificateServiceServer,
			}
			_, err := s.SendCertificateToGateway(tt.args.in0, tt.args.req)

			if (err != nil) && !tt.wantErr {
				t.Errorf("SendCertificateToGateway() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("SendCertificateToGateway() err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestCertificatePubSubServer_SyncCertificateToProvider(t *testing.T) {

	cache := MewMemoryCache()

	type fields struct {
		cache                    *MemoryCache
		mu                       sync.Mutex
		gateways                 map[string]pb.CertificateService_GatewaySubscribeServer
		CertificateServiceServer pb.CertificateServiceServer
	}
	type args struct {
		in0   context.Context
		certs *pb.CertificateList
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "run syncCertificateToProvider success",
			fields: struct {
				cache                    *MemoryCache
				mu                       sync.Mutex
				gateways                 map[string]pb.CertificateService_GatewaySubscribeServer
				CertificateServiceServer pb.CertificateServiceServer
			}{
				cache:                    cache,
				mu:                       sync.Mutex{},
				gateways:                 make(map[string]pb.CertificateService_GatewaySubscribeServer),
				CertificateServiceServer: NewCertificatePubSubServer(),
			},
			args: struct {
				in0   context.Context
				certs *pb.CertificateList
			}{in0: context.Background(), certs: &pb.CertificateList{Certs: []*pb.Cert{
				{Id: uint64(utils.Id()), Domain: "test.anson.com", Certificate: "", PrivateKey: "", IssuerCertificate: "", Target: "https://anson.test.com"},
			}}},
			want:    1,
			wantErr: false,
		},
		{
			name: "run syncCertificateToProvider error",
			fields: struct {
				cache                    *MemoryCache
				mu                       sync.Mutex
				gateways                 map[string]pb.CertificateService_GatewaySubscribeServer
				CertificateServiceServer pb.CertificateServiceServer
			}{
				cache:                    cache,
				mu:                       sync.Mutex{},
				gateways:                 make(map[string]pb.CertificateService_GatewaySubscribeServer),
				CertificateServiceServer: NewCertificatePubSubServer(),
			},
			args: struct {
				in0   context.Context
				certs *pb.CertificateList
			}{in0: context.Background(), certs: &pb.CertificateList{Certs: []*pb.Cert{
				{Id: uint64(utils.Id()), Domain: "test.anson.com", Certificate: "", PrivateKey: "", IssuerCertificate: "", Target: "https://anson.test.com"},
			}}},
			want:    2,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CertificatePubSubServer{
				cache:                    tt.fields.cache,
				mu:                       tt.fields.mu,
				gateways:                 tt.fields.gateways,
				CertificateServiceServer: tt.fields.CertificateServiceServer,
			}
			_, err := s.SyncCertificateToProvider(tt.args.in0, tt.args.certs)

			if err != nil {
				t.Errorf("SyncCertificateToProvider() error = %v", err)
				return
			}

			if len(tt.fields.cache.GetAll().Certs) != tt.want && !tt.wantErr {
				t.Errorf("The number of certificates must be %d, result: [%d]", tt.want, len(tt.fields.cache.GetAll().Certs))
				return
			}
		})
	}
}

func TestNewCertificatePubSubServer(t *testing.T) {
	tests := []struct {
		name string
		want *CertificatePubSubServer
	}{
		{
			name: "can new CertificatePubSubServer",
			want: NewCertificatePubSubServer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCertificatePubSubServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCertificatePubSubServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamInterceptor(t *testing.T) {
	type args struct {
		conf *conf.Config
		key  string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "run StreamInterceptor success",
			args: struct {
				conf *conf.Config
				key  string
			}{
				conf: &conf.Config{
					Secret: "123121",
					Logger: struct {
						Level    string
						Mode     string
						Path     string
						KeepDays int
						MaxSize  int
					}{
						Level: "debug", Mode: "console", Path: "", KeepDays: 0, MaxSize: 0,
					},
					Sync: struct {
						Target     string
						GrpcServer struct{ Port string }
					}{
						Target: "",
						GrpcServer: struct{ Port string }{
							Port: "10086",
						},
					},
				}, key: "Authorization"},
			want: nil,
		},
		{
			name: "run StreamInterceptor 404",
			args: struct {
				conf *conf.Config
				key  string
			}{
				conf: &conf.Config{
					Secret: "123121",
					Logger: struct {
						Level    string
						Mode     string
						Path     string
						KeepDays int
						MaxSize  int
					}{
						Level: "debug", Mode: "console", Path: "", KeepDays: 0, MaxSize: 0,
					},
					Sync: struct {
						Target     string
						GrpcServer struct{ Port string }
					}{
						Target: "",
						GrpcServer: struct{ Port string }{
							Port: "10086",
						},
					},
				}, key: "errorKey"},
			want: errs.UnAuthorizationException,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StreamInterceptor(tt.args.conf)

			ctx := metadata.NewIncomingContext(
				context.Background(),
				metadata.Pairs(tt.args.key, tt.args.conf.Secret))

			err := got("anson", wrapServerStream(ctx), &grpc.StreamServerInfo{}, func(srv any, stream grpc.ServerStream) error {
				return nil
			})

			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("err: [%s], tt.want: [%s]", err, tt.want)
			}
		})
	}
}

func TestUnaryInterceptor(t *testing.T) {
	type args struct {
		conf *conf.Config
		key  string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "run UnaryInterceptor success",
			args: struct {
				conf *conf.Config
				key  string
			}{
				conf: &conf.Config{
					Secret: "123121",
					Logger: struct {
						Level    string
						Mode     string
						Path     string
						KeepDays int
						MaxSize  int
					}{
						Level: "debug", Mode: "console", Path: "", KeepDays: 0, MaxSize: 0,
					},
					Sync: struct {
						Target     string
						GrpcServer struct{ Port string }
					}{
						Target: "",
						GrpcServer: struct{ Port string }{
							Port: "10086",
						},
					},
				},
				key: "Authorization",
			},
			want: nil,
		},
		{
			name: "run UnaryInterceptor 404",
			args: struct {
				conf *conf.Config
				key  string
			}{
				conf: &conf.Config{
					Secret: "123121",
					Logger: struct {
						Level    string
						Mode     string
						Path     string
						KeepDays int
						MaxSize  int
					}{
						Level: "debug", Mode: "console", Path: "", KeepDays: 0, MaxSize: 0,
					},
					Sync: struct {
						Target     string
						GrpcServer struct{ Port string }
					}{
						Target: "",
						GrpcServer: struct{ Port string }{
							Port: "10086",
						},
					},
				},
				key: "errorKey",
			},
			want: errs.UnAuthorizationException,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UnaryInterceptor(tt.args.conf)

			ctx := metadata.NewIncomingContext(
				context.Background(),
				metadata.Pairs(tt.args.key, tt.args.conf.Secret))

			_, err := got(ctx, "anson", &grpc.UnaryServerInfo{}, func(ctx context.Context, req any) (any, error) {
				return req, nil
			})

			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("err: [%s], tt.want: [%s]", err, tt.want)
			}
		})
	}
}

func Test_getTokenFromCtx(t *testing.T) {

	str := "123456"

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		wantToken *string
		wantErr   bool
	}{
		{
			name: "can get token",
			args: struct{ ctx context.Context }{ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.Pairs("Authorization", "123456"))},
			wantToken: &str,
			wantErr:   false,
		},
		{
			name: "can not get token and throw 404",
			args: struct{ ctx context.Context }{ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.Pairs("error", "123456"))},
			wantToken: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := getTokenFromCtx(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTokenFromCtx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !reflect.DeepEqual(err, errs.UnAuthorizationException) {
				t.Errorf("err: [%s], wantErr: [%s]", err, errs.UnAuthorizationException)
				return
			}
			if !reflect.DeepEqual(gotToken, tt.wantToken) {
				t.Errorf("getTokenFromCtx() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}
