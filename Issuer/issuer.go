package issuer

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	xhttp "github.com/zeromicro/x/http"
	"issuer-gateway/issuer/internal/config"
	"issuer-gateway/issuer/internal/database"
	"issuer-gateway/issuer/internal/errs"
	"issuer-gateway/issuer/internal/handler"
	"issuer-gateway/issuer/internal/svc"
	"issuer-gateway/issuer/middleware"
	"issuer-gateway/pkg/acme"
	"net/http"
)

func Run(conPath string) {

	var c config.Config
	conf.MustLoad(conPath, &c)

	database.Init(&c)

	server := rest.MustNewServer(c.Issuer.RestConf,
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

	fmt.Printf("Starting server at %s:%d...\n", c.Issuer.Host, c.Issuer.Port)

	server.Start()
}
