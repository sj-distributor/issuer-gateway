package middleware

import (
	"github.com/pygzfei/issuer-gateway/issuer/internal/svc"
	"github.com/pygzfei/issuer-gateway/pkg/acme"
	"github.com/pygzfei/issuer-gateway/pkg/acme/providers"
	"github.com/pygzfei/issuer-gateway/pkg/errs"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
	"strings"
)

func AcceptAcmeChallenge(serverCtx *svc.ServiceContext) rest.Route {

	provider := serverCtx.AcmeProvider.(*acme.AcmeProvider)

	return rest.Route{
		Method: http.MethodGet,
		Path:   "/.well-known/acme-challenge/:token",
		Handler: func(w http.ResponseWriter, r *http.Request) {

			token := strings.Replace(r.URL.Path, "/.well-known/acme-challenge/", "", -1)

			originalDomain := r.Header.Get("X-Forwarded-Host")

			logx.Info(
				logx.Field("Token", token),
				logx.Field("X-Forwarded-Host", originalDomain),
				logx.Field("RemoteAddr", r.RemoteAddr),
				logx.Field("Method", r.Method),
				logx.Field("Host", r.Host),
				logx.Field("URL", r.URL),
				logx.Field("Header", r.Header),
			)

			value, ok := provider.MemoryCache.Load(originalDomain)

			logx.Info(
				logx.Field("value", value),
				logx.Field("ok", ok),
			)

			if ok {
				domainMapping := value.(providers.DomainMapping)

				if domainMapping.Token == token {

					w.Header().Set("Content-Type", "text/plain")
					// keyAuth
					_, err := w.Write([]byte(domainMapping.KeyAuth))
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					// domain
					logx.Infof("[%s] Served key authentication", originalDomain)
					return
				}

			}
			http.Error(w, errs.NotFoundException.Error(), http.StatusNotFound)
		},
	}
}
