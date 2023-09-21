package gateway

import (
	"cert-gateway/gateway/internal/cache"
	"cert-gateway/gateway/internal/config"
	"cert-gateway/gateway/internal/handler"
	"cert-gateway/gateway/internal/syncx"
	"cert-gateway/utils"
	"fmt"
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
