package cert

import (
	"net/http"

	"cert-gateway/issuer/internal/logic/cert"
	"cert-gateway/issuer/internal/svc"
	"cert-gateway/issuer/internal/types"
	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"
)

func AddCertFormUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddCertFormUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		if err := validator.New().Struct(&req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := cert.NewAddCertFormUploadLogic(r.Context(), svcCtx)
		resp, err := l.AddCertFormUpload(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
