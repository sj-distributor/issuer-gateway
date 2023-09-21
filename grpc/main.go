package grpcServer

import (
	"cert-gateway/grpc/config"
	"cert-gateway/pkg/driver"
	"cert-gateway/utils"
	"strings"
)

func Run(confPath string) {
	var c conf.Config
	utils.MustLoad(&confPath, &c)
	if strings.ToUpper(c.Sync.Target) == "GRPC" {
		driver.NewGrpcServiceAndListen(c.Sync.Grpc.Addr)
	}
}
