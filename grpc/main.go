package grpcServer

import (
	"github.com/pygzfei/issuer-gateway/grpc/config"
	"github.com/pygzfei/issuer-gateway/pkg/driver"
	"github.com/pygzfei/issuer-gateway/pkg/logger"
	"github.com/pygzfei/issuer-gateway/utils"
	"strings"
)

func Run(confPath string) {
	var c conf.Config
	utils.MustLoad(&confPath, &c)
	if strings.ToUpper(c.Sync.Target) == "GRPC" {
		logger.Init(c.Env)
		driver.NewGrpcServiceAndListen(c.Sync.GrpcServer.Port)
	}
}
