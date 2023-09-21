package gateway

import (
	"fmt"
	"github.com/pygzfei/issuer-gateway/gateway/internal/cache"
	"github.com/pygzfei/issuer-gateway/gateway/internal/config"
	"github.com/pygzfei/issuer-gateway/gateway/internal/handler"
	"github.com/pygzfei/issuer-gateway/gateway/internal/syncx"
	"github.com/pygzfei/issuer-gateway/utils"
	"log"
	"time"
)

func Run(confPath string) {

	utils.MustLoad(&confPath, config.C)

	cache.Init(config.C)

	go syncx.Init(config.C)

	fmt.Println("HTTP server listening on :80")
	handler.Http(config.C)

	fmt.Println("HTTPS server listening on :443")
	if err := utils.GraceFul(time.Minute, handler.Https()).ListenAndServeTLS("", ""); err != nil {
		log.Fatalln(err)
	}
}
