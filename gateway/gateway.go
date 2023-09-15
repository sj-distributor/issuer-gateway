package main

import (
	"cert-gateway/gateway/internal/cache"
	"cert-gateway/gateway/internal/config"
	"cert-gateway/gateway/internal/handler"
	"cert-gateway/gateway/internal/syncx"
	"cert-gateway/utils"
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {

	var configFile = flag.String("f", "etc/config.yaml", "the config file")

	utils.MustLoad(configFile, config.C)
	cache.Init(config.C)

	go syncx.Init(config.C)

	fmt.Println("HTTP server listening on :80")
	handler.Http(config.C)

	fmt.Println("HTTPS server listening on :443")
	if err := utils.GraceFul(time.Minute, handler.Https()).ListenAndServeTLS("", ""); err != nil {
		log.Fatalln(err)
	}
}
