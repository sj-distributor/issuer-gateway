package grpcServer

import (
	"github.com/pygzfei/issuer-gateway/grpc/config"
	"github.com/pygzfei/issuer-gateway/pkg/driver"
	"github.com/pygzfei/issuer-gateway/pkg/logger"
	"github.com/pygzfei/issuer-gateway/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

func Run(confPath string) {
	var c conf.Config
	utils.MustLoad(&confPath, &c)
	if strings.ToUpper(c.Sync.Target) == "GRPC" {
		logger.Init(logx.LogConf{
			Level:       c.Logger.Level,
			Mode:        c.Logger.Mode,
			Path:        c.Logger.Path,
			KeepDays:    c.Logger.KeepDays,
			MaxSize:     c.Logger.MaxSize,
			ServiceName: "GrpcServer",
		})
		driver.NewGrpcServiceAndListen(c.Sync.GrpcServer.Port)
	}
}
