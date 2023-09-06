package main

import (
	"cert-gateway/cert/internal/config"
	"cert-gateway/cert/internal/database"
	"cert-gateway/cert/internal/errs"
	"cert-gateway/cert/internal/handler"
	"cert-gateway/cert/internal/svc"
	"cert-gateway/cert/middleware"
	"cert-gateway/cert/pkg/cache"
	"cert-gateway/pkg/acme"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	xhttp "github.com/zeromicro/x/http"
	"log"
	"net/http"
)

var configFile = flag.String("f", "etc/cert-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	database.Init(&c)

	if err := cache.Init(database.DB(), &c); err != nil {
		log.Fatalln(err)
	}

	server := rest.MustNewServer(c.RestConf,
		middleware.Cors(),
		rest.WithNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			xhttp.JsonBaseResponseCtx(r.Context(), w, errs.NotFoundException)
		})),
	)

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
