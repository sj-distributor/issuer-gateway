package issuer

import (
	"github.com/pygzfei/issuer-gateway/issuer/internal/config"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database"
	"github.com/pygzfei/issuer-gateway/issuer/internal/handler"
	"github.com/pygzfei/issuer-gateway/issuer/internal/scheduler"
	"github.com/pygzfei/issuer-gateway/issuer/internal/svc"
	"github.com/pygzfei/issuer-gateway/issuer/middleware"
	"github.com/pygzfei/issuer-gateway/pkg/errs"
	"github.com/pygzfei/issuer-gateway/pkg/logger"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	xhttp "github.com/zeromicro/x/http"
	"net/http"
)

func Run(conPath string) {

	var c config.Config
	conf.MustLoad(conPath, &c)

	logger.Init(logx.LogConf{
		Level:       c.Logger.Level,
		Mode:        c.Logger.Mode,
		Path:        c.Logger.Path,
		KeepDays:    c.Logger.KeepDays,
		MaxSize:     c.Logger.MaxSize,
		ServiceName: "Issuer",
	})

	database.Init(&c)

	server := rest.MustNewServer(c.Issuer.RestConf,
		middleware.Cors(),
		rest.WithNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			xhttp.JsonBaseResponseCtx(r.Context(), w, errs.NotFoundException)
		})),
	)

	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	server.AddRoute(middleware.AcceptAcmeChallenge(ctx))

	scheduler.NewScheduler(&c, ctx.SyncProvider, ctx.AcmeProvider)

	logx.Infof("Starting server at %s:%d...", c.Issuer.Host, c.Issuer.Port)

	server.Start()
}
