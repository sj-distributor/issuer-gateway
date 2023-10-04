package gateway

import (
	"fmt"
	"github.com/pygzfei/issuer-gateway/gateway/internal/cache"
	"github.com/pygzfei/issuer-gateway/gateway/internal/config"
	"github.com/pygzfei/issuer-gateway/gateway/internal/handler"
	"github.com/pygzfei/issuer-gateway/gateway/internal/syncx"
	"github.com/pygzfei/issuer-gateway/pkg/logger"
	"github.com/pygzfei/issuer-gateway/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"time"
)

func Run(confPath string) {
	var c = config.C
	utils.MustLoad(&confPath, c)

	logger.Init(logx.LogConf{
		Level:       c.Logger.Level,
		Mode:        c.Logger.Mode,
		Path:        c.Logger.Path,
		KeepDays:    c.Logger.KeepDays,
		MaxSize:     c.Logger.MaxSize,
		ServiceName: "Gateway",
	})

	cache.Init(config.C)

	go syncx.Init(config.C)

	fmt.Println("HTTP server listening on :80")
	handler.Http(config.C)

	fmt.Println("HTTPS server listening on :443")
	if err := utils.GraceFul(time.Minute, handler.Https()).ListenAndServeTLS("", ""); err != nil {
		log.Fatalln(err)
	}
}
