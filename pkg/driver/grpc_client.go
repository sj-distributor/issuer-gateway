package driver

import (
	"cert-gateway/grpc/pb"
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
)

type GrpcClient struct {
	client pb.CertificateServiceClient
	ctx    context.Context
	conn   *grpc.ClientConn
}

func NewGrpcClient(addr string, ctx context.Context) *GrpcClient {

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &GrpcClient{
		ctx:    ctx,
		client: pb.NewCertificateServiceClient(conn),
		conn:   conn,
	}
}

func (c *GrpcClient) GatewaySubscribe(ip string, onMegReceived OnMessageReceived, onErrReceiving ...OnErrReceiving) error {
	stream, err := c.client.GatewaySubscribe(c.ctx, &pb.SubscribeRequest{LocalIp: ip})
	if err != nil {
		log.Fatalf("Subscribe failed: %v", err)
		return err
	}

	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				if len(onErrReceiving) > 0 {
					onErrReceiving[0](err)
				}
			}
			log.Printf("Received message: %v\n", message.Certs)
			onMegReceived(message.Certs)
		}
	}()

	return nil
}

func (c *GrpcClient) SendCertificateToGateway(localId string) error {
	_, err := c.client.SendCertificateToGateway(c.ctx, &pb.SubscribeRequest{LocalIp: localId})
	return err
}

func (c *GrpcClient) SyncCertificateToProvider(certificateList *pb.CertificateList) error {
	_, err := c.client.SyncCertificateToProvider(c.ctx, certificateList)
	return err
}
