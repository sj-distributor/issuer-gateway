package grpcServer

import (
	"issuer-gateway/grpc/config"
	"issuer-gateway/pkg/driver"
	"issuer-gateway/utils"
	"strings"
)

func Run(confPath string) {
	var c conf.Config
	utils.MustLoad(&confPath, &c)
	if strings.ToUpper(c.Sync.Target) == "GRPC" {
		driver.NewGrpcServiceAndListen(c.Sync.GrpcServer.Port)
	}
}
