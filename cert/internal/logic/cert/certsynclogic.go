package cert

import (
	"cert-gateway/cert/pkg/cache"
	"context"

	"cert-gateway/cert/internal/svc"
	"cert-gateway/cert/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CertSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	cache  *cache.MemoryCache
}

func NewCertSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CertSyncLogic {
	return &CertSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		cache:  cache.CertCache,
	}
}

func (l *CertSyncLogic) CertSync(req *types.CertSyncReq) (resp *types.CertSyncResp, err error) {
	changed := l.cache.CaptureChanged(req.Maximum)
	return &types.CertSyncResp{
		Certs: *changed,
	}, nil
}
