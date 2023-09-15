package driver

import (
	"cert-gateway/bus/pb"
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
)

type GrpcClient struct {
	client pb.PubSubServiceClient
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
		client: pb.NewPubSubServiceClient(conn),
		conn:   conn,
	}
}

func (c *GrpcClient) Subscribe(ip string, onMegReceived OnMessageReceived, onErrReceiving ...OnErrReceiving) error {
	stream, err := c.client.Subscribe(c.ctx, &pb.SubscribeRequest{SubscriberId: ip})
	if err != nil {
		log.Fatalf("Subscribe failed: %v", err)
		return err
	}

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
		log.Printf("Received message: %s\n", message.Message)
		onMegReceived(message.Message)
	}

	return nil
}

func (c *GrpcClient) Publish(msg string) error {
	_, err := c.client.Publish(c.ctx, &pb.PublishRequest{Message: msg})
	return err
}
