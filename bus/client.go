package main

import (
	"cert-gateway/bus/pb"
	"cert-gateway/pkg/driver"
	"context"
	"google.golang.org/grpc/metadata"
	"log"
)

func main() {

	md := metadata.Pairs("Authorization", "Bearer "+"66d2e42661bc292f8237b4736a423a36")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	grpcClient := driver.NewGrpcClient("192.168.101.2:9527", ctx)

	var certs []*pb.Cert
	certs = append(certs, &pb.Cert{
		Id:     123123,
		Domain: "absob.123.com",
	})
	err := grpcClient.SyncCertificateToProvider(&pb.CertificateList{Certs: certs})
	if err != nil {
		log.Println(err)
	}

	err = grpcClient.SendCertificateToGateway("10086")

	if err != nil {
		log.Println(err)
	}
}
