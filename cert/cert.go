package main

import (
	"cert-gateway/cert/internal/config"
	"cert-gateway/cert/internal/database"
	"cert-gateway/cert/internal/handler"
	"cert-gateway/cert/internal/svc"
	"cert-gateway/cert/pkg/acme"
	"cert-gateway/cert/pkg/cache"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

var configFile = flag.String("f", "etc/cert-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	database.Init(&c)

	cache.Init(database.DB(), &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/.well-known/acme-challenge/:token",
		Handler: acme.AcceptChallenge(),
	})

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
