package cert

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pygzfei/issuer-gateway/issuer/internal/logic/cert"
	"github.com/pygzfei/issuer-gateway/issuer/internal/svc"
	"github.com/pygzfei/issuer-gateway/issuer/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"
)

func GetCertsPagingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCertsPagingReq
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		if err := validator.New().Struct(&req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := cert.NewGetCertsPagingLogic(r.Context(), svcCtx)
		resp, err := l.GetCertsPaging(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
