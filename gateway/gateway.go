package gateway

import (
	"fmt"
	"issuer-gateway/gateway/internal/cache"
	"issuer-gateway/gateway/internal/config"
	"issuer-gateway/gateway/internal/handler"
	"issuer-gateway/gateway/internal/syncx"
	"issuer-gateway/utils"
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
