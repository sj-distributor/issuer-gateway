package main

import (
	"cert-gateway/gateway/internal/cache"
	"cert-gateway/gateway/internal/configs"
	"cert-gateway/gateway/internal/handler"
	"cert-gateway/utils"
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {

	var configFile = flag.String("f", "internal/configs/config.yaml", "the config file")

	utils.MustLoad(configFile, configs.C)

	cache.Init(configs.C)

	fmt.Println(configs.C)

	fmt.Println("HTTP server listening on :80")
	handler.Http()

	fmt.Println("HTTPS server listening on :443")
	if err := utils.GraceFul(time.Minute, handler.Https()).ListenAndServeTLS("", ""); err != nil {
		log.Fatalln(err)
	}
}
