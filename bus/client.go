package main

import (
	"cert-gateway/pkg/driver"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
)

func main() {

	md := metadata.Pairs("Authorization", "Bearer "+"66d2e42661bc292f8237b4736a423a36")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	grpcClient := driver.NewGrpcClient("192.167.167.167:9527", ctx)

	err := grpcClient.Subscribe("192873912309",
		func(msg string) {
			fmt.Println(fmt.Sprintf("receiving msg %s", msg))
		}, func(err error) {

		})

	if err != nil {
		log.Println(err)
	}

	err = grpcClient.Publish("来自client的问候")

	if err != nil {
		log.Println(err)
	}

	//
	//conn, err := grpc.Dial("192.167.167.167:9527", grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//defer conn.Close()
	//
	//client := pb.NewPubSubServiceClient(conn)
	//
	//md := metadata.Pairs("Authorization", "Bearer "+"66d2e42661bc292f8237b4736a423a36")
	//ctx := metadata.NewOutgoingContext(context.Background(), md)
	//
	//// 订阅消息
	//_, err = client.Publish(ctx, &pb.PublishRequest{Message: "我是anson仔123123123"})
	//if err != nil {
	//	log.Fatalf("Subscribe failed: %v", err)
	//}

}
